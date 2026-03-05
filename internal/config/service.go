package config

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/deahtstroke/protheon/internal/db"
)

type Service interface {
	CreateConfig(context.Context, string, string) (*db.Config, error)
	ListConfigs(context.Context) error
}

type ConfigService struct {
	Repository *db.Repository
}

func NewService(repo *db.Repository) Service {
	return &ConfigService{
		Repository: repo,
	}
}

func (s *ConfigService) ListConfigs(ctx context.Context) error {
	configs, err := s.Repository.GetAllConfigs(ctx)
	if err != nil {
		return err
	}

	headers := []string{"CONFIG ID", "ALIAS", "CREATED AT"}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	defer w.Flush()

	header := strings.Join(headers, "\t")
	fmt.Fprintln(w, header)

	for _, config := range configs {
		t := time.Unix(config.CreatedAt, 0)
		fmt.Fprintf(w, "%s\t%s\t%s\n", config.Id, config.Alias, t.Format(time.RFC1123))
	}

	return nil
}

func (s *ConfigService) CreateConfig(ctx context.Context, format, alias string) (*db.Config, error) {
	aliasExists, err := s.Repository.ExistsByAlias(ctx, db.ExistsByAliasParams{Alias: alias})
	if err != nil {
		return nil, err
	}

	if aliasExists {
		return nil, errors.New("Already existing configuration with the given alias")
	}

	formatConfig, err := ResolveConfig(format)
	if err != nil {
		return nil, err
	}

	configId, err := GenerateId(16)
	if err != nil {
		return nil, fmt.Errorf("Error generating config ID: %s", err)
	}

	path, err := InitializeConfig(configId, getUserPreferredEditor(), formatConfig)
	if err != nil {
		return nil, err
	}

	entity, err := s.Repository.CreateConfig(ctx, db.CreateConfigParams{Id: configId, Alias: alias, Path: path})
	if err != nil {
		return nil, fmt.Errorf("Error creating config: %s", err)
	}

	return &entity, nil
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
