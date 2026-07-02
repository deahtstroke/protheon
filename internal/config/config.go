package config

import (
	"fmt"
	"strings"

	"github.com/deahtstroke/protheon/internal/core/engine"
)

// ConfigFormat uses the strategy design pattern
// Matches the corresponding config extension to its format
type ConfigFormat interface {
	// Get the config template, the default text the user sees when
	// editing the configuration in their selected text editor
	GetTemplate() string

	// Unmarshals the content of the ETLConfig to a Go struct
	ToETLConfig([]byte) (*engine.RunOpts, error)

	// The pattern that's used by the created temporary file before
	// atomically writing it to the config directory
	TempPattern() string

	// Sets the configuration name that's used in the atomic write to the config
	// directory
	SetConfigName(id string) string
}

func ResolveFormat(format string) (ConfigFormat, error) {
	switch strings.ToLower(format) {
	case "yaml", "yml", ".yaml", ".yml":
		return &YamlFormat{}, nil
	case "toml", ".toml":
		return &TomlFormat{}, nil
	default:
		return nil, fmt.Errorf("Unable to recognize configuration format %s", format)
	}
}
