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
	teamDetails, err := h.db.GetTeamByID(ctxDB, uint32(id))
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

	var getStatsFromDB func(context.Context, models.TeamStatsFilter) (models.TeamStatsInPeriod, error)
	var cacheKey string

	if filter.Season != "" {
		getStatsFromDB = h.db.GetTeamStatsBySeason
		cacheKey = "team:" + fmt.Sprint(id) + ":stats:season:" + filter.Season
	} else if (filter.ToDate != time.Time{} || filter.FromDate != time.Time{}) {
		getStatsFromDB = h.db.GetTeamStatsByDates
		cacheKey = "team:" + fmt.Sprint(id) + ":stats:dates:" + fmt.Sprintf("%s_%s", filter.FromDate, filter.ToDate)
	} else {
		getStatsFromDB = h.db.GetTeamFullStats
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
	teamStats, err := h.db.GetStandingsByTeamAndSeason(ctxMainDB, uint32(id), seasonParam)
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
	match, err := h.db.GetTeamNextGame(ctxMainDB, uint32(id))
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
	teamPlayers, err := h.db.GetTeamPlayersBySeason(ctxMainDB, uint32(id), seasonParam)
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

	cacheKey := fmt.Sprintf("team:%d:fixtures:limit:%d:offset:%d", id, limit, offset)

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
	matches, err := h.db.GetTeamLastGames(ctxMainDB, uint32(id), uint32(limit), uint32(offset))
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
	manager, err := h.db.GetCurrentManagerByTeam(ctxMainDB, uint32(id))
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
	teamPlayers, err := h.db.GetTeamPlayersWithStatsBySeason(ctxMainDB, uint32(id), seasonParam)
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
