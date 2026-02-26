package db

import (
	"database/sql"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func Connect(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(dir)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, err
	}

	if err := goose.Up(conn, "migrations"); err != nil {
		return nil, err
	}

	return conn, err
}
