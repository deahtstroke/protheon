package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
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

func runServer() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	defer ch.Close()
	args := amqp091.Table{"x-max-length": int32(1000)}
	q, err := ch.QueueDeclare(
		"pgcr_jobs",
		true,
		false,
		false,
		false,
		args,
	)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Server shutting down...")
		case t := <-ticker.C:
			job := Job{
				ID:   rand.Int(),
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
				log.Printf("Failed to publish job: %+v", job)
			} else {
				log.Printf("Pulbished job: %+v\n", job)
			}
		}
	}
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

func runWorker(rabbitUrl string) {
	conn, err := dialWithRetry(rabbitUrl, 5, 1*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Channel open failed: %v", err)
	}
	args := amqp091.Table{"x-max-length": int32(1000)}
	q, err := ch.QueueDeclare("pgcr_jobs", true, false, false, false, args)
	if err != nil {
		log.Fatalf("Queue declared failed: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Consume failed: %v", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Println("Worker started. Press Ctrl+C to stop gracefully.")
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down worker..e")
			return
		case d, ok := <-msgs:
			if !ok {
				log.Println("Channel close, exiting...")
				return
			}
			var job Job
			if err := json.Unmarshal(d.Body, &job); err != nil {
				log.Printf("Bad job json: %v", err)
				continue
			}
			log.Printf("Worker got job: %+v", job)
			d.Ack(false)
			log.Printf("Job %d done", job.ID)
		}
	}
}

func main() {
	role := flag.String("role", "worker", "server or worker")
	defaultURL := os.Getenv("RABBIT_URL")
	rabbitFlag := flag.String("rabbit", defaultURL, "amqp url (e.g., amqp://user:pass@manager:5672/")
	flag.Parse()

	switch *role {
	case "server":
		runServer()
	case "worker":
		runWorker(*rabbitFlag)
	default:
		fmt.Println("Unknown role, use --role=server or --role=worker")
		os.Exit(1)
	}
}
