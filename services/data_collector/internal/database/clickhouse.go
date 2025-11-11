package database

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync/atomic"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/narroworb/data_collector/internal/collector"
	"github.com/narroworb/data_collector/internal/models"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClichouseDB struct {
	conn  driver.Conn
	maxID struct {
		athlete                 uint32
		manager                 uint32
		team                    uint32
		tournament              uint32
		match                   uint32
		footballPlayerMatchStat uint32
		footballGoalieMatchStat uint32
		footballTeamMatchStat   uint32
		teamPerformance         uint32
	}
}

func NewClickhouseDB() (*ClichouseDB, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: "sports_data",
			Username: "default",
			Password: "",
		},
		Settings: map[string]interface{}{
			"max_execution_time": 60,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error in connect to Clickhouse: %v", err)
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("error in Ping to Clickhouse: %v", err)
	}
	c := &ClichouseDB{
		conn: conn,
	}
	if err := c.initCounters(); err != nil {
		return nil, fmt.Errorf("error in init counters: %v", err)
	}

	return c, nil
}

func (c *ClichouseDB) initCounters() error {
	queries := map[*uint32]string{
		&c.maxID.manager:                 "SELECT ifNull(max(manager_id), 0) FROM Managers",
		&c.maxID.team:                    "SELECT ifNull(max(team_id), 0) FROM Teams",
		&c.maxID.match:                   "SELECT ifNull(max(match_id), 0) FROM Matches",
		&c.maxID.athlete:                 "SELECT ifNull(max(athlete_id), 0) FROM Athletes",
		&c.maxID.tournament:              "SELECT ifNull(max(tournament_id), 0) FROM Tournaments",
		&c.maxID.footballPlayerMatchStat: "SELECT ifNull(max(stat_id), 0) FROM Football_Player_Match_Stats",
		&c.maxID.footballGoalieMatchStat: "SELECT ifNull(max(stat_id), 0) FROM Football_Goalie_Match_Stats",
		&c.maxID.footballTeamMatchStat:   "SELECT ifNull(max(stat_id), 0) FROM Football_Team_Match_Stats",
		&c.maxID.teamPerformance:         "SELECT ifNull(max(performance_id), 0) FROM Team_Tournament_Performances",
	}

	for ptr, q := range queries {
		var val uint32
		if err := c.conn.QueryRow(context.Background(), q).Scan(&val); err != nil {
			return err
		}
		atomic.StoreUint32(ptr, val)
	}
	return nil
}

func (c *ClichouseDB) Close() {
	c.conn.Close()
}

func (c *ClichouseDB) GetUnactualTournamentsAndTours(ctx context.Context) ([]collector.UnactualTournamentsAndTours, error) {
	rows, err := c.conn.Query(ctx, "SELECT DISTINCT t.name, t.season, m.round FROM Tournaments t INNER JOIN Matches m ON t.tournament_id=m.tournament_id WHERE m.status='Not started' AND date < NOW() - INTERVAL 3 HOURS ORDER BY t.name, t.season, m.round;")
	if err != nil {
		return nil, fmt.Errorf("error in query to db: %+v", err)
	}

	defer rows.Close()
	res := make([]collector.UnactualTournamentsAndTours, 0, 10)

	for rows.Next() {
		var leagueName, season string
		var tour uint16

		if err := rows.Scan(&leagueName, &season, &tour); err != nil {
			return nil, fmt.Errorf("error in scan query: %+v", err)
		}

		if len(res) == 0 || (res[len(res)-1].LeagueName != leagueName || res[len(res)-1].Season != season) {
			res = append(res, collector.UnactualTournamentsAndTours{
				LeagueName: leagueName,
				Season:     season,
				Tours:      make([]uint16, 0),
			})
		}

		res[len(res)-1].Tours = append(res[len(res)-1].Tours, tour)
	}

	return res, nil
}

