package rabbitmq

import (
	"context"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type Publisher interface {
	Publish(ctx context.Context, body []byte) error
}

type RabbitPublisher struct {
	url     string
	Queue   *amqp.Queue
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewPublisherCtx(ctx context.Context, url, queueName string) (*RabbitPublisher, error) {
	conn, err := dialWithRetry(ctx, amqp.Dial, url, 5, 1*time.Second)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitPublisher{
		url:     url,
		Conn:    conn,
		Channel: ch,
		Queue:   &q,
	}, nil
}

func (p *RabbitPublisher) Publish(ctx context.Context, body []byte) error {
	if p.Channel == nil || p.Queue == nil {
		return errors.New("Publisher not initialized")
	}

	return p.Channel.PublishWithContext(ctx, "", p.Queue.Name, false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		})
}
