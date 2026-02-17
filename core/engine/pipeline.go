package engine

import (
	"io"

	"github.com/deahtstroke/protheon/core/input"
)

type Executor interface {
	Execute() error
}

type ExecutorConfig struct {
	Path       string
	Compress   string
	Script     string
	Datasource string
	Table      string
}

type ProtheonExecutor struct {
	input.Decoder
}

func NewProtheonExecutor(cfg ExecutorConfig) *ProtheonExecutor {
	extractor := input.FileExtractor{
		Compress: cfg.Compress,
		Path:     cfg.Path,
	}

	rc, err := extractor.Open()
	if err != nil {
		return nil
	}

	decoder := input.NewJSONLDecoder(rc)
	return &ProtheonExecutor{
		Decoder: decoder,
	}
}

func (e *ProtheonExecutor) Execute() error {
	for {
		_, err := e.Decoder.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

