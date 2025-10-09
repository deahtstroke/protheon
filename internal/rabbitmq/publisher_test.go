package rabbitmq

import (
	"context"
	"fmt"
	"testing"

	rabbitcontainer "github.com/testcontainers/testcontainers-go/modules/rabbitmq"
)

func TestPublishSuccess(t *testing.T) {
	ctx := context.Background()
	rabbitmqContainer, err := rabbitcontainer.Run(
		ctx,
		"rabbitmq:3.7.25-management-alpine",
		rabbitcontainer.WithAdminUsername("protheon"),
		rabbitcontainer.WithAdminPassword("password"))

	if err != nil {
		t.Fatalf("Error running test container: %v", err)
	}

	host, err := rabbitmqContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Error getting hostname of container: %v", err)
	}

	port, err := rabbitmqContainer.MappedPort(ctx, "5672")
	if err != nil {
		t.Fatalf("Error getting mapped port from container: %v", err)
	}

	url := fmt.Sprintf("amqp://protheon:password@%s:%s/", host, port.Port())

	publisher, err := NewPublisherCtx(ctx, url, "pgcr_jobs")
	if err != nil {
		t.Fatalf("Error creating publisher: %v", err)
	}

	err = publisher.Publish(ctx, []byte(string("Hello World!")))
	if err != nil {
		t.Fatalf("Publishing failed: %v", err)
	}
}

func TestPublishContainerTerminated(t *testing.T) {
	ctx := context.Background()
	rabbitmqContainer, err := rabbitcontainer.Run(
		ctx,
		"rabbitmq:3.7.25-management-alpine",
		rabbitcontainer.WithAdminUsername("protheon"),
		rabbitcontainer.WithAdminPassword("password"))
	if err != nil {
		t.Fatalf("Error running test container: %v", err)
	}

	host, err := rabbitmqContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Error getting hostname of container: %v", err)
	}

	port, err := rabbitmqContainer.MappedPort(ctx, "5672")
	if err != nil {
		t.Fatalf("Error getting mapped port from container: %v", err)
	}

	url := fmt.Sprintf("amqp://protheon:password@%s:%s/", host, port.Port())

	publisher, err := NewPublisherCtx(ctx, url, "pgcr_jobs")
	if err != nil {
		t.Fatalf("Error creating publisher: %v", err)
	}

	err = rabbitmqContainer.Terminate(ctx)
	if err != nil {
		t.Fatalf("Error while terminating rabbitmq container: %v", err)
	}

	err = publisher.Publish(ctx, []byte(string("Hello World!")))
	if err == nil {
		t.Fatalf("Expecting error, found none")
	}
}

func TestPublishFailureChannelClosed(t *testing.T) {
	ctx := context.Background()
	rabbitmqContainer, err := rabbitcontainer.Run(
		ctx,
		"rabbitmq:3.7.25-management-alpine",
		rabbitcontainer.WithAdminUsername("protheon"),
		rabbitcontainer.WithAdminPassword("password"))
	if err != nil {
		t.Fatalf("Error running test container: %v", err)
	}

	port, err := rabbitmqContainer.MappedPort(ctx, "5672")
	if err != nil {
		t.Fatalf("Error getting mapped port from container: %v", err)
	}

	host, err := rabbitmqContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Error getting host from container: %v", err)
	}

	url := fmt.Sprintf("amqp://protheon:password@%s:%s/", host, port.Port())

	publisher, err := NewPublisherCtx(ctx, url, "pgcr_jobs")
	if err != nil {
		t.Fatalf("Error creating publisher: %v", err)
	}

	err = publisher.Channel.Close()
	if err != nil {
		t.Fatalf("Error closing rabbitmq channel")
	}

	err = publisher.Publish(ctx, []byte(string("Hello world!")))
	if err == nil {
		t.Fatalf("Expecting error, found none")
	}
}
