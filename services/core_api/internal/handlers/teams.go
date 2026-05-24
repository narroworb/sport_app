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

func (h *HandlerRepo) GetTeamDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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

	cacheKey := "team:" + fmt.Sprint(id) + ":details"
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
	teamDetails, err := h.adb.GetTeamByID(ctxDB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "team not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetTeamDetails: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, teamDetails)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetTeamDetails: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetTeamDetails: %v\n", err)
	}
}

func (h *HandlerRepo) GetTeamStats(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty id parameter in query"})
		return
	}

	var filter models.TeamStatsFilter

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id parameter in query"})
		return
	}

	filter.TeamID = uint32(id)

	seasonParam := r.URL.Query().Get("season")
	if seasonParam != "" {
		filter.Season = validateAndFormatSeason(seasonParam)
		if filter.Season == "" {
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid season parameter in query, use YYYY/YYYY or YYYY-YYYY"})
			return
		}
	}

	dateFromStr := r.URL.Query().Get("dateFrom")
	dateToStr := r.URL.Query().Get("dateTo")

	var dateFrom, dateTo time.Time
	var hasDateFrom, hasDateTo bool

	if dateFromStr != "" {
		dateFrom, err = parseDate(dateFromStr)
		if err != nil {
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date_from format, use YYYY-MM-DD"})
			return
		}
		hasDateFrom = true
	}

	if dateToStr != "" {
		dateTo, err = parseDate(dateToStr)
		if err != nil {
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date_to format, use YYYY-MM-DD"})
			return
		}
		hasDateTo = true
	}

	// Валидация и установка дат
	if hasDateFrom && hasDateTo {
		if !dateTo.After(dateFrom) {
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "date_to must be greater than date_from"})
			return
		}
		filter.FromDate = dateFrom
		filter.ToDate = dateTo
	} else if hasDateFrom {
		// Только date_from: от этой даты до сегодня
		filter.FromDate = dateFrom
		filter.ToDate = time.Now()
	} else if hasDateTo {
		// Только date_to: с 2014 года до этой даты
		filter.FromDate, _ = time.Parse("2006-01-02", "2014-01-01")
		filter.ToDate = dateTo
	}

	var getStatsFromDB func(context.Context, models.TeamStatsFilter) (models.TeamStatsInPeriod, error)
	var cacheKey string

	if filter.Season != "" {
		getStatsFromDB = h.adb.GetTeamStatsBySeason
		cacheKey = "team:" + fmt.Sprint(id) + ":stats:season:" + filter.Season
	} else if hasDateFrom || hasDateTo { // Исправлено: проверяем наличие любого фильтра дат
		getStatsFromDB = h.adb.GetTeamStatsByDates
		cacheKey = "team:" + fmt.Sprint(id) + ":stats:dates:" + fmt.Sprintf("%s_%s",
			filter.FromDate.Format("2006-01-02"),
			filter.ToDate.Format("2006-01-02"))
	} else {
		getStatsFromDB = h.adb.GetTeamFullStats
		cacheKey = "team:" + fmt.Sprint(id) + ":stats:full"
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
	teamStats, err := getStatsFromDB(ctxMainDB, filter)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "team stats not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetTeamStats: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, teamStats)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetTeamStats: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetTeamStats: %v\n", err)
		fmt.Printf("%+v\n", teamStats)
	}
}

func (h *HandlerRepo) GetTeamStandings(w http.ResponseWriter, r *http.Request) {
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

	seasonParam := r.URL.Query().Get("season")
	if seasonParam != "" {
		seasonParam = validateAndFormatSeason(seasonParam)
		if seasonParam == "" {
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid season parameter in query, use YYYY/YYYY or YYYY-YYYY"})
			return
		}
	} else {
		seasonParam = getCurrentSeason()
	}

	cacheKey := fmt.Sprintf("team:%d:stats:season:%s", id, seasonParam)

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
	teamStats, err := h.adb.GetStandingsByTeamAndSeason(ctxMainDB, uint32(id), seasonParam)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "team standings not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetStandingsByTeamAndSeason: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, teamStats)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetStandingsByTeamAndSeason: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetStandingsByTeamAndSeason: %v\n", err)
		fmt.Printf("%+v\n", teamStats)
	}
}

func (h *HandlerRepo) GetTeamNextGame(w http.ResponseWriter, r *http.Request) {
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

	cacheKey := fmt.Sprintf("team:%d:nextgame", id)

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
	match, err := h.adb.GetTeamNextGame(ctxMainDB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "team next game not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetTeamNextGame: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, match)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetTeamNextGame: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetTeamNextGame: %v\n", err)
		fmt.Printf("%+v\n", match)
	}
}

func (h *HandlerRepo) GetTeamPlayers(w http.ResponseWriter, r *http.Request) {
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

	seasonParam := r.URL.Query().Get("season")
	if seasonParam != "" {
		seasonParam = validateAndFormatSeason(seasonParam)
		if seasonParam == "" {
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid season parameter in query, use YYYY/YYYY or YYYY-YYYY"})
			return
		}
	} else {
		seasonParam = getCurrentSeason()
	}

	cacheKey := fmt.Sprintf("team:%d:players:season:%s", id, seasonParam)

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
	teamPlayers, err := h.adb.GetTeamPlayersBySeason(ctxMainDB, uint32(id), seasonParam)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "team players not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetTeamPlayersBySeason: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, teamPlayers)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetTeamPlayersBySeason: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetTeamPlayersBySeason: %v\n", err)
		fmt.Printf("%+v\n", teamPlayers)
	}
}

