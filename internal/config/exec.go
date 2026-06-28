package config

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/deahtstroke/protheon/internal/core/engine"
)

// InitializeConfig creates a configuration based on the specified format passed
// it creates a temporary file in the current working directory to achieve an atomic write
//
// Once the user finishes editing the temp file its subsequently written to the global configuration
// directory for Protheon in the form of ~/${XDG_DATA_HOME}/protheon
//
// The return value of this function is the destination of where the file was written
func InitializeConfig(id string, format ConfigFormat) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	tempFile, err := os.CreateTemp(cwd, format.TempPattern())
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())

	editedContent, err := editTempFile(*tempFile, []byte(format.GetTemplate()))
	if err != nil {
		return "", err
	}

	cfg, err := format.ToETLConfig(editedContent)
	if err != nil {
		return "", err
	}

	if err := engine.VerifyFields(cfg); err != nil {
		return "", err
	}

	cfgName := format.SetConfigName(id)
	destination := filepath.Join(GlobalConfigDataDir(), cfgName)
	return destination, atomicWrite(destination, editedContent)
}

// EditConfig, like the name suggests, edits an already-existing Protheon configuration
// given the path to the existing configuration and the format, either YAML or TOML
//
// It opens up the configuration in the user's preferred text editor and does an atomic
// write whenever the user finishes editing
//
// The return value of this function is the original destination of the edited configuration
func EditConfig(configPath string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	format, err := ResolveFormat(filepath.Ext(configPath))
	if err != nil {
		return err
	}

	src, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	tempFile, err := os.CreateTemp(cwd, format.TempPattern())
	if err != nil {
		return err
	}

	defer os.Remove(tempFile.Name())

	content, err := editTempFile(*tempFile, src)
	if err != nil {
		return err
	}

	if err := validateConfig(content, format); err != nil {
		return err
	}

	return atomicWrite(configPath, content)
}

func editTempFile(tempFile os.File, content []byte) ([]byte, error) {
	if _, err := tempFile.Write(content); err != nil {
		return nil, err
	}

	cmd := exec.Command(getUserPreferredEditor(), tempFile.Name())

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	editedContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return nil, err
	}

	return editedContent, nil
}

func validateConfig(source []byte, format ConfigFormat) error {
	etlConfig, err := format.ToETLConfig(source)
	if err != nil {
		return err
	}
	return engine.VerifyFields(etlConfig)
}
