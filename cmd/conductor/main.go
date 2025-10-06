package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/deahtstroke/protheon/internal/api"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rabbitmq/amqp091-go"
)

type Job struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
}

type Worker struct {
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	IP       string `json:"ip"`
}

var (
	workers   = make(map[string]Worker)
	workersMu sync.Mutex
)

func registerWorker(w http.ResponseWriter, r *http.Request) {
	var req api.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	id := uuid.New().String()
	worker := Worker{
		ID:       id,
		Hostname: req.Hostname,
		OS:       req.OS,
		IP:       host,
	}

	workersMu.Lock()
	workers[id] = worker
	workersMu.Unlock()

	resp := api.RegisterResponse{
		ID:                id,
		HeartbeatInterval: 30,
		QueueName:         "jobs",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	fmt.Printf("[Conductor] Worker registered: %+v\n", worker)
}

func receiveHeartbeat(w http.ResponseWriter, r *http.Request) {
	var req api.HeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	workersMu.Lock()
	_, ok := workers[req.ID]
	workersMu.Unlock()

	if !ok {
		log.Printf("Unable to find worker with Id [%s]", req.ID)
	}
}

func dialWithRetry(ctx context.Context, url string, attempts int, backoff time.Duration) (*amqp091.Connection, error) {
	var conn *amqp091.Connection
	var err error
	for i := range attempts {
		select {
		case <-ctx.Done():
			log.Printf("Context cancelled, dialing cancelled...")
			return nil, errors.New("Dialing cancelled")
		default:
			conn, err = amqp091.Dial(url)
			if err == nil {
				return conn, nil
			}
			log.Printf("amqp dial failed (attempt %d/%d): %v", i+1, attempts, err)
			time.Sleep(backoff)
			backoff *= 2
		}
	}
	return nil, err
}

func EmitJobs(ctx context.Context) {
	conn, err := dialWithRetry(ctx, "amqp://protheon:secretpassword@localhost:5672/", 5, 1*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	defer conn.Close()
	log.Printf("Connected to RabbitMQ successfully!")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"pgcr_jobs",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error declaring queue: %v", err)
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	log.Printf("Starting job dispatch operation...")
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			log.Println("Server shutting down...")
			return
		case t := <-ticker.C:
			job := Job{
				ID:   i,
				Data: fmt.Sprintf("Job created at %s", t.Format(time.RFC3339)),
			}
			msg, _ := json.Marshal(job)
			err = ch.PublishWithContext(ctx,
				"",
				q.Name,
				false,
				false, amqp091.Publishing{
					DeliveryMode: amqp091.Persistent,
					ContentType:  "application/json",
					Body:         msg,
				})

			if err != nil {
				log.Printf("Failed to publish job: %+v: %v", job, err)
			} else {
				log.Printf("Published job: %+v\n", job)
			}
		}
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	r := mux.NewRouter()
	r.HandleFunc("/mind/register", registerWorker).Methods("POST")
	r.HandleFunc("/mind/heartbet", receiveHeartbeat).Methods("POST")

	go EmitJobs(ctx)

	log.Printf("HTTP server listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
