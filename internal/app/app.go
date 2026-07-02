package app

import (
	"context"
	"database/sql"
	"os"

	"github.com/deahtstroke/protheon/internal/config"
	repo "github.com/deahtstroke/protheon/internal/db"
)

type Protheon struct {
	ConfigService config.Service
	GlobalCtx     context.Context
}

func New(ctx context.Context, db *sql.DB, debug bool) (*Protheon, error) {
	app := &Protheon{}
	dbPath := config.GlobalDatabaseUrl()

	if err := ensureDir(dbPath); err != nil {
		return nil, err
	}

	if err := ensureDir(config.GlobalConfigDataDir()); err != nil {
		return nil, err
	}

	repo := repo.NewRepository(db)
	configService := config.NewService(repo)

	app.ConfigService = configService
	app.GlobalCtx = ctx

	return app, nil
}

func ensureDir(dir string) error {
	if _, err := os.Stat(config.GlobalConfigDataDir()); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0o775)
	}
	return nil
}
