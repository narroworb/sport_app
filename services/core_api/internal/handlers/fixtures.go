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
)

func (h *HandlerRepo) GetFixtureDetails(w http.ResponseWriter, r *http.Request) {
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

	cacheKey := "fixture:" + fmt.Sprint(id) + ":details"
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
	fixtureDetails, err := h.db.GetMatchByID(ctxDB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "fixture not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetFixtureDetails: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, fixtureDetails)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetFixtureDetails: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetFixtureDetails: %v\n", err)
	}
}

func (h *HandlerRepo) GetFixturesByDate(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	date, err := parseDate(dateStr)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date parameter in query"})
		return
	}

	cacheKey := fmt.Sprintf("fixtures:%s", date)

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
	matches, err := h.db.GetMatchesByDate(ctxMainDB, date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "fixtures not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetMatchesByDate: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, matches)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetMatchesByDate: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetMatchesByDate: %v\n", err)
		fmt.Printf("%+v\n", matches)
	}
}

func (h *HandlerRepo) GetFixturePlayersStats(w http.ResponseWriter, r *http.Request) {
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

	cacheKey := fmt.Sprintf("fixtures:%d:stats:players", id)

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
	stats, err := h.db.GetMatchPlayersStats(ctxMainDB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "fixture stats not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetMatchPlayersStats: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, stats)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetMatchPlayersStats: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetMatchPlayersStats: %v\n", err)
		fmt.Printf("%+v\n", stats)
	}
}

func (h *HandlerRepo) GetFixtureGoaliesStats(w http.ResponseWriter, r *http.Request) {
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

	cacheKey := fmt.Sprintf("fixtures:%d:stats:goalies", id)

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
	stats, err := h.db.GetMatchGoaliesStats(ctxMainDB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "fixture stats not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetMatchGoalieStats: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, stats)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetMatchGoalieStats: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetMatchGoalieStats: %v\n", err)
		fmt.Printf("%+v\n", stats)
	}
}

func (h *HandlerRepo) GetFixtureTeamsStats(w http.ResponseWriter, r *http.Request) {
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

	cacheKey := fmt.Sprintf("fixtures:%d:stats:teams", id)

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
	stats, err := h.db.GetMatchTeamsStats(ctxMainDB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "fixture stats not found"})
			return
		}

		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "try later"})
		log.Printf("error in GetMatchTeamsStats: %v\n", err)
		return
	}

	resp, err := h.writeJSON(w, http.StatusOK, stats)
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), cacheTimeout)
			defer cancel()

			if err := h.cacheDB.Set(ctx, cacheKey, resp, cacheTTL); err != nil {
				log.Printf("error in set to cache from GetMatchTeamsStats: %v\n", err)
			}
		}()
	} else {
		log.Printf("error in writeJSON from GetMatchTeamsStats: %v\n", err)
		fmt.Printf("%+v\n", stats)
	}
}
