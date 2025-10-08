package rabbitmq

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestDialWithRetrySuccess(t *testing.T) {
	attempts := 0
	fakeDial := func(url string) (*amqp091.Connection, error) {
		attempts++
		if attempts < 3 {
			return nil, errors.New("failure")
		}
		return &amqp091.Connection{}, nil
	}

	ctx := context.Background()
	conn, err := dialWithRetry(ctx, fakeDial, "fake-url", 5, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	if conn == nil {
		t.Fatal("Expected non-nill connection")
	}
}
