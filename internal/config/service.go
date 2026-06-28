package config

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/deahtstroke/protheon/internal/db"
	protheonErrors "github.com/deahtstroke/protheon/internal/errors"
)

type EditConfigurationOptions struct {
	TransformationPath string
	Input              *ConfigurationInput
	Datasource         *ConfigurationDatasource
}

type ConfigurationInput struct {
	Path      string
	Extension string
}

type ConfigurationDatasource struct {
	Url   string
	Table string
}

func (o EditConfigurationOptions) ToKeyValues() map[string]any {
	var kvs map[string]any = make(map[string]any)

	if o.Datasource.Table != "" {
		kvs["datasource.table"] = o.Datasource.Table
	}

	if o.Datasource.Url != "" {
		kvs["datasource.url"] = o.Datasource.Url
	}

	if o.Input.Extension != "" {
		kvs["input.extension"] = o.Input.Extension
	}

	if o.Input.Path != "" {
		kvs["input.path"] = o.Input.Path
	}

	if o.TransformationPath != "" {
		kvs["transformations"] = o.TransformationPath
	}

	return kvs
}

type Service interface {
	CreateConfig(context.Context, string, string) error
	DeleteConfigs(context.Context, []string) error
	ListConfigs(context.Context) ([]db.Config, error)
	EditConfig(context.Context, string, EditConfigurationOptions) error
}

type ConfigService struct {
	Repository *db.Repository
}

func NewService(repo *db.Repository) Service {
	return &ConfigService{
		Repository: repo,
	}
}

func (s *ConfigService) ListConfigs(ctx context.Context) ([]db.Config, error) {
	configs, err := s.Repository.GetAllConfigs(ctx)
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func (s *ConfigService) CreateConfig(ctx context.Context, format, alias string) error {
	aliasExists, err := s.Repository.ExistsByAlias(ctx, db.ExistsByAliasParams{Alias: alias})
	if err != nil {
		return err
	}

	if aliasExists {
		return errors.New("Already existing configuration with the given alias")
	}

	formatConfig, err := ResolveFormat(format)
	if err != nil {
		return err
	}

	configId, err := GenerateId(16)
	if err != nil {
		return fmt.Errorf("Error generating config ID: %s", err)
	}

	path, err := InitializeConfig(configId, formatConfig)
	if err != nil {
		return err
	}

	entity, err := s.Repository.CreateConfig(ctx, db.CreateConfigParams{Id: configId, Alias: alias, Path: path})
	if err != nil {
		return fmt.Errorf("Error creating config: %s", err)
	}

	if entity.Alias != "" {
		fmt.Println(entity.Alias)
	} else {
		fmt.Println(entity.Id)
	}

	return nil
}

func (s *ConfigService) DeleteConfigs(ctx context.Context, configIds []string) error {
	var idsNotFound []string

	for _, id := range configIds {
		config, err := s.Repository.GetConfigPathByAliasOrId(ctx, id)
		if err != nil {
			return err
		}

		if err = os.Remove(config.Path); err != nil {
			return err
		}

		err = s.Repository.DeleteByIdOrAlias(ctx, fmt.Sprint(id))
		if err == nil {
			fmt.Println(id)
		}

		var notFound *protheonErrors.NotFoundError
		if errors.As(err, &notFound) {
			idsNotFound = append(idsNotFound, notFound.Identifier)
			continue
		}

		return err
	}

	if len(idsNotFound) > 0 {
		fmt.Fprintf(os.Stderr, "\nNo config found for the following identifiers:\n\n")
		for _, id := range idsNotFound {
			fmt.Fprintf(os.Stderr, "	- %s\n", id)
		}
	}

	return nil
}

func (s *ConfigService) EditConfig(ctx context.Context, id string, opts EditConfigurationOptions) error {
	config, err := s.Repository.GetConfigPathByAliasOrId(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("No matching config with Id/Alias %s", id)
	}

	if err != nil {
		return err
	}

	if kvs := opts.ToKeyValues(); len(kvs) > 0 {
		extension := filepath.Ext(config.Path)
		switch extension {
		case ".toml":
			if err := MutateToml(config.Path, kvs); err != nil {
				return err
			}
		case ".yaml", ".yml":
			if err := MutateYaml(config.Path, kvs); err != nil {
				return err
			}
		}
	} else if err := EditConfig(config.Path); err != nil {
		return err
	}

	if config.Alias != "" {
		fmt.Print(config.Alias)
	} else {
		fmt.Print(config.Id)
	}
	return nil
}

func getUserPreferredEditor() string {
	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}

	if editor == "" {
		if runtime.GOOS == "windows" {
			return "notepad"
		} else {
			return "vi"
		}
	}

	return editor
}

func GenerateId(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
