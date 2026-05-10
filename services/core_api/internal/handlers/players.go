package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/narroworb/core_api/internal/middleware"
	"github.com/narroworb/core_api/internal/models"
)

func (h *HandlerRepo) GetPlayerDetails(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty id parameter in query"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id parameter in query"})
		return
	}

	cacheKey := "player:" + fmt.Sprint(id) + ":details"
	ctxCache, cancelCache := context.WithTimeout(r.Context(), cacheTimeout)
	defer cancelCache()
	cached, err := h.cacheDB.Get(ctxCache, cacheKey)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cached))
		return
	}

	ctxDB, cancelDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelDB()
	playerDetails, err := h.adb.GetPlayerByID(ctxDB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "player not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetPlayerDetails: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, playerDetails)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetPlayerDetails: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetPlayerDetails: %v\n", err)
	}
}

func (h *HandlerRepo) GetPlayerStats(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty id parameter in query"})
		return
	}

	var filter models.PlayerStatsFilter

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id parameter in query"})
		return
	}

	ctxDB, cancelDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelDB()
	playerPosition, err := h.adb.GetPlayerPositionByID(ctxDB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "player not found"})
			return
		}
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetPlayerPositionByID: %v\n", err)
		return
	}

	filter.PlayerID = uint32(id)

	seasonParam := r.URL.Query().Get("season")
	if seasonParam != "" {
		filter.Season = validateAndFormatSeason(seasonParam)
		if filter.Season == "" {
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid season parameter in query, use YYYY/YYYY or YYYY-YYYY"})
			return
		}
	}

	dateFromStr := r.URL.Query().Get("date_from")
	dateToStr := r.URL.Query().Get("date_to")

	dateFrom, _ := parseDate(dateFromStr)
	dateTo, _ := parseDate(dateToStr)

	if (dateFrom != time.Time{} || dateTo != time.Time{}) {
		if !dateTo.After(dateFrom) {
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date parameters in query, date_from should be less than date_to"})
			return
		}
		if (dateFrom != time.Time{}) {
			filter.FromDate = dateFrom
		} else {
			filter.FromDate, _ = time.Parse("2006-01-02", "2014-01-01")
		}
		if (dateTo != time.Time{}) {
			filter.ToDate = dateTo
		} else {
			filter.ToDate = time.Now()
		}
	}

	var getStatsFromDB func(context.Context, models.PlayerStatsFilter) (models.PlayerStatsInPeriod, error)
	var cacheKey string

	if filter.Season != "" {
		if playerPosition == "G" {
			getStatsFromDB = h.adb.GetGoalieStatsBySeason
		} else {
			getStatsFromDB = h.adb.GetPlayerStatsBySeason
		}
		cacheKey = "player:" + fmt.Sprint(id) + ":stats:season:" + filter.Season
	} else if (filter.ToDate != time.Time{} || filter.FromDate != time.Time{}) {
		if playerPosition == "G" {
			getStatsFromDB = h.adb.GetGoalieStatsByDates
		} else {
			getStatsFromDB = h.adb.GetPlayerStatsByDates
		}
		cacheKey = "player:" + fmt.Sprint(id) + ":stats:dates:" + fmt.Sprintf("%s_%s", filter.FromDate, filter.ToDate)
	} else {
		if playerPosition == "G" {
			getStatsFromDB = h.adb.GetGoalieFullStats
		} else {
			getStatsFromDB = h.adb.GetPlayerFullStats
		}
		cacheKey = "player:" + fmt.Sprint(id) + ":stats:full"
	}

	ctxCache, cancelCache := context.WithTimeout(r.Context(), cacheTimeout)
	defer cancelCache()
	cached, err := h.cacheDB.Get(ctxCache, cacheKey)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cached))
		return
	}

	ctxMainDB, cancelMainDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelMainDB()
	playerStats, err := getStatsFromDB(ctxMainDB, filter)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "player stats not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetPlayerStats: %v\n", err)
		return
	}

	playerStats = completePlayerStats(playerStats)

	resp, err := h.writeJSON(w, http.StatusOK, playerStats)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetPlayerStats: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetPlayerStats: %v\n", err)
		fmt.Printf("%+v\n", playerStats)
	}
}

