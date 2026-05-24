package models

import "time"

type Player struct {
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	ID            uint32    `json:"athlete_id"`
	Position      string    `json:"position"`
	DateOfBirth   time.Time `json:"date_of_birth"`
	Height        uint16    `json:"height"`
	PreferredFoot string    `json:"preffered_foot"`
	Nation        Country   `json:"nation"`
	CurrentStatus string    `json:"current_status,omitempty"`
	URLPhoto      string    `json:"url_photo,omitempty"`
}

type Country struct {
	Name    string `json:"name,omitempty"`
	URLFlag string `json:"url_flag,omitempty"`
}

type Match struct {
	ID          uint32     `json:"match_id"`
	Date        time.Time  `json:"date"`
	HomeTeam    Team       `json:"home_team"`
	AwayTeam    Team       `json:"away_team"`
	HomeGoals   int16      `json:"home_team_score"`
	AwayGoals   int16      `json:"away_team_score"`
	Round       uint16     `json:"round"`
	HomeManager Manager    `json:"home_team_manager,omitempty"`
	AwayManager Manager    `json:"away_team_manager,omitempty"`
	Status      string     `json:"status"`
	Tournament  Tournament `json:"tournament,omitempty"`
}

type ShortMatch struct {
	ID         uint32     `json:"match_id"`
	Date       time.Time  `json:"date"`
	HomeTeam   Team       `json:"home_team"`
	AwayTeam   Team       `json:"away_team"`
	HomeGoals  int16      `json:"home_team_score"`
	AwayGoals  int16      `json:"away_team_score"`
	Round      uint16     `json:"round"`
	Status     string     `json:"status"`
	Tournament Tournament `json:"tournament,omitempty"`
}

type Tournament struct {
	ID      uint32  `json:"tournament_id"`
	Name    string  `json:"name,omitempty"`
	Country Country `json:"country"`
	Season  string  `json:"season,omitempty"`
	URLLogo string  `json:"url_logo,omitempty"`
}

type Team struct {
	Name    string `json:"name,omitempty"`
	ID      uint32 `json:"team_id"`
	URLLogo string `json:"url_logo,omitempty"`
}

type PlayerStatsInMatch struct {
	ID      uint32 `json:"stat_id"`
	IDMatch uint32 `json:"match_id"`
	Player  Player `json:"athlete"`

	StartPlayer bool `json:"start_player"`

	Rating float64 `json:"rating"`

	MinutesPlayed uint16 `json:"minutes_played"`
	Goals         uint8  `json:"goals"`
	Assists       uint8  `json:"assists"`

	BlockedShots  uint8 `json:"blocked_shots"`
	Interceptions uint8 `json:"interceptions"`
	TotalTackles  uint8 `json:"total_tackles"`
	DribbledPast  uint8 `json:"dribbled_past"`

	Duels     uint8 `json:"duels"`
	DuelsWon  uint8 `json:"duels_won"`
	Fouls     uint8 `json:"fouls"`
	WasFouled uint8 `json:"was_fouled"`

	PassAttempts   uint16 `json:"pass_attempts"`
	CompletePasses uint16 `json:"complete_passes"`
	KeyPasses      uint8  `json:"key_passes"`

	ShotsOnTarget    uint8 `json:"shots_on_target"`
	TotalShots       uint8 `json:"total_shots"`
	DribbleAttempts  uint8 `json:"dribble_attempts"`
	CompleteDribbles uint8 `json:"complete_dribbles"`

	PenaltyScored uint8 `json:"penalty_scored"`
	PenaltyMissed uint8 `json:"penalty_missed"`

	YellowCards uint8 `json:"yellow_cards"`
	RedCards    uint8 `json:"red_cards"`

	Captain        bool `json:"captain,omitempty"`
	HomeTeamPlayer bool
}

