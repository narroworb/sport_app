package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/narroworb/core_api/internal/models"
)

var ErrProtectedTournament = errors.New("protected tournament cannot be modified by users")

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
		teamPerformance         uint32
	}
}

func NewClickhouseDB() (*ClichouseDB, error) {
	dbAddr := os.Getenv("DATABASE_ADDR")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{dbAddr},
		Auth: clickhouse.Auth{
			Database: dbName,
			Username: dbUser,
			Password: dbPass,
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
		return nil, fmt.Errorf("error initializing id counters: %v", err)
	}

	return c, nil
}

func (c *ClichouseDB) Close() {
	c.conn.Close()
}

func (c *ClichouseDB) initCounters() error {
	queries := map[*uint32]string{
		&c.maxID.athlete:                 "SELECT ifNull(max(athlete_id), 0) FROM Athletes",
		&c.maxID.manager:                 "SELECT ifNull(max(manager_id), 0) FROM Managers",
		&c.maxID.team:                    "SELECT ifNull(max(team_id), 0) FROM Teams",
		&c.maxID.tournament:              "SELECT ifNull(max(tournament_id), 0) FROM Tournaments",
		&c.maxID.match:                   "SELECT ifNull(max(match_id), 0) FROM Matches",
		&c.maxID.footballPlayerMatchStat: "SELECT ifNull(max(stat_id), 0) FROM Football_Player_Match_Stats",
		&c.maxID.footballGoalieMatchStat: "SELECT ifNull(max(stat_id), 0) FROM Football_Goalie_Match_Stats",
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

func (c *ClichouseDB) NextAthleteID() uint32 {
	return atomic.AddUint32(&c.maxID.athlete, 1)
}

func (c *ClichouseDB) NextManagerID() uint32 {
	return atomic.AddUint32(&c.maxID.manager, 1)
}

func (c *ClichouseDB) NextTeamID() uint32 {
	return atomic.AddUint32(&c.maxID.team, 1)
}

func (c *ClichouseDB) NextTournamentID() uint32 {
	return atomic.AddUint32(&c.maxID.tournament, 1)
}

func (c *ClichouseDB) NextMatchID() uint32 {
	return atomic.AddUint32(&c.maxID.match, 1)
}

func (c *ClichouseDB) NextFootballPlayerMatchStatID() uint32 {
	return atomic.AddUint32(&c.maxID.footballPlayerMatchStat, 1)
}

func (c *ClichouseDB) NextFootballGoalieMatchStatID() uint32 {
	return atomic.AddUint32(&c.maxID.footballGoalieMatchStat, 1)
}

func (c *ClichouseDB) NextTeamTournamentPerformanceID() uint32 {
	return atomic.AddUint32(&c.maxID.teamPerformance, 1)
}

func (c *ClichouseDB) GetPlayerByID(ctx context.Context, id uint32) (models.Player, error) {
	row := c.conn.QueryRow(ctx, "SELECT first_name, last_name, position, date_of_birth, height, preferred_foot_or_handness, a.nation, current_status, photo_url, flag_url FROM Athletes a INNER JOIN Athlete_Photos ap ON a.athlete_id=ap.athlete_id INNER JOIN National_Flags nf ON a.nation=nf.nation WHERE a.athlete_id=$1", id)

	var player models.Player
	if err := row.Scan(&player.FirstName, &player.LastName, &player.Position, &player.DateOfBirth, &player.Height, &player.PreferredFoot, &player.Nation.Name, &player.CurrentStatus, &player.URLPhoto, &player.Nation.URLFlag); err != nil {
		return models.Player{}, err
	}
	player.ID = id

	return player, nil
}

func (c *ClichouseDB) GetPlayerPositionByID(ctx context.Context, id uint32) (string, error) {
	row := c.conn.QueryRow(ctx, "SELECT position FROM Athletes WHERE athlete_id=$1", id)

	var pos string
	if err := row.Scan(&pos); err != nil {
		return "", err
	}

	return pos, nil
}

func (c *ClichouseDB) GetManagerByID(ctx context.Context, id uint32) (models.Manager, error) {
	row := c.conn.QueryRow(ctx, "SELECT first_name, last_name, m.nation, yellow_cards, red_cards, photo_url, flag_url FROM Managers m INNER JOIN Manager_Photos mp ON m.manager_id=mp.manager_id INNER JOIN National_Flags nf ON m.nation=nf.nation WHERE m.manager_id=$1", id)

	var manager models.Manager
	if err := row.Scan(&manager.FirstName, &manager.LastName, &manager.Nation.Name, &manager.YellowCards, &manager.RedCards, &manager.URLPhoto, &manager.Nation.URLFlag); err != nil {
		return models.Manager{}, err
	}
	manager.ID = id

	return manager, nil
}

func (c *ClichouseDB) GetTeamByID(ctx context.Context, id uint32) (models.Team, error) {
	row := c.conn.QueryRow(ctx, "SELECT name, logo_url FROM Teams t INNER JOIN Team_Logos tl ON t.team_id=tl.team_id WHERE t.team_id=$1", id)

	var team models.Team
	if err := row.Scan(&team.Name, &team.URLLogo); err != nil {
		return models.Team{}, err
	}
	team.ID = id

	return team, nil
}

func (c *ClichouseDB) GetMatchByID(ctx context.Context, id uint32) (models.Match, error) {
	row := c.conn.QueryRow(ctx, "SELECT Matches.tournament_id, date, home_team_id, away_team_id, home_score, away_score, round, home_team_manager_id, away_team_manager_id, status, logo_url FROM Matches INNER JOIN Tournaments ON Matches.tournament_id = Tournaments.tournament_id INNER JOIN Tournament_Logos ON Tournaments.name = Tournament_Logos.tournament_name  WHERE match_id=$1", id)

	var match models.Match
	if err := row.Scan(&match.Tournament.ID, &match.Date, &match.HomeTeam.ID, &match.AwayTeam.ID, &match.HomeGoals, &match.AwayGoals, &match.Round, &match.HomeManager.ID, &match.AwayManager.ID, &match.Status, &match.Tournament.URLLogo); err != nil {
		return models.Match{}, err
	}
	match.ID = id

	return match, nil
}

func (c *ClichouseDB) GetTournamentByID(ctx context.Context, id uint32) (models.Tournament, error) {
	row := c.conn.QueryRow(ctx, "SELECT name, t.country, flag_url, season, logo_url FROM Tournaments t INNER JOIN National_Flags nf ON t.country=nf.nation INNER JOIN Tournament_Logos tl ON t.name=tl.tournament_name WHERE tournament_id=$1", id)

	var tournament models.Tournament
	if err := row.Scan(&tournament.Name, &tournament.Country.Name, &tournament.Country.URLFlag, &tournament.Season, &tournament.URLLogo); err != nil {
		return models.Tournament{}, err
	}
	tournament.ID = id

	return tournament, nil
}

func (c *ClichouseDB) GetAllTournaments(ctx context.Context) (map[string]models.TournamentWithSeason, error) {
	rows, err := c.conn.Query(ctx, "SELECT t.tournament_id, t.name, season, logo_url FROM Tournaments t LEFT JOIN Tournament_Logos tl ON t.name=tl.tournament_name ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tournaments := make(map[string]models.TournamentWithSeason)
	for rows.Next() {
		var name, season string
		var logo sql.NullString
		var id uint32
		if err := rows.Scan(&id, &name, &season, &logo); err != nil {
			return map[string]models.TournamentWithSeason{}, err
		}

		if _, exists := tournaments[name]; !exists {
			tournaments[name] = models.TournamentWithSeason{
				Name: name,
				Seasons: []struct {
					Season string "json:\"season,omitempty\""
					ID     uint32 "json:\"tournament_id\""
				}{},
				URLLogo: logo.String,
			}
		}
		entry := tournaments[name]
		entry.Seasons = append(entry.Seasons, struct {
			Season string "json:\"season,omitempty\""
			ID     uint32 "json:\"tournament_id\""
		}{Season: season, ID: id})
		tournaments[name] = entry
	}

	return tournaments, nil
}

func (c *ClichouseDB) GetPlayerStatsBySeason(ctx context.Context, filter models.PlayerStatsFilter) (models.PlayerStatsInPeriod, error) {
	row := c.conn.QueryRow(ctx, `
	SELECT 
		COUNT(*) AS total_matches,
		SUM(start_player),
		COUNTIf(minutes_played != 0) as played_matches,
		AVGIf(rating, rating != 0),
		SUM(minutes_played),
		SUM(goals),
		SUM(assists),
		SUM(blocked_shots),
		SUM(interceptions),
		SUM(total_tackles),
		SUM(dribbled_past),
		SUM(duels),
		SUM(duels_won),
		SUM(fouls),
		SUM(was_fouled),
		SUM(pass_attempts),
		SUM(complete_passes),
		SUM(key_passes),
		SUM(shot_on_target),
		SUM(total_shots),
		SUM(dribble_attempts),
		SUM(complete_dribbles),
		SUM(penalty_scored),
		SUM(penalty_missed),
		SUM(yellow_cards),
		SUM(red_cards),
		SUM(captain),
		FROM Football_Player_Match_Stats s 
	INNER JOIN Matches m ON s.match_id=m.match_id 
	INNER JOIN Tournaments t ON t.tournament_id=m.tournament_id 
	WHERE athlete_id=$1 AND season=$2`,
		filter.PlayerID, filter.Season)

	var stats models.PlayerStatsInPeriod
	if err := row.Scan(&stats.TotalMatches, &stats.StartPlayer, &stats.MatchesPlayed, &stats.AvgRating, &stats.MinutesPlayed, &stats.Goals, &stats.Assists, &stats.BlockedShots, &stats.Interceptions, &stats.TotalTackles, &stats.DribbledPast, &stats.Duels,
		&stats.DuelsWon, &stats.Fouls, &stats.WasFouled, &stats.PassAttempts, &stats.CompletePasses, &stats.KeyPasses, &stats.ShotsOnTarget, &stats.TotalShots, &stats.DribbleAttempts, &stats.CompleteDribbles, &stats.PenaltyScored, &stats.PenaltyMissed,
		&stats.YellowCards, &stats.RedCards, &stats.CaptainTimes); err != nil {
		return models.PlayerStatsInPeriod{}, err
	}
	stats.IDPlayer = filter.PlayerID

	return stats, nil
}

func (c *ClichouseDB) GetPlayerStatsByDates(ctx context.Context, filter models.PlayerStatsFilter) (models.PlayerStatsInPeriod, error) {
	row := c.conn.QueryRow(ctx, `
	SELECT 
		COUNT(*) AS total_matches,
		SUM(start_player),
		COUNTIf(minutes_played != 0) as played_matches,
		AVGIf(rating, rating != 0),
		SUM(minutes_played),
		SUM(goals),
		SUM(assists),
		SUM(blocked_shots),
		SUM(interceptions),
		SUM(total_tackles),
		SUM(dribbled_past),
		SUM(duels),
		SUM(duels_won),
		SUM(fouls),
		SUM(was_fouled),
		SUM(pass_attempts),
		SUM(complete_passes),
		SUM(key_passes),
		SUM(shot_on_target),
		SUM(total_shots),
		SUM(dribble_attempts),
		SUM(complete_dribbles),
		SUM(penalty_scored),
		SUM(penalty_missed),
		SUM(yellow_cards),
		SUM(red_cards),
		SUM(captain),
		FROM Football_Player_Match_Stats s 
	INNER JOIN Matches m ON s.match_id=m.match_id 
	WHERE athlete_id=$1 AND date >= $2 AND date <= $3`,
		filter.PlayerID, filter.FromDate, filter.ToDate)

	var stats models.PlayerStatsInPeriod
	if err := row.Scan(&stats.TotalMatches, &stats.StartPlayer, &stats.MatchesPlayed, &stats.AvgRating, &stats.MinutesPlayed, &stats.Goals, &stats.Assists, &stats.BlockedShots, &stats.Interceptions, &stats.TotalTackles, &stats.DribbledPast, &stats.Duels,
		&stats.DuelsWon, &stats.Fouls, &stats.WasFouled, &stats.PassAttempts, &stats.CompletePasses, &stats.KeyPasses, &stats.ShotsOnTarget, &stats.TotalShots, &stats.DribbleAttempts, &stats.CompleteDribbles, &stats.PenaltyScored, &stats.PenaltyMissed,
		&stats.YellowCards, &stats.RedCards, &stats.CaptainTimes); err != nil {
		return models.PlayerStatsInPeriod{}, err
	}
	stats.IDPlayer = filter.PlayerID

	return stats, nil
}

func (c *ClichouseDB) GetPlayerFullStats(ctx context.Context, filter models.PlayerStatsFilter) (models.PlayerStatsInPeriod, error) {
	row := c.conn.QueryRow(ctx, `
	SELECT 
		COUNT(*) AS total_matches,
		SUM(start_player),
		COUNTIf(minutes_played != 0) as played_matches,
		AVGIf(rating, rating != 0),
		SUM(minutes_played),
		SUM(goals),
		SUM(assists),
		SUM(blocked_shots),
		SUM(interceptions),
		SUM(total_tackles),
		SUM(dribbled_past),
		SUM(duels),
		SUM(duels_won),
		SUM(fouls),
		SUM(was_fouled),
		SUM(pass_attempts),
		SUM(complete_passes),
		SUM(key_passes),
		SUM(shot_on_target),
		SUM(total_shots),
		SUM(dribble_attempts),
		SUM(complete_dribbles),
		SUM(penalty_scored),
		SUM(penalty_missed),
		SUM(yellow_cards),
		SUM(red_cards),
		SUM(captain),
		FROM Football_Player_Match_Stats
	WHERE athlete_id=$1`,
		filter.PlayerID)

	var stats models.PlayerStatsInPeriod
	if err := row.Scan(&stats.TotalMatches, &stats.StartPlayer, &stats.MatchesPlayed, &stats.AvgRating, &stats.MinutesPlayed, &stats.Goals, &stats.Assists, &stats.BlockedShots, &stats.Interceptions, &stats.TotalTackles, &stats.DribbledPast, &stats.Duels,
		&stats.DuelsWon, &stats.Fouls, &stats.WasFouled, &stats.PassAttempts, &stats.CompletePasses, &stats.KeyPasses, &stats.ShotsOnTarget, &stats.TotalShots, &stats.DribbleAttempts, &stats.CompleteDribbles, &stats.PenaltyScored, &stats.PenaltyMissed,
		&stats.YellowCards, &stats.RedCards, &stats.CaptainTimes); err != nil {
		return models.PlayerStatsInPeriod{}, err
	}
	stats.IDPlayer = filter.PlayerID

	return stats, nil
}

func (c *ClichouseDB) GetGoalieStatsBySeason(ctx context.Context, filter models.PlayerStatsFilter) (models.PlayerStatsInPeriod, error) {
	row := c.conn.QueryRow(ctx, `
	SELECT 
		COUNT(*) AS total_matches,
		SUM(start_player),
		COUNTIf(minutes_played != 0) as played_matches,
		AVGIf(rating, rating != 0),
		SUM(minutes_played),
		SUM(goals),
		SUM(assists),
		SUM(goals_conceded),
		SUM(saves),
		SUM(pass_attempts),
		SUM(complete_passes),
		SUM(key_passes),
		SUM(penalty_saved),
		SUM(penalty_conceded),
		SUM(fouls),
		SUM(was_fouled),
		SUM(yellow_cards),
		SUM(red_cards),
		SUM(captain),
		FROM Football_Goalie_Match_Stats s 
	INNER JOIN Matches m ON s.match_id=m.match_id 
	INNER JOIN Tournaments t ON t.tournament_id=m.tournament_id 
	WHERE athlete_id=$1 AND season=$2`,
		filter.PlayerID, filter.Season)

	var stats models.PlayerStatsInPeriod
	if err := row.Scan(&stats.TotalMatches, &stats.StartPlayer, &stats.MatchesPlayed, &stats.AvgRating, &stats.MinutesPlayed, &stats.Goals, &stats.Assists, &stats.GoalsConceded, &stats.Saves,
		&stats.PassAttempts, &stats.CompletePasses, &stats.KeyPasses, &stats.PenaltySaved, &stats.PenaltyConceded, &stats.Fouls, &stats.WasFouled, &stats.YellowCards, &stats.RedCards, &stats.CaptainTimes); err != nil {
		return models.PlayerStatsInPeriod{}, err
	}
	stats.IDPlayer = filter.PlayerID

	return stats, nil
}

func (c *ClichouseDB) GetGoalieStatsByDates(ctx context.Context, filter models.PlayerStatsFilter) (models.PlayerStatsInPeriod, error) {
	row := c.conn.QueryRow(ctx, `
	SELECT 
		COUNT(*) AS total_matches,
		SUM(start_player),
		COUNTIf(minutes_played != 0) as played_matches,
		AVGIf(rating, rating != 0),
		SUM(minutes_played),
		SUM(goals),
		SUM(assists),
		SUM(goals_conceded),
		SUM(saves),
		SUM(pass_attempts),
		SUM(complete_passes),
		SUM(key_passes),
		SUM(penalty_saved),
		SUM(penalty_conceded),
		SUM(fouls),
		SUM(was_fouled),
		SUM(yellow_cards),
		SUM(red_cards),
		SUM(captain),
		FROM Football_Goalie_Match_Stats s 
	INNER JOIN Matches m ON s.match_id=m.match_id 
	WHERE athlete_id=$1 AND date >= $2 AND date <= $3`,
		filter.PlayerID, filter.FromDate, filter.ToDate)

	var stats models.PlayerStatsInPeriod
	if err := row.Scan(&stats.TotalMatches, &stats.StartPlayer, &stats.MatchesPlayed, &stats.AvgRating, &stats.MinutesPlayed, &stats.Goals, &stats.Assists, &stats.GoalsConceded, &stats.Saves,
		&stats.PassAttempts, &stats.CompletePasses, &stats.KeyPasses, &stats.PenaltySaved, &stats.PenaltyConceded, &stats.Fouls, &stats.WasFouled, &stats.YellowCards, &stats.RedCards, &stats.CaptainTimes); err != nil {
		return models.PlayerStatsInPeriod{}, err
	}
	stats.IDPlayer = filter.PlayerID

	return stats, nil
}

func (c *ClichouseDB) GetGoalieFullStats(ctx context.Context, filter models.PlayerStatsFilter) (models.PlayerStatsInPeriod, error) {
	row := c.conn.QueryRow(ctx, `
	SELECT 
		COUNT(*) AS total_matches,
		SUM(start_player),
		COUNTIf(minutes_played != 0) as played_matches,
		AVGIf(rating, rating != 0),
		SUM(minutes_played),
		SUM(goals),
		SUM(assists),
		SUM(goals_conceded),
		SUM(saves),
		SUM(pass_attempts),
		SUM(complete_passes),
		SUM(key_passes),
		SUM(penalty_saved),
		SUM(penalty_conceded),
		SUM(fouls),
		SUM(was_fouled),
		SUM(yellow_cards),
		SUM(red_cards),
		SUM(captain),
		FROM Football_Goalie_Match_Stats s
	WHERE athlete_id=$1`,
		filter.PlayerID, filter.FromDate, filter.ToDate)

	var stats models.PlayerStatsInPeriod
	if err := row.Scan(&stats.TotalMatches, &stats.StartPlayer, &stats.MatchesPlayed, &stats.AvgRating, &stats.MinutesPlayed, &stats.Goals, &stats.Assists, &stats.GoalsConceded, &stats.Saves,
		&stats.PassAttempts, &stats.CompletePasses, &stats.KeyPasses, &stats.PenaltySaved, &stats.PenaltyConceded, &stats.Fouls, &stats.WasFouled, &stats.YellowCards, &stats.RedCards, &stats.CaptainTimes); err != nil {
		return models.PlayerStatsInPeriod{}, err
	}
	stats.IDPlayer = filter.PlayerID

	return stats, nil
}

func (c *ClichouseDB) GetPlayerFixtures(ctx context.Context, id, limit, offset uint32, season string) ([]models.PlayerMatch, error) {
	baseQuery := `
    SELECT 
        s.match_id,
        date,
        home_team_id,
        ht.name,
        htl.logo_url,
        away_team_id,
        at.name,
        atl.logo_url,
        home_score,
        away_score,
        round,
        status,
        m.tournament_id,
        t.name,
        goals,
        assists,
        red_cards,
        toFloat64(rating) as rating,
        minutes_played
    FROM Football_Player_Match_Stats s 
    INNER JOIN Matches m ON s.match_id=m.match_id 
    INNER JOIN Teams ht ON m.home_team_id=ht.team_id
    INNER JOIN Team_Logos htl ON m.home_team_id = htl.team_id
    INNER JOIN Teams at ON m.away_team_id=at.team_id
    INNER JOIN Team_Logos atl ON m.away_team_id = atl.team_id
    INNER JOIN Tournaments t ON m.tournament_id=t.tournament_id
    WHERE athlete_id=$1 and season = $2
    
    UNION ALL
    
    SELECT 
        s.match_id,
        date,
        home_team_id,
        ht.name,
        htl.logo_url,
        away_team_id,
        at.name,
        atl.logo_url,
        home_score,
        away_score,
        round,
        status,
        m.tournament_id,
        t.name,
        goals,
        assists,
        red_cards,
        toFloat64(rating) as rating,
        minutes_played
    FROM Football_Goalie_Match_Stats s 
    INNER JOIN Matches m ON s.match_id=m.match_id 
    INNER JOIN Teams ht ON m.home_team_id=ht.team_id
    INNER JOIN Team_Logos htl ON m.home_team_id = htl.team_id
    INNER JOIN Teams at ON m.away_team_id=at.team_id
    INNER JOIN Team_Logos atl ON m.away_team_id = atl.team_id
    INNER JOIN Tournaments t ON m.tournament_id=t.tournament_id
    WHERE athlete_id=$1 and season = $2
    `

	var rows driver.Rows
	var err error

	// Оборачиваем в подзапрос для правильной пагинации
	if limit > 0 {
		rows, err = c.conn.Query(ctx, `
            SELECT * FROM (
                `+baseQuery+`
            ) AS combined
            ORDER BY date DESC
            LIMIT $3 OFFSET $4
        `, id, season, limit, offset)
	} else {
		rows, err = c.conn.Query(ctx, `
            SELECT * FROM (
                `+baseQuery+`
            ) AS combined
            ORDER BY date DESC
        `, id, season)
		limit = 100
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]models.PlayerMatch, 0, limit)

	for rows.Next() {
		var match models.PlayerMatch
		if err := rows.Scan(&match.ID, &match.Date, &match.HomeTeam.ID, &match.HomeTeam.Name, &match.HomeTeam.URLLogo,
			&match.AwayTeam.ID, &match.AwayTeam.Name, &match.AwayTeam.URLLogo,
			&match.HomeGoals, &match.AwayGoals, &match.Round, &match.Status,
			&match.Tournament.ID, &match.Tournament.Name, &match.PlayerGoals, &match.PlayerAssists,
			&match.PlayerRedCards, &match.PlayerRating, &match.PlayerMinutesPlayed); err != nil {
			return nil, err
		}
		match.PlayerID = id
		matches = append(matches, match)
	}

	return matches, nil
}

func parseDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err == nil {
		return parsed, nil
	}
	return time.Parse(time.RFC3339, value)
}

type createPlayerPayload struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Position      string `json:"position"`
	DateOfBirth   string `json:"date_of_birth"`
	Height        uint16 `json:"height"`
	PreferredFoot string `json:"preffered_foot"`
	Nation        struct {
		Name    string `json:"name"`
		URLFlag string `json:"url_flag,omitempty"`
	} `json:"nation"`
	CurrentStatus string `json:"current_status"`
	PhotoURL      string `json:"url_photo,omitempty"`
}

