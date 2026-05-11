package handlers

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/narroworb/core_api/internal/models"
)

const (
	cacheTimeout = 1 * time.Second
	dbTimeout    = 3 * time.Second
	cacheTTL     = 15 * time.Minute
)

type HandlerRepo struct {
	adb      AnalyticDatabaseInterface
	tdb      TransactionDatabaseInterface
	cacheDB  CacheInterface
	searchDB SearchDatabaseInterface
}

func NewHandlerRepo(adb AnalyticDatabaseInterface, tdb TransactionDatabaseInterface, cacheDB CacheInterface, searchDB SearchDatabaseInterface) *HandlerRepo {
	return &HandlerRepo{adb: adb, tdb: tdb, cacheDB: cacheDB, searchDB: searchDB}
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
	query := r.URL.Query().Get("q")
	if strings.TrimSpace(query) == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "query parameter 'q' is required",
		})
		return
	}

	entityType := r.URL.Query().Get("type")
	if entityType == "" {
		entityType = "all"
	}

	page := int32(1)
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.ParseInt(pageStr, 10, 32); err == nil {
			page = int32(p)
		}
	}

	pageSize := int32(10)
	if sizeStr := r.URL.Query().Get("size"); sizeStr != "" {
		if s, err := strconv.ParseInt(sizeStr, 10, 32); err == nil && s > 0 && s <= 100 {
			pageSize = int32(s)
		}
	}

	filters := models.SearchFilters{
		Query:    query,
		Position: r.URL.Query().Get("position"),
		Nation:   r.URL.Query().Get("nation"),
		Season:   r.URL.Query().Get("season"),
		Page:     page,
		PageSize: pageSize,
	}

	cacheKey := fmt.Sprintf("search:%s", hashSearchFilters(filters))

	ctx, cancel := context.WithTimeout(r.Context(), cacheTimeout)
	defer cancel()

	if cached, err := h.cacheDB.Get(ctx, cacheKey); err == nil {
		var cachedResults []models.SearchResult
		if err := json.Unmarshal([]byte(cached), &cachedResults); err == nil {
			h.writeJSON(w, http.StatusOK, models.SearchResponse{
				Query:   query,
				Total:   int64(len(cachedResults)),
				Results: cachedResults,
			})
			return
		}
	}

	ctx, cancel = context.WithTimeout(r.Context(), dbTimeout)
	defer cancel()

	results, err := h.searchDB.Search(ctx, query, entityType, filters)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "search failed",
		})
		return
	}

	response := models.SearchResponse{
		Query:   query,
		Total:   int64(len(results)),
		Results: results,
	}

	if resultJSON, err := json.Marshal(results); err == nil {
		ctx, cancel := context.WithTimeout(r.Context(), cacheTimeout)
		defer cancel()
		_ = h.cacheDB.Set(ctx, cacheKey, resultJSON, cacheTTL)
	}

	h.writeJSON(w, http.StatusOK, response)
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

func hashSearchFilters(filters models.SearchFilters) string {
	data := fmt.Sprintf("%s:%s:%s:%s:%d:%d",
		filters.Query,
		filters.Position,
		filters.Nation,
		filters.Season,
		filters.Page,
		filters.PageSize,
	)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}
