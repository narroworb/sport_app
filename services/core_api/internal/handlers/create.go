package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/narroworb/core_api/internal/database"
	"github.com/narroworb/core_api/internal/middleware"
	"github.com/narroworb/core_api/internal/models"
)

func genID() uint32 {
	return uint32(time.Now().UnixNano() & 0xffffffff)
}

func (h *HandlerRepo) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey)
	if uid == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var p models.Player
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&p); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	// persist to transaction DB
	body, err := json.Marshal(p)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "serialization failed"})
		return
	}

	userID := uid.(int64)
	id, err := h.adb.CreatePlayer(r.Context(), userID, body)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "db insert failed"})
		return
	}

	p.ID = uint32(id)
	w.Header().Set("Location", "/api/player/"+fmt.Sprint(p.ID))
	h.writeJSON(w, http.StatusCreated, p)
}

func (h *HandlerRepo) CreateTeam(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey)
	if uid == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var t models.Team
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&t); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	body, err := json.Marshal(t)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "serialization failed"})
		return
	}
	userID := uid.(int64)
	id, err := h.adb.CreateTeam(r.Context(), userID, body)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "db insert failed"})
		return
	}
	t.ID = uint32(id)
	h.writeJSON(w, http.StatusCreated, t)
}

func (h *HandlerRepo) CreateTournament(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey)
	if uid == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var t models.Tournament
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&t); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	body, err := json.Marshal(t)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "serialization failed"})
		return
	}
	userID := uid.(int64)
	id, err := h.adb.CreateTournament(r.Context(), userID, body)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "db insert failed"})
		return
	}
	t.ID = uint32(id)
	h.writeJSON(w, http.StatusCreated, t)
}

func (h *HandlerRepo) CreateManager(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey)
	if uid == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var m models.Manager
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&m); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	body, err := json.Marshal(m)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "serialization failed"})
		return
	}
	userID := uid.(int64)
	id, err := h.adb.CreateManager(r.Context(), userID, body)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "db insert failed"})
		return
	}
	m.ID = uint32(id)
	h.writeJSON(w, http.StatusCreated, m)
}

func (h *HandlerRepo) CreateFixture(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey)
	if uid == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var m models.Match
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&m); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	body, err := json.Marshal(m)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "serialization failed"})
		return
	}
	userID := uid.(int64)
	id, err := h.adb.CreateFixture(r.Context(), userID, body)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "db insert failed"})
		return
	}
	m.ID = uint32(id)
	h.writeJSON(w, http.StatusCreated, m)
}

type MatchStatsPayload struct {
	Players []models.PlayerStatsInMatch `json:"players,omitempty"`
	Goalies []models.GoalieStatsInMatch `json:"goalies,omitempty"`
	Teams   models.TeamMatchStats       `json:"teams,omitempty"`
}

func (h *HandlerRepo) CreateMatchStats(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey)
	if uid == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var p MatchStatsPayload
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&p); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	body, err := json.Marshal(p)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "serialization failed"})
		return
	}
	userID := uid.(int64)
	id, err := h.adb.CreateMatchStats(r.Context(), userID, body)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "db insert failed"})
		return
	}
	// return created id
	h.writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

type TablePayload struct {
	TournamentID uint32            `json:"tournament_id"`
	Rows         []models.TableRow `json:"rows"`
}

func (h *HandlerRepo) CreateTournamentTable(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey)
	if uid == nil {
		h.writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var p TablePayload
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&p); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	body, err := json.Marshal(p)
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "serialization failed"})
		return
	}
	userID := uid.(int64)
	id, err := h.adb.CreateTournamentTable(r.Context(), userID, body)
	if err != nil {
		if errors.Is(err, database.ErrProtectedTournament) {
			h.writeJSON(w, http.StatusForbidden, map[string]string{"error": "tournament is protected"})
			return
		}
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "db insert failed"})
		return
	}
	h.writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}