type createMatchPayload struct {
	Date        string            `json:"date"`
	HomeTeam    models.Team       `json:"home_team"`
	AwayTeam    models.Team       `json:"away_team"`
	Tournament  models.Tournament `json:"tournament"`
	HomeGoals   int16             `json:"home_team_score"`
	AwayGoals   int16             `json:"away_team_score"`
	Round       uint16            `json:"round"`
	HomeManager models.Manager    `json:"home_team_manager"`
	AwayManager models.Manager    `json:"away_team_manager"`
	Status      string            `json:"status"`
}

type createMatchStatsPayload struct {
	Players []models.PlayerStatsInMatch `json:"players,omitempty"`
	Goalies []models.GoalieStatsInMatch `json:"goalies,omitempty"`
	Teams   models.TeamMatchStats       `json:"teams,omitempty"`
}

type createTournamentTablePayload struct {
	TournamentID uint32            `json:"tournament_id"`
	Rows         []models.TableRow `json:"rows"`
}

func (c *ClichouseDB) ensureFlag(ctx context.Context, nation, flagURL string) error {
	if nation == "" {
		return nil
	}
	var count uint64
	if err := c.conn.QueryRow(ctx, "SELECT count() FROM National_Flags WHERE nation=$1", nation).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return c.conn.Exec(ctx, `INSERT INTO National_Flags (nation, flag_url) VALUES ($1, $2)`, nation, flagURL)
}

func (c *ClichouseDB) ensureAthletePhoto(ctx context.Context, athleteID uint32, photoURL string) error {
	var count uint64
	if err := c.conn.QueryRow(ctx, "SELECT count() FROM Athlete_Photos WHERE athlete_id=$1", athleteID).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return c.conn.Exec(ctx, `INSERT INTO Athlete_Photos (athlete_id, photo_url) VALUES ($1, $2)`, athleteID, photoURL)
}

func (c *ClichouseDB) ensureManagerPhoto(ctx context.Context, managerID uint32, photoURL string) error {
	var count uint64
	if err := c.conn.QueryRow(ctx, "SELECT count() FROM Manager_Photos WHERE manager_id=$1", managerID).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return c.conn.Exec(ctx, `INSERT INTO Manager_Photos (manager_id, photo_url) VALUES ($1, $2)`, managerID, photoURL)
}

func (c *ClichouseDB) ensureTeamLogo(ctx context.Context, teamID uint32, logoURL string) error {
	var count uint64
	if err := c.conn.QueryRow(ctx, "SELECT count() FROM Team_Logos WHERE team_id=$1", teamID).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return c.conn.Exec(ctx, `INSERT INTO Team_Logos (team_id, logo_url) VALUES ($1, $2)`, teamID, logoURL)
}

func (c *ClichouseDB) ensureTournamentLogo(ctx context.Context, tournamentName, logoURL string) error {
	if tournamentName == "" {
		return nil
	}
	var count uint64
	if err := c.conn.QueryRow(ctx, "SELECT count() FROM Tournament_Logos WHERE tournament_name=$1", tournamentName).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return c.conn.Exec(ctx, `INSERT INTO Tournament_Logos (tournament_name, logo_url) VALUES ($1, $2)`, tournamentName, logoURL)
}

func (c *ClichouseDB) CreatePlayer(ctx context.Context, userID int64, payload []byte) (int64, error) {
	var p createPlayerPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return 0, err
	}
	dateOfBirth, err := parseDate(p.DateOfBirth)
	if err != nil {
		return 0, err
	}
	id := c.NextAthleteID()
	if err := c.conn.Exec(ctx,
		`INSERT INTO Athletes (athlete_id, sport_id, first_name, last_name, position, date_of_birth, height, preferred_foot_or_handness, nation, current_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		id, 1, p.FirstName, p.LastName, p.Position, dateOfBirth, p.Height, p.PreferredFoot, p.Nation.Name, p.CurrentStatus); err != nil {
		return 0, err
	}
	if err := c.ensureFlag(ctx, p.Nation.Name, p.Nation.URLFlag); err != nil {
		return 0, err
	}
	if err := c.ensureAthletePhoto(ctx, id, p.PhotoURL); err != nil {
		return 0, err
	}
	return int64(id), nil
}

func (c *ClichouseDB) CreateTeam(ctx context.Context, userID int64, payload []byte) (int64, error) {
	var t models.Team
	if err := json.Unmarshal(payload, &t); err != nil {
		return 0, err
	}
	id := c.NextTeamID()
	if err := c.conn.Exec(ctx, `INSERT INTO Teams (team_id, name) VALUES ($1, $2)`, id, t.Name); err != nil {
		return 0, err
	}
	if err := c.ensureTeamLogo(ctx, id, t.URLLogo); err != nil {
		return 0, err
	}
	return int64(id), nil
}

func (c *ClichouseDB) CreateTournament(ctx context.Context, userID int64, payload []byte) (int64, error) {
	var t models.Tournament
	if err := json.Unmarshal(payload, &t); err != nil {
		return 0, err
	}
	id := c.NextTournamentID()
	if err := c.conn.Exec(ctx,
		`INSERT INTO Tournaments (tournament_id, sport_id, name, country, season)
		VALUES ($1, $2, $3, $4, $5)`,
		id, 1, t.Name, t.Country.Name, t.Season); err != nil {
		return 0, err
	}
	if err := c.ensureFlag(ctx, t.Country.Name, t.Country.URLFlag); err != nil {
		return 0, err
	}
	if err := c.ensureTournamentLogo(ctx, t.Name, t.URLLogo); err != nil {
		return 0, err
	}
	return int64(id), nil
}

func (c *ClichouseDB) CreateManager(ctx context.Context, userID int64, payload []byte) (int64, error) {
	var m models.Manager
	if err := json.Unmarshal(payload, &m); err != nil {
		return 0, err
	}
	id := c.NextManagerID()
	if err := c.conn.Exec(ctx,
		`INSERT INTO Managers (manager_id, first_name, last_name, nation, yellow_cards, red_cards)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		id, m.FirstName, m.LastName, m.Nation.Name, m.YellowCards, m.RedCards); err != nil {
		return 0, err
	}
	if err := c.ensureFlag(ctx, m.Nation.Name, m.Nation.URLFlag); err != nil {
		return 0, err
	}
	if err := c.ensureManagerPhoto(ctx, id, m.URLPhoto); err != nil {
		return 0, err
	}
	return int64(id), nil
}

