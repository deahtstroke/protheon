package rabbitmq

import (
	"context"
	"errors"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type DialFunc func(url string) (*amqp.Connection, error)

func dialWithRetry(ctx context.Context, dial DialFunc, url string, attempts int, backoff time.Duration) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error
	for i := range attempts {
		select {
		case <-ctx.Done():
			log.Printf("Context cancelled, dialing cancelled...")
			return nil, errors.New("Dialing cancelled")
		default:
			conn, err = dial(url)
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
