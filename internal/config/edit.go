package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/deahtstroke/tast"
	"go.yaml.in/yaml/v4"
)

func MutateToml(path string, pairs map[string]any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	doc, err := tast.ParseBytes(data)
	if err != nil {
		return err
	}

	for k, v := range pairs {
		segments := strings.Split(k, ".")
		if len(segments) > 1 {
			tableDef := strings.Join(segments[:len(segments)-1], ".")
			keyDef := segments[len(segments)-1]

			t, ok := doc.Table(tableDef)
			if !ok {
				return fmt.Errorf("Error finding table %s", tableDef)
			}

			if err := t.Set(keyDef, v); err != nil {
				return err
			}
		} else {
			kv, ok := doc.FindKey(segments[0])
			if !ok {
				return fmt.Errorf("Error finding key-value with key %s", segments[0])
			}

			if err := kv.Set(v); err != nil {
				return err
			}
		}
	}

	return doc.Save(path)
}

// MutateYaml applies the given key/value assignments to the YAML document
// at the path, preserving the document's existing structure.
//
// Each entry in keyValuePairs has the form `<key>=<value>` where key can
// be a dotted entry. For example, "database.host=localhost" sets the
// "host" field of the database mapping to "localhost"
//
// MutateYaml returns an error if path cannot be read, if any entry in
// keyValuePairs is malformed, if a key path does not resolve to an existing
// node in the document, or if the updated document cannot be marshaled or
// written back to path
func MutateYaml(path string, pairs map[string]any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return err
	}

	for k, v := range pairs {
		keys := strings.Split(k, ".")
		if err := setYAMLValue(root.Content[0], keys, v); err != nil {
			return fmt.Errorf("set '%s': %w", k, err)
		}
	}

	out, err := yaml.Marshal(root.Content[0])
	if err != nil {
		return fmt.Errorf("Marshal yaml: %v", err)
	}

	return atomicWrite(path, out)
}

func setYAMLValue(node *yaml.Node, keys []string, value any) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("Expected a mapping node, got %v", node.Kind)
	}

	for i := 0; i < len(node.Content)-1; i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		if keyNode.Value != keys[0] {
			continue
		}

		if len(keys) == 1 {
			valueNode.Value = fmt.Sprintf("%v", value)
			valueNode.Kind = yaml.ScalarNode
			valueNode.Tag = "!!str"
			return nil
		}

		return setYAMLValue(valueNode, keys[1:], value)
	}

	return fmt.Errorf("key '%s' not found", keys[0])
}

// AtomicWrite tries to create a temporary file in the same directory
// as the provided path using a file pattern. It then
// writes out the provided bytes to the temporary file, closes it, and
// renames it to the oriignal path
//
// After everything is done the temporary file is removed
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