func (c *ClichouseDB) CreateFixture(ctx context.Context, userID int64, payload []byte) (int64, error) {
	var m createMatchPayload
	if err := json.Unmarshal(payload, &m); err != nil {
		return 0, err
	}
	matchDate, err := parseDate(m.Date)
	if err != nil {
		return 0, err
	}
	if m.Tournament.ID < 100 {
		return 0, fmt.Errorf("%w: tournament_id %d is protected", ErrProtectedTournament, m.Tournament.ID)
	}
	id := c.NextMatchID()
	if err := c.conn.Exec(ctx,
		`INSERT INTO Matches (match_id, tournament_id, date, home_team_id, away_team_id, home_score, away_score, round, home_team_manager_id, away_team_manager_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		id, m.Tournament.ID, matchDate, m.HomeTeam.ID, m.AwayTeam.ID, m.HomeGoals, m.AwayGoals, m.Round, m.HomeManager.ID, m.AwayManager.ID, m.Status); err != nil {
		return 0, err
	}
	return int64(id), nil
}

func (c *ClichouseDB) CreateMatchStats(ctx context.Context, userID int64, payload []byte) (int64, error) {
	var p createMatchStatsPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return 0, err
	}

	var firstID int64
	for _, stats := range p.Players {
		id := c.NextFootballPlayerMatchStatID()
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
		if err := c.conn.Exec(ctx,
			`INSERT INTO Football_Player_Match_Stats (stat_id, match_id, athlete_id, start_player, rating, minutes_played, goals, assists, blocked_shots, interceptions, total_tackles,
				dribbled_past, duels, duels_won, fouls, was_fouled, pass_attempts, complete_passes, key_passes, shot_on_target, total_shots, dribble_attempts, complete_dribbles,
				penalty_scored, penalty_missed, yellow_cards, red_cards, captain, home_team_player)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29)`,
			id, stats.IDMatch, stats.Player.ID, startPlayer, fmt.Sprintf("%.1f", stats.Rating), stats.MinutesPlayed, stats.Goals, stats.Assists, stats.BlockedShots, stats.Interceptions,
			stats.TotalTackles, stats.DribbledPast, stats.Duels, stats.DuelsWon, stats.Fouls, stats.WasFouled, stats.PassAttempts, stats.CompletePasses,
			stats.KeyPasses, stats.ShotsOnTarget, stats.TotalShots, stats.DribbleAttempts, stats.CompleteDribbles, stats.PenaltyScored, stats.PenaltyMissed,
			stats.YellowCards, stats.RedCards, captain, homeTeamPlayer); err != nil {
			return 0, err
		}
		if firstID == 0 {
			firstID = int64(id)
		}
	}
	for _, stats := range p.Goalies {
		id := c.NextFootballGoalieMatchStatID()
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
		if err := c.conn.Exec(ctx,
			`INSERT INTO Football_Goalie_Match_Stats (stat_id, match_id, athlete_id, start_player, rating, minutes_played, goals, assists, goals_conceded, saves, pass_attempts,
				complete_passes, key_passes, penalty_saved, penalty_conceded, fouls, was_fouled, yellow_cards, red_cards, captain, home_team_player)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`,
			id, stats.IDMatch, stats.Player.ID, startPlayer, fmt.Sprintf("%.1f", stats.Rating), stats.MinutesPlayed, stats.Goals, stats.Assists, stats.GoalsConceded, stats.Saves,
			stats.PassAttempts, stats.CompletePasses, stats.KeyPasses, stats.PenaltySaved, stats.PenaltyConceded, stats.Fouls, stats.WasFouled,
			stats.YellowCards, stats.RedCards, captain, homeTeamPlayer); err != nil {
			return 0, err
		}
		if firstID == 0 {
			firstID = int64(id)
		}
	}
	if p.Teams.IDMatch != 0 {
		if err := c.conn.Exec(ctx,
			`INSERT INTO Football_Team_Match_Stats (match_id, shots_on_goal_home, shots_on_goal_away, total_shots_home, total_shots_away, blocked_shots_home, blocked_shots_away, fouls_home, fouls_away,
				corner_kicks_home, corner_kicks_away, ball_possession_home, ball_possession_away, yellow_cards_home, yellow_cards_away, red_cards_home, red_cards_away,
				total_passes_home, total_passes_away, complete_passes_home, complete_passes_away, offsides_home, offsides_away, shots_inside_box_home, shots_inside_box_away)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)`,
			p.Teams.IDMatch, p.Teams.ShotsOnGoalHome, p.Teams.ShotsOnGoalAway, p.Teams.TotalShotsHome, p.Teams.TotalShotsAway, p.Teams.BlockedShotsHome, p.Teams.BlockedShotsAway,
			p.Teams.FoulsHome, p.Teams.FoulsAway, p.Teams.CornerKicksHome, p.Teams.CornerKicksAway, p.Teams.BallPossessionHome, p.Teams.BallPossessionAway,
			p.Teams.YellowCardsHome, p.Teams.YellowCardsAway, p.Teams.RedCardsHome, p.Teams.RedCardsAway, p.Teams.TotalPassesHome, p.Teams.TotalPassesAway,
			p.Teams.CompletePassesHome, p.Teams.CompletePassesAway, p.Teams.OffsidesHome, p.Teams.OffsidesAway, p.Teams.ShotsInsideBoxHome, p.Teams.ShotsInsideBoxAway); err != nil {
			return 0, err
		}
	}
	return firstID, nil
}

func (c *ClichouseDB) CreateTournamentTable(ctx context.Context, userID int64, payload []byte) (int64, error) {
	var p createTournamentTablePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return 0, err
	}
	if p.TournamentID == 0 {
		return 0, fmt.Errorf("tournament_id is required")
	}
	if p.TournamentID < 100 {
		return 0, fmt.Errorf("%w: tournament_id %d is protected", ErrProtectedTournament, p.TournamentID)
	}
	var firstID int64
	for _, row := range p.Rows {
		if row.Team.ID == 0 {
			continue
		}
		id := c.NextTeamTournamentPerformanceID()
		if err := c.conn.Exec(ctx,
			`INSERT INTO Team_Tournament_Performances (performance_id, tournament_id, team_id, points, position, games_played, wins, draws, losses, goals_scored, goals_conceded)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			id, p.TournamentID, row.Team.ID, row.Points, row.Pos, row.Matches, row.Wins, row.Draws, row.Losses, row.ScoresFor, row.ScoresAgainst); err != nil {
			return 0, err
		}
		if firstID == 0 {
			firstID = int64(id)
		}
	}
	return firstID, nil
}

