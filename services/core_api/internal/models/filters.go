package models

import "time"

type PlayerStatsFilter struct {
	PlayerID uint32
	Season   string
	FromDate time.Time
	ToDate   time.Time
}

type TeamStatsFilter struct {
	TeamID   uint32
	Season   string
	FromDate time.Time
	ToDate   time.Time
}

type SearchFilters struct {
	Query    string
	Position string
	Nation   string
	Season   string
	Page     int32
	PageSize int32
}

type SearchResult struct {
	Type  string      `json:"type"`
	Score float32     `json:"score"`
	Data  interface{} `json:"data"`
}

type SearchResponse struct {
	Query   string         `json:"query"`
	Total   int64          `json:"total"`
	Results []SearchResult `json:"results"`
}