func completePlayerStats(stats models.PlayerStatsInPeriod) models.PlayerStatsInPeriod {
	if stats.MinutesPlayed == 0 {
		return stats
	}
	stats.GoalsPer90 = float32(stats.Goals) / float32(stats.MinutesPlayed)
	stats.AssistsPer90 = float32(stats.Assists) / float32(stats.MinutesPlayed)
	stats.GoalsConcededPer90 = float32(stats.GoalsConceded) / float32(stats.MinutesPlayed)
	stats.SavesPer90 = float32(stats.Saves) / float32(stats.MinutesPlayed)
	stats.BlockedShotsPer90 = float32(stats.BlockedShots) / float32(stats.MinutesPlayed)
	stats.InterceptionsPer90 = float32(stats.Interceptions) / float32(stats.MinutesPlayed)
	stats.TotalTacklesPer90 = float32(stats.TotalTackles) / float32(stats.MinutesPlayed)
	stats.DribbledPastPer90 = float32(stats.DribbledPast) / float32(stats.MinutesPlayed)
	stats.DuelsPer90 = float32(stats.Duels) / float32(stats.MinutesPlayed)
	stats.DuelsWonPer90 = float32(stats.DuelsWon) / float32(stats.MinutesPlayed)
	if stats.PenaltySaved+stats.PenaltyConceded != 0 {
		stats.PenaltySavedPercent = float32(stats.PenaltySaved) / float32(stats.PenaltySaved+stats.PenaltyConceded)
	}
	stats.FoulsPer90 = float32(stats.Fouls) / float32(stats.MinutesPlayed)
	stats.WasFouledPer90 = float32(stats.WasFouled) / float32(stats.MinutesPlayed)
	stats.PassAttemptsPer90 = float32(stats.PassAttempts) / float32(stats.MinutesPlayed)
	stats.CompletePassesPer90 = float32(stats.CompletePasses) / float32(stats.MinutesPlayed)
	stats.KeyPassesPer90 = float32(stats.KeyPasses) / float32(stats.MinutesPlayed)
	stats.ShotsOnTargetPer90 = float32(stats.ShotsOnTarget) / float32(stats.MinutesPlayed)
	stats.TotalShotsPer90 = float32(stats.TotalShots) / float32(stats.MinutesPlayed)
	stats.DribbleAttemptsPer90 = float32(stats.DribbleAttempts) / float32(stats.MinutesPlayed)
	stats.CompleteDribblesPer90 = float32(stats.CompleteDribbles) / float32(stats.MinutesPlayed)
	if stats.DribbleAttempts != 0 {
		stats.DribbleAccuracy = float32(stats.CompleteDribbles) / float32(stats.DribbleAttempts)
	}
	if stats.PenaltyMissed+stats.PenaltyScored != 0 {
		stats.PenaltyAccuracy = float32(stats.PenaltyScored) / float32(stats.PenaltyMissed+stats.PenaltyScored)
	}

	return stats
}

func (h *HandlerRepo) GetPlayerFixtures(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty id parameter in query"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id parameter in query"})
		return
	}

	var limit, offset int

	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid limit parameter in query, use only positive digits"})
			return
		}
	}

	offsetStr := r.URL.Query().Get("offset")
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 1 {
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid offset parameter in query, use only positive digits"})
			return
		}
	}

	if offset > 0 && limit == 0 {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "offset parameter in query is set, but limit parameter empty"})
		return
	}

	cacheKey := fmt.Sprintf("player:%d:fixtures:limit:%d:offset:%d", id, limit, offset)

	ctxCache, cancelCache := context.WithTimeout(r.Context(), cacheTimeout)
	defer cancelCache()
	cached, err := h.cacheDB.Get(ctxCache, cacheKey)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cached))
		return
	}

	ctxMainDB, cancelMainDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelMainDB()
	playerFixtures, err := h.adb.GetPlayerFixtures(ctxMainDB, uint32(id), uint32(limit), uint32(offset))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "player fixtures not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetPlayerFixtures: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, playerFixtures)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetPlayerFixtures: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetPlayerFixtures: %v\n", err)
		fmt.Printf("%+v\n", playerFixtures)
	}
}