func (c *ClichouseDB) GetPlayerTeams(ctx context.Context, id uint32) ([]models.TeamSeason, error) {
	rows, err := c.conn.Query(ctx, `
	WITH last_match AS 
	(SELECT f.athlete_id, first_name, last_name, nation, home_team_player, f.match_id AS match_id FROM Football_Player_Match_Stats AS f 
	INNER JOIN Matches AS m ON f.match_id=m.match_id 
	INNER JOIN Athletes AS a ON a.athlete_id=f.athlete_id 
	WHERE f.athlete_id = $1 
	UNION ALL
	SELECT f.athlete_id, first_name, last_name, nation, home_team_player, f.match_id AS match_id FROM Football_Goalie_Match_Stats AS f 
	INNER JOIN Matches AS m ON f.match_id=m.match_id 
	INNER JOIN Athletes AS a ON a.athlete_id=f.athlete_id 
	WHERE f.athlete_id = $1 
	ORDER BY f.match_id DESC) 
	SELECT t.team_id, t.name, logo_url, tr.season FROM Matches m
	INNER JOIN last_match ON last_match.match_id=m.match_id 
	INNER JOIN Teams t ON t.team_id=home_team_id
	INNER JOIN Team_Logos tl ON t.team_id=tl.team_id
	INNER JOIN Tournaments tr ON tr.tournament_id=m.tournament_id
	WHERE home_team_player
	UNION DISTINCT
	SELECT t.team_id, t.name, logo_url, tr.season FROM Matches m
	INNER JOIN last_match ON last_match.match_id=m.match_id 
	INNER JOIN Teams t ON t.team_id=away_team_id
	INNER JOIN Team_Logos tl ON t.team_id=tl.team_id
	INNER JOIN Tournaments tr ON tr.tournament_id=m.tournament_id
	WHERE home_team_player=0;`,
		id)
	if err != nil {
		return nil, err
	}

	teams := make([]models.TeamSeason, 0, 5)

	for rows.Next() {
		var team models.TeamSeason
		if err := rows.Scan(&team.ID, &team.Name, &team.URLLogo, &team.Season); err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	return teams, nil
}

func (c *ClichouseDB) GetTeamStatsBySeason(ctx context.Context, filter models.TeamStatsFilter) (models.TeamStatsInPeriod, error) {
	row := c.conn.QueryRow(ctx, `
	WITH all_matches_stats AS (
		SELECT 
		m.match_id,
		home_team_id,
		away_team_id,
		home_score AS goals,
		away_score AS goals_conceded,
		IF(home_score > away_score, 3, IF(home_score < away_score, 0, 1)) AS points,
		IF(home_score > away_score, 1, 0) AS is_win,
		shots_on_goal_home AS shots_on_goal,
		shots_on_goal_away AS shots_on_goal_conceded,
		total_shots_home AS total_shots,
		total_shots_away AS total_shots_conceded,
		blocked_shots_home AS blocked_shots,
		blocked_shots_away AS blocked_shots_conceded,
		fouls_home AS fouls,
		fouls_away AS fouls_conceded,
		corner_kicks_home AS corner_kicks,
		corner_kicks_away AS corner_kicks_conceded,
		ball_possession_home AS ball_possession,
		yellow_cards_home AS yellow_cards,
		yellow_cards_away AS yellow_cards_conceded,
		red_cards_home AS red_cards,
		red_cards_away AS red_cards_conceded,
		total_passes_home AS total_passes,
		total_passes_away AS total_passes_conceded,
		complete_passes_home AS complete_passes,
		complete_passes_away AS complete_passes_conceded, 
		offsides_home AS offsides,
		offsides_away AS offsides_conceded,
		shots_inside_box_home AS shots_inside_box,
		shots_inside_box_away AS shots_inside_box_conceded
		FROM Matches m 
		INNER JOIN Football_Team_Match_Stats s
		ON m.match_id=s.match_id
		INNER JOIN Tournaments t ON t.tournament_id=m.tournament_id 
		WHERE m.home_team_id=$1 AND status='Ended' and season=$2
		
		UNION ALL
		
		SELECT 
		m.match_id,
		home_team_id,
		away_team_id,
		home_score AS goals_conceded,
		away_score AS goals,
		IF(away_score > home_score, 3, IF(away_score < home_score, 0, 1)) AS points,
		IF(away_score > home_score, 1, 0) AS is_win,
		shots_on_goal_away AS shots_on_goal,
		shots_on_goal_home AS shots_on_goal_conceded,
		total_shots_away AS total_shots,
		total_shots_home AS total_shots_conceded,
		blocked_shots_away AS blocked_shots,
		blocked_shots_home AS blocked_shots_conceded,
		fouls_away AS fouls,
		fouls_home AS fouls_conceded,
		corner_kicks_away AS corner_kicks,
		corner_kicks_home AS corner_kicks_conceded,
		ball_possession_away AS ball_possession,
		yellow_cards_away AS yellow_cards,
		yellow_cards_home AS yellow_cards_conceded,
		red_cards_away AS red_cards,
		red_cards_home AS red_cards_conceded,
		total_passes_away AS total_passes,
		total_passes_home AS total_passes_conceded,
		complete_passes_away AS complete_passes,
		complete_passes_home AS complete_passes_conceded, 
		offsides_away AS offsides,
		offsides_home AS offsides_conceded,
		shots_inside_box_away AS shots_inside_box,
		shots_inside_box_home AS shots_inside_box_conceded
		FROM Matches m 
		INNER JOIN Football_Team_Match_Stats s
		ON m.match_id=s.match_id
		INNER JOIN Tournaments t ON t.tournament_id=m.tournament_id 
		WHERE m.away_team_id=$1 AND status='Ended' and season=$2
	)
	SELECT 
	COUNT(*) AS total_matches,
	AVG(is_win) as average_win_rate,
	AVG(points) as average_points,
	AVG(ball_possession) as average_ball_possession,
	SUM(goals),
	AVG(goals),
	SUM(goals_conceded),
	AVG(goals_conceded),
	SUM(shots_on_goal),
	AVG(shots_on_goal),
	SUM(shots_on_goal_conceded),
	AVG(shots_on_goal_conceded),
	SUM(total_shots),
	AVG(total_shots),
	SUM(total_shots_conceded),
	AVG(total_shots_conceded),
	SUM(shots_inside_box),
	AVG(shots_inside_box),
	SUM(shots_inside_box_conceded),
	AVG(shots_inside_box_conceded),
	SUM(blocked_shots),
	AVG(blocked_shots),
	SUM(fouls),
	AVG(fouls),
	SUM(fouls_conceded),
	AVG(fouls_conceded),
	SUM(corner_kicks),
	AVG(corner_kicks),
	SUM(corner_kicks_conceded),
	AVG(corner_kicks_conceded),
	SUM(yellow_cards),
	AVG(yellow_cards),
	SUM(yellow_cards_conceded),
	AVG(yellow_cards_conceded),
	SUM(red_cards),
	AVG(red_cards),
	SUM(red_cards_conceded),
	AVG(red_cards_conceded),
	SUM(total_passes),
	AVG(total_passes),
	SUM(total_passes_conceded),
	AVG(total_passes_conceded),
	SUM(complete_passes),
	AVG(complete_passes),
	SUM(complete_passes_conceded),
	AVG(complete_passes_conceded),
	if(isNaN(AVG(complete_passes/total_passes)), 0, AVG(complete_passes/total_passes)),
	if(isNaN(AVG(complete_passes_conceded/total_passes_conceded)), 0, AVG(complete_passes_conceded/total_passes_conceded)),
	SUM(offsides),
	AVG(offsides),
	SUM(offsides_conceded),
	AVG(offsides_conceded)
	FROM all_matches_stats`,
		filter.TeamID, filter.Season)

	var stats models.TeamStatsInPeriod
	if err := row.Scan(&stats.TotalMatches, &stats.AvgWinRate, &stats.AvgPoints, &stats.AvgBallPossession, &stats.Goals, &stats.GoalsPer90, &stats.GoalsConceded, &stats.GoalsConcededPer90,
		&stats.ShotsOnGoal, &stats.ShotsOnGoalPer90, &stats.ShotsOnGoalConceded, &stats.ShotsOnGoalConcededPer90, &stats.TotalShots, &stats.TotalShotsPer90,
		&stats.TotalShotsConceded, &stats.TotalShotsConcededPer90, &stats.ShotsInsideBox, &stats.ShotsInsideBoxPer90, &stats.ShotsInsideBoxConceded, &stats.ShotsInsideBoxConcededPer90,
		&stats.BlockedShots, &stats.BlockedShotsPer90, &stats.Fouls, &stats.FoulsPer90, &stats.WasFouled, &stats.WasFouledPer90, &stats.CornerKicks, &stats.CornerKicksPer90,
		&stats.CornerKicksConceded, &stats.CornerKicksConcededPer90, &stats.YellowCards, &stats.YellowCardsPer90, &stats.YellowCardsOpp, &stats.YellowCardsOppPer90,
		&stats.RedCards, &stats.RedCardsPer90, &stats.RedCardsOpp, &stats.RedCardsOppPer90, &stats.PassAttempts, &stats.PassAttemptsPer90, &stats.PassAttemptsOpp, &stats.PassAttemptsOppPer90,
		&stats.CompletePasses, &stats.CompletePassesPer90, &stats.CompletePassesOpp, &stats.CompletePassesOppPer90, &stats.AvgPassAccuracy, &stats.AvgPassAccuracyOpp,
		&stats.Offsides, &stats.OffsidesPer90, &stats.OffsidesOpp, &stats.OffsidesOppPer90); err != nil {
		return models.TeamStatsInPeriod{}, err
	}
	stats.IDTeam = filter.TeamID

	return stats, nil
}

func (c *ClichouseDB) GetTeamStatsByDates(ctx context.Context, filter models.TeamStatsFilter) (models.TeamStatsInPeriod, error) {
	row := c.conn.QueryRow(ctx, `
	WITH all_matches_stats AS (
		SELECT 
		m.match_id,
		home_team_id,
		away_team_id,
		home_score AS goals,
		away_score AS goals_conceded,
		IF(home_score > away_score, 3, IF(home_score < away_score, 0, 1)) AS points,
		IF(home_score > away_score, 1, 0) AS is_win,
		shots_on_goal_home AS shots_on_goal,
		shots_on_goal_away AS shots_on_goal_conceded,
		total_shots_home AS total_shots,
		total_shots_away AS total_shots_conceded,
		blocked_shots_home AS blocked_shots,
		blocked_shots_away AS blocked_shots_conceded,
		fouls_home AS fouls,
		fouls_away AS fouls_conceded,
		corner_kicks_home AS corner_kicks,
		corner_kicks_away AS corner_kicks_conceded,
		ball_possession_home AS ball_possession,
		yellow_cards_home AS yellow_cards,
		yellow_cards_away AS yellow_cards_conceded,
		red_cards_home AS red_cards,
		red_cards_away AS red_cards_conceded,
		total_passes_home AS total_passes,
		total_passes_away AS total_passes_conceded,
		complete_passes_home AS complete_passes,
		complete_passes_away AS complete_passes_conceded, 
		offsides_home AS offsides,
		offsides_away AS offsides_conceded,
		shots_inside_box_home AS shots_inside_box,
		shots_inside_box_away AS shots_inside_box_conceded
		FROM Matches m 
		INNER JOIN Football_Team_Match_Stats s
		ON m.match_id=s.match_id
		WHERE m.home_team_id=$1 AND status='Ended' AND date >= $2 AND date <= $3
		
		UNION ALL
		
		SELECT 
		m.match_id,
		home_team_id,
		away_team_id,
		home_score AS goals_conceded,
		away_score AS goals,
		IF(away_score > home_score, 3, IF(away_score < home_score, 0, 1)) AS points,
		IF(away_score > home_score, 1, 0) AS is_win,
		shots_on_goal_away AS shots_on_goal,
		shots_on_goal_home AS shots_on_goal_conceded,
		total_shots_away AS total_shots,
		total_shots_home AS total_shots_conceded,
		blocked_shots_away AS blocked_shots,
		blocked_shots_home AS blocked_shots_conceded,
		fouls_away AS fouls,
		fouls_home AS fouls_conceded,
		corner_kicks_away AS corner_kicks,
		corner_kicks_home AS corner_kicks_conceded,
		ball_possession_away AS ball_possession,
		yellow_cards_away AS yellow_cards,
		yellow_cards_home AS yellow_cards_conceded,
		red_cards_away AS red_cards,
		red_cards_home AS red_cards_conceded,
		total_passes_away AS total_passes,
		total_passes_home AS total_passes_conceded,
		complete_passes_away AS complete_passes,
		complete_passes_home AS complete_passes_conceded, 
		offsides_away AS offsides,
		offsides_home AS offsides_conceded,
		shots_inside_box_away AS shots_inside_box,
		shots_inside_box_home AS shots_inside_box_conceded
		FROM Matches m 
		INNER JOIN Football_Team_Match_Stats s
		ON m.match_id=s.match_id
		WHERE m.away_team_id=$1 AND status='Ended' AND date >= $2 AND date <= $3
	)
	SELECT 
	COUNT(*) AS total_matches,
	AVG(is_win) as average_win_rate,
	AVG(points) as average_points,
	AVG(ball_possession) as average_ball_possession,
	SUM(goals),
	AVG(goals),
	SUM(goals_conceded),
	AVG(goals_conceded),
	SUM(shots_on_goal),
	AVG(shots_on_goal),
	SUM(shots_on_goal_conceded),
	AVG(shots_on_goal_conceded),
	SUM(total_shots),
	AVG(total_shots),
	SUM(total_shots_conceded),
	AVG(total_shots_conceded),
	SUM(shots_inside_box),
	AVG(shots_inside_box),
	SUM(shots_inside_box_conceded),
	AVG(shots_inside_box_conceded),
	SUM(blocked_shots),
	AVG(blocked_shots),
	SUM(fouls),
	AVG(fouls),
	SUM(fouls_conceded),
	AVG(fouls_conceded),
	SUM(corner_kicks),
	AVG(corner_kicks),
	SUM(corner_kicks_conceded),
	AVG(corner_kicks_conceded),
	SUM(yellow_cards),
	AVG(yellow_cards),
	SUM(yellow_cards_conceded),
	AVG(yellow_cards_conceded),
	SUM(red_cards),
	AVG(red_cards),
	SUM(red_cards_conceded),
	AVG(red_cards_conceded),
	SUM(total_passes),
	AVG(total_passes),
	SUM(total_passes_conceded),
	AVG(total_passes_conceded),
	SUM(complete_passes),
	AVG(complete_passes),
	SUM(complete_passes_conceded),
	AVG(complete_passes_conceded),
	if(isNaN(AVG(complete_passes/total_passes)), 0, AVG(complete_passes/total_passes)),
	if(isNaN(AVG(complete_passes_conceded/total_passes_conceded)), 0, AVG(complete_passes_conceded/total_passes_conceded)),
	SUM(offsides),
	AVG(offsides),
	SUM(offsides_conceded),
	AVG(offsides_conceded)
	FROM all_matches_stats`,
		filter.TeamID, filter.FromDate, filter.ToDate)

	var stats models.TeamStatsInPeriod
	if err := row.Scan(&stats.TotalMatches, &stats.AvgWinRate, &stats.AvgPoints, &stats.AvgBallPossession, &stats.Goals, &stats.GoalsPer90, &stats.GoalsConceded, &stats.GoalsConcededPer90,
		&stats.ShotsOnGoal, &stats.ShotsOnGoalPer90, &stats.ShotsOnGoalConceded, &stats.ShotsOnGoalConcededPer90, &stats.TotalShots, &stats.TotalShotsPer90,
		&stats.TotalShotsConceded, &stats.TotalShotsConcededPer90, &stats.ShotsInsideBox, &stats.ShotsInsideBoxPer90, &stats.ShotsInsideBoxConceded, &stats.ShotsInsideBoxConcededPer90,
		&stats.BlockedShots, &stats.BlockedShotsPer90, &stats.Fouls, &stats.FoulsPer90, &stats.WasFouled, &stats.WasFouledPer90, &stats.CornerKicks, &stats.CornerKicksPer90,
		&stats.CornerKicksConceded, &stats.CornerKicksConcededPer90, &stats.YellowCards, &stats.YellowCardsPer90, &stats.YellowCardsOpp, &stats.YellowCardsOppPer90,
		&stats.RedCards, &stats.RedCardsPer90, &stats.RedCardsOpp, &stats.RedCardsOppPer90, &stats.PassAttempts, &stats.PassAttemptsPer90, &stats.PassAttemptsOpp, &stats.PassAttemptsOppPer90,
		&stats.CompletePasses, &stats.CompletePassesPer90, &stats.CompletePassesOpp, &stats.CompletePassesOppPer90, &stats.AvgPassAccuracy, &stats.AvgPassAccuracyOpp,
		&stats.Offsides, &stats.OffsidesPer90, &stats.OffsidesOpp, &stats.OffsidesOppPer90); err != nil {
		return models.TeamStatsInPeriod{}, err
	}
	stats.IDTeam = filter.TeamID

	return stats, nil
}

func (c *ClichouseDB) GetTeamFullStats(ctx context.Context, filter models.TeamStatsFilter) (models.TeamStatsInPeriod, error) {
	row := c.conn.QueryRow(ctx, `
	WITH all_matches_stats AS (
		SELECT 
		m.match_id,
		home_team_id,
		away_team_id,
		home_score AS goals,
		away_score AS goals_conceded,
		IF(home_score > away_score, 3, IF(home_score < away_score, 0, 1)) AS points,
		IF(home_score > away_score, 1, 0) AS is_win,
		shots_on_goal_home AS shots_on_goal,
		shots_on_goal_away AS shots_on_goal_conceded,
		total_shots_home AS total_shots,
		total_shots_away AS total_shots_conceded,
		blocked_shots_home AS blocked_shots,
		blocked_shots_away AS blocked_shots_conceded,
		fouls_home AS fouls,
		fouls_away AS fouls_conceded,
		corner_kicks_home AS corner_kicks,
		corner_kicks_away AS corner_kicks_conceded,
		ball_possession_home AS ball_possession,
		yellow_cards_home AS yellow_cards,
		yellow_cards_away AS yellow_cards_conceded,
		red_cards_home AS red_cards,
		red_cards_away AS red_cards_conceded,
		total_passes_home AS total_passes,
		total_passes_away AS total_passes_conceded,
		complete_passes_home AS complete_passes,
		complete_passes_away AS complete_passes_conceded, 
		offsides_home AS offsides,
		offsides_away AS offsides_conceded,
		shots_inside_box_home AS shots_inside_box,
		shots_inside_box_away AS shots_inside_box_conceded
		FROM Matches m 
		INNER JOIN Football_Team_Match_Stats s
		ON m.match_id=s.match_id 
		WHERE m.home_team_id=$1 AND status='Ended'
		
		UNION ALL
		
		SELECT 
		m.match_id,
		home_team_id,
		away_team_id,
		home_score AS goals_conceded,
		away_score AS goals,
		IF(away_score > home_score, 3, IF(away_score < home_score, 0, 1)) AS points,
		IF(away_score > home_score, 1, 0) AS is_win,
		shots_on_goal_away AS shots_on_goal,
		shots_on_goal_home AS shots_on_goal_conceded,
		total_shots_away AS total_shots,
		total_shots_home AS total_shots_conceded,
		blocked_shots_away AS blocked_shots,
		blocked_shots_home AS blocked_shots_conceded,
		fouls_away AS fouls,
		fouls_home AS fouls_conceded,
		corner_kicks_away AS corner_kicks,
		corner_kicks_home AS corner_kicks_conceded,
		ball_possession_away AS ball_possession,
		yellow_cards_away AS yellow_cards,
		yellow_cards_home AS yellow_cards_conceded,
		red_cards_away AS red_cards,
		red_cards_home AS red_cards_conceded,
		total_passes_away AS total_passes,
		total_passes_home AS total_passes_conceded,
		complete_passes_away AS complete_passes,
		complete_passes_home AS complete_passes_conceded, 
		offsides_away AS offsides,
		offsides_home AS offsides_conceded,
		shots_inside_box_away AS shots_inside_box,
		shots_inside_box_home AS shots_inside_box_conceded
		FROM Matches m 
		INNER JOIN Football_Team_Match_Stats s
		ON m.match_id=s.match_id
		WHERE m.away_team_id=$1 AND status='Ended'
	)
	SELECT 
	COUNT(*) AS total_matches,
	AVG(is_win) as average_win_rate,
	AVG(points) as average_points,
	AVG(ball_possession) as average_ball_possession,
	SUM(goals),
	AVG(goals),
	SUM(goals_conceded),
	AVG(goals_conceded),
	SUM(shots_on_goal),
	AVG(shots_on_goal),
	SUM(shots_on_goal_conceded),
	AVG(shots_on_goal_conceded),
	SUM(total_shots),
	AVG(total_shots),
	SUM(total_shots_conceded),
	AVG(total_shots_conceded),
	SUM(shots_inside_box),
	AVG(shots_inside_box),
	SUM(shots_inside_box_conceded),
	AVG(shots_inside_box_conceded),
	SUM(blocked_shots),
	AVG(blocked_shots),
	SUM(fouls),
	AVG(fouls),
	SUM(fouls_conceded),
	AVG(fouls_conceded),
	SUM(corner_kicks),
	AVG(corner_kicks),
	SUM(corner_kicks_conceded),
	AVG(corner_kicks_conceded),
	SUM(yellow_cards),
	AVG(yellow_cards),
	SUM(yellow_cards_conceded),
	AVG(yellow_cards_conceded),
	SUM(red_cards),
	AVG(red_cards),
	SUM(red_cards_conceded),
	AVG(red_cards_conceded),
	SUM(total_passes),
	AVG(total_passes),
	SUM(total_passes_conceded),
	AVG(total_passes_conceded),
	SUM(complete_passes),
	AVG(complete_passes),
	SUM(complete_passes_conceded),
	AVG(complete_passes_conceded),
	if(isNaN(AVG(complete_passes/total_passes)), 0, AVG(complete_passes/total_passes)),
	if(isNaN(AVG(complete_passes_conceded/total_passes_conceded)), 0, AVG(complete_passes_conceded/total_passes_conceded)),
	SUM(offsides),
	AVG(offsides),
	SUM(offsides_conceded),
	AVG(offsides_conceded)
	FROM all_matches_stats`,
		filter.TeamID)

	var stats models.TeamStatsInPeriod
	if err := row.Scan(&stats.TotalMatches, &stats.AvgWinRate, &stats.AvgPoints, &stats.AvgBallPossession, &stats.Goals, &stats.GoalsPer90, &stats.GoalsConceded, &stats.GoalsConcededPer90,
		&stats.ShotsOnGoal, &stats.ShotsOnGoalPer90, &stats.ShotsOnGoalConceded, &stats.ShotsOnGoalConcededPer90, &stats.TotalShots, &stats.TotalShotsPer90,
		&stats.TotalShotsConceded, &stats.TotalShotsConcededPer90, &stats.ShotsInsideBox, &stats.ShotsInsideBoxPer90, &stats.ShotsInsideBoxConceded, &stats.ShotsInsideBoxConcededPer90,
		&stats.BlockedShots, &stats.BlockedShotsPer90, &stats.Fouls, &stats.FoulsPer90, &stats.WasFouled, &stats.WasFouledPer90, &stats.CornerKicks, &stats.CornerKicksPer90,
		&stats.CornerKicksConceded, &stats.CornerKicksConcededPer90, &stats.YellowCards, &stats.YellowCardsPer90, &stats.YellowCardsOpp, &stats.YellowCardsOppPer90,
		&stats.RedCards, &stats.RedCardsPer90, &stats.RedCardsOpp, &stats.RedCardsOppPer90, &stats.PassAttempts, &stats.PassAttemptsPer90, &stats.PassAttemptsOpp, &stats.PassAttemptsOppPer90,
		&stats.CompletePasses, &stats.CompletePassesPer90, &stats.CompletePassesOpp, &stats.CompletePassesOppPer90, &stats.AvgPassAccuracy, &stats.AvgPassAccuracyOpp,
		&stats.Offsides, &stats.OffsidesPer90, &stats.OffsidesOpp, &stats.OffsidesOppPer90); err != nil {
		return models.TeamStatsInPeriod{}, err
	}
	stats.IDTeam = filter.TeamID

	return stats, nil
}

func (c *ClichouseDB) GetStandingsByTeamAndSeason(ctx context.Context, teamID uint32, season string) ([]models.TableRow, error) {
	rows, err := c.conn.Query(ctx, `
	SELECT ttp.team_id, name, logo_url, points, position, games_played, wins, draws, losses, goals_scored, goals_conceded
	FROM Team_Tournament_Performances ttp
	INNER JOIN Teams t
	ON ttp.team_id=t.team_id
	INNER JOIN Team_Logos tl
	ON tl.team_id=t.team_id
	WHERE ttp.tournament_id IN (
	SELECT t.tournament_id FROM Team_Tournament_Performances ttp2
	INNER JOIN Tournaments t
	ON ttp2.tournament_id=t.tournament_id
	WHERE t.season = $1 AND ttp2.team_id=$2)
	ORDER BY position`,
		season, teamID)
	if err != nil {
		return nil, err
	}

	table := make([]models.TableRow, 0, 20)

	for rows.Next() {
		var row models.TableRow
		if err := rows.Scan(&row.Team.ID, &row.Team.Name, &row.Team.URLLogo, &row.Points, &row.Pos, &row.Matches, &row.Wins,
			&row.Draws, &row.Losses, &row.ScoresFor, &row.ScoresAgainst); err != nil {
			return nil, err
		}

		table = append(table, row)
	}

	return table, nil
}

func (c *ClichouseDB) GetTeamNextGame(ctx context.Context, id uint32) (models.ShortMatch, error) {
	row := c.conn.QueryRow(ctx, `SELECT 
	m.match_id, ht.team_id, ht.name, htl.logo_url, at.team_id, at.name, atl.logo_url, date, round, t.name, tl.logo_url FROM Matches m
	INNER JOIN Teams ht ON ht.team_id=m.home_team_id
	INNER JOIN Team_Logos htl ON ht.team_id=htl.team_id
	INNER JOIN Teams at ON at.team_id=m.away_team_id
	INNER JOIN Team_Logos atl ON at.team_id=atl.team_id
	INNER JOIN Tournaments t ON t.tournament_id=m.tournament_id
	INNER JOIN Tournament_Logos tl ON t.name=tl.tournament_name
	WHERE (home_team_id = $1 OR away_team_id = $1) AND date >= now()
	ORDER BY date
	LIMIT 1
	`, id)

	var match models.ShortMatch
	if err := row.Scan(&match.ID, &match.HomeTeam.ID, &match.HomeTeam.Name, &match.HomeTeam.URLLogo, &match.AwayTeam.ID, &match.AwayTeam.Name, &match.AwayTeam.URLLogo, &match.Date, &match.Round, &match.Tournament.Name, &match.Tournament.URLLogo); err != nil {
		return models.ShortMatch{}, err
	}

	return match, nil
}

func (c *ClichouseDB) GetTeamPlayersBySeason(ctx context.Context, teamID uint32, season string) (map[string][]models.Player, error) {
	rows, err := c.conn.Query(ctx, `
	WITH all_players AS (
		SELECT DISTINCT athlete_id FROM Football_Player_Match_Stats s
		INNER JOIN Matches m ON m.match_id = s.match_id
		INNER JOIN Tournaments t ON m.tournament_id = t.tournament_id
		WHERE ((home_team_id=$1 AND home_team_player) OR (away_team_id=$1 AND NOT home_team_player)) AND t.season = $2
		UNION ALL
		SELECT DISTINCT athlete_id FROM Football_Goalie_Match_Stats s
		INNER JOIN Matches m ON m.match_id = s.match_id
		INNER JOIN Tournaments t ON m.tournament_id = t.tournament_id
		WHERE ((home_team_id=$1 AND home_team_player) OR (away_team_id=$1 AND NOT home_team_player)) AND t.season = $2
	)
	SELECT 
		ap.athlete_id, first_name, last_name, a.nation, flag_url, position, height, preferred_foot_or_handness, current_status, date_of_birth, photo_url
	FROM all_players ap
	INNER JOIN Athletes a ON ap.athlete_id=a.athlete_id
	INNER JOIN Athlete_Photos p ON p.athlete_id=a.athlete_id
	INNER JOIN National_Flags nf ON a.nation=nf.nation
	ORDER BY position `,
		teamID, season)
	if err != nil {
		return nil, err
	}

	squad := make(map[string][]models.Player, 4)
	squad["G"] = make([]models.Player, 0, 5)
	squad["D"] = make([]models.Player, 0, 15)
	squad["M"] = make([]models.Player, 0, 15)
	squad["F"] = make([]models.Player, 0, 10)

	for rows.Next() {
		var player models.Player
		if err := rows.Scan(&player.ID, &player.FirstName, &player.LastName, &player.Nation.Name, &player.Nation.URLFlag, &player.Position, &player.Height, &player.PreferredFoot, &player.CurrentStatus, &player.DateOfBirth, &player.URLPhoto); err != nil {
			return nil, err
		}

		squad[player.Position] = append(squad[player.Position], player)
	}

	return squad, nil
}

func (c *ClichouseDB) GetTeamGames(ctx context.Context, teamID uint32, limit, offset uint32, season string) ([]models.ShortMatch, error) {
	rows, err := c.conn.Query(ctx, `
	SELECT 
	match_id, ht.team_id, ht.name, htl.logo_url, at.team_id, at.name, atl.logo_url, home_score, away_score, date, round, t.name, tl.logo_url, status FROM Matches m
	INNER JOIN Teams ht ON ht.team_id=m.home_team_id
	INNER JOIN Team_Logos htl ON ht.team_id=htl.team_id
	INNER JOIN Teams at ON at.team_id=m.away_team_id
	INNER JOIN Team_Logos atl ON at.team_id=atl.team_id
	INNER JOIN Tournaments t ON t.tournament_id=m.tournament_id
	INNER JOIN Tournament_Logos tl ON t.name=tl.tournament_name
	WHERE (home_team_id = $1 OR away_team_id = $1) and t.season=$2
	ORDER BY date DESC
	LIMIT $3
	OFFSET $4`,
		teamID, season, limit, offset)
	if err != nil {
		return nil, err
	}

	matches := make([]models.ShortMatch, 0, limit)

	for rows.Next() {
		var match models.ShortMatch
		if err := rows.Scan(&match.ID, &match.HomeTeam.ID, &match.HomeTeam.Name, &match.HomeTeam.URLLogo, &match.AwayTeam.ID, &match.AwayTeam.Name, &match.AwayTeam.URLLogo, &match.HomeGoals, &match.AwayGoals, &match.Date, &match.Round, &match.Tournament.Name, &match.Tournament.URLLogo, &match.Status); err != nil {
			return nil, err
		}

		matches = append(matches, match)
	}

	return matches, nil
}

func (c *ClichouseDB) GetCurrentManagerByTeam(ctx context.Context, id uint32) (models.Manager, error) {
	row := c.conn.QueryRow(ctx, `
	WITH last_matches AS (
	SELECT home_team_manager_id AS manager_id, date FROM Matches 
	WHERE home_team_id = $1 AND status != 'Not started'
	UNION ALL
	SELECT away_team_manager_id AS manager_id, date FROM Matches 
	WHERE away_team_id = $1 AND status != 'Not started'
	)
	
	SELECT 
	m.manager_id, first_name, last_name, m.nation, nf.flag_url, mp.photo_url FROM last_matches lm
	INNER JOIN Managers m ON m.manager_id=lm.manager_id
	LEFT JOIN National_Flags nf ON nf.nation=m.nation
	LEFT JOIN Manager_Photos mp ON mp.manager_id=m.manager_id
	ORDER BY date DESC
	LIMIT 1
	`, id)

	var manager models.Manager
	if err := row.Scan(&manager.ID, &manager.FirstName, &manager.LastName, &manager.Nation.Name, &manager.Nation.URLFlag, &manager.URLPhoto); err != nil {
		return models.Manager{}, err
	}

	return manager, nil
}

func (c *ClichouseDB) GetTeamPlayersWithStatsBySeason(ctx context.Context, teamID uint32, season string) (map[string][]models.PlayerWithStats, error) {
	rows, err := c.conn.Query(ctx, `
	WITH all_players AS (
		SELECT athlete_id, SUM(goals) AS sum_goals, SUM(assists) AS sum_assists, SUM(minutes_played) AS sum_minutes, COUNT(*) AS total_matches, COUNTIf(minutes_played != 0) as played_matches, AVGIf(rating, rating != 0) AS avg_rating, SUM(start_player) as starts FROM Football_Player_Match_Stats s
		INNER JOIN Matches m ON m.match_id = s.match_id
		INNER JOIN Tournaments t ON m.tournament_id = t.tournament_id
		WHERE ((home_team_id=$1 AND home_team_player) OR (away_team_id=$1 AND NOT home_team_player)) AND t.season = $2
		GROUP BY athlete_id
		UNION ALL
		SELECT athlete_id, SUM(goals) AS sum_goals, SUM(assists) AS sum_assists, SUM(minutes_played) AS sum_minutes, COUNT(*) AS total_matches, COUNTIf(minutes_played != 0) as played_matches, AVGIf(rating, rating != 0) AS avg_rating, SUM(start_player) as starts FROM Football_Goalie_Match_Stats s
		INNER JOIN Matches m ON m.match_id = s.match_id
		INNER JOIN Tournaments t ON m.tournament_id = t.tournament_id
		WHERE ((home_team_id=$1 AND home_team_player) OR (away_team_id=$1 AND NOT home_team_player)) AND t.season = $2
		GROUP BY athlete_id
	)
	SELECT 
		ap.athlete_id, first_name, last_name, a.nation, flag_url, position, height, preferred_foot_or_handness, current_status, date_of_birth, photo_url,
		sum_goals, sum_assists, sum_minutes, played_matches, total_matches, if(isNaN(avg_rating), 0, avg_rating) AS avg_rating, starts
	FROM all_players ap
	INNER JOIN Athletes a ON ap.athlete_id=a.athlete_id
	INNER JOIN Athlete_Photos p ON p.athlete_id=a.athlete_id
	INNER JOIN National_Flags nf ON a.nation=nf.nation
	ORDER BY position `,
		teamID, season)
	if err != nil {
		return nil, err
	}

	squad := make(map[string][]models.PlayerWithStats, 4)
	squad["G"] = make([]models.PlayerWithStats, 0, 5)
	squad["D"] = make([]models.PlayerWithStats, 0, 15)
	squad["M"] = make([]models.PlayerWithStats, 0, 15)
	squad["F"] = make([]models.PlayerWithStats, 0, 10)

	for rows.Next() {
		var player models.PlayerWithStats
		if err := rows.Scan(&player.ID, &player.FirstName, &player.LastName, &player.Nation.Name, &player.Nation.URLFlag, &player.Position, &player.Height, &player.PreferredFoot, &player.CurrentStatus, &player.DateOfBirth, &player.URLPhoto,
			&player.Goals, &player.Assists, &player.MinutesPlayed, &player.MatchesPlayed, &player.TotalMatches, &player.AvgRating, &player.StartPlayer); err != nil {
			return nil, err
		}

		squad[player.Position] = append(squad[player.Position], player)
	}

	return squad, nil
}

func (c *ClichouseDB) GetTournamentTableByID(ctx context.Context, tournamentID uint32) ([]models.TableRow, error) {
	rows, err := c.conn.Query(ctx, `
	SELECT ttp.team_id, name, logo_url, points, position, games_played, wins, draws, losses, goals_scored, goals_conceded
	FROM Team_Tournament_Performances ttp
	INNER JOIN Teams t
	ON ttp.team_id=t.team_id
	INNER JOIN Team_Logos tl
	ON tl.team_id=t.team_id
	WHERE ttp.tournament_id=$1
	ORDER BY position`,
		tournamentID)
	if err != nil {
		return nil, err
	}

	table := make([]models.TableRow, 0, 20)

	for rows.Next() {
		var row models.TableRow
		if err := rows.Scan(&row.Team.ID, &row.Team.Name, &row.Team.URLLogo, &row.Points, &row.Pos, &row.Matches, &row.Wins,
			&row.Draws, &row.Losses, &row.ScoresFor, &row.ScoresAgainst); err != nil {
			return nil, err
		}

		table = append(table, row)
	}

	return table, nil
}

func (c *ClichouseDB) GetTournamentFixturesByID(ctx context.Context, id uint32) (map[uint16][]models.ShortMatch, error) {
	rows, err := c.conn.Query(ctx, `
	SELECT 
		m.match_id,
		date,
		home_team_id,
		ht.name,
		htl.logo_url,
		away_team_id,
		at.name,
		atl.logo_url,
		home_score,
		away_score,
		round,
		status,
		t.tournament_id,
		t.name
		FROM Matches m 
	INNER JOIN Teams ht ON m.home_team_id=ht.team_id
	INNER JOIN Team_Logos htl ON m.home_team_id = htl.team_id
	INNER JOIN Teams at ON m.away_team_id=at.team_id
	INNER JOIN Team_Logos atl ON m.away_team_id = atl.team_id
	INNER JOIN Tournaments t ON t.tournament_id = m.tournament_id
	WHERE t.tournament_id = $1
	ORDER BY round, m.date DESC
	`, id)
	if err != nil {
		return nil, err
	}

	matches := make(map[uint16][]models.ShortMatch, 40)

	for rows.Next() {
		var match models.ShortMatch
		if err := rows.Scan(&match.ID, &match.Date, &match.HomeTeam.ID, &match.HomeTeam.Name, &match.HomeTeam.URLLogo, &match.AwayTeam.ID, &match.AwayTeam.Name, &match.AwayTeam.URLLogo,
			&match.HomeGoals, &match.AwayGoals, &match.Round, &match.Status, &match.Tournament.ID, &match.Tournament.Name); err != nil {
			return nil, err
		}

		if _, ok := matches[match.Round]; !ok {
			matches[match.Round] = make([]models.ShortMatch, 0, 10)
		}

		matches[match.Round] = append(matches[match.Round], match)
	}

	return matches, nil
}

func (c *ClichouseDB) GetTournamentTeamsStatsByID(ctx context.Context, id uint32) ([]models.TeamStatsInPeriod, error) {
	rows, err := c.conn.Query(ctx, `
	WITH all_matches_stats AS (
		SELECT 
		m.match_id,
		home_team_id AS team_id,
		home_score AS goals,
		away_score AS goals_conceded,
		IF(home_score > away_score, 3, IF(home_score < away_score, 0, 1)) AS points,
		IF(home_score > away_score, 1, 0) AS is_win,
		shots_on_goal_home AS shots_on_goal,
		shots_on_goal_away AS shots_on_goal_conceded,
		total_shots_home AS total_shots,
		total_shots_away AS total_shots_conceded,
		blocked_shots_home AS blocked_shots,
		blocked_shots_away AS blocked_shots_conceded,
		fouls_home AS fouls,
		fouls_away AS fouls_conceded,
		corner_kicks_home AS corner_kicks,
		corner_kicks_away AS corner_kicks_conceded,
		ball_possession_home AS ball_possession,
		yellow_cards_home AS yellow_cards,
		yellow_cards_away AS yellow_cards_conceded,
		red_cards_home AS red_cards,
		red_cards_away AS red_cards_conceded,
		total_passes_home AS total_passes,
		total_passes_away AS total_passes_conceded,
		complete_passes_home AS complete_passes,
		complete_passes_away AS complete_passes_conceded, 
		offsides_home AS offsides,
		offsides_away AS offsides_conceded,
		shots_inside_box_home AS shots_inside_box,
		shots_inside_box_away AS shots_inside_box_conceded
		FROM Matches m 
		INNER JOIN Football_Team_Match_Stats s
		ON m.match_id=s.match_id
		WHERE tournament_id=$1 AND status='Ended'
		
		UNION ALL
		
		SELECT 
		m.match_id,
		away_team_id AS team_id,
		home_score AS goals_conceded,
		away_score AS goals,
		IF(away_score > home_score, 3, IF(away_score < home_score, 0, 1)) AS points,
		IF(away_score > home_score, 1, 0) AS is_win,
		shots_on_goal_away AS shots_on_goal,
		shots_on_goal_home AS shots_on_goal_conceded,
		total_shots_away AS total_shots,
		total_shots_home AS total_shots_conceded,
		blocked_shots_away AS blocked_shots,
		blocked_shots_home AS blocked_shots_conceded,
		fouls_away AS fouls,
		fouls_home AS fouls_conceded,
		corner_kicks_away AS corner_kicks,
		corner_kicks_home AS corner_kicks_conceded,
		ball_possession_away AS ball_possession,
		yellow_cards_away AS yellow_cards,
		yellow_cards_home AS yellow_cards_conceded,
		red_cards_away AS red_cards,
		red_cards_home AS red_cards_conceded,
		total_passes_away AS total_passes,
		total_passes_home AS total_passes_conceded,
		complete_passes_away AS complete_passes,
		complete_passes_home AS complete_passes_conceded, 
		offsides_away AS offsides,
		offsides_home AS offsides_conceded,
		shots_inside_box_away AS shots_inside_box,
		shots_inside_box_home AS shots_inside_box_conceded
		FROM Matches m 
		INNER JOIN Football_Team_Match_Stats s
		ON m.match_id=s.match_id
		WHERE tournament_id=$1 AND status='Ended'
	)
	SELECT 
	team_id,
	COUNT(*) AS total_matches,
	AVG(is_win) as average_win_rate,
	AVG(points) as average_points,
	AVG(ball_possession) as average_ball_possession,
	SUM(goals),
	AVG(goals),
	SUM(goals_conceded),
	AVG(goals_conceded),
	SUM(shots_on_goal),
	AVG(shots_on_goal),
	SUM(shots_on_goal_conceded),
	AVG(shots_on_goal_conceded),
	SUM(total_shots),
	AVG(total_shots),
	SUM(total_shots_conceded),
	AVG(total_shots_conceded),
	SUM(shots_inside_box),
	AVG(shots_inside_box),
	SUM(shots_inside_box_conceded),
	AVG(shots_inside_box_conceded),
	SUM(blocked_shots),
	AVG(blocked_shots),
	SUM(fouls),
	AVG(fouls),
	SUM(fouls_conceded),
	AVG(fouls_conceded),
	SUM(corner_kicks),
	AVG(corner_kicks),
	SUM(corner_kicks_conceded),
	AVG(corner_kicks_conceded),
	SUM(yellow_cards),
	AVG(yellow_cards),
	SUM(yellow_cards_conceded),
	AVG(yellow_cards_conceded),
	SUM(red_cards),
	AVG(red_cards),
	SUM(red_cards_conceded),
	AVG(red_cards_conceded),
	SUM(total_passes),
	AVG(total_passes),
	SUM(total_passes_conceded),
	AVG(total_passes_conceded),
	SUM(complete_passes),
	AVG(complete_passes),
	SUM(complete_passes_conceded),
	AVG(complete_passes_conceded),
	AVG(if(total_passes > 0, complete_passes/total_passes, 0)) AS avg_pass_completion_rate,
	AVG(if(total_passes_conceded > 0, complete_passes_conceded/total_passes_conceded, 0)) AS avg_pass_completion_conceded_rate,
	SUM(offsides),
	AVG(offsides),
	SUM(offsides_conceded),
	AVG(offsides_conceded)
	FROM all_matches_stats
	GROUP BY team_id
	ORDER BY average_win_rate DESC`,
		id)
	if err != nil {
		return nil, err
	}

	allStats := make([]models.TeamStatsInPeriod, 0, 20)

	for rows.Next() {
		var stats models.TeamStatsInPeriod
		if err := rows.Scan(&stats.IDTeam, &stats.TotalMatches, &stats.AvgWinRate, &stats.AvgPoints, &stats.AvgBallPossession, &stats.Goals, &stats.GoalsPer90, &stats.GoalsConceded, &stats.GoalsConcededPer90,
			&stats.ShotsOnGoal, &stats.ShotsOnGoalPer90, &stats.ShotsOnGoalConceded, &stats.ShotsOnGoalConcededPer90, &stats.TotalShots, &stats.TotalShotsPer90,
			&stats.TotalShotsConceded, &stats.TotalShotsConcededPer90, &stats.ShotsInsideBox, &stats.ShotsInsideBoxPer90, &stats.ShotsInsideBoxConceded, &stats.ShotsInsideBoxConcededPer90,
			&stats.BlockedShots, &stats.BlockedShotsPer90, &stats.Fouls, &stats.FoulsPer90, &stats.WasFouled, &stats.WasFouledPer90, &stats.CornerKicks, &stats.CornerKicksPer90,
			&stats.CornerKicksConceded, &stats.CornerKicksConcededPer90, &stats.YellowCards, &stats.YellowCardsPer90, &stats.YellowCardsOpp, &stats.YellowCardsOppPer90,
			&stats.RedCards, &stats.RedCardsPer90, &stats.RedCardsOpp, &stats.RedCardsOppPer90, &stats.PassAttempts, &stats.PassAttemptsPer90, &stats.PassAttemptsOpp, &stats.PassAttemptsOppPer90,
			&stats.CompletePasses, &stats.CompletePassesPer90, &stats.CompletePassesOpp, &stats.CompletePassesOppPer90, &stats.AvgPassAccuracy, &stats.AvgPassAccuracyOpp,
			&stats.Offsides, &stats.OffsidesPer90, &stats.OffsidesOpp, &stats.OffsidesOppPer90); err != nil {
			return nil, err
		}

		allStats = append(allStats, stats)
	}

	return allStats, nil
}

func (c *ClichouseDB) GetTournamentPlayersStatsByID(ctx context.Context, id, limit, offset uint32) ([]models.PlayerStatsInPeriod, error) {
	rows, err := c.conn.Query(ctx, `
	SELECT 
		athlete_id,
		COUNT(*) AS total_matches,
		SUM(start_player),
		COUNTIf(minutes_played != 0) as played_matches,
		AVGIf(rating, rating != 0) as avg_rating,
		SUM(minutes_played) AS full_minutes_played,
		SUM(goals),
		SUM(assists),
		SUM(blocked_shots),
		SUM(interceptions),
		SUM(total_tackles),
		SUM(dribbled_past),
		SUM(duels),
		SUM(duels_won),
		SUM(fouls),
		SUM(was_fouled),
		SUM(pass_attempts),
		SUM(complete_passes),
		SUM(key_passes),
		SUM(shot_on_target),
		SUM(total_shots),
		SUM(dribble_attempts),
		SUM(complete_dribbles),
		SUM(penalty_scored),
		SUM(penalty_missed),
		SUM(yellow_cards),
		SUM(red_cards),
		SUM(captain)
		FROM Football_Player_Match_Stats s
	INNER JOIN Matches m ON s.match_id=m.match_id 
	WHERE tournament_id=$1
	GROUP BY athlete_id
	ORDER BY avg_rating DESC, full_minutes_played DESC
	LIMIT $2
	OFFSET $3`,
		id, limit, offset)
	if err != nil {
		return nil, err
	}

	allStats := make([]models.PlayerStatsInPeriod, 0, limit)

	for rows.Next() {
		var stats models.PlayerStatsInPeriod
		if err := rows.Scan(&stats.IDPlayer, &stats.TotalMatches, &stats.StartPlayer, &stats.MatchesPlayed, &stats.AvgRating, &stats.MinutesPlayed, &stats.Goals, &stats.Assists, &stats.BlockedShots, &stats.Interceptions, &stats.TotalTackles, &stats.DribbledPast, &stats.Duels,
			&stats.DuelsWon, &stats.Fouls, &stats.WasFouled, &stats.PassAttempts, &stats.CompletePasses, &stats.KeyPasses, &stats.ShotsOnTarget, &stats.TotalShots, &stats.DribbleAttempts, &stats.CompleteDribbles, &stats.PenaltyScored, &stats.PenaltyMissed,
			&stats.YellowCards, &stats.RedCards, &stats.CaptainTimes); err != nil {
			return nil, err
		}

		allStats = append(allStats, stats)
	}

	return allStats, nil
}

func (c *ClichouseDB) GetTournamentTableGraphByID(ctx context.Context, tournamentID uint32) ([][]models.ShortTableRow, error) {
	rows, err := c.conn.Query(ctx, `
	WITH all_matches_stats AS (
		SELECT 
		m.match_id,
		home_team_id AS team_id,
		home_score AS goals,
		away_score AS goals_conceded,
		round,
		IF(home_score > away_score, 3, IF(home_score < away_score, 0, 1)) AS points
		FROM Matches m 
		WHERE tournament_id=$1 AND status='Ended'
		
		UNION ALL
		
		SELECT 
		m.match_id,
		away_team_id AS team_id,
		home_score AS goals_conceded,
		away_score AS goals,
		round,
		IF(away_score > home_score, 3, IF(away_score < home_score, 0, 1)) AS points
		FROM Matches m 
		WHERE tournament_id=$1 AND status='Ended'
		),
	teams_with_stats AS (
		SELECT 
		team_id,
		round,
		SUM(points) OVER (PARTITION BY team_id ORDER BY round) AS total_points_after_round,
		SUM(goals-goals_conceded) OVER (PARTITION BY team_id ORDER BY round) AS total_goals_after_round
		FROM all_matches_stats
	)
		SELECT 
		tws.team_id,
		name,
		logo_url,
		round,
		ROW_NUMBER() OVER (PARTITION BY round ORDER BY total_points_after_round DESC, total_goals_after_round DESC) AS position,
		total_points_after_round,
		total_goals_after_round
		FROM teams_with_stats tws
		INNER JOIN Teams t
		ON tws.team_id=t.team_id
		INNER JOIN Team_Logos tl
		ON tl.team_id=t.team_id
		ORDER BY round, position
		`,
		tournamentID)
	if err != nil {
		return nil, err
	}

	tableGraph := make([][]models.ShortTableRow, 0, 40)

	for rows.Next() {
		var row models.ShortTableRow
		if err := rows.Scan(&row.Team.ID, &row.Team.Name, &row.Team.URLLogo, &row.Matches, &row.Pos, &row.Points, &row.Difference); err != nil {
			return nil, err
		}

		if len(tableGraph) == 0 || len(tableGraph) < int(row.Matches) {
			tableGraph = append(tableGraph, make([]models.ShortTableRow, 0, 20))
		}

		tableGraph[len(tableGraph)-1] = append(tableGraph[len(tableGraph)-1], row)
	}

	return tableGraph, nil
}

func (c *ClichouseDB) GetManagerStatsBySeason(ctx context.Context, id uint32) (map[string]models.ManagerStatsInPeriod, error) {
	rows, err := c.conn.Query(ctx, `
	WITH all_matches_stats AS (
		SELECT 
		m.match_id,
		t.season,
		home_team_id,
		away_team_id,
		home_score AS goals,
		away_score AS goals_conceded,
		IF(home_score > away_score, 3, IF(home_score < away_score, 0, 1)) AS points,
		IF(home_score > away_score, 1, 0) AS is_win,
		shots_on_goal_home AS shots_on_goal,
		shots_on_goal_away AS shots_on_goal_conceded,
		total_shots_home AS total_shots,
		total_shots_away AS total_shots_conceded,
		blocked_shots_home AS blocked_shots,
		blocked_shots_away AS blocked_shots_conceded,
		fouls_home AS fouls,
		fouls_away AS fouls_conceded,
		corner_kicks_home AS corner_kicks,
		corner_kicks_away AS corner_kicks_conceded,
		ball_possession_home AS ball_possession,
		total_passes_home AS total_passes,
		total_passes_away AS total_passes_conceded,
		complete_passes_home AS complete_passes,
		complete_passes_away AS complete_passes_conceded, 
		offsides_home AS offsides,
		offsides_away AS offsides_conceded,
		shots_inside_box_home AS shots_inside_box,
		shots_inside_box_away AS shots_inside_box_conceded
		FROM Matches m 
		INNER JOIN Football_Team_Match_Stats s
		ON m.match_id=s.match_id
		INNER JOIN Tournaments t ON t.tournament_id=m.tournament_id 
		WHERE m.home_team_manager_id=$1 AND status='Ended'
		
		UNION ALL
		
		SELECT 
		m.match_id,
		t.season,
		home_team_id,
		away_team_id,
		home_score AS goals_conceded,
		away_score AS goals,
		IF(away_score > home_score, 3, IF(away_score < home_score, 0, 1)) AS points,
		IF(away_score > home_score, 1, 0) AS is_win,
		shots_on_goal_away AS shots_on_goal,
		shots_on_goal_home AS shots_on_goal_conceded,
		total_shots_away AS total_shots,
		total_shots_home AS total_shots_conceded,
		blocked_shots_away AS blocked_shots,
		blocked_shots_home AS blocked_shots_conceded,
		fouls_away AS fouls,
		fouls_home AS fouls_conceded,
		corner_kicks_away AS corner_kicks,
		corner_kicks_home AS corner_kicks_conceded,
		ball_possession_away AS ball_possession,
		total_passes_away AS total_passes,
		total_passes_home AS total_passes_conceded,
		complete_passes_away AS complete_passes,
		complete_passes_home AS complete_passes_conceded, 
		offsides_away AS offsides,
		offsides_home AS offsides_conceded,
		shots_inside_box_away AS shots_inside_box,
		shots_inside_box_home AS shots_inside_box_conceded
		FROM Matches m 
		INNER JOIN Football_Team_Match_Stats s
		ON m.match_id=s.match_id
		INNER JOIN Tournaments t ON t.tournament_id=m.tournament_id 
		WHERE m.away_team_manager_id=$1 AND status='Ended'
	)
	SELECT 
	season, 
	COUNT(*) AS total_matches,
	AVG(is_win) as average_win_rate,
	AVG(points) as average_points,
	AVG(ball_possession) as average_ball_possession,
	SUM(goals),
	AVG(goals),
	SUM(goals_conceded),
	AVG(goals_conceded),
	SUM(shots_on_goal),
	AVG(shots_on_goal),
	SUM(shots_on_goal_conceded),
	AVG(shots_on_goal_conceded),
	SUM(total_shots),
	AVG(total_shots),
	SUM(total_shots_conceded),
	AVG(total_shots_conceded),
	SUM(shots_inside_box),
	AVG(shots_inside_box),
	SUM(shots_inside_box_conceded),
	AVG(shots_inside_box_conceded),
	SUM(blocked_shots),
	AVG(blocked_shots),
	SUM(fouls),
	AVG(fouls),
	SUM(fouls_conceded),
	AVG(fouls_conceded),
	SUM(corner_kicks),
	AVG(corner_kicks),
	SUM(corner_kicks_conceded),
	AVG(corner_kicks_conceded),
	SUM(total_passes),
	AVG(total_passes),
	SUM(total_passes_conceded),
	AVG(total_passes_conceded),
	SUM(complete_passes),
	AVG(complete_passes),
	SUM(complete_passes_conceded),
	AVG(complete_passes_conceded),
	if(isNaN(AVG(complete_passes/total_passes)), 0, AVG(complete_passes/total_passes)),
	if(isNaN(AVG(complete_passes_conceded/total_passes_conceded)), 0, AVG(complete_passes_conceded/total_passes_conceded)),
	SUM(offsides),
	AVG(offsides),
	SUM(offsides_conceded),
	AVG(offsides_conceded)
	FROM all_matches_stats
	GROUP BY season
	ORDER BY season DESC`,
		id)
	if err != nil {
		return nil, err
	}

	allStats := make(map[string]models.ManagerStatsInPeriod)

	for rows.Next() {
		var stats models.ManagerStatsInPeriod
		var season string
		if err := rows.Scan(&season, &stats.TotalMatches, &stats.WinPercentage, &stats.AvgPoints, &stats.AvgBallPossession, &stats.Goals, &stats.GoalsPer90, &stats.GoalsConceded, &stats.GoalsConcededPer90,
			&stats.ShotsOnGoal, &stats.ShotsOnGoalPer90, &stats.ShotsOnGoalConceded, &stats.ShotsOnGoalConcededPer90, &stats.TotalShots, &stats.TotalShotsPer90,
			&stats.TotalShotsConceded, &stats.TotalShotsConcededPer90, &stats.ShotsInsideBox, &stats.ShotsInsideBoxPer90, &stats.ShotsInsideBoxConceded, &stats.ShotsInsideBoxConcededPer90,
			&stats.BlockedShots, &stats.BlockedShotsPer90, &stats.Fouls, &stats.FoulsPer90, &stats.WasFouled, &stats.WasFouledPer90, &stats.CornerKicks, &stats.CornerKicksPer90,
			&stats.CornerKicksConceded, &stats.CornerKicksConcededPer90, &stats.PassAttempts, &stats.PassAttemptsPer90, &stats.PassAttemptsOpp, &stats.PassAttemptsOppPer90,
			&stats.CompletePasses, &stats.CompletePassesPer90, &stats.CompletePassesOpp, &stats.CompletePassesOppPer90, &stats.AvgPassAccuracy, &stats.AvgPassAccuracyOpp,
			&stats.Offsides, &stats.OffsidesPer90, &stats.OffsidesOpp, &stats.OffsidesOppPer90); err != nil {
			return nil, err
		}
		stats.IDManager = id
		allStats[season] = stats
	}

	return allStats, nil
}

func (c *ClichouseDB) GetManagerStatsByTeam(ctx context.Context, id uint32) (map[string]models.ManagerStatsInPeriod, error) {
	rows, err := c.conn.Query(ctx, `
	WITH all_matches_stats AS (
		SELECT 
		m.match_id,
		t.name,
		home_team_id,
		away_team_id,
		home_score AS goals,
		away_score AS goals_conceded,
		IF(home_score > away_score, 3, IF(home_score < away_score, 0, 1)) AS points,
		IF(home_score > away_score, 1, 0) AS is_win,
		shots_on_goal_home AS shots_on_goal,
		shots_on_goal_away AS shots_on_goal_conceded,
		total_shots_home AS total_shots,
		total_shots_away AS total_shots_conceded,
		blocked_shots_home AS blocked_shots,
		blocked_shots_away AS blocked_shots_conceded,
		fouls_home AS fouls,
		fouls_away AS fouls_conceded,
		corner_kicks_home AS corner_kicks,
		corner_kicks_away AS corner_kicks_conceded,
		ball_possession_home AS ball_possession,
		total_passes_home AS total_passes,
		total_passes_away AS total_passes_conceded,
		complete_passes_home AS complete_passes,
		complete_passes_away AS complete_passes_conceded, 
		offsides_home AS offsides,
		offsides_away AS offsides_conceded,
		shots_inside_box_home AS shots_inside_box,
		shots_inside_box_away AS shots_inside_box_conceded
		FROM Matches m 
		INNER JOIN Football_Team_Match_Stats s
		ON m.match_id=s.match_id
		INNER JOIN Teams t ON t.team_id=m.home_team_id 
		WHERE m.home_team_manager_id=$1 AND status='Ended'
		
		UNION ALL
		
		SELECT 
		m.match_id,
		t.name,
		home_team_id,
		away_team_id,
		home_score AS goals_conceded,
		away_score AS goals,
		IF(away_score > home_score, 3, IF(away_score < home_score, 0, 1)) AS points,
		IF(away_score > home_score, 1, 0) AS is_win,
		shots_on_goal_away AS shots_on_goal,
		shots_on_goal_home AS shots_on_goal_conceded,
		total_shots_away AS total_shots,
		total_shots_home AS total_shots_conceded,
		blocked_shots_away AS blocked_shots,
		blocked_shots_home AS blocked_shots_conceded,
		fouls_away AS fouls,
		fouls_home AS fouls_conceded,
		corner_kicks_away AS corner_kicks,
		corner_kicks_home AS corner_kicks_conceded,
		ball_possession_away AS ball_possession,
		total_passes_away AS total_passes,
		total_passes_home AS total_passes_conceded,
		complete_passes_away AS complete_passes,
		complete_passes_home AS complete_passes_conceded, 
		offsides_away AS offsides,
		offsides_home AS offsides_conceded,
		shots_inside_box_away AS shots_inside_box,
		shots_inside_box_home AS shots_inside_box_conceded
		FROM Matches m 
		INNER JOIN Football_Team_Match_Stats s
		ON m.match_id=s.match_id
		INNER JOIN Teams t ON t.team_id=m.away_team_id 
		WHERE m.away_team_manager_id=$1 AND status='Ended'
	)
	SELECT 
	name, 
	COUNT(*) AS total_matches,
	AVG(is_win) as average_win_rate,
	AVG(points) as average_points,
	AVG(ball_possession) as average_ball_possession,
	SUM(goals),
	AVG(goals),
	SUM(goals_conceded),
	AVG(goals_conceded),
	SUM(shots_on_goal),
	AVG(shots_on_goal),
	SUM(shots_on_goal_conceded),
	AVG(shots_on_goal_conceded),
	SUM(total_shots),
	AVG(total_shots),
	SUM(total_shots_conceded),
	AVG(total_shots_conceded),
	SUM(shots_inside_box),
	AVG(shots_inside_box),
	SUM(shots_inside_box_conceded),
	AVG(shots_inside_box_conceded),
	SUM(blocked_shots),
	AVG(blocked_shots),
	SUM(fouls),
	AVG(fouls),
	SUM(fouls_conceded),
	AVG(fouls_conceded),
	SUM(corner_kicks),
	AVG(corner_kicks),
	SUM(corner_kicks_conceded),
	AVG(corner_kicks_conceded),
	SUM(total_passes),
	AVG(total_passes),
	SUM(total_passes_conceded),
	AVG(total_passes_conceded),
	SUM(complete_passes),
	AVG(complete_passes),
	SUM(complete_passes_conceded),
	AVG(complete_passes_conceded),
	if(isNaN(AVG(complete_passes/total_passes)), 0, AVG(complete_passes/total_passes)),
	if(isNaN(AVG(complete_passes_conceded/total_passes_conceded)), 0, AVG(complete_passes_conceded/total_passes_conceded)),
	SUM(offsides),
	AVG(offsides),
	SUM(offsides_conceded),
	AVG(offsides_conceded)
	FROM all_matches_stats
	GROUP BY name
	ORDER BY total_matches DESC`,
		id)
	if err != nil {
		return nil, err
	}

	allStats := make(map[string]models.ManagerStatsInPeriod)

	for rows.Next() {
		var stats models.ManagerStatsInPeriod
		var team string
		if err := rows.Scan(&team, &stats.TotalMatches, &stats.WinPercentage, &stats.AvgPoints, &stats.AvgBallPossession, &stats.Goals, &stats.GoalsPer90, &stats.GoalsConceded, &stats.GoalsConcededPer90,
			&stats.ShotsOnGoal, &stats.ShotsOnGoalPer90, &stats.ShotsOnGoalConceded, &stats.ShotsOnGoalConcededPer90, &stats.TotalShots, &stats.TotalShotsPer90,
			&stats.TotalShotsConceded, &stats.TotalShotsConcededPer90, &stats.ShotsInsideBox, &stats.ShotsInsideBoxPer90, &stats.ShotsInsideBoxConceded, &stats.ShotsInsideBoxConcededPer90,
			&stats.BlockedShots, &stats.BlockedShotsPer90, &stats.Fouls, &stats.FoulsPer90, &stats.WasFouled, &stats.WasFouledPer90, &stats.CornerKicks, &stats.CornerKicksPer90,
			&stats.CornerKicksConceded, &stats.CornerKicksConcededPer90, &stats.PassAttempts, &stats.PassAttemptsPer90, &stats.PassAttemptsOpp, &stats.PassAttemptsOppPer90,
			&stats.CompletePasses, &stats.CompletePassesPer90, &stats.CompletePassesOpp, &stats.CompletePassesOppPer90, &stats.AvgPassAccuracy, &stats.AvgPassAccuracyOpp,
			&stats.Offsides, &stats.OffsidesPer90, &stats.OffsidesOpp, &stats.OffsidesOppPer90); err != nil {
			return nil, err
		}
		stats.IDManager = id
		allStats[team] = stats
	}

	return allStats, nil
}

func (c *ClichouseDB) GetManagerFullStats(ctx context.Context, id uint32) (map[string]models.ManagerStatsInPeriod, error) {
	row := c.conn.QueryRow(ctx, `
	WITH all_matches_stats AS (
		SELECT 
		m.match_id,
		home_team_id,
		away_team_id,
		home_score AS goals,
		away_score AS goals_conceded,
		IF(home_score > away_score, 3, IF(home_score < away_score, 0, 1)) AS points,
		IF(home_score > away_score, 1, 0) AS is_win,
		shots_on_goal_home AS shots_on_goal,
		shots_on_goal_away AS shots_on_goal_conceded,
		total_shots_home AS total_shots,
		total_shots_away AS total_shots_conceded,
		blocked_shots_home AS blocked_shots,
		blocked_shots_away AS blocked_shots_conceded,
		fouls_home AS fouls,
		fouls_away AS fouls_conceded,
		corner_kicks_home AS corner_kicks,
		corner_kicks_away AS corner_kicks_conceded,
		ball_possession_home AS ball_possession,
		total_passes_home AS total_passes,
		total_passes_away AS total_passes_conceded,
		complete_passes_home AS complete_passes,
		complete_passes_away AS complete_passes_conceded, 
		offsides_home AS offsides,
		offsides_away AS offsides_conceded,
		shots_inside_box_home AS shots_inside_box,
		shots_inside_box_away AS shots_inside_box_conceded
		FROM Matches m 
		INNER JOIN Football_Team_Match_Stats s
		ON m.match_id=s.match_id
		WHERE m.home_team_manager_id=$1 AND status='Ended'
		
		UNION ALL
		
		SELECT 
		m.match_id,
		home_team_id,
		away_team_id,
		home_score AS goals_conceded,
		away_score AS goals,
		IF(away_score > home_score, 3, IF(away_score < home_score, 0, 1)) AS points,
		IF(away_score > home_score, 1, 0) AS is_win,
		shots_on_goal_away AS shots_on_goal,
		shots_on_goal_home AS shots_on_goal_conceded,
		total_shots_away AS total_shots,
		total_shots_home AS total_shots_conceded,
		blocked_shots_away AS blocked_shots,
		blocked_shots_home AS blocked_shots_conceded,
		fouls_away AS fouls,
		fouls_home AS fouls_conceded,
		corner_kicks_away AS corner_kicks,
		corner_kicks_home AS corner_kicks_conceded,
		ball_possession_away AS ball_possession,
		total_passes_away AS total_passes,
		total_passes_home AS total_passes_conceded,
		complete_passes_away AS complete_passes,
		complete_passes_home AS complete_passes_conceded, 
		offsides_away AS offsides,
		offsides_home AS offsides_conceded,
		shots_inside_box_away AS shots_inside_box,
		shots_inside_box_home AS shots_inside_box_conceded
		FROM Matches m 
		INNER JOIN Football_Team_Match_Stats s
		ON m.match_id=s.match_id
		WHERE m.away_team_manager_id=$1 AND status='Ended'
	)
	SELECT 
	COUNT(*) AS total_matches,
	AVG(is_win) as average_win_rate,
	AVG(points) as average_points,
	AVG(ball_possession) as average_ball_possession,
	SUM(goals),
	AVG(goals),
	SUM(goals_conceded),
	AVG(goals_conceded),
	SUM(shots_on_goal),
	AVG(shots_on_goal),
	SUM(shots_on_goal_conceded),
	AVG(shots_on_goal_conceded),
	SUM(total_shots),
	AVG(total_shots),
	SUM(total_shots_conceded),
	AVG(total_shots_conceded),
	SUM(shots_inside_box),
	AVG(shots_inside_box),
	SUM(shots_inside_box_conceded),
	AVG(shots_inside_box_conceded),
	SUM(blocked_shots),
	AVG(blocked_shots),
	SUM(fouls),
	AVG(fouls),
	SUM(fouls_conceded),
	AVG(fouls_conceded),
	SUM(corner_kicks),
	AVG(corner_kicks),
	SUM(corner_kicks_conceded),
	AVG(corner_kicks_conceded),
	SUM(total_passes),
	AVG(total_passes),
	SUM(total_passes_conceded),
	AVG(total_passes_conceded),
	SUM(complete_passes),
	AVG(complete_passes),
	SUM(complete_passes_conceded),
	AVG(complete_passes_conceded),
	if(isNaN(AVG(complete_passes/total_passes)), 0, AVG(complete_passes/total_passes)),
	if(isNaN(AVG(complete_passes_conceded/total_passes_conceded)), 0, AVG(complete_passes_conceded/total_passes_conceded)),
	SUM(offsides),
	AVG(offsides),
	SUM(offsides_conceded),
	AVG(offsides_conceded)
	FROM all_matches_stats`,
		id)

	allStats := make(map[string]models.ManagerStatsInPeriod)

	var stats models.ManagerStatsInPeriod
	if err := row.Scan(&stats.TotalMatches, &stats.WinPercentage, &stats.AvgPoints, &stats.AvgBallPossession, &stats.Goals, &stats.GoalsPer90, &stats.GoalsConceded, &stats.GoalsConcededPer90,
		&stats.ShotsOnGoal, &stats.ShotsOnGoalPer90, &stats.ShotsOnGoalConceded, &stats.ShotsOnGoalConcededPer90, &stats.TotalShots, &stats.TotalShotsPer90,
		&stats.TotalShotsConceded, &stats.TotalShotsConcededPer90, &stats.ShotsInsideBox, &stats.ShotsInsideBoxPer90, &stats.ShotsInsideBoxConceded, &stats.ShotsInsideBoxConcededPer90,
		&stats.BlockedShots, &stats.BlockedShotsPer90, &stats.Fouls, &stats.FoulsPer90, &stats.WasFouled, &stats.WasFouledPer90, &stats.CornerKicks, &stats.CornerKicksPer90,
		&stats.CornerKicksConceded, &stats.CornerKicksConcededPer90, &stats.PassAttempts, &stats.PassAttemptsPer90, &stats.PassAttemptsOpp, &stats.PassAttemptsOppPer90,
		&stats.CompletePasses, &stats.CompletePassesPer90, &stats.CompletePassesOpp, &stats.CompletePassesOppPer90, &stats.AvgPassAccuracy, &stats.AvgPassAccuracyOpp,
		&stats.Offsides, &stats.OffsidesPer90, &stats.OffsidesOpp, &stats.OffsidesOppPer90); err != nil {
		return nil, err
	}
	stats.IDManager = id
	allStats["full"] = stats

	return allStats, nil
}

func (c *ClichouseDB) GetManagerCardsByID(ctx context.Context, id uint32) (uint16, uint16, error) {
	row := c.conn.QueryRow(ctx, "SELECT yellow_cards, red_cards FROM Managers WHERE manager_id=$1", id)

	var yellow, red uint16
	if err := row.Scan(&yellow, &red); err != nil {
		return 0, 0, err
	}
	return yellow, red, nil
}

func (c *ClichouseDB) GetManagerTeams(ctx context.Context, id uint32) ([]models.ManagerTeams, error) {
	rows, err := c.conn.Query(ctx, `
	WITH all_teams_seasons AS (
		SELECT DISTINCT 
		season,
		home_team_id AS team_id
		FROM Matches m 
		INNER JOIN Tournaments t ON t.tournament_id = m.tournament_id
		WHERE m.home_team_manager_id=$1 AND status='Ended'
		
		UNION ALL
		
		SELECT DISTINCT 
		season,
		away_team_id AS team_id
		FROM Matches m 
		INNER JOIN Tournaments t ON t.tournament_id = m.tournament_id
		WHERE m.away_team_manager_id=$1 AND status='Ended'
	)
	SELECT DISTINCT
	season,
	t.team_id,
	name,
	logo_url
	FROM all_teams_seasons a
	INNER JOIN Teams t ON t.team_id=a.team_id
	INNER JOIN Team_Logos tl ON tl.team_id=a.team_id
	ORDER BY name, season`,
		id)
	if err != nil {
		return nil, err
	}

	allTeams := make([]models.ManagerTeams, 0, 10)

	for rows.Next() {
		var season, team, logoURL string
		var teamID uint32

		if err := rows.Scan(&season, &teamID, &team, &logoURL); err != nil {
			return nil, err
		}
		allTeams = append(allTeams, models.ManagerTeams{
			Season: season,
			Teams:  models.Team{Name: team, URLLogo: logoURL, ID: teamID},
		})
	}

	return allTeams, nil
}

func (c *ClichouseDB) GetManagerFixtures(ctx context.Context, id, limit, offset uint32, season string) ([]models.ManagerMatch, error) {
	// Базовый запрос без ORDER BY и LIMIT/OFFSET
	baseQuery := `
    SELECT 
        m.match_id,
        date,
        home_team_id,
        ht.name,
        htl.logo_url,
        away_team_id,
        at.name,
        atl.logo_url,
        home_score,
        away_score,
        round,
        status,
        t.tournament_id,
        t.name
    FROM Matches m 
    INNER JOIN Teams ht ON m.home_team_id = ht.team_id
    INNER JOIN Team_Logos htl ON m.home_team_id = htl.team_id
    INNER JOIN Teams at ON m.away_team_id = at.team_id
    INNER JOIN Team_Logos atl ON m.away_team_id = atl.team_id
    INNER JOIN Tournaments t ON m.tournament_id = t.tournament_id
    WHERE (home_team_manager_id = $1 OR away_team_manager_id = $1) AND t.season = $2
    `

	var rows driver.Rows
	var err error

	if limit > 0 {
		// Оборачиваем в подзапрос для правильной пагинации
		rows, err = c.conn.Query(ctx, `
            SELECT * FROM (
                `+baseQuery+`
            ) AS combined
            ORDER BY date DESC
            LIMIT $3 OFFSET $4
        `, id, season, limit, offset)
	} else {
		rows, err = c.conn.Query(ctx, `
            SELECT * FROM (
                `+baseQuery+`
            ) AS combined
            ORDER BY date DESC
        `, id, season)
		limit = 100
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]models.ManagerMatch, 0, limit)

	for rows.Next() {
		var match models.ManagerMatch
		if err := rows.Scan(
			&match.ID,
			&match.Date,
			&match.HomeTeam.ID,
			&match.HomeTeam.Name,
			&match.HomeTeam.URLLogo,
			&match.AwayTeam.ID,
			&match.AwayTeam.Name,
			&match.AwayTeam.URLLogo,
			&match.HomeGoals,
			&match.AwayGoals,
			&match.Round,
			&match.Status,
			&match.Tournament.ID,
			&match.Tournament.Name,
		); err != nil {
			return nil, err
		}
		match.ManagerID = id
		matches = append(matches, match)
	}

	return matches, nil
}

