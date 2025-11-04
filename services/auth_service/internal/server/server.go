package server

import (
	"net/http"
	"os"

	"github.com/narroworb/auth_service/internal/db"
	"github.com/narroworb/auth_service/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	database := db.Connect()
	handler := handlers.Handler{
		DB: database,
	}

	r := chi.NewRouter()

	r.Post("/register", handler.Register)
	r.Post("/login", handler.Login)
	r.Get("/me", handler.Me)

	http.ListenAndServe(":"+port, r)
}
