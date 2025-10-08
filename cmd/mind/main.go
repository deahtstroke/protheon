package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/deahtstroke/protheon/internal/api"
	"github.com/rabbitmq/amqp091-go"
)

type Job struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
}

func sendHeartbeat(managerURL, workerID string, jobsDone int, lastJob time.Time, start time.Time) {
	hb := api.HeartbeatRequest{
		ID:          workerID,
		JobsDone:    jobsDone,
		LastJobTime: time.Now(),
		Uptime:      time.Since(start).String(),
	}

	body, err := json.Marshal(hb)
	resp, err := http.Post(managerURL+"/heartbeat", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("⚠️ heartbeat failed: %v", err)
		return
	}

	resp.Body.Close()
}

func startHeartbeat(ctx context.Context, managerURL, workerID string, jobsDone *int, lastJob *time.Time, start *time.Time) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ctx.Done():
			log.Println("Hearbeat shutting down...")
			return
		case <-ticker.C:
			sendHeartbeat(managerURL, workerID, *jobsDone, *lastJob, *start)
		}
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

func register(serverAddr string) (*api.RegisterResponse, error) {
	log.Printf("Registering with Conductor mind@[%s]...", serverAddr)
	registerURL := fmt.Sprintf("http://%s:8080/mind/register", serverAddr)
	hostname, _ := os.Hostname()
	os := runtime.GOOS
	registerRequest := api.RegisterRequest{
		Hostname: hostname,
		OS:       os,
	}

	data, err := json.Marshal(registerRequest)
	if err != nil {
		log.Fatalf("Failed to marshal register request: %v", err)
		return nil, err
	}
	resp, err := http.Post(registerURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("Failed to register with conductor: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Registration failed: status %d", resp.StatusCode)
		return nil, err
	}

	var regResp api.RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&regResp); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
		return nil, err
	}

	return &regResp, nil
}

func DoJobs(ctx context.Context, serverAddr string) {
	rabbitURL := fmt.Sprintf("amqp://protheon:secretpassword@%s:5672/", serverAddr)
	conn, err := dialWithRetry(ctx, rabbitURL, 5, 1*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Printf("Connected to Host RabbitMQ!")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Channel open failed: %v", err)
	}
	q, err := ch.QueueDeclare(
		"pgcr_jobs",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		log.Fatalf("Queue declared failed: %v", err)
	}

	err = ch.Qos(5, 0, false)
	if err != nil {
		log.Fatalf("Error setting QoS for Rabbit channel: %v", err)
	}
	msgs, err := ch.ConsumeWithContext(
		ctx,
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Fatalf("Consume failed: %v", err)
	}

	log.Println("Worker started. Press Ctrl+C to stop gracefully.")
	for {
		select {
		case <-ctx.Done():
			log.Print("Worker shutting down...")
			return
		case d, ok := <-msgs:
			if !ok {
				log.Printf("Msgs channel is closed, retrying...")
				time.Sleep(2 * time.Second)
				continue
			}
			var job Job
			if err := json.Unmarshal(d.Body, &job); err != nil {
				log.Printf("Bad job json: %v", err)
				continue
			}
			d.Ack(false)
			log.Printf("Job %+v done", job)
		}
	}
}

func main() {
	serverAddr := flag.String("server-addr", "", "Server URL")
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	start := time.Now()

	if *serverAddr == "" {
		fmt.Println("Unknown server url, use --server-addr")
		os.Exit(1)
	}

	resp, err := register(*serverAddr)
	if err != nil {
		fmt.Printf("Error register worker: %v", err)
		os.Exit(1)
	}

	go startHeartbeat(ctx, fmt.Sprintf("%s:8080", *serverAddr), resp.ID, nil, nil, &start)
	go DoJobs(ctx, *serverAddr)

	<-ctx.Done()
	log.Printf("Shutting down worker gracefully")
}