func (c *ClichouseDB) GetMatchesByDate(ctx context.Context, date time.Time) ([]models.Match, error) {
	rows, err := c.conn.Query(ctx, `
	SELECT 
	m.match_id, ht.team_id, ht.name, htl.logo_url, at.team_id, at.name, atl.logo_url, date, round, t.name, tl.logo_url, home_score, away_score, status FROM Matches m
	INNER JOIN Teams ht ON ht.team_id=m.home_team_id
	INNER JOIN Team_Logos htl ON ht.team_id=htl.team_id
	INNER JOIN Teams at ON at.team_id=m.away_team_id
	INNER JOIN Team_Logos atl ON at.team_id=atl.team_id
	INNER JOIN Tournaments t ON t.tournament_id=m.tournament_id
	INNER JOIN Tournament_Logos tl ON t.name=tl.tournament_name
	WHERE toDate(date) = $1
	ORDER BY m.tournament_id, date DESC
	`,
		date)
	if err != nil {
		return nil, err
	}

	matches := make([]models.Match, 0, 10)

	for rows.Next() {
		var match models.Match
		if err := rows.Scan(&match.ID, &match.HomeTeam.ID, &match.HomeTeam.Name, &match.HomeTeam.URLLogo, &match.AwayTeam.ID, &match.AwayTeam.Name, &match.AwayTeam.URLLogo, &match.Date, &match.Round, &match.Tournament.Name, &match.Tournament.URLLogo,
			&match.HomeGoals, &match.AwayGoals, &match.Status); err != nil {
			return nil, err
		}

		matches = append(matches, match)
	}

	return matches, nil
}

