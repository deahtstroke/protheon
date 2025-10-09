package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/deahtstroke/protheon/internal/api"
	"github.com/deahtstroke/protheon/internal/producer"
	"github.com/deahtstroke/protheon/internal/rabbitmq"
	"github.com/gorilla/mux"
)

func main() {
	url := flag.String("url", "", "AMQPS url to running rabbitmq instance")
	fileRoot := flag.String("files", "", "URI of where to scan for files")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	rabbitPublisher, err := rabbitmq.NewPublisherCtx(ctx, *url, "pgcr")
	if err != nil {
		log.Fatalf("Error creating rabbitmq publisher: %v", err)
	}

	pgcrProducer := producer.NewPgcrProducer(*fileRoot)

	r := mux.NewRouter()
	r.HandleFunc("/mind/register", api.RegisterWorker).Methods("POST")
	r.HandleFunc("/mind/heartbet", api.ReceiveHeartbeat).Methods("POST")

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Printf("HTTP server listening on :8080")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Print("Server exited gracefully")
}
