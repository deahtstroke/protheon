package config

import (
	"fmt"

	"github.com/deahtstroke/protheon/internal/core/engine"
	"go.yaml.in/yaml/v4"
)

type YamlFormat struct{}

var _ ConfigFormat = (*YamlFormat)(nil)

func (f *YamlFormat) GetTemplate() string {
	return yamlTemplate
}

func (f *YamlFormat) ToETLConfig(b []byte) (*engine.RunOpts, error) {
	var cfg engine.RunOpts
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (f *YamlFormat) TempPattern() string {
	return "protheon-*.yaml"
}

func (f *YamlFormat) SetConfigName(id string) string {
	return fmt.Sprintf("%s-protheon.yaml", id)
}
