package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/narroworb/data_collector/internal/collector"
	"github.com/narroworb/data_collector/internal/database"
	"github.com/narroworb/data_collector/internal/scheduler"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	db, err := database.NewClickhouseDB()
	if err != nil {
		log.Fatalf("error in create db: %v", err)
	}
	log.Println("Connected to db")
	defer db.Close()
	updater := collector.NewUpdater(db)
	log.Println("Connected to updater")

	scheduler := scheduler.NewScheduler(updater.StartUpdate)

	log.Println("scheduler created")

	go scheduler.Start(time.Hour * 1)

	r := chi.NewRouter()

	r.Post("/run-update", func(w http.ResponseWriter, r *http.Request) {
		scheduler.RunNow()
		w.WriteHeader(http.StatusAccepted)
	})

	go func() {
		if err := http.ListenAndServe(":8090", r); err != nil {
			signalChan <- syscall.SIGINT
		}
	}()
	log.Println("server listen on localhost:8090")

	<-signalChan

	log.Println("collector stopped gracefully")
}
