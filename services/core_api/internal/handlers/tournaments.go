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

func (h *HandlerRepo) GetTournamentDetails(w http.ResponseWriter, r *http.Request) {
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

	cacheKey := "tournament:" + fmt.Sprint(id) + ":details"
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
	tournamentDetails, err := h.adb.GetTournamentByID(ctxDB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "tournament not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetTournamentDetails: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, tournamentDetails)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetTournamentDetails: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetTournamentDetails: %v\n", err)
	}
}

func (h *HandlerRepo) GetTournamentTable(w http.ResponseWriter, r *http.Request) {
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

	cacheKey := fmt.Sprintf("tournament:%d:table", id)
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
	tournamentTable, err := h.adb.GetTournamentTableByID(ctxDB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "tournament table not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetTournamentTableByIDAndSeason: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, tournamentTable)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetTournamentTableByIDAndSeason: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetTournamentTableByIDAndSeason: %v\n", err)
	}
}

func (h *HandlerRepo) GetTournamentFixtures(w http.ResponseWriter, r *http.Request) {
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

	cacheKey := fmt.Sprintf("tournament:%d:fixtures", id)
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
	tournamentFixtures, err := h.adb.GetTournamentFixturesByID(ctxDB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "tournament fixtures not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetTournamentFixturesByID: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, tournamentFixtures)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetTournamentFixturesByID: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetTournamentFixturesByID: %v\n", err)
	}
}

func (h *HandlerRepo) GetTournamentTeamsStats(w http.ResponseWriter, r *http.Request) {
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

	cacheKey := fmt.Sprintf("tournament:%d:stats:teams", id)
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
	tournamentStats, err := h.adb.GetTournamentTeamsStatsByID(ctxDB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "tournament team stats not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetTournamentTeamsStatsByID: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, tournamentStats)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetTournamentTeamsStatsByID: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetTournamentTeamsStatsByID: %v\n", err)
	}
}

func (h *HandlerRepo) GetTournamentPlayersStats(w http.ResponseWriter, r *http.Request) {
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

	if limit == 0 {
		limit = 10
	}

	cacheKey := fmt.Sprintf("tournament:%d:stats:players:limit:%d:offset:%d", id, limit, offset)
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
	tournamentStats, err := h.adb.GetTournamentPlayersStatsByID(ctxDB, uint32(id), uint32(limit), uint32(offset))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "tournament player stats not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetTournamentPlayersStatsByID: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, tournamentStats)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetTournamentPlayersStatsByID: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetTournamentPlayersStatsByID: %v\n", err)
	}
}

func (h *HandlerRepo) GetTournamentTableGraph(w http.ResponseWriter, r *http.Request) {
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

	cacheKey := fmt.Sprintf("tournament:%d:table:graph", id)
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
	tournamentTableGraph, err := h.adb.GetTournamentTableGraphByID(ctxDB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "tournament table graph not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetTournamentTableGraphByID: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, tournamentTableGraph)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetTournamentTableGraphByID: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetTournamentTableGraphByID: %v\n", err)
	}
}

func (h *HandlerRepo) GetFavouriteTournaments(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	cacheKey := fmt.Sprintf("user:%d:tournaments", userID)

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
	tournamentIDs, err := h.tdb.GetFavoriteTournamentIDs(ctxTransactDB, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "favorite tournaments not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetFavoriteTournamentIDs: %v\n", err)
		return
	}

	tournaments := make([]models.Tournament, 0, len(tournamentIDs))

	for _, id := range tournamentIDs {
		ctxAnalyticDB, cancelAnalyticDB := context.WithTimeout(r.Context(), dbTimeout)
		defer cancelAnalyticDB()
		tournament, err := h.adb.GetTournamentByID(ctxAnalyticDB, id)
		if err != nil {
			h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
			log.Printf("error in GetTournamentByID: %v\n", err)
			return
		}

		tournaments = append(tournaments, tournament)
	}

	resp, err := h.writeJSON(w, http.StatusOK, tournaments)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetFavoriteTournaments: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetFavoriteTournaments: %v\n", err)
		fmt.Printf("%+v\n", tournaments)
	}
}

func (h *HandlerRepo) SetFavouriteTournament(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty id parameter in query"})
		return
	}

	tournamentID, err := strconv.Atoi(idStr)
	if err != nil || tournamentID < 1 {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id parameter in query"})
		return
	}

	ctxDB, cancelDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelDB()
	_, err = h.adb.GetTournamentByID(ctxDB, uint32(tournamentID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "tournament not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetTournamentByID: %v\n", err)
		return
	}

	ctxTransactDB, cancelTransactDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelTransactDB()
	err = h.tdb.SetFavoriteTournamentByID(ctxTransactDB, userID, int64(tournamentID))
	if err != nil {
		if err.Error() == fmt.Sprintf("tournament with ID=%d already favorite for user with ID=%d", tournamentID, userID) {
			h.writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in SetFavoriteTournamentByID: %v\n", err)
		return
	}

	cacheKey := fmt.Sprintf("user:%d:tournaments", userID)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
		defer cancel()
		if err := h.cacheDB.Del(ctx, cacheKey); err != nil {
			log.Printf("error deleting cache for user %d: %v\n", userID, err)
		}
	}()

	_, err = h.writeJSON(w, http.StatusCreated, nil)
	if err != nil {
		log.Printf("error in writeJSON from SetFavoriteTournament: %v\n", err)
	}
}

func (h *HandlerRepo) DeleteFavouriteTournament(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)

	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty id parameter in query"})
		return
	}

	tournamentID, err := strconv.Atoi(idStr)
	if err != nil || tournamentID < 1 {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id parameter in query"})
		return
	}

	ctxDB, cancelDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelDB()
	_, err = h.adb.GetTournamentByID(ctxDB, uint32(tournamentID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "tournament not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetTournamentByID: %v\n", err)
		return
	}

	ctxTransactDB, cancelTransactDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelTransactDB()
	err = h.tdb.DeleteFavoriteTournamentByID(ctxTransactDB, userID, int64(tournamentID))
	if err != nil {
		if err.Error() == fmt.Sprintf("tournament with ID=%d already not favorite for user with ID=%d", tournamentID, userID) {
			h.writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in DeleteFavoriteTournamentByID: %v\n", err)
		return
	}

	cacheKey := fmt.Sprintf("user:%d:tournaments", userID)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
		defer cancel()
		if err := h.cacheDB.Del(ctx, cacheKey); err != nil {
			log.Printf("error deleting cache for user %d: %v\n", userID, err)
		}
	}()

	_, err = h.writeJSON(w, http.StatusCreated, nil)
	if err != nil {
		log.Printf("error in writeJSON from DeleteFavoriteTournament: %v\n", err)
	}
}

func (h *HandlerRepo) GetAllTournaments(w http.ResponseWriter, r *http.Request) {
	cacheKey := fmt.Sprintf("alltournaments")

	ctxCache, cancelCache := context.WithTimeout(r.Context(), cacheTimeout)
	defer cancelCache()
	cached, err := h.cacheDB.Get(ctxCache, cacheKey)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cached))
		return
	}

	ctxAnalyticDB, cancelAnalyticDB := context.WithTimeout(r.Context(), dbTimeout)
	defer cancelAnalyticDB()
	tournaments, err := h.adb.GetAllTournaments(ctxAnalyticDB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "favorite tournaments not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetAllTournaments: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, tournaments)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetAllTournaments: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetAllTournaments: %v\n", err)
		fmt.Printf("%+v\n", tournaments)
	}
}