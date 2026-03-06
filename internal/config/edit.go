package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v4"
)

func MutateToml(path string, keyValuePairs []string) error {
	kvs, err := parseKeyValues(keyValuePairs)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

}

func MutateYaml(path string, keyValuePairs []string) error {
	kvs, err := parseKeyValues(keyValuePairs)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return err
	}

	for k, v := range kvs {
		keys := strings.Split(k, ".")
		if err := setYAMLValue(root.Content[0], keys, v); err != nil {
			return fmt.Errorf("set '%s': %w", k, err)
		}
	}

	out, err := yaml.Marshal(&root)
	if err != nil {
		return fmt.Errorf("Marshal yaml: %v", err)
	}

	return atomicWrite(path, out)
}

func atomicWrite(path string, out []byte) error {
	currDir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(currDir, "*-protheon.temp")
	if err != nil {
		return err
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}()

	if _, err := tmpFile.Write(out); err != nil {
		return err
	}

	if err := tmpFile.Sync(); err != nil {
		return err
	}

	// Windows only: Cannot rename an open file
	if err := tmpFile.Close(); err != nil {
		return err
	}

	return os.Rename(tmpFile.Name(), path)
}

func setYAMLValue(node *yaml.Node, keys []string, value string) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("Expected a mapping node, got %v", node.Kind)
	}

	for i := 0; i < len(node.Content)-1; i++ {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		if keyNode.Value != keys[0] {
			continue
		}

		if len(keys) == 1 {
			valueNode.Value = value
			valueNode.Kind = yaml.ScalarNode
			valueNode.Tag = "!!str"
			return nil
		}

		return setYAMLValue(valueNode, keys[1:], value)
	}

	return fmt.Errorf("key '%s' not found", keys[0])
}

func parseKeyValues(args []string) (map[string]string, error) {
	kvs := make(map[string]string)

	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("Invalid argument '%s': Expected format key=value", arg)
		}

		key, value := parts[0], parts[1]
		if key == "" {
			return nil, fmt.Errorf("Invalid argument '%s': Key cannot be empty", arg)
		}

		if _, exists := kvs[key]; exists {
			return nil, fmt.Errorf("Duplicate key '%s': Each key may only be set once", key)
		}

		kvs[key] = value
	}

	return kvs, nil
}