type GoalieStatsInMatch struct {
	ID      uint32 `json:"stat_id"`
	IDMatch uint32 `json:"match_id"`
	Player  Player `json:"athlete"`

	StartPlayer bool `json:"start_player"`

	Rating float64 `json:"rating"`

	MinutesPlayed uint16 `json:"minutes_played"`
	Goals         uint8  `json:"goals"`
	Assists       uint8  `json:"assists"`
	GoalsConceded uint8  `json:"goals_conceded"`

	Saves          uint8  `json:"saves"`
	PassAttempts   uint16 `json:"pass_attempts"`
	CompletePasses uint16 `json:"complete_passes"`
	KeyPasses      uint8  `json:"key_passes"`

	PenaltySaved    uint8 `json:"penalty_saved"`
	PenaltyConceded uint8 `json:"penalty_conceded"`

	Fouls     uint8 `json:"fouls"`
	WasFouled uint8 `json:"was_fouled"`

	YellowCards uint8 `json:"yellow_cards"`
	RedCards    uint8 `json:"red_cards"`

	Captain        bool `json:"captain,omitempty"`
	HomeTeamPlayer bool
}

type TeamMatchStats struct {
	IDMatch            uint32 `json:"match_id"`
	ShotsOnGoalHome    uint16 `json:"shots_on_goal_home_team"`
	ShotsOnGoalAway    uint16 `json:"shots_on_goal_away_team"`
	TotalShotsHome     uint16 `json:"total_shots_home_team"`
	TotalShotsAway     uint16 `json:"total_shots_away_team"`
	BlockedShotsHome   uint16 `json:"blocked_shots_home_team"`
	BlockedShotsAway   uint16 `json:"blocked_shots_away_team"`
	FoulsHome          uint16 `json:"fouls_home_team"`
	FoulsAway          uint16 `json:"fouls_away_team"`
	CornerKicksHome    uint16 `json:"corner_kicks_home_team"`
	CornerKicksAway    uint16 `json:"corner_kicks_away_team"`
	BallPossessionHome uint8  `json:"ball_possession_home_team"`
	BallPossessionAway uint8  `json:"ball_possession_away_team"`
	YellowCardsHome    uint8  `json:"yellow_cards_home_team"`
	YellowCardsAway    uint8  `json:"yellow_cards_away_team"`
	RedCardsHome       uint8  `json:"red_cards_home_team"`
	RedCardsAway       uint8  `json:"red_cards_away_team"`
	TotalPassesHome    uint16 `json:"total_passes_home_team"`
	TotalPassesAway    uint16 `json:"total_passes_away_team"`
	CompletePassesHome uint16 `json:"complete_passes_home_team"`
	CompletePassesAway uint16 `json:"complete_passes_away_team"`
	OffsidesHome       uint8  `json:"offsides_home_team"`
	OffsidesAway       uint8  `json:"offsides_away_team"`
	ShotsInsideBoxHome uint8  `json:"shots_inside_box_home_team"`
	ShotsInsideBoxAway uint8  `json:"shots_inside_box_away_team"`
}

type Manager struct {
	FirstName   string  `json:"first_name,omitempty"`
	LastName    string  `json:"last_name,omitempty"`
	ID          uint32  `json:"manager_id"`
	Nation      Country `json:"nation"`
	YellowCards uint16  `json:"total_yellow_cards,omitempty"`
	RedCards    uint16  `json:"total_red_cards,omitempty"`
	URLPhoto    string  `json:"url_photo,omitempty"`
}