func (c *ClichouseDB) GetMatchPlayersStats(ctx context.Context, id uint32) (map[string][]models.PlayerStatsInMatch, error) {
	rows, err := c.conn.Query(ctx, `
	SELECT 
		stat_id, match_id, s.athlete_id, first_name, last_name, position, start_player, toFloat64(rating), minutes_played, goals, assists, blocked_shots, interceptions, total_tackles,
		dribbled_past, duels, duels_won, fouls, was_fouled, pass_attempts, complete_passes, key_passes, shot_on_target, total_shots, dribble_attempts, complete_dribbles,
		penalty_scored, penalty_missed, yellow_cards, red_cards, captain, home_team_player FROM Football_Player_Match_Stats s
	INNER JOIN Athletes a ON a.athlete_id=s.athlete_id
	WHERE match_id = $1
	ORDER BY position, minutes_played DESC, rating DESC
	`,
		id)
	if err != nil {
		return nil, err
	}

	allStats := make(map[string][]models.PlayerStatsInMatch, 2)
	allStats["home_team"] = make([]models.PlayerStatsInMatch, 0, 20)
	allStats["away_team"] = make([]models.PlayerStatsInMatch, 0, 20)

	for rows.Next() {
		var stats models.PlayerStatsInMatch
		if err := rows.Scan(&stats.ID, &stats.IDMatch, &stats.Player.ID, &stats.Player.FirstName, &stats.Player.LastName, &stats.Player.Position, &stats.StartPlayer, &stats.Rating, &stats.MinutesPlayed, &stats.Goals, &stats.Assists, &stats.BlockedShots,
			&stats.Interceptions, &stats.TotalTackles, &stats.DribbledPast, &stats.Duels, &stats.DuelsWon, &stats.Fouls, &stats.WasFouled, &stats.PassAttempts, &stats.CompletePasses, &stats.KeyPasses,
			&stats.ShotsOnTarget, &stats.TotalShots, &stats.DribbleAttempts, &stats.CompleteDribbles, &stats.PenaltyScored, &stats.PenaltyMissed, &stats.YellowCards, &stats.RedCards, &stats.Captain, &stats.HomeTeamPlayer); err != nil {
			return nil, err
		}

		if stats.HomeTeamPlayer {
			allStats["home_team"] = append(allStats["home_team"], stats)
		} else {
			allStats["away_team"] = append(allStats["away_team"], stats)
		}
	}

	return allStats, nil
}

