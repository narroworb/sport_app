package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/narroworb/core_api/internal/models"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClichouseDB struct {
	conn driver.Conn
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

	return c, nil
}

func (c *ClichouseDB) Close() {
	c.conn.Close()
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
	row := c.conn.QueryRow(ctx, "SELECT tournament_id, date, home_team_id, away_team_id, home_score, away_score, round, home_team_manager_id, away_team_manager_id, status FROM Matches WHERE match_id=$1", id)

	var match models.Match
	if err := row.Scan(&match.Tournament.ID, &match.Date, &match.HomeTeam.ID, &match.AwayTeam.ID, &match.HomeGoals, &match.AwayGoals, &match.Round, &match.HomeManager.ID, &match.AwayManager.ID, &match.Status); err != nil {
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

func (c *ClichouseDB) GetPlayerFixtures(ctx context.Context, id, limit, offset uint32) ([]models.PlayerMatch, error) {
	var rows driver.Rows
	var err error
	if limit > 0 {
		rows, err = c.conn.Query(ctx, `
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
	WHERE athlete_id=$1
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
	WHERE athlete_id=$1
	ORDER BY m.date DESC
	LIMIT $2
	OFFSET $3`,
			id, limit, offset)
	} else {
		rows, err = c.conn.Query(ctx, `
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
	WHERE athlete_id=$1
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
	WHERE athlete_id=$1
	ORDER BY m.date DESC`, id)
		limit = 100
	}

	if err != nil {
		return nil, err
	}

	matches := make([]models.PlayerMatch, 0, limit)

	for rows.Next() {
		var match models.PlayerMatch
		if err := rows.Scan(&match.ID, &match.Date, &match.HomeTeam.ID, &match.HomeTeam.Name, &match.HomeTeam.URLLogo, &match.AwayTeam.ID, &match.AwayTeam.Name, &match.AwayTeam.URLLogo,
			&match.HomeGoals, &match.AwayGoals, &match.Round, &match.Status, &match.Tournament.ID, &match.Tournament.Name, &match.PlayerGoals, &match.PlayerAssists,
			&match.PlayerRedCards, &match.PlayerRating, &match.PlayerMinutesPlayed); err != nil {
			return nil, err
		}
		match.PlayerID = id

		matches = append(matches, match)
	}

	return matches, nil
}

func (c *ClichouseDB) GetPlayerTeams(ctx context.Context, id uint32) ([]models.Team, error) {
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
	SELECT t.team_id, name, logo_url FROM Matches
	INNER JOIN last_match ON last_match.match_id=Matches.match_id 
	INNER JOIN Teams t ON t.team_id=home_team_id
	INNER JOIN Team_Logos tl ON t.team_id=tl.team_id
	WHERE home_team_player
	UNION DISTINCT
	SELECT t.team_id, name, logo_url FROM Matches
	INNER JOIN last_match ON last_match.match_id=Matches.match_id 
	INNER JOIN Teams t ON t.team_id=away_team_id
	INNER JOIN Team_Logos tl ON t.team_id=tl.team_id
	WHERE home_team_player=0;`,
		id)
	if err != nil {
		return nil, err
	}

	teams := make([]models.Team, 0, 5)

	for rows.Next() {
		var team models.Team
		if err := rows.Scan(&team.ID, &team.Name, &team.URLLogo); err != nil {
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
	AVG(complete_passes/total_passes),
	AVG(complete_passes_conceded/total_passes_conceded),
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
	AVG(complete_passes/total_passes),
	AVG(complete_passes_conceded/total_passes_conceded),
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
	AVG(complete_passes/total_passes),
	AVG(complete_passes_conceded/total_passes_conceded),
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
	ht.team_id, ht.name, htl.logo_url, at.team_id, at.name, atl.logo_url, date, round, t.name, tl.logo_url FROM Matches m
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
	if err := row.Scan(&match.HomeTeam.ID, &match.HomeTeam.Name, &match.HomeTeam.URLLogo, &match.AwayTeam.ID, &match.AwayTeam.Name, &match.AwayTeam.URLLogo, &match.Date, &match.Round, &match.Tournament.Name, &match.Tournament.URLLogo); err != nil {
		return models.ShortMatch{}, err
	}
	match.ID = id

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

func (c *ClichouseDB) GetTeamLastGames(ctx context.Context, teamID uint32, limit, offset uint32) ([]models.ShortMatch, error) {
	rows, err := c.conn.Query(ctx, `
	SELECT 
	match_id, ht.team_id, ht.name, htl.logo_url, at.team_id, at.name, atl.logo_url, date, round, t.name, tl.logo_url, status FROM Matches m
	INNER JOIN Teams ht ON ht.team_id=m.home_team_id
	INNER JOIN Team_Logos htl ON ht.team_id=htl.team_id
	INNER JOIN Teams at ON at.team_id=m.away_team_id
	INNER JOIN Team_Logos atl ON at.team_id=atl.team_id
	INNER JOIN Tournaments t ON t.tournament_id=m.tournament_id
	INNER JOIN Tournament_Logos tl ON t.name=tl.tournament_name
	WHERE (home_team_id = $1 OR away_team_id = $1) AND date <= now()
	ORDER BY date DESC
	LIMIT $2
	OFFSET $3`,
		teamID, limit, offset)
	if err != nil {
		return nil, err
	}

	matches := make([]models.ShortMatch, 0, limit)

	for rows.Next() {
		var match models.ShortMatch
		if err := rows.Scan(&match.ID, &match.HomeTeam.ID, &match.HomeTeam.Name, &match.HomeTeam.URLLogo, &match.AwayTeam.ID, &match.AwayTeam.Name, &match.AwayTeam.URLLogo, &match.Date, &match.Round, &match.Tournament.Name, &match.Tournament.URLLogo, &match.Status); err != nil {
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
	AVG(complete_passes/total_passes),
	AVG(complete_passes_conceded/total_passes_conceded),
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
	AVG(complete_passes/total_passes),
	AVG(complete_passes_conceded/total_passes_conceded),
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
	AVG(complete_passes/total_passes),
	AVG(complete_passes_conceded/total_passes_conceded),
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

func (c *ClichouseDB) GetManagerTeams(ctx context.Context, id uint32) (map[string][]string, error) {
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
	name
	FROM all_teams_seasons a
	INNER JOIN Teams t ON t.team_id=a.team_id
	ORDER BY name, season`,
		id)
	if err != nil {
		return nil, err
	}

	allTeams := make(map[string][]string)

	for rows.Next() {
		var season, team string

		if err := rows.Scan(&season, &team); err != nil {
			return nil, err
		}

		if len(allTeams[team]) == 0 {
			allTeams[team] = make([]string, 0, 4)
		}
		allTeams[team] = append(allTeams[team], season)
	}

	return allTeams, nil
}

func (c *ClichouseDB) GetManagerFixtures(ctx context.Context, id, limit, offset uint32) ([]models.ManagerMatch, error) {
	var rows driver.Rows
	var err error
	if limit > 0 {
		rows, err = c.conn.Query(ctx, `
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
	INNER JOIN Tournaments t ON m.tournament_id=t.tournament_id
	WHERE home_team_manager_id=$1 OR away_team_manager_id=$1
	ORDER BY m.date DESC
	LIMIT $2
	OFFSET $3`,
			id, limit, offset)
	} else {
		rows, err = c.conn.Query(ctx, `
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
	INNER JOIN Tournaments t ON m.tournament_id=t.tournament_id
	WHERE home_team_manager_id=$1 OR away_team_manager_id=$1
	ORDER BY m.date DESC`, id)
		limit = 100
	}

	if err != nil {
		return nil, err
	}

	matches := make([]models.ManagerMatch, 0, limit)

	for rows.Next() {
		var match models.ManagerMatch
		if err := rows.Scan(&match.ID, &match.Date, &match.HomeTeam.ID, &match.HomeTeam.Name, &match.HomeTeam.URLLogo, &match.AwayTeam.ID, &match.AwayTeam.Name, &match.AwayTeam.URLLogo,
			&match.HomeGoals, &match.AwayGoals, &match.Round, &match.Status, &match.Tournament.ID, &match.Tournament.Name); err != nil {
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
	ht.team_id, ht.name, htl.logo_url, at.team_id, at.name, atl.logo_url, date, round, t.name, tl.logo_url, home_score, away_score, status FROM Matches m
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
		if err := rows.Scan(&match.HomeTeam.ID, &match.HomeTeam.Name, &match.HomeTeam.URLLogo, &match.AwayTeam.ID, &match.AwayTeam.Name, &match.AwayTeam.URLLogo, &match.Date, &match.Round, &match.Tournament.Name, &match.Tournament.URLLogo,
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
