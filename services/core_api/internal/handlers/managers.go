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
	"github.com/narroworb/core_api/internal/middleware"
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
	managerDetails, err := h.adb.GetManagerByID(ctxDB, uint32(id))
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
		getStatsFromDB = h.adb.GetManagerStatsBySeason
		cacheKey = fmt.Sprintf("manager:%d:stats:by_season", id)
	} else if byTeam == "true" {
		getStatsFromDB = h.adb.GetManagerStatsByTeam
		cacheKey = fmt.Sprintf("manager:%d:stats:by_team", id)
	} else {
		getStatsFromDB = h.adb.GetManagerFullStats
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
		yellowCards, redCards, err := h.adb.GetManagerCardsByID(ctxMainDB, uint32(id))
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
	managerTeams, err := h.adb.GetManagerTeams(ctxMainDB, uint32(id))
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

	cacheKey := fmt.Sprintf("manager:%d:fixtures:limit:%d:offset:%d:season:%s", id, limit, offset, seasonParam)

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
	managerFixtures, err := h.adb.GetManagerFixtures(ctxMainDB, uint32(id), uint32(limit), uint32(offset), seasonParam)
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

func (h *HandlerRepo) GetFavouriteManagers(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	cacheKey := fmt.Sprintf("user:%d:managers", userID)

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
	managerIDs, err := h.tdb.GetFavoriteManagersIDs(ctxTransactDB, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "favorite managers not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetFavoriteManagersIDs: %v\n", err)
		return
	}

	managers := make([]models.Manager, 0, len(managerIDs))

	for _, id := range managerIDs {
		ctxAnalyticDB, cancelAnalyticDB := context.WithTimeout(r.Context(), dbTimeout)
		defer cancelAnalyticDB()
		manager, err := h.adb.GetManagerByID(ctxAnalyticDB, id)
		if err != nil {
			h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
			log.Printf("error in GetManagerByID: %v\n", err)
			return
		}

		managers = append(managers, manager)
	}

	resp, err := h.writeJSON(w, http.StatusOK, managers)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetFavoriteManagers: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetFavoriteManagers: %v\n", err)
		fmt.Printf("%+v\n", managers)
	}
}

func (h *HandlerRepo) SetFavouriteManager(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty id parameter in query"})
		return
	}

	managerID, err := strconv.Atoi(idStr)
	if err != nil || managerID < 1 {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id parameter in query"})
		return
	}

	ctxDB, cancelDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelDB()
	_, err = h.adb.GetManagerByID(ctxDB, uint32(managerID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "manager not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetManagerByID: %v\n", err)
		return
	}

	ctxTransactDB, cancelTransactDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelTransactDB()
	err = h.tdb.SetFavoriteManagerByID(ctxTransactDB, userID, int64(managerID))
	if err != nil {
		if err.Error() == fmt.Sprintf("manager with ID=%d already favorite for user with ID=%d", managerID, userID) {
			h.writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in SetFavoriteManagerByID: %v\n", err)
		return
	}

	cacheKey := fmt.Sprintf("user:%d:managers", userID)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
		defer cancel()
		if err := h.cacheDB.Del(ctx, cacheKey); err != nil {
			log.Printf("error deleting cache for user %d: %v\n", userID, err)
		}
	}()

	_, err = h.writeJSON(w, http.StatusCreated, nil)
	if err != nil {
		log.Printf("error in writeJSON from SetFavoriteManager: %v\n", err)
	}
}

func (h *HandlerRepo) DeleteFavouriteManager(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty id parameter in query"})
		return
	}

	managerID, err := strconv.Atoi(idStr)
	if err != nil || managerID < 1 {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id parameter in query"})
		return
	}

	ctxDB, cancelDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelDB()
	_, err = h.adb.GetManagerByID(ctxDB, uint32(managerID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "manager not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetManagerByID: %v\n", err)
		return
	}

	ctxTransactDB, cancelTransactDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelTransactDB()
	err = h.tdb.DeleteFavoriteManagerByID(ctxTransactDB, userID, int64(managerID))
	if err != nil {
		if err.Error() == fmt.Sprintf("manager with ID=%d already not favorite for user with ID=%d", managerID, userID) {
			h.writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in DeleteFavoriteManagerByID: %v\n", err)
		return
	}

	cacheKey := fmt.Sprintf("user:%d:managers", userID)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
		defer cancel()
		if err := h.cacheDB.Del(ctx, cacheKey); err != nil {
			log.Printf("error deleting cache for user %d: %v\n", userID, err)
		}
	}()

	_, err = h.writeJSON(w, http.StatusOK, nil)
	if err != nil {
		log.Printf("error in writeJSON from DeleteFavoriteManager: %v\n", err)
	}
}
