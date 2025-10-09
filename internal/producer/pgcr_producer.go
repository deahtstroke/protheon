package producer

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	bungie "github.com/deahtstroke/protheon/internal/bungie/types"
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
	Source    string
	Publisher rabbitmq.Publisher
}

func NewPgcrProducer(source string, publisher rabbitmq.Publisher) Producer {
	return &PgcrProducer{
		Source:    source,
		Publisher: publisher,
	}
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
			var pgcr bungie.PGCR
			err = json.Unmarshal(scanner.Bytes(), &pgcr)
			if err != nil {
				return err
			}

			err = pp.Publisher.Publish(ctx, scanner.Bytes())
			if err != nil {
				log.Printf("Error publishing pgcr [%s]: %v", pgcr.ActivityDetails.InstanceID, err)
				return err
			}
		}
	}

	return nil
}
