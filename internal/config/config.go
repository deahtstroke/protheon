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

	// Validates the content of the temporary file before writing it anywhere
	RetrieveETLConfig([]byte) (*engine.ETLConfig, error)

	// The pattern that's used by the created temporary file before
	// atomically writing it to the config directory
	TempPattern() string

	// Sets the configuration name that's used in the atomic write to the config
	// directory
	SetConfigName(id string) string
}

func ResolveConfig(format string) (ConfigFormat, error) {
	switch strings.ToLower(format) {
	case "yaml", "yml":
		return &YamlFormat{}, nil
	case "toml":
		return &TomlFormat{}, nil
	default:
		return nil, fmt.Errorf("Unable to recognize configuration format %s", format)
	}
}
