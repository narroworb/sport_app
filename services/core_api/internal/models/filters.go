package models

import "time"

type PlayerStatsFilter struct {
	PlayerID uint32
	Season   string
	FromDate time.Time
	ToDate   time.Time
}

type TeamStatsFilter struct {
	TeamID uint32
	Season   string
	FromDate time.Time
	ToDate   time.Time
}
