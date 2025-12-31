package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/narroworb/data_collector/internal/collector"
	"github.com/narroworb/data_collector/internal/database"
	"github.com/narroworb/data_collector/internal/kafka/consumer"
	"github.com/narroworb/data_collector/internal/kafka/producer"
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

	producer := producer.NewKafkaProducer()
	log.Println("Connected to kafka producer")
	defer producer.Close()

	consumer := consumer.NewKafkaConsumer("1", db)
	log.Println("Connected to kafka consumer")
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Start to complete previous messages in kafka")
	if err := consumer.DrainAndProcess(ctx, time.Second*2); err != nil {
		log.Fatalf("error in drainandprocess: %v", err)
	}
	log.Println("Completed previous messages in kafka")

	go func() {
		log.Println("Kafka consumer is listening")
		if err := consumer.Start(ctx); err != nil {
			log.Fatalf("error in work consumer: %v", err)
		}
	}()

	updater := collector.NewUpdater(db, producer)
	log.Println("Connected to updater")

	schedulerStats := scheduler.NewScheduler(updater.StartUpdate)

	log.Println("scheduler by stats created")

	go schedulerStats.Start(time.Hour * 1)

	schedulerWithoutStats := scheduler.NewScheduler(updater.StartUpdateWithoutStatistics)

	log.Println("scheduler without stats created")

	go schedulerWithoutStats.Start(time.Hour * 2)

	r := chi.NewRouter()

	r.Post("/run-update", func(w http.ResponseWriter, r *http.Request) {
		schedulerStats.RunNow()
		w.WriteHeader(http.StatusAccepted)
	})

	r.Post("/run-update-without-stats", func(w http.ResponseWriter, r *http.Request) {
		schedulerWithoutStats.RunNow()
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
