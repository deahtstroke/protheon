package config

import (
	"fmt"

	"github.com/deahtstroke/protheon/internal/core/engine"
	"github.com/pelletier/go-toml/v2"
)

var _ ConfigFormat = (*TomlFormat)(nil)

type TomlFormat struct{}

func (f *TomlFormat) GetTemplate() string {
	return tomlTemplate
}

func (f *TomlFormat) ToETLConfig(b []byte) (*engine.RunOpts, error) {
	var cfg engine.RunOpts
	if err := toml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (f *TomlFormat) TempPattern() string {
	return "protheon-*.toml"
}

func (f *TomlFormat) SetConfigName(id string) string {
	return fmt.Sprintf("%s-protheon.toml", id)
}
