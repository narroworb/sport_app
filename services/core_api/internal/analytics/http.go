package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	analyticsv1 "github.com/narroworb/core_api/gen/analytics/v1"
)

type HTTPHandler struct {
	client analyticsv1.AnalyticsServiceClient
}

func NewHTTPHandler(client analyticsv1.AnalyticsServiceClient) *HTTPHandler {
	return &HTTPHandler{client: client}
}

func (h *HTTPHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/health", h.Health)
	r.Get("/match_win_probabilities", h.GetMatchWinProbabilities)
	r.Get("/team_form", h.GetTeamFormIndex)
	r.Get("/player_similarity", h.GetPlayerSimilarityTopK)
	return r
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func parseUint32(q string, name string) (uint32, error) {
	q = strings.TrimSpace(q)
	if q == "" {
		return 0, fmt.Errorf("missing %s", name)
	}
	v, err := strconv.ParseUint(q, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid %s", name)
	}
	return uint32(v), nil
}

func formatSeason(season string) string {
	season = strings.TrimSpace(season)
	if season == "" {
		return ""
	}
	season = strings.ReplaceAll(season, "-", "/")
	season = strings.ReplaceAll(season, "\\", "/")
	parts := strings.Split(season, "/")
	if len(parts) == 2 && len(parts[0]) == 4 && len(parts[1]) == 4 {
		return parts[0] + "/" + parts[1]
	}
	return ""
}

func (h *HTTPHandler) Health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	resp, err := h.client.Health(ctx, &analyticsv1.HealthRequest{})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{
			"error": err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *HTTPHandler) GetMatchWinProbabilities(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	homeID, err := parseUint32(q.Get("home_team_id"), "home_team_id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	awayID, err := parseUint32(q.Get("away_team_id"), "away_team_id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	season := formatSeason(q.Get("season"))
	if season == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid season"})
		return
	}

	matchesBack, _ := strconv.ParseUint(strings.TrimSpace(q.Get("matches_back")), 10, 32)
	if matchesBack == 0 {
		matchesBack = 10
	}
	maxGoals, _ := strconv.ParseUint(strings.TrimSpace(q.Get("max_goals")), 10, 32)
	if maxGoals == 0 {
		maxGoals = 10
	}

	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	resp, err := h.client.GetMatchWinProbabilities(ctx, &analyticsv1.GetMatchWinProbabilitiesRequest{
		HomeTeamId:  homeID,
		AwayTeamId:  awayID,
		Season:      season,
		MatchesBack: uint32(matchesBack),
		MaxGoals:    uint32(maxGoals),
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *HTTPHandler) GetTeamFormIndex(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	teamID, err := parseUint32(q.Get("team_id"), "team_id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	season := formatSeason(q.Get("season"))
	if season == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid season"})
		return
	}

	matchesBack, _ := strconv.ParseUint(strings.TrimSpace(q.Get("matches_back")), 10, 32)
	if matchesBack == 0 {
		matchesBack = 10
	}
	halfLife, _ := strconv.ParseFloat(strings.TrimSpace(q.Get("half_life_matches")), 64)
	if halfLife <= 0 {
		halfLife = 5.0
	}

	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	resp, err := h.client.GetTeamFormIndex(ctx, &analyticsv1.GetTeamFormIndexRequest{
		TeamId:          teamID,
		Season:          season,
		MatchesBack:     uint32(matchesBack),
		HalfLifeMatches: halfLife,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *HTTPHandler) GetPlayerSimilarityTopK(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	playerID, err := parseUint32(q.Get("player_id"), "player_id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	season := formatSeason(q.Get("season"))
	if season == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid season"})
		return
	}

	topK, _ := strconv.ParseUint(strings.TrimSpace(q.Get("top_k")), 10, 32)
	if topK == 0 {
		topK = 10
	}
	minMinutes, _ := strconv.ParseUint(strings.TrimSpace(q.Get("min_minutes")), 10, 32)
	if minMinutes == 0 {
		minMinutes = 600
	}

	ctx, cancel := context.WithTimeout(r.Context(), 6*time.Second)
	defer cancel()

	resp, err := h.client.GetPlayerSimilarityTopK(ctx, &analyticsv1.GetPlayerSimilarityTopKRequest{
		PlayerId:   playerID,
		Season:     season,
		TopK:       uint32(topK),
		MinMinutes: uint32(minMinutes),
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}
