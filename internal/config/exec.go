package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/deahtstroke/protheon/internal/core/engine"
)

func InitializeConfig(id, editor string, format ConfigFormat) (string, error) {
	editor = getUserPreferredEditor()

	tempFile, err := os.CreateTemp("", format.TempPattern())
	if err != nil {
		return "", fmt.Errorf("Error creating a temporary file: %v", err)
	}

	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte(format.GetTemplate()))
	if err != nil {
		return "", fmt.Errorf("Error writing to tempfile %s: %s", tempFile.Name(), err)
	}

	cmd := exec.Command(editor, tempFile.Name())

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if err = cmd.Run(); err != nil {
		return "", err
	}

	// Validate config before saving it anywhere
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return "", err
	}

	cfg, err := format.RetrieveETLConfig(content)
	if err != nil {
		return "", err
	}

	if err := engine.VerifyFields(cfg); err != nil {
		return "", err
	}

	cfgName := format.SetConfigName(id)
	destination := filepath.Join(GlobalConfigDataDir(), cfgName)
	err = os.Rename(tempFile.Name(), destination)
	if err != nil {
		return "", fmt.Errorf("Error performing atomic write for file %s: %v", cfgName, err)
	}

	return destination, nil
}
