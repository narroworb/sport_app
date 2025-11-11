package models

import "time"

type Match struct {
	ID          int
	IDAppDB  uint32
	Date        time.Time
	HomeTeam    *Team
	AwayTeam    *Team
	HomeGoals   uint16
	AwayGoals   uint16
	Round       uint16
	HomeManager *Manager
	AwayManager *Manager
	Status      string
}

type Season struct {
	Year  string
	Table map[uint8]TableRow
	Tours map[uint16][]Match
	Teams map[string]*Team
}

type TableRow struct {
	Team          *Team
	Points        uint16
	Pos           uint8
	Matches       uint16
	Wins          uint16
	Loses         uint16
	Draws         uint16
	ScoresFor     uint16
	ScoresAgainst uint16
}

type Competition struct {
	Country string
	Seasons map[string]Season
}

type Team struct {
	Name string
	ID   uint32
}

type Player struct {
	FirstName     string
	LastName      string
	ID            uint32
	Position      string
	DateOfBirth   time.Time
	Height        uint16
	PreferredFoot string
	Nation        string
	CurrentStatus string
}

type PlayerStatsInMatch struct {
	ID       uint32
	IDMatch  uint32
	IDPlayer uint32

	StartPlayer bool

	Rating float32

	MinutesPlayed uint16
	Goals         uint8
	Assists       uint8

	BlockedShots  uint8
	Interceptions uint8
	TotalTackles  uint8
	DribbledPast  uint8

	Duels     uint8
	DuelsWon  uint8
	Fouls     uint8
	WasFouled uint8

	PassAttempts   uint16
	CompletePasses uint16
	KeyPasses      uint8

	ShotsOnTarget    uint8
	TotalShots       uint8
	DribbleAttempts  uint8
	CompleteDribbles uint8

	PenaltyScored uint8
	PenaltyMissed uint8

	YellowCards uint8
	RedCards    uint8

	Captain        bool
	HomeTeamPlayer bool
}

type GoalieStatsInMatch struct {
	ID       uint32
	IDMatch  uint32
	IDPlayer uint32

	StartPlayer bool

	Rating float32

	MinutesPlayed uint16
	Goals         uint8
	Assists       uint8
	GoalsConceded uint8

	Saves          uint8
	PassAttempts   uint16
	CompletePasses uint16
	KeyPasses      uint8

	PenaltySaved    uint8
	PenaltyConceded uint8

	Fouls     uint8
	WasFouled uint8

	YellowCards uint8
	RedCards    uint8

	Captain        bool
	HomeTeamPlayer bool
}

type TeamMatchStats struct {
	IDMatch            uint32
	ShotsOnGoalHome    uint16
	ShotsOnGoalAway    uint16
	TotalShotsHome     uint16
	TotalShotsAway     uint16
	BlockedShotsHome   uint16
	BlockedShotsAway   uint16
	FoulsHome          uint16
	FoulsAway          uint16
	CornerKicksHome    uint16
	CornerKicksAway    uint16
	BallPossessionHome uint8
	BallPossessionAway uint8
	YellowCardsHome    uint8
	YellowCardsAway    uint8
	RedCardsHome       uint8
	RedCardsAway       uint8
	TotalPassesHome    uint16
	TotalPassesAway    uint16
	CompletePassesHome uint16
	CompletePassesAway uint16
	OffsidesHome       uint8
	OffsidesAway       uint8
	ShotsInsideBoxHome uint8
	ShotsInsideBoxAway uint8
}

type Manager struct {
	FirstName   string
	LastName    string
	ID          uint32
	Nation      string
	YellowCards uint16
	RedCards    uint16
}