func (h *HandlerRepo) GetPlayerTeams(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty id parameter in query"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id parameter in query"})
		return
	}

	cacheKey := fmt.Sprintf("player:%d:teams", id)

	ctxCache, cancelCache := context.WithTimeout(r.Context(), cacheTimeout)
	defer cancelCache()
	cached, err := h.cacheDB.Get(ctxCache, cacheKey)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cached))
		return
	}

	ctxMainDB, cancelMainDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelMainDB()
	playerTeams, err := h.adb.GetPlayerTeams(ctxMainDB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "player teams not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetPlayerTeams: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, playerTeams)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetPlayerTeams: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetPlayerTeams: %v\n", err)
		fmt.Printf("%+v\n", playerTeams)
	}
}

func (h *HandlerRepo) GetFavouritePlayers(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	cacheKey := fmt.Sprintf("user:%d:players", userID)

	ctxCache, cancelCache := context.WithTimeout(r.Context(), cacheTimeout)
	defer cancelCache()
	cached, err := h.cacheDB.Get(ctxCache, cacheKey)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cached))
		return
	}

	ctxTransactDB, cancelTransactDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelTransactDB()
	playerIDs, err := h.tdb.GetFavoritePlayersIDs(ctxTransactDB, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "favorite players not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetFavoritePlayersIDs: %v\n", err)
		return
	}

	players := make([]models.Player, 0, len(playerIDs))

	for _, id := range playerIDs {
		ctxAnalyticDB, cancelAnalyticDB := context.WithTimeout(r.Context(), dbTimeout)
		defer cancelAnalyticDB()
		player, err := h.adb.GetPlayerByID(ctxAnalyticDB, id)
		if err != nil {
			h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
			log.Printf("error in GetPlayerByID: %v\n", err)
			return
		}

		players = append(players, player)
	}

	resp, err := h.writeJSON(w, http.StatusOK, players)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetFavoritePlayers: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetFavoritePlayers: %v\n", err)
		fmt.Printf("%+v\n", players)
	}
}

func (h *HandlerRepo) SetFavouritePlayer(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty id parameter in query"})
		return
	}

	athleteID, err := strconv.Atoi(idStr)
	if err != nil || athleteID < 1 {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id parameter in query"})
		return
	}

	ctxDB, cancelDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelDB()
	_, err = h.adb.GetPlayerByID(ctxDB, uint32(athleteID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "player not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetPlayerByID: %v\n", err)
		return
	}

	ctxTransactDB, cancelTransactDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelTransactDB()
	err = h.tdb.SetFavoritePlayerByID(ctxTransactDB, userID, int64(athleteID))
	if err != nil {
		if err.Error() == fmt.Sprintf("player with ID=%d already favorite for user with ID=%d", athleteID, userID) {
			h.writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in SetFavoritePlayerByID: %v\n", err)
		return
	}

	_, err = h.writeJSON(w, http.StatusCreated, nil)
	if err != nil {
		log.Printf("error in writeJSON from SetFavoritePlayer: %v\n", err)
	}
}

func (h *HandlerRepo) DeleteFavouritePlayer(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty id parameter in query"})
		return
	}

	athleteID, err := strconv.Atoi(idStr)
	if err != nil || athleteID < 1 {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id parameter in query"})
		return
	}

	ctxDB, cancelDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelDB()
	_, err = h.adb.GetPlayerByID(ctxDB, uint32(athleteID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "player not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetPlayerByID: %v\n", err)
		return
	}

	ctxTransactDB, cancelTransactDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelTransactDB()
	err = h.tdb.DeleteFavoritePlayerByID(ctxTransactDB, userID, int64(athleteID))
	if err != nil {
		if err.Error() == fmt.Sprintf("player with ID=%d already not favorite for user with ID=%d", athleteID, userID) {
			h.writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in DeleteFavoritePlayerByID: %v\n", err)
		return
	}

	_, err = h.writeJSON(w, http.StatusOK, nil)
	if err != nil {
		log.Printf("error in writeJSON from DeleteFavoritePlayer: %v\n", err)
	}
}
