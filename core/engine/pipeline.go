package engine

import (
	"io"
	"log"

	"github.com/deahtstroke/protheon/core/input"
	"github.com/deahtstroke/protheon/core/load"
	"github.com/deahtstroke/protheon/core/transform"
)

type EngineConfig struct {
	InputPath       string
	Compress        string
	Transformations string
	Format          string
	Datasource      string
	Table           string
}

type ProtheonEngine struct {
	input.Decoder
	transform.Transformer
	load.Loader
	CleanupFuncs []func() error
}

func NewProtheonEngine(cfg EngineConfig) (*ProtheonEngine, error) {
	executor := &ProtheonEngine{}

	opts := []input.FileExtractorOpt{}
	switch cfg.Compress {
	case "zstd":
		opts = append(opts, input.WithZstd())
	case "gzip":
		opts = append(opts, input.WithGzip())
	}

	extractor, err := input.NewFileExtractor(cfg.InputPath, opts...)

	rc, err := extractor.Open()
	if err != nil {
		return nil, err
	}

	decoder := input.NewJSONLDecoder(rc)
	executor.Decoder = decoder
	executor.CleanupFuncs = append(executor.CleanupFuncs, decoder.Close)

	transformer := transform.NewLuaTransformer(cfg.Transformations)
	executor.Transformer = transformer
	executor.CleanupFuncs = append(executor.CleanupFuncs, transformer.Close)

	loader, err := load.NewSqlLoader(cfg.Datasource, cfg.Table)
	if err != nil {
		return nil, err
	}

	executor.Loader = loader
	executor.CleanupFuncs = append(executor.CleanupFuncs, loader.Close)

	return executor, nil
}

func (e *ProtheonEngine) Cleanup() {
	for _, clean := range e.CleanupFuncs {
		err := clean()
		if err != nil {
			log.Printf("There was an error executing a cleanup function: %v", err)
			continue
		}
	}
}

func (e *ProtheonEngine) Run() error {
	for {
		v, err := e.Decoder.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		t, err := e.Transformer.Transform(v)
		if err != nil {
			return err
		}

		if err := e.Loader.Load(t); err != nil {
			return err
		}
	}

	return nil
}