type PlayerStatsInPeriod struct {
	IDPlayer uint32 `json:"athlete_id"`

	TotalMatches  uint64  `json:"total_matches"`
	MatchesPlayed uint64  `json:"matches_played"`
	StartPlayer   uint64  `json:"start_player"`
	AvgRating     float64 `json:"avg_rating"`

	MinutesPlayed uint64  `json:"minutes_played"`
	Goals         uint64  `json:"goals"`
	GoalsPer90    float32 `json:"goals_per_90"`

	Assists      uint64  `json:"assists"`
	AssistsPer90 float32 `json:"assists_per_90"`

	GoalsConceded      uint64  `json:"goals_conceded"`
	GoalsConcededPer90 float32 `json:"goals_conceded_per_90"`

	Saves      uint64  `json:"saves"`
	SavesPer90 float32 `json:"saves_per_90"`

	BlockedShots      uint64  `json:"blocked_shots"`
	BlockedShotsPer90 float32 `json:"blocked_shots_per_90"`

	Interceptions      uint64  `json:"interceptions"`
	InterceptionsPer90 float32 `json:"interceptions_per_90"`

	TotalTackles      uint64  `json:"total_tackles"`
	TotalTacklesPer90 float32 `json:"total_tackles_per_90"`

	DribbledPast      uint64  `json:"dribbled_past"`
	DribbledPastPer90 float32 `json:"dribbled_past_per_90"`

	Duels      uint64  `json:"duels"`
	DuelsPer90 float32 `json:"duels_per_90"`

	DuelsWon      uint64  `json:"duels_won"`
	DuelsWonPer90 float32 `json:"duels_won_per_90"`

	PenaltySaved        uint64  `json:"penalty_saved"`
	PenaltyConceded     uint64  `json:"penalty_conceded"`
	PenaltySavedPercent float32 `json:"penalty_saves_percent"`

	Fouls      uint64  `json:"fouls"`
	FoulsPer90 float32 `json:"fouls_per_90"`

	WasFouled      uint64  `json:"was_fouled"`
	WasFouledPer90 float32 `json:"was_fouled_per_90"`

	PassAttempts      uint64  `json:"pass_attempts"`
	PassAttemptsPer90 float32 `json:"pass_attempts_per_90"`

	CompletePasses      uint64  `json:"complete_passes"`
	CompletePassesPer90 float32 `json:"complete_passes_per_90"`

	KeyPasses      uint64  `json:"key_passes"`
	KeyPassesPer90 float32 `json:"key_passes_per_90"`

	ShotsOnTarget      uint64  `json:"shots_on_target"`
	ShotsOnTargetPer90 float32 `json:"shots_on_target_per_90"`

	TotalShots      uint64  `json:"total_shots"`
	TotalShotsPer90 float32 `json:"total_shots_per_90"`

	DribbleAttempts      uint64  `json:"dribble_attempts"`
	DribbleAttemptsPer90 float32 `json:"dribble_attempts_per90"`

	CompleteDribbles      uint64  `json:"complete_dribbles"`
	CompleteDribblesPer90 float32 `json:"complete_dribbles_per_90"`

	DribbleAccuracy float32 `json:"dribble_accuracy"`

	PenaltyScored   uint64  `json:"penalty_scored"`
	PenaltyMissed   uint64  `json:"penalty_missed"`
	PenaltyAccuracy float32 `json:"penalty_accuracy"`

	YellowCards uint64 `json:"yellow_cards"`
	RedCards    uint64 `json:"red_cards"`

	CaptainTimes uint64 `json:"captain_times,omitempty"`
}

type PlayerMatch struct {
	ID                  uint32     `json:"match_id"`
	PlayerID            uint32     `json:"athlete_id"`
	Date                time.Time  `json:"date"`
	HomeTeam            Team       `json:"home_team"`
	AwayTeam            Team       `json:"away_team"`
	HomeGoals           int16      `json:"home_team_score"`
	AwayGoals           int16      `json:"away_team_score"`
	Round               uint16     `json:"round"`
	Status              string     `json:"status"`
	Tournament          Tournament `json:"tournament"`
	PlayerGoals         uint8      `json:"athlete_goals"`
	PlayerAssists       uint8      `json:"athlete_assists"`
	PlayerRedCards      uint8      `json:"athlete_red_cards"`
	PlayerRating        float64    `json:"athlete_rating"`
	PlayerMinutesPlayed uint16     `json:"athlete_minutes_played"`
}

