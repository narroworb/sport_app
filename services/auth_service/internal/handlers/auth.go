package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/narroworb/auth_service/internal/jwt"
	"github.com/narroworb/auth_service/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	DB *sqlx.DB
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid input data", http.StatusBadRequest)
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	_, err := h.DB.Exec(`INSERT INTO users (username, password_hash) VALUES($1, $2);`, req.Username, hash)
	if err != nil {
		http.Error(w, "username already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "registered",
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid input data", http.StatusBadRequest)
		return
	}

	var user models.User

	err := h.DB.Get(&user, `SELECT id, username, password_hash, created_at FROM users WHERE username=$1;`, req.Username)
	if err != nil {
		log.Printf("error in search user: %v", err)
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		http.Error(w, "password incorrect", http.StatusUnauthorized)
		return
	}

	token, _ := jwt.GenerateToken(user.Username, user.ID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "no token", http.StatusUnauthorized)
		return
	}

	username, userID, err := jwt.ParseToken(authHeader)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"username": username,
		"id":       userID,
	})
}
