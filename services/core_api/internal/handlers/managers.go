package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/narroworb/core_api/internal/models"
)

func (h *HandlerRepo) GetManagerDetails(w http.ResponseWriter, r *http.Request) {
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

	cacheKey := "manager:" + fmt.Sprint(id) + ":details"
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
	managerDetails, err := h.db.GetManagerByID(ctxDB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "manager not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetManagerDetails: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, managerDetails)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetManagerDetails: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetManagerDetails: %v\n", err)
	}
}

func (h *HandlerRepo) GetManagerStats(w http.ResponseWriter, r *http.Request) {
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

	byTeam := r.URL.Query().Get("by_team")
	if byTeam != "true" && byTeam != "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid by_team parameter in query, use boolean value"})
		return
	}

	bySeason := r.URL.Query().Get("by_season")
	if bySeason != "true" && bySeason != "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid by_season parameter in query, use boolean value"})
		return
	}

	if bySeason == "true" && byTeam == "true" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid query, there can be only one grouping criterion"})
		return
	}

	var getStatsFromDB func(context.Context, uint32) (map[string]models.ManagerStatsInPeriod, error)
	var cacheKey string

	if bySeason == "true" {
		getStatsFromDB = h.db.GetManagerStatsBySeason
		cacheKey = fmt.Sprintf("manager:%d:stats:by_season", id)
	} else if byTeam == "true" {
		getStatsFromDB = h.db.GetManagerStatsByTeam
		cacheKey = fmt.Sprintf("manager:%d:stats:by_team", id)
	} else {
		getStatsFromDB = h.db.GetManagerFullStats
		cacheKey = fmt.Sprintf("manager:%d:stats:full", id)
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
	managerStats, err := getStatsFromDB(ctxMainDB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "manager stats not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetManagerStats: %v\n", err)
		return
	}

	if cacheKey == fmt.Sprintf("manager:%d:stats:full", id) {
		ctxMainDB, cancelMainDB := context.WithTimeout(r.Context(), dbTimeout)
		defer cancelMainDB()
		yellowCards, redCards, err := h.db.GetManagerCardsByID(ctxMainDB, uint32(id))
		if err != nil {
			h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
			log.Printf("error in GetManagerCardsByID: %v\n", err)
			return
		}
		stats := managerStats["full"]

		stats.YellowCards = yellowCards
		stats.RedCards = redCards

		managerStats["full"] = stats
	}

	resp, err := h.writeJSON(w, http.StatusOK, managerStats)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetManagerStats: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetManagerStats: %v\n", err)
		fmt.Printf("%+v\n", managerStats)
	}
}

func (h *HandlerRepo) GetManagerTeams(w http.ResponseWriter, r *http.Request) {
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

	cacheKey := fmt.Sprintf("manager:%d:teams", id)

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
	managerTeams, err := h.db.GetManagerTeams(ctxMainDB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "manager teams not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetManagerTeams: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, managerTeams)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetManagerTeams: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetManagerTeams: %v\n", err)
		fmt.Printf("%+v\n", managerTeams)
	}
}

func (h *HandlerRepo) GetManagerFixtures(w http.ResponseWriter, r *http.Request) {
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

	cacheKey := fmt.Sprintf("manager:%d:fixtures:limit:%d:offset:%d", id, limit, offset)

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
	managerFixtures, err := h.db.GetManagerFixtures(ctxMainDB, uint32(id), uint32(limit), uint32(offset))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "manager fixtures not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetManagerFixtures: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, managerFixtures)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetManagerFixtures: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetManagerFixtures: %v\n", err)
		fmt.Printf("%+v\n", managerFixtures)
	}
}
