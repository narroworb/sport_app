package main

import (
	"log"

	"github.com/narroworb/data_collector/internal/collector"
	"github.com/narroworb/data_collector/internal/database"
)

func main() {
	db, err := database.NewClickhouseDB()
	if err != nil {
		log.Fatalf("error in create db: %v", err)
	}
	defer db.Close()
	updater := collector.NewUpdater(db)

	updater.StartUpdate()
}