type TeamStatsInPeriod struct {
	IDTeam uint32 `json:"team_id"`

	TotalMatches uint64  `json:"total_matches"`
	AvgWinRate   float64 `json:"avg_win_rate"`
	AvgPoints    float64 `json:"avg_points"`

	AvgBallPossession float64 `json:"average_ball_possession"`

	Goals      int64   `json:"goals"`
	GoalsPer90 float64 `json:"goals_per_90"`

	GoalsConceded      int64   `json:"goals_conceded"`
	GoalsConcededPer90 float64 `json:"goals_conceded_per_90"`

	ShotsOnGoal      uint64  `json:"shots_on_goal"`
	ShotsOnGoalPer90 float64 `json:"shots_on_goal_per_90"`

	ShotsOnGoalConceded      uint64  `json:"shots_on_goal_conceded"`
	ShotsOnGoalConcededPer90 float64 `json:"shots_on_goal_conceded_per_90"`

	TotalShots      uint64  `json:"total_shots"`
	TotalShotsPer90 float64 `json:"total_shots_per_90"`

	TotalShotsConceded      uint64  `json:"total_shots_conceded"`
	TotalShotsConcededPer90 float64 `json:"total_shots_conceded_per_90"`

	ShotsInsideBox      uint64  `json:"shots_inside_box"`
	ShotsInsideBoxPer90 float64 `json:"shots_inside_box_per_90"`

	ShotsInsideBoxConceded      uint64  `json:"shots_inside_box_conceded"`
	ShotsInsideBoxConcededPer90 float64 `json:"shots_inside_box_conceded_per_90"`

	BlockedShots      uint64  `json:"blocked_shots"`
	BlockedShotsPer90 float64 `json:"blocked_shots_per_90"`

	Fouls      uint64  `json:"fouls"`
	FoulsPer90 float64 `json:"fouls_per_90"`

	WasFouled      uint64  `json:"was_fouled"`
	WasFouledPer90 float64 `json:"was_fouled_per_90"`

	CornerKicks      uint64  `json:"corner_kicks"`
	CornerKicksPer90 float64 `json:"corner_kicks_per_90"`

	CornerKicksConceded      uint64  `json:"corner_kicks_conceded"`
	CornerKicksConcededPer90 float64 `json:"corner_kicks_conceded_per_90"`

	YellowCards      uint64  `json:"yellow_cards"`
	YellowCardsPer90 float64 `json:"yellow_cards_per_90"`

	YellowCardsOpp      uint64  `json:"yellow_cards_opponent"`
	YellowCardsOppPer90 float64 `json:"yellow_cards_opponent_per_90"`

	RedCards      uint64  `json:"red_cards"`
	RedCardsPer90 float64 `json:"red_cards_per_90"`

	RedCardsOpp      uint64  `json:"red_cards_opponent"`
	RedCardsOppPer90 float64 `json:"red_cards_opponent_per_90"`

	PassAttempts      uint64  `json:"pass_attempts"`
	PassAttemptsPer90 float64 `json:"pass_attempts_per_90"`

	PassAttemptsOpp      uint64  `json:"pass_attempts_opponent"`
	PassAttemptsOppPer90 float64 `json:"pass_attempts_opponent_per_90"`

	CompletePasses      uint64  `json:"complete_passes"`
	CompletePassesPer90 float64 `json:"complete_passes_per_90"`

	CompletePassesOpp      uint64  `json:"complete_passes_opp"`
	CompletePassesOppPer90 float64 `json:"complete_passes_opp_per_90"`

	AvgPassAccuracy    float64 `json:"average_pass_accuracy"`
	AvgPassAccuracyOpp float64 `json:"average_pass_accuracy_opponent"`

	Offsides      uint64  `json:"offsides"`
	OffsidesPer90 float64 `json:"offsides_per_90"`

	OffsidesOpp      uint64  `json:"offsides_opponent"`
	OffsidesOppPer90 float64 `json:"offsides_opponent_per_90"`
}

type TableRow struct {
	Team          Team   `json:"team"`
	Points        uint16 `json:"points"`
	Pos           uint8  `json:"position"`
	Matches       uint16 `json:"matches_played"`
	Wins          uint16 `json:"wins"`
	Draws         uint16 `json:"draws"`
	Losses        uint16 `json:"losses"`
	ScoresFor     uint16 `json:"goals_scored"`
	ScoresAgainst uint16 `json:"goals_conceded"`
}

type PlayerWithStats struct {
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	ID            uint32    `json:"athlete_id"`
	Position      string    `json:"position"`
	DateOfBirth   time.Time `json:"date_of_birth"`
	Height        uint16    `json:"height"`
	PreferredFoot string    `json:"preffered_foot"`
	Nation        Country   `json:"nation"`
	CurrentStatus string    `json:"current_status"`
	URLPhoto      string    `json:"url_photo"`
	TotalMatches  uint64    `json:"total_matches"`
	MatchesPlayed uint64    `json:"matches_played"`
	StartPlayer   uint64    `json:"start_player"`
	AvgRating     float64   `json:"avg_rating"`
	MinutesPlayed uint64    `json:"minutes_played"`
	Goals         uint64    `json:"goals"`
	Assists       uint64    `json:"assists"`
}

type ShortTableRow struct {
	Team       Team   `json:"team"`
	Points     uint64 `json:"points"`
	Pos        uint64 `json:"position"`
	Matches    uint16 `json:"matches_played"`
	Difference int64  `json:"goal_difference"`
}