func (h *HandlerRepo) GetTeamLastGames(w http.ResponseWriter, r *http.Request) {
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

	seasonParam := r.URL.Query().Get("season")
	if seasonParam != "" {
		seasonParam = validateAndFormatSeason(seasonParam)
		if seasonParam == "" {
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid season parameter in query, use YYYY/YYYY or YYYY-YYYY"})
			return
		}
	} else {
		seasonParam = getCurrentSeason()
	}

	cacheKey := fmt.Sprintf("team:%d:fixtures:limit:%d:offset:%d:season:%s", id, limit, offset, seasonParam)

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
	matches, err := h.adb.GetTeamGames(ctxMainDB, uint32(id), uint32(limit), uint32(offset), seasonParam)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "team last games not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetTeamLastGames: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, matches)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetTeamLastGames: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetTeamLastGames: %v\n", err)
		fmt.Printf("%+v\n", matches)
	}
}

func (h *HandlerRepo) GetTeamManager(w http.ResponseWriter, r *http.Request) {
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

	cacheKey := fmt.Sprintf("team:%d:manager", id)

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
	manager, err := h.adb.GetCurrentManagerByTeam(ctxMainDB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || manager.ID == 1 {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "team manager not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetCurrentManagerByTeam: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, manager)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetCurrentManagerByTeam: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetCurrentManagerByTeam: %v\n", err)
		fmt.Printf("%+v\n", manager)
	}
}

func (h *HandlerRepo) GetTeamPlayersWithStats(w http.ResponseWriter, r *http.Request) {
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

	seasonParam := r.URL.Query().Get("season")
	if seasonParam != "" {
		seasonParam = validateAndFormatSeason(seasonParam)
		if seasonParam == "" {
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid season parameter in query, use YYYY/YYYY or YYYY-YYYY"})
			return
		}
	} else {
		seasonParam = getCurrentSeason()
	}

	cacheKey := fmt.Sprintf("team:%d:players_with_stats:season:%s", id, seasonParam)

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
	teamPlayers, err := h.adb.GetTeamPlayersWithStatsBySeason(ctxMainDB, uint32(id), seasonParam)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "team players not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetTeamPlayersWithStatsBySeason: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, teamPlayers)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetTeamPlayersWithStatsBySeason: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetTeamPlayersWithStatsBySeason: %v\n", err)
		fmt.Printf("%+v\n", teamPlayers)
	}
}

func (h *HandlerRepo) GetFavouriteTeams(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	cacheKey := fmt.Sprintf("user:%d:teams", userID)

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
	teamIDs, err := h.tdb.GetFavoriteTeamsIDs(ctxTransactDB, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "favorite teams not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetFavoriteTeamsIDs: %v\n", err)
		return
	}

	teams := make([]models.Team, 0, len(teamIDs))

	for _, id := range teamIDs {
		ctxAnalyticDB, cancelAnalyticDB := context.WithTimeout(r.Context(), dbTimeout)
		defer cancelAnalyticDB()
		team, err := h.adb.GetTeamByID(ctxAnalyticDB, id)
		if err != nil {
			h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
			log.Printf("error in GetTeamByID: %v\n", err)
			return
		}

		teams = append(teams, team)
	}

	resp, err := h.writeJSON(w, http.StatusOK, teams)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetFavoriteTeams: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetFavoriteTeams: %v\n", err)
		fmt.Printf("%+v\n", teams)
	}
}

func (h *HandlerRepo) SetFavouriteTeam(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty id parameter in query"})
		return
	}

	teamID, err := strconv.Atoi(idStr)
	if err != nil || teamID < 1 {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id parameter in query"})
		return
	}

	ctxDB, cancelDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelDB()
	_, err = h.adb.GetTeamByID(ctxDB, uint32(teamID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "team not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GEtTeamByID: %v\n", err)
		return
	}

	ctxTransactDB, cancelTransactDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelTransactDB()
	err = h.tdb.SetFavoriteTeamByID(ctxTransactDB, userID, int64(teamID))
	if err != nil {
		if err.Error() == fmt.Sprintf("team with ID=%d already favorite for user with ID=%d", teamID, userID) {
			h.writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in SetFavoriteTeamByID: %v\n", err)
		return
	}

	cacheKey := fmt.Sprintf("user:%d:teams", userID)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
		defer cancel()
		if err := h.cacheDB.Del(ctx, cacheKey); err != nil {
			log.Printf("error deleting cache for user %d: %v\n", userID, err)
		}
	}()

	_, err = h.writeJSON(w, http.StatusCreated, nil)
	if err != nil {
		log.Printf("error in writeJSON from SetFavoriteTeam: %v\n", err)
	}
}

func (h *HandlerRepo) DeleteFavouriteTeam(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty id parameter in query"})
		return
	}

	teamID, err := strconv.Atoi(idStr)
	if err != nil || teamID < 1 {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id parameter in query"})
		return
	}

	ctxDB, cancelDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelDB()
	_, err = h.adb.GetTeamByID(ctxDB, uint32(teamID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "team not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GEtTeamByID: %v\n", err)
		return
	}

	ctxTransactDB, cancelTransactDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelTransactDB()
	err = h.tdb.DeleteFavoriteTeamByID(ctxTransactDB, userID, int64(teamID))
	if err != nil {
		if err.Error() == fmt.Sprintf("team with ID=%d already not favorite for user with ID=%d", teamID, userID) {
			h.writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in DeleteFavoriteTeamByID: %v\n", err)
		return
	}

	cacheKey := fmt.Sprintf("user:%d:teams", userID)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
		defer cancel()
		if err := h.cacheDB.Del(ctx, cacheKey); err != nil {
			log.Printf("error deleting cache for user %d: %v\n", userID, err)
		}
	}()

	_, err = h.writeJSON(w, http.StatusOK, nil)
	if err != nil {
		log.Printf("error in writeJSON from DeleteFavoriteTeam: %v\n", err)
	}
}
