package engine

import (
	"io"
	"log"

	"github.com/deahtstroke/protheon/core/input"
	"github.com/deahtstroke/protheon/core/load"
	"github.com/deahtstroke/protheon/core/transform"
)

type ETLConfig struct {
	Input           Input      `toml:"input" yaml:"input"`
	Transformations string     `toml:"transformations" yaml:"transformations"`
	Datasource      Datasource `toml:"datasource" yaml:"datasource"`
}

type Input struct {
	Path      string `toml:"path" yaml:"path"`
	Extension string `toml:"extension" yaml:"extension"`
	Compress  string `toml:"compress" yaml:"compress"`
}

type Datasource struct {
	URL   string `toml:"url" yaml:"url"`
	Table string `toml:"table" yaml:"table"`
}

type ProtheonEngine struct {
	input.Decoder
	transform.Transformer
	load.Loader
	CleanupFuncs []func() error
}

func NewEngine() (*ProtheonEngine, error) {
	e := &ProtheonEngine{}
	return e, nil
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

func (e *ProtheonEngine) Run(cfg ETLConfig) error {
	opts := []input.FileExtractorOpt{}

	switch cfg.Input.Compress {
	case "zstd":
		opts = append(opts, input.WithZstd())
	case "gzip":
		opts = append(opts, input.WithGzip())
	default:
	}

	extractor, err := input.NewFileExtractor(cfg.Input.Path, opts...)
	if err != nil {
		return err
	}

	rc, err := extractor.Open()
	if err != nil {
		return err
	}

	decoder := input.NewJSONLDecoder(rc)
	e.Decoder = decoder
	e.CleanupFuncs = append(e.CleanupFuncs, decoder.Close)

	transformer := transform.NewLuaTransformer(cfg.Transformations)
	e.Transformer = transformer
	e.CleanupFuncs = append(e.CleanupFuncs, transformer.Close)
	loader, err := load.NewSqlLoader(cfg.Datasource.URL, cfg.Datasource.Table)
	if err != nil {
		return err
	}
	e.Loader = loader
	e.CleanupFuncs = append(e.CleanupFuncs, loader.Close)

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
