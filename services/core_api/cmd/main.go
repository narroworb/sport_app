package main

import (
	"log"
	"os"

	"github.com/narroworb/core_api/internal/database"
	"github.com/narroworb/core_api/internal/server"
)

func main() {
	db, err := database.NewClickhouseDB()
	if err != nil {
		log.Fatalf("error in connect clickhouse: %v", err)
	}
	log.Println("Connected to Clickhouse")
	defer db.Close()

	redis, err := database.NewRedis()
	if err != nil {
		log.Fatalf("error in connect redis: %v", err)
	}
	log.Println("Connected to Redis")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	srv := server.NewServerRepo(db, redis, port)

	if err := srv.Run(); err != nil {
		log.Fatalf("error in runnig server: %v", err)
	}
}
