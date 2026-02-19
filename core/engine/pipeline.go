package engine

import (
	"io"

	"github.com/deahtstroke/protheon/core/input"
	"github.com/deahtstroke/protheon/core/load"
	"github.com/deahtstroke/protheon/core/transform"
)

type Executor interface {
	Execute() error
}

type ExecutorConfig struct {
	Path       string
	Compress   string
	Script     string
	Format     string
	Datasource string
	Table      string
}

type ProtheonExecutor struct {
	input.Decoder
	transform.Transformer
	load.Loader
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
	transformer := transform.NewLuaTransformer(cfg.Script)

	loader := load.NewSqlLoader(cfg.Datasource, cfg.Table)
	return &ProtheonExecutor{
		Decoder:     decoder,
		Transformer: transformer,
		Loader:      loader,
	}
}

func (e *ProtheonExecutor) Execute() error {
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
