package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	cacheTimeout = 1 * time.Second
	dbTimeout    = 3 * time.Second
	cacheTTL     = 15 * time.Minute
)

type HandlerRepo struct {
	adb     AnalyticDatabaseInterface
	tdb     TransactionDatabaseInterface
	cacheDB CacheInterface
}

func NewHandlerRepo(adb AnalyticDatabaseInterface, tdb TransactionDatabaseInterface, cacheDB CacheInterface) *HandlerRepo {
	return &HandlerRepo{adb: adb, tdb: tdb, cacheDB: cacheDB}
}

func (h *HandlerRepo) writeJSON(w http.ResponseWriter, status int, data any) ([]byte, error) {
	w.Header().Set("Content-Type", "application/json")

	if data == nil {
		data = map[string]struct{}{}
	}

	resp, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
		return nil, err
	}

	w.WriteHeader(status)
	_, err = w.Write(resp)
	return resp, err
}

func (h *HandlerRepo) Search(w http.ResponseWriter, r *http.Request) {
}

func validateAndFormatSeason(season string) string {
	season = strings.TrimSpace(season)

	if strings.Contains(season, "/") {
		parts := strings.Split(season, "/")
		if len(parts) == 2 && len(parts[0]) == 4 && len(parts[1]) == 4 {
			return parts[0] + "/" + parts[1]
		}
	}

	if strings.Contains(season, "-") {
		parts := strings.Split(season, "-")
		if len(parts) == 2 && len(parts[0]) == 4 && len(parts[1]) == 4 {
			return parts[0] + "/" + parts[1]
		}
	}

	if strings.Contains(season, "\\") {
		parts := strings.Split(season, "\\")
		if len(parts) == 2 && len(parts[0]) == 4 && len(parts[1]) == 4 {
			return parts[0] + "/" + parts[1]
		}
	}

	return ""
}

func getCurrentSeason() string {
	timeNow := time.Now()
	curYear := timeNow.Year()

	if timeNow.Month() > time.June {
		return fmt.Sprintf("%d/%d", curYear, curYear+1)
	} else {
		return fmt.Sprintf("%d/%d", curYear-1, curYear)
	}
}

func parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("empty string")
	}

	formats := []string{"2006-01-02", "02.01.2006", "2006/01/02"}
	for _, format := range formats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return date, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse this format of date: %s", dateStr)
}