func (c *ClichouseDB) GetFootballTournamentID(ctx context.Context, name, season string) (uint32, error) {
	row := c.conn.QueryRow(ctx,
		`SELECT tournament_id FROM Tournaments
		WHERE sport_id=1 AND season=$1 AND name=$2`,
		season, name)
	var id uint32
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (c *ClichouseDB) GetFootballTeamID(ctx context.Context, name string) (uint32, error) {
	row := c.conn.QueryRow(ctx,
		`SELECT team_id FROM Teams
		WHERE name=$1`,
		name)
	var id uint32
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (c *ClichouseDB) GetFootballManagerID(ctx context.Context, manager *models.Manager) (uint32, error) {
	row := c.conn.QueryRow(ctx,
		`SELECT manager_id FROM Managers
		WHERE first_name=$1 AND last_name=$2`,
		manager.FirstName, manager.LastName)
	var id uint32
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	manager.ID = id
	return id, nil
}

func (c *ClichouseDB) InsertFootballManager(ctx context.Context, manager *models.Manager) (uint32, error) {
	id := c.nextFootballManagerID()
	err := c.conn.Exec(ctx,
		`INSERT INTO Managers (manager_id, first_name, last_name, nation, yellow_cards, red_cards)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		id, manager.FirstName, manager.LastName, manager.Nation, manager.YellowCards, manager.RedCards)
	if err != nil {
		return 0, err
	}
	manager.ID = id
	return id, nil
}

func (c *ClichouseDB) nextFootballManagerID() uint32 {
	return atomic.AddUint32(&c.maxID.manager, 1)
}

func (c *ClichouseDB) GetFootballMatchID(ctx context.Context, match *models.Match, tournamentID uint32) (uint32, error) {
	row := c.conn.QueryRow(ctx,
		`SELECT match_id FROM Matches
		WHERE tournament_id=$1 AND date=$2 AND home_team_id=$3 AND away_team_id=$4 AND round=$5`,
		tournamentID, match.Date, match.HomeTeam.ID, match.AwayTeam.ID, match.Round)
	var id uint32
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	match.IDAppDB = id
	return id, nil
}

func (c *ClichouseDB) GetFootballNotPlayedMatchID(ctx context.Context, match *models.Match, tournamentID uint32) (uint32, error) {
	row := c.conn.QueryRow(ctx,
		`SELECT match_id FROM Matches
		WHERE tournament_id=$1 AND status='Not started' AND home_team_id=$2 AND away_team_id=$3 AND round=$4`,
		tournamentID, match.HomeTeam.ID, match.AwayTeam.ID, match.Round)
	var id uint32
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	match.IDAppDB = id
	return id, nil
}

func (c *ClichouseDB) GetFootballPlayedMatchID(ctx context.Context, match *models.Match, tournamentID uint32) (uint32, error) {
	row := c.conn.QueryRow(ctx,
		`SELECT match_id FROM Matches
		WHERE tournament_id=$1 AND status='Ended' AND home_team_id=$2 AND away_team_id=$3 AND round=$4`,
		tournamentID, match.HomeTeam.ID, match.AwayTeam.ID, match.Round)
	var id uint32
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	match.IDAppDB = id
	return id, nil
}

func (c *ClichouseDB) GetFootballMatchStatus(ctx context.Context, matchID uint32) (string, error) {
	row := c.conn.QueryRow(ctx,
		`SELECT status FROM Matches
		WHERE match_id=$1`,
		matchID)
	var status string
	if err := row.Scan(&status); err != nil {
		return "", err
	}
	return status, nil
}

func (c *ClichouseDB) UpdateFootballMatch(ctx context.Context, match *models.Match) error {
	err := c.conn.Exec(ctx,
		`ALTER TABLE Matches UPDATE home_score=$2, away_score=$3, home_team_manager_id=$4, away_team_manager_id=$5, status=$6, date=$7 WHERE match_id=$1`,
		match.IDAppDB, match.HomeGoals, match.AwayGoals, match.HomeManager.ID, match.AwayManager.ID, match.Status, match.Date)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClichouseDB) InsertFootballMatch(ctx context.Context, match *models.Match, tournamentID uint32) (uint32, error) {
	id := c.nextFootballMatchID()

	err := c.conn.Exec(ctx,
		`INSERT INTO Matches (match_id, tournament_id, date, home_team_id, away_team_id, home_score, away_score, round, home_team_manager_id, away_team_manager_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		id, tournamentID, match.Date, match.HomeTeam.ID, match.AwayTeam.ID, match.HomeGoals, match.AwayGoals, match.Round, match.HomeManager.ID, match.AwayManager.ID, match.Status)
	if err != nil {
		return 0, err
	}
	match.IDAppDB = id
	return id, nil
}

func (c *ClichouseDB) nextFootballMatchID() uint32 {
	return atomic.AddUint32(&c.maxID.match, 1)
}

func (c *ClichouseDB) GetFootballPlayerID(ctx context.Context, name string, position string, height uint16) (uint32, error) {
	firstName := strings.Split(name, " ")[0]
	lastName := strings.TrimSpace(strings.Join(strings.Split(name, " ")[1:], " "))
	if lastName == "" {
		lastName = firstName
		firstName = ""
	}
	row := c.conn.QueryRow(ctx,
		`SELECT athlete_id FROM Athletes
		WHERE sport_id=1 AND first_name=$1 AND last_name=$2 AND position=$3 AND height=$4`,
		firstName, lastName, position, height)
	var id uint32
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (c *ClichouseDB) InsertFootballPlayer(ctx context.Context, player *models.Player) (uint32, error) {
	id := c.nextFootballPlayerID()

	err := c.conn.Exec(ctx,
		`INSERT INTO Athletes (athlete_id, sport_id, first_name, last_name, position, date_of_birth, height, preferred_foot_or_handness, nation, current_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		id, 1, player.FirstName, player.LastName, player.Position, player.DateOfBirth, player.Height, player.PreferredFoot, player.Nation, player.CurrentStatus)
	if err != nil {
		return 0, err
	}
	player.ID = id
	return id, nil
}

func (c *ClichouseDB) nextFootballPlayerID() uint32 {
	return atomic.AddUint32(&c.maxID.athlete, 1)
}

func (c *ClichouseDB) IncrementYellowCardsManager(ctx context.Context, managerID uint32) error {
	// row := c.conn.QueryRow(ctx,
	// 	`SELECT yellow_cards FROM Managers
	// 	WHERE manager_id=$1`,
	// 	managerID)
	// var yellow_cards uint16
	// if err := row.Scan(&yellow_cards); err != nil {
	// 	return err
	// }
	// yellow_cards++
	// err := c.conn.Exec(ctx,
	// 	`ALTER TABLE Managers UPDATE yellow_cards=$1 WHERE manager_id=$2`,
	// 	yellow_cards, managerID)
	// if err != nil {
	// 	return err
	// }

	err := c.conn.Exec(ctx,
		`ALTER TABLE Managers UPDATE yellow_cards=yellow_cards+1 WHERE manager_id=$1`,
		managerID)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClichouseDB) IncrementRedCardsManager(ctx context.Context, managerID uint32) error {
	// row := c.conn.QueryRow(ctx,
	// 	`SELECT red_cards FROM Managers
	// 	WHERE manager_id=$1`,
	// 	managerID)
	// var red_cards int
	// if err := row.Scan(&red_cards); err != nil {
	// 	return err
	// }
	// red_cards++
	// err := c.conn.Exec(ctx,
	// 	`ALTER TABLE Managers UPDATE red_cards=$1 WHERE manager_id=$2`,
	// 	red_cards, managerID)
	// if err != nil {
	// 	return err
	// }

	err := c.conn.Exec(ctx,
		`ALTER TABLE Managers UPDATE red_cards=red_cards+1 WHERE manager_id=$1`,
		managerID)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClichouseDB) InsertFootballMatchStats(ctx context.Context, stats models.TeamMatchStats, matchID uint32) (uint32, error) {
	id := c.nextFootballTeamMatchStatID()

	err := c.conn.Exec(ctx,
		`INSERT INTO Football_Team_Match_Stats (stat_id, match_id, shots_on_goal_home, shots_on_goal_away, total_shots_home, total_shots_away, blocked_shots_home, blocked_shots_away, fouls_home, fouls_away,
		corner_kicks_home, corner_kicks_away, ball_possession_home, ball_possession_away, yellow_cards_home, yellow_cards_away, red_cards_home, red_cards_away, total_passes_home, total_passes_away,
		complete_passes_home, complete_passes_away, offsides_home, offsides_away, shots_inside_box_home, shots_inside_box_away)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)`,
		id, matchID, stats.ShotsOnGoalHome, stats.ShotsOnGoalAway, stats.TotalShotsHome, stats.TotalShotsAway, stats.BlockedShotsHome, stats.BlockedShotsAway, stats.FoulsHome, stats.FoulsAway,
		stats.CornerKicksHome, stats.CornerKicksAway, stats.BallPossessionHome, stats.BallPossessionAway, stats.YellowCardsHome, stats.YellowCardsAway, stats.RedCardsHome, stats.RedCardsAway,
		stats.TotalPassesHome, stats.TotalPassesAway, stats.CompletePassesHome, stats.CompletePassesAway, stats.OffsidesHome, stats.OffsidesAway, stats.ShotsInsideBoxHome, stats.ShotsInsideBoxAway)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *ClichouseDB) nextFootballTeamMatchStatID() uint32 {
	return atomic.AddUint32(&c.maxID.footballTeamMatchStat, 1)
}

func (c *ClichouseDB) GetFootballMatchStats(ctx context.Context, stats models.TeamMatchStats, matchID uint32) (uint32, error) {
	row := c.conn.QueryRow(ctx,
		`SELECT stat_id FROM Football_Team_Match_Stats WHERE match_id=$1`,
		matchID)
	var id uint32
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (c *ClichouseDB) InsertFootballPlayerMatchStats(ctx context.Context, stats models.PlayerStatsInMatch, matchID uint32) (uint32, error) {
	var startPlayer, captain, homeTeamPlayer uint8
	if stats.StartPlayer {
		startPlayer = 1
	}
	if stats.Captain {
		captain = 1
	}
	if stats.HomeTeamPlayer {
		homeTeamPlayer = 1
	}

	id := c.nextFootballPlayerMatchStatID()

	err := c.conn.Exec(ctx,
		`INSERT INTO Football_Player_Match_Stats (match_id, athlete_id, start_player, rating, minutes_played, goals, assists, blocked_shots, interceptions, total_tackles,
		dribbled_past, duels, duels_won, fouls, was_fouled, pass_attempts, complete_passes, key_passes, shot_on_target, total_shots, dribble_attempts, complete_dribbles,
		penalty_scored, penalty_missed, yellow_cards, red_cards, captain, home_team_player)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29)`,
		id, matchID, stats.IDPlayer, startPlayer, fmt.Sprintf("%.1f", stats.Rating), stats.MinutesPlayed, stats.Goals, stats.Assists, stats.BlockedShots, stats.Interceptions,
		stats.TotalTackles, stats.DribbledPast, stats.Duels, stats.DuelsWon, stats.Fouls, stats.WasFouled, stats.PassAttempts, stats.CompletePasses,
		stats.KeyPasses, stats.ShotsOnTarget, stats.TotalShots, stats.DribbleAttempts, stats.CompleteDribbles, stats.PenaltyScored, stats.PenaltyMissed,
		stats.YellowCards, stats.RedCards, captain, homeTeamPlayer)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *ClichouseDB) nextFootballPlayerMatchStatID() uint32 {
	return atomic.AddUint32(&c.maxID.footballPlayerMatchStat, 1)
}

func (c *ClichouseDB) GetFootballPlayerMatchStats(ctx context.Context, stats models.PlayerStatsInMatch, matchID uint32) (uint32, error) {
	row := c.conn.QueryRow(ctx,
		`SELECT stat_id FROM Football_Player_Match_Stats WHERE match_id=$1 AND athlete_id=$2`,
		matchID, stats.IDPlayer)
	var id uint32
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (c *ClichouseDB) InsertFootballGoalieMatchStats(ctx context.Context, stats models.GoalieStatsInMatch, matchID uint32) (uint32, error) {
	var startPlayer, captain, homeTeamPlayer uint8
	if stats.StartPlayer {
		startPlayer = 1
	}
	if stats.Captain {
		captain = 1
	}
	if stats.HomeTeamPlayer {
		homeTeamPlayer = 1
	}

	id := c.nextFootballGoalieMatchStatID()

	err := c.conn.Exec(ctx,
		`INSERT INTO Football_Goalie_Match_Stats (match_id, athlete_id, start_player, rating, minutes_played, goals, assists, goals_conceded, saves, pass_attempts,
		complete_passes, key_passes, penalty_saved, penalty_conceded, fouls, was_fouled, yellow_cards, red_cards, captain, home_team_player)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)`,
		id, matchID, stats.IDPlayer, startPlayer, fmt.Sprintf("%.1f", stats.Rating), stats.MinutesPlayed, stats.Goals, stats.Assists, stats.GoalsConceded, stats.Saves,
		stats.PassAttempts, stats.CompletePasses, stats.KeyPasses, stats.PenaltySaved, stats.PenaltyConceded, stats.Fouls, stats.WasFouled, stats.YellowCards,
		stats.RedCards, captain, homeTeamPlayer)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *ClichouseDB) nextFootballGoalieMatchStatID() uint32 {
	return atomic.AddUint32(&c.maxID.footballGoalieMatchStat, 1)
}

func (c *ClichouseDB) GetFootballGoalieMatchStats(ctx context.Context, stats models.GoalieStatsInMatch, matchID uint32) (uint32, error) {
	row := c.conn.QueryRow(ctx,
		`SELECT stat_id FROM Football_Goalie_Match_Stats WHERE match_id=$1 AND athlete_id=$2`,
		matchID, stats.IDPlayer)
	var id uint32
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (c *ClichouseDB) InsertFootballGoalieMatchStatsBatch(ctx context.Context, statsBatch map[uint32]*models.GoalieStatsInMatch, matchID uint32) error {
	if len(statsBatch) == 0 {
		return fmt.Errorf("empty batch")
	}

	batchInsert, err := c.conn.PrepareBatch(ctx, `INSERT INTO Football_Goalie_Match_Stats (stat_id, match_id, athlete_id, start_player, rating, minutes_played, goals, assists, goals_conceded, saves,
	pass_attempts, complete_passes, key_passes, penalty_saved, penalty_conceded, fouls, was_fouled, yellow_cards, red_cards, captain, home_team_player)`)
	if err != nil {
		return fmt.Errorf("ошибка при подготовке батча: %v", err)
	}

	for _, stats := range statsBatch {
		if stats.IDPlayer < 1 {
			log.Printf("try inserting goalie stats with ID=0. matchID: %d, stats: %+v", matchID, *stats)
			continue
		}

		var startPlayer, captain, homeTeamPlayer uint8
		if stats.StartPlayer {
			startPlayer = 1
		}
		if stats.Captain {
			captain = 1
		}
		if stats.HomeTeamPlayer {
			homeTeamPlayer = 1
		}

		if err := batchInsert.Append(c.nextFootballGoalieMatchStatID(), matchID, stats.IDPlayer, startPlayer, fmt.Sprintf("%.1f", stats.Rating), stats.MinutesPlayed, stats.Goals, stats.Assists, stats.GoalsConceded, stats.Saves,
			stats.PassAttempts, stats.CompletePasses, stats.KeyPasses, stats.PenaltySaved, stats.PenaltyConceded, stats.Fouls, stats.WasFouled, stats.YellowCards,
			stats.RedCards, captain, homeTeamPlayer); err != nil {
			return fmt.Errorf("ошибка при добавлении в батч: %v", err)
		}
	}

	if err := batchInsert.Send(); err != nil {
		return err
	}
	return nil
}

func (c *ClichouseDB) InsertFootballPlayerMatchStatsBatch(ctx context.Context, statsBatch map[uint32]*models.PlayerStatsInMatch, matchID uint32) error {
	if len(statsBatch) == 0 {
		return fmt.Errorf("empty batch")
	}

	batchInsert, err := c.conn.PrepareBatch(ctx, `INSERT INTO Football_Player_Match_Stats (stat_id, match_id, athlete_id, start_player, rating, minutes_played, goals, assists, blocked_shots, interceptions, total_tackles,
		dribbled_past, duels, duels_won, fouls, was_fouled, pass_attempts, complete_passes, key_passes, shot_on_target, total_shots, dribble_attempts, complete_dribbles,
		penalty_scored, penalty_missed, yellow_cards, red_cards, captain, home_team_player)`)
	if err != nil {
		return fmt.Errorf("ошибка при подготовке батча: %v", err)
	}

	for _, stats := range statsBatch {
		if stats.IDPlayer < 1 {
			log.Printf("try inserting player stats with ID=0. matchID: %d, stats: %+v", matchID, *stats)
			continue
		}

		var startPlayer, captain, homeTeamPlayer uint8
		if stats.StartPlayer {
			startPlayer = 1
		}
		if stats.Captain {
			captain = 1
		}
		if stats.HomeTeamPlayer {
			homeTeamPlayer = 1
		}

		if err := batchInsert.Append(c.nextFootballPlayerMatchStatID(), matchID, stats.IDPlayer, startPlayer, fmt.Sprintf("%.1f", stats.Rating), stats.MinutesPlayed, stats.Goals, stats.Assists, stats.BlockedShots, stats.Interceptions,
			stats.TotalTackles, stats.DribbledPast, stats.Duels, stats.DuelsWon, stats.Fouls, stats.WasFouled, stats.PassAttempts, stats.CompletePasses,
			stats.KeyPasses, stats.ShotsOnTarget, stats.TotalShots, stats.DribbleAttempts, stats.CompleteDribbles, stats.PenaltyScored, stats.PenaltyMissed,
			stats.YellowCards, stats.RedCards, captain, homeTeamPlayer); err != nil {
			return fmt.Errorf("ошибка при добавлении в батч: %v", err)
		}
	}

	if err := batchInsert.Send(); err != nil {
		return err
	}
	return nil
}

func (c *ClichouseDB) GetCountPlayersStatsByMatchID(ctx context.Context, matchID uint32) (uint64, error) {
	row := c.conn.QueryRow(ctx,
		`SELECT COUNT(*) FROM Football_Player_Match_Stats
		WHERE match_id=$1`,
		matchID)
	var cnt uint64
	if err := row.Scan(&cnt); err != nil {
		return 0, err
	}
	return cnt, nil
}