func (c *ClichouseDB) GetMatchGoaliesStats(ctx context.Context, id uint32) (map[string][]models.GoalieStatsInMatch, error) {
	rows, err := c.conn.Query(ctx, `
	SELECT 
		stat_id, match_id, a.athlete_id, first_name, last_name, position, start_player, toFloat64(rating), minutes_played, goals, assists, goals_conceded, saves, pass_attempts,
		complete_passes, key_passes, penalty_saved, penalty_conceded, fouls, was_fouled, yellow_cards, red_cards, captain, home_team_player FROM Football_Goalie_Match_Stats s
	INNER JOIN Athletes a ON a.athlete_id=s.athlete_id
	WHERE match_id = $1
	ORDER BY minutes_played DESC, rating DESC
	`,
		id)
	if err != nil {
		return nil, err
	}

	allStats := make(map[string][]models.GoalieStatsInMatch, 2)
	allStats["home_team"] = make([]models.GoalieStatsInMatch, 0, 5)
	allStats["away_team"] = make([]models.GoalieStatsInMatch, 0, 5)

	for rows.Next() {
		var stats models.GoalieStatsInMatch
		if err := rows.Scan(&stats.ID, &stats.IDMatch, &stats.Player.ID, &stats.Player.FirstName, &stats.Player.LastName, &stats.Player.Position, &stats.StartPlayer, &stats.Rating, &stats.MinutesPlayed, &stats.Goals, &stats.Assists, &stats.GoalsConceded,
			&stats.Saves, &stats.PassAttempts, &stats.CompletePasses, &stats.KeyPasses, &stats.PenaltySaved, &stats.PenaltyConceded, &stats.Fouls, &stats.WasFouled, &stats.YellowCards, &stats.RedCards,
			&stats.Captain, &stats.HomeTeamPlayer); err != nil {
			return nil, err
		}

		if stats.HomeTeamPlayer {
			allStats["home_team"] = append(allStats["home_team"], stats)
		} else {
			allStats["away_team"] = append(allStats["away_team"], stats)
		}
	}

	return allStats, nil
}

