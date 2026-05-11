package main

import (
	"log"
	"os"

	"github.com/narroworb/core_api/internal/database"
	"github.com/narroworb/core_api/internal/elasticsearch"
	"github.com/narroworb/core_api/internal/server"
)

func main() {
	clickhouseDB, err := database.NewClickhouseDB()
	if err != nil {
		log.Fatalf("error in connect clickhouse: %v", err)
	}
	log.Println("Connected to Clickhouse")
	defer clickhouseDB.Close()

	redis, err := database.NewRedis()
	if err != nil {
		log.Fatalf("error in connect redis: %v", err)
	}
	log.Println("Connected to Redis")

	postgresDB, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("error in connect postgres: %v", err)
	}
	log.Println("Connected to Postgres")

	esClient, err := elasticsearch.NewElasticsearch(clickhouseDB)
	if err != nil {
		log.Fatalf("error in connect elasticsearch: %v", err)
	}
	log.Println("Connected to Elasticsearch")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	srv := server.NewServerRepo(clickhouseDB, postgresDB, redis, esClient, port)

	if err := srv.Run(); err != nil {
		log.Fatalf("error in runnig server: %v", err)
	}
}
