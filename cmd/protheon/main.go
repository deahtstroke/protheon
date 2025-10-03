package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Job struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
}

func dialWithRetry(url string, attempts int, backoff time.Duration) (*amqp091.Connection, error) {
	var conn *amqp091.Connection
	var err error
	for i := range attempts {
		conn, err = amqp091.Dial(url)
		if err == nil {
			return conn, nil
		}
		log.Printf("amqp dial failed (attempt %d/%d): %v", i+1, attempts, err)
		time.Sleep(backoff)
		backoff *= 2
	}
	return nil, err
}

func runServer(ctx context.Context) {
	conn, err := dialWithRetry("amqp://protheon:secretpassword@localhost:5672/", 5, 1*time.Second)
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

func runWorker(ctx context.Context, rabbitUrl string) {
	conn, err := dialWithRetry(rabbitUrl, 5, 1*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

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
	role := flag.String("role", "worker", "server or worker")
	defaultURL := os.Getenv("RABBIT_URL")
	rabbitFlag := flag.String("rabbit", defaultURL, "amqp url (e.g., amqp://user:pass@manager:5672/")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	switch *role {
	case "server":
		runServer(ctx)
	case "worker":
		runWorker(ctx, *rabbitFlag)
	default:
		fmt.Println("Unknown role, use --role=server or --role=worker")
		os.Exit(1)
	}
}