type ManagerStatsInPeriod struct {
	IDManager uint32 `json:"manager_id"`

	TotalMatches  uint64  `json:"total_matches"`
	WinPercentage float64 `json:"win_percentage"`
	AvgPoints     float64 `json:"avg_points"`

	AvgBallPossession float64 `json:"average_ball_possession"`

	YellowCards uint16 `json:"yellow_cards"`
	RedCards    uint16 `json:"red_cards"`

	Goals      int64   `json:"goals"`
	GoalsPer90 float64 `json:"goals_per_90"`

	GoalsConceded      int64   `json:"goals_conceded"`
	GoalsConcededPer90 float64 `json:"goals_conceded_per_90"`

	ShotsOnGoal      uint64  `json:"shots_on_goal"`
	ShotsOnGoalPer90 float64 `json:"shots_on_goal_per_90"`

	ShotsOnGoalConceded      uint64  `json:"shots_on_goal_conceded"`
	ShotsOnGoalConcededPer90 float64 `json:"shots_on_goal_conceded_per_90"`

	TotalShots      uint64  `json:"total_shots"`
	TotalShotsPer90 float64 `json:"total_shots_per_90"`

	TotalShotsConceded      uint64  `json:"total_shots_conceded"`
	TotalShotsConcededPer90 float64 `json:"total_shots_conceded_per_90"`

	ShotsInsideBox      uint64  `json:"shots_inside_box"`
	ShotsInsideBoxPer90 float64 `json:"shots_inside_box_per_90"`

	ShotsInsideBoxConceded      uint64  `json:"shots_inside_box_conceded"`
	ShotsInsideBoxConcededPer90 float64 `json:"shots_inside_box_conceded_per_90"`

	BlockedShots      uint64  `json:"blocked_shots"`
	BlockedShotsPer90 float64 `json:"blocked_shots_per_90"`

	Fouls      uint64  `json:"fouls"`
	FoulsPer90 float64 `json:"fouls_per_90"`

	WasFouled      uint64  `json:"was_fouled"`
	WasFouledPer90 float64 `json:"was_fouled_per_90"`

	CornerKicks      uint64  `json:"corner_kicks"`
	CornerKicksPer90 float64 `json:"corner_kicks_per_90"`

	CornerKicksConceded      uint64  `json:"corner_kicks_conceded"`
	CornerKicksConcededPer90 float64 `json:"corner_kicks_conceded_per_90"`

	PassAttempts      uint64  `json:"pass_attempts"`
	PassAttemptsPer90 float64 `json:"pass_attempts_per_90"`

	PassAttemptsOpp      uint64  `json:"pass_attempts_opponent"`
	PassAttemptsOppPer90 float64 `json:"pass_attempts_opponent_per_90"`

	CompletePasses      uint64  `json:"complete_passes"`
	CompletePassesPer90 float64 `json:"complete_passes_per_90"`

	CompletePassesOpp      uint64  `json:"complete_passes_opp"`
	CompletePassesOppPer90 float64 `json:"complete_passes_opp_per_90"`

	AvgPassAccuracy    float64 `json:"average_pass_accuracy"`
	AvgPassAccuracyOpp float64 `json:"average_pass_accuracy_opponent"`

	Offsides      uint64  `json:"offsides"`
	OffsidesPer90 float64 `json:"offsides_per_90"`

	OffsidesOpp      uint64  `json:"offsides_opponent"`
	OffsidesOppPer90 float64 `json:"offsides_opponent_per_90"`
}

type ManagerMatch struct {
	ID         uint32     `json:"match_id"`
	ManagerID  uint32     `json:"manager_id"`
	Date       time.Time  `json:"date"`
	HomeTeam   Team       `json:"home_team"`
	AwayTeam   Team       `json:"away_team"`
	HomeGoals  int16      `json:"home_team_score"`
	AwayGoals  int16      `json:"away_team_score"`
	Round      uint16     `json:"round"`
	Status     string     `json:"status"`
	Tournament Tournament `json:"tournament"`
}

type TeamSeason struct {
	ID      uint32 `json:"team_id"`
	Name    string `json:"name,omitempty"`
	URLLogo string `json:"url_logo,omitempty"`
	Season  string `json:"season,omitempty"`
}

type ManagerTeams struct {
	Season string `json:"season"`
	Teams  Team   `json:"team"`
}
