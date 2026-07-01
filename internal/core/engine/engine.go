package engine

import (
	"errors"
	"io"
	"log"
	"strings"

	"github.com/deahtstroke/protheon/internal/core/input"
	"github.com/deahtstroke/protheon/internal/core/load"
	"github.com/deahtstroke/protheon/internal/core/transform"
)

type RunOpts struct {
	Input           InputOpts      `toml:"input" yaml:"input"`
	Transformations string         `toml:"transformations" yaml:"transformations"`
	Datasource      DatasourceOpts `toml:"datasource" yaml:"datasource"`
}

type InputOpts struct {
	Path       string `toml:"path" yaml:"path"`
	SampleSize int64  `toml:"sampleSize" yaml:"sampleSize"`
	Extension  string `toml:"extension" yaml:"extension"`
	Compress   string `toml:"compress" yaml:"compress"`
}

type DatasourceOpts struct {
	URL   string `toml:"url" yaml:"url"`
	Table string `toml:"table" yaml:"table"`
}

type Engine struct {
	input.Decoder
	transform.Transformer
	load.Loader
	CleanupFuncs []func() error
}

func NewEngine() (*Engine, error) {
	e := &Engine{}
	return e, nil
}

func (e *Engine) Cleanup() {
	for _, clean := range e.CleanupFuncs {
		err := clean()
		if err != nil {
			log.Printf("There was an error executing a cleanup function: %v", err)
			continue
		}
	}
}

func (e *Engine) Run(cfg RunOpts, dryRun bool) error {
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

	if dryRun {
		e.Loader = load.NewDryRunLoader(nil)
	} else {
		loader, err := load.NewSqlLoader(cfg.Datasource.URL, cfg.Datasource.Table)
		if err != nil {
			return err
		}
		e.Loader = loader
		e.CleanupFuncs = append(e.CleanupFuncs, loader.Close)
	}

	processed := int64(0)
	for {
		if cfg.Input.SampleSize > 0 && processed >= cfg.Input.SampleSize {
			break
		}

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

		processed++
	}

	return nil
}

// Validates the minimum required fields for a proper ETL pipeline run
// Note: This method only tests whether the fields are present, if the
// value itself is not a valid or is poorly formatted then this method
// won't catch it
func VerifyFields(cfg *RunOpts) error {
	errs := []string{}
	if cfg.Input.Path == "" {
		errs = append(errs, "Input path not provided")
	}

	if cfg.Input.Extension == "" {
		errs = append(errs, "Input extension not provided")
	}

	if cfg.Transformations == "" {
		errs = append(errs, "Transformations not provided or empty")
	}

	if cfg.Datasource.URL == "" {
		errs = append(errs, "Datasource URL is not defined")
	}
	if cfg.Datasource.Table == "" {
		errs = append(errs, "Datasource destination table is not defined")
	}

	if len(errs) > 0 {
		var builder strings.Builder
		builder.WriteString("There were several errors when validating the config file: \n")
		for _, msg := range errs {
			builder.WriteString("\n • " + msg)
		}

		return errors.New(builder.String())
	}

	return nil
}
