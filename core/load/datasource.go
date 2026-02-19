package load

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type Loader interface {
	Load(input any) error
	Close() error
}

type SqlLoader struct {
	URL   string
	Db    *sql.DB
	Table string
}

func NewSqlLoader(url, table string) (*SqlLoader, error) {
	var db *sql.DB
	var err error

	switch {
	case strings.Contains(url, "postgres://"):
		db, err = sql.Open("postgres", url)
		if err != nil {
			return nil, err
		}
	case strings.Contains(url, "mysql://"):
		db, err = sql.Open("mysql", url)
		if err != nil {
			return nil, err
		}
	}

	return &SqlLoader{
		URL:   url,
		Db:    db,
		Table: table,
	}, nil
}

func (s *SqlLoader) Close() error {
	return s.Db.Close()
}

func (l *SqlLoader) Load(input any) error {
	data, ok := input.(map[string]any)
	if !ok {
		return fmt.Errorf("Cannot infer the type of the underlying input as `map[string]any`")
	}

	return InsertAny(l.Db, l.Table, data)
}

func InsertAny(db *sql.DB, tableName string, data map[string]any) error {
	var columns []string
	var placeholders []string
	var values []any

	i := 1
	for col, val := range data {
		columns = append(columns, col)
		values = append(values, val)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		i++
	}

	query := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s);`,
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	_, err := db.Exec(query, values...)
	return err
}
