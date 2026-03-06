package config

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/deahtstroke/protheon/internal/db"
	protheonErrors "github.com/deahtstroke/protheon/internal/errors"
)

type Service interface {
	CreateConfig(context.Context, string, string) error
	DeleteConfigs(context.Context, []string) error
	ListConfigs(context.Context) error
	EditConfig(context.Context, string, []string) error
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

func (s *ConfigService) CreateConfig(ctx context.Context, format, alias string) error {
	aliasExists, err := s.Repository.ExistsByAlias(ctx, db.ExistsByAliasParams{Alias: alias})
	if err != nil {
		return err
	}

	if aliasExists {
		return errors.New("Already existing configuration with the given alias")
	}

	formatConfig, err := ResolveConfig(format)
	if err != nil {
		return err
	}

	configId, err := GenerateId(16)
	if err != nil {
		return fmt.Errorf("Error generating config ID: %s", err)
	}

	path, err := InitializeConfig(configId, getUserPreferredEditor(), formatConfig)
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
		path, err := s.Repository.GetConfigPathByAliasOrId(ctx, id)
		if err != nil {
			return err
		}

		if err = os.Remove(path); err != nil {
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

func (s *ConfigService) EditConfig(ctx context.Context, identifier string, keyValues []string) error {
	path, err := s.Repository.GetConfigPathByAliasOrId(ctx, identifier)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("No matching config for %s", identifier)
	}

	if err != nil {
		return err
	}

	if len(keyValues) > 0 {
		return MutateYaml(path, keyValues)
	} else {

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
