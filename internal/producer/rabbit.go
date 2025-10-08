package producer

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deahtstroke/protheon/internal/rabbitmq"
	"github.com/klauspost/compress/zstd"
)

const (
	MAX_CAPACITY = 64 * 1048 * 1048
)

type Producer interface {
	Produce(ctx context.Context) error
}

type PgcrProducer struct {
	Source   string
	Rabbitmq rabbitmq.Publisher
}

func (pp *PgcrProducer) Produce(ctx context.Context) error {
	file, err := os.Open(pp.Source)
	if err != nil {
		return fmt.Errorf("Error opening source [%s]: %v", pp.Source, err)
	}

	defer file.Close()

	bufReader := bufio.NewReader(file)
	decoder, err := zstd.NewReader(bufReader)
	if err != nil {
		return fmt.Errorf("Error creating ZSTD reader for source [%s]: %v", pp.Source, err)
	}

	buf := make([]byte, MAX_CAPACITY)
	scanner := bufio.NewScanner(decoder)
	scanner.Buffer(buf, MAX_CAPACITY)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			log.Printf("[Producer] Context cancelled. Finishing producing PGCRs...")
			return nil
		default:
		}
	}

	return nil
}