func (c *ClichouseDB) GetMatchTeamsStats(ctx context.Context, id uint32) (models.TeamMatchStats, error) {
	row := c.conn.QueryRow(ctx, `
	SELECT 
		shots_on_goal_home, shots_on_goal_away, total_shots_home, total_shots_away, blocked_shots_home, blocked_shots_away, fouls_home, fouls_away,
		corner_kicks_home, corner_kicks_away, ball_possession_home, ball_possession_away, yellow_cards_home, yellow_cards_away, red_cards_home, red_cards_away, total_passes_home, total_passes_away,
		complete_passes_home, complete_passes_away, offsides_home, offsides_away, shots_inside_box_home, shots_inside_box_away FROM Football_Team_Match_Stats s
	WHERE match_id = $1
	`,
		id)

	var stats models.TeamMatchStats
	if err := row.Scan(&stats.ShotsOnGoalHome, &stats.ShotsOnGoalAway, &stats.TotalShotsHome, &stats.TotalShotsAway, &stats.BlockedShotsHome, &stats.BlockedShotsAway, &stats.FoulsHome, &stats.FoulsAway,
		&stats.CornerKicksHome, &stats.CornerKicksAway, &stats.BallPossessionHome, &stats.BallPossessionAway, &stats.YellowCardsHome, &stats.YellowCardsAway, &stats.RedCardsHome, &stats.RedCardsAway, &stats.TotalPassesHome, &stats.TotalPassesAway,
		&stats.CompletePassesHome, &stats.CompletePassesAway, &stats.OffsidesHome, &stats.OffsidesAway, &stats.ShotsInsideBoxHome, &stats.ShotsInsideBoxAway); err != nil {
		return stats, err
	}
	stats.IDMatch = id

	return stats, nil
}

func (c *ClichouseDB) GetUnindexedPlayers(ctx context.Context) ([]models.Player, error) {
	rows, err := c.conn.Query(ctx, `
	SELECT a.athlete_id, a.first_name, a.last_name, a.position, a.date_of_birth, a.height, a.preferred_foot_or_handness, a.nation, a.current_status, ap.photo_url, nf.flag_url FROM Athletes a 
	INNER JOIN Athlete_Photos ap ON a.athlete_id=ap.athlete_id 
	INNER JOIN National_Flags nf ON a.nation=nf.nation 
	WHERE is_indexed=0;
	`)
	if err != nil {
		return nil, err
	}

	players := make([]models.Player, 0)
	for rows.Next() {
		var player models.Player
		if err := rows.Scan(&player.ID, &player.FirstName, &player.LastName, &player.Position, &player.DateOfBirth, &player.Height, &player.PreferredFoot, &player.Nation.Name, &player.CurrentStatus, &player.URLPhoto, &player.Nation.URLFlag); err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return players, nil
}

func (c *ClichouseDB) GetUnindexedTeams(ctx context.Context) ([]models.Team, error) {
	rows, err := c.conn.Query(ctx, "SELECT t.team_id, t.name, tl.logo_url FROM Teams t INNER JOIN Team_Logos tl ON t.team_id=tl.team_id WHERE is_indexed=0;")
	if err != nil {
		return nil, err
	}

	teams := make([]models.Team, 0)
	for rows.Next() {
		var team models.Team
		if err := rows.Scan(&team.ID, &team.Name, &team.URLLogo); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func (c *ClichouseDB) GetUnindexedManagers(ctx context.Context) ([]models.Manager, error) {
	rows, err := c.conn.Query(ctx, "SELECT m.manager_id, m.first_name, last_name, m.nation, mp.photo_url, nf.flag_url FROM Managers m INNER JOIN Manager_Photos mp ON m.manager_id=mp.manager_id INNER JOIN National_Flags nf ON m.nation=nf.nation WHERE is_indexed=0;")
	if err != nil {
		return nil, err
	}
	managers := make([]models.Manager, 0)
	for rows.Next() {
		var manager models.Manager
		if err := rows.Scan(&manager.ID, &manager.FirstName, &manager.LastName, &manager.Nation.Name, &manager.URLPhoto, &manager.Nation.URLFlag); err != nil {
			return nil, err
		}
		managers = append(managers, manager)
	}
	return managers, nil
}

func (c *ClichouseDB) GetUnindexedTournaments(ctx context.Context) ([]models.Tournament, error) {
	rows, err := c.conn.Query(ctx, "SELECT t.tournament_id, t.name, t.country, nf.flag_url, t.season, tl.logo_url FROM Tournaments t INNER JOIN National_Flags nf ON t.country=nf.nation INNER JOIN Tournament_Logos tl ON t.name=tl.tournament_name WHERE is_indexed=0;")
	if err != nil {
		return nil, err
	}
	tournaments := make([]models.Tournament, 0)
	for rows.Next() {
		var tournament models.Tournament
		if err := rows.Scan(&tournament.ID, &tournament.Name, &tournament.Country.Name, &tournament.Country.URLFlag, &tournament.Season, &tournament.URLLogo); err != nil {
			return nil, err
		}
		tournaments = append(tournaments, tournament)
	}
	return tournaments, nil
}

func (c *ClichouseDB) UpdateBatchPlayersIndexedStatus(ctx context.Context, ids []uint32) error {
	return c.conn.Exec(ctx, "ALTER TABLE Athletes UPDATE is_indexed=1 WHERE athlete_id IN $1", ids)
}

func (c *ClichouseDB) UpdateBatchManagersIndexedStatus(ctx context.Context, ids []uint32) error {
	return c.conn.Exec(ctx, "ALTER TABLE Managers UPDATE is_indexed=1 WHERE manager_id IN $1", ids)
}

func (c *ClichouseDB) UpdateBatchTeamsIndexedStatus(ctx context.Context, ids []uint32) error {
	return c.conn.Exec(ctx, "ALTER TABLE Teams UPDATE is_indexed=1 WHERE team_id IN $1", ids)
}

func (c *ClichouseDB) UpdateBatchTournamentsIndexedStatus(ctx context.Context, ids []uint32) error {
	return c.conn.Exec(ctx, "ALTER TABLE Tournaments UPDATE is_indexed=1 WHERE tournament_id IN $1", ids)
}
