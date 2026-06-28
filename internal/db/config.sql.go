package db

import (
	"context"

	"github.com/deahtstroke/protheon/internal/errors"
)

const createConfg = `
INSERT INTO config (
	id,
	path,
	alias,
	created_at,
	updated_at
) VALUES (
	?, ?, ?, strftime('%s', 'now'), strftime('%s', 'now')
)
RETURNING id, path, alias, created_at, updated_at
`

type CreateConfigParams struct {
	Id    string `json:"id"`
	Path  string `json:"path"`
	Alias string `json:"alias"`
}

func (r *Repository) CreateConfig(ctx context.Context, args CreateConfigParams) (Config, error) {
	row := r.db.QueryRowContext(ctx, createConfg, args.Id, args.Path, args.Alias)

	var c Config
	err := row.Scan(
		&c.Id,
		&c.Path,
		&c.Alias,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	return c, err
}

const existsByAlias = `
SELECT EXISTS (
	SELECT 1
	FROM config c
	WHERE alias = ? AND alias != ''
);
`

type ExistsByAliasParams struct {
	Alias string `json:"alias"`
}

func (r *Repository) ExistsByAlias(ctx context.Context, args ExistsByAliasParams) (bool, error) {
	row := r.db.QueryRowContext(ctx, existsByAlias, args.Alias)

	var exists bool
	err := row.Scan(&exists)

	return exists, err
}

const getAllConfis = `
 SELECT id, alias, created_at
 FROM config c
`

func (r *Repository) GetAllConfigs(ctx context.Context) ([]Config, error) {
	rows, err := r.db.QueryContext(ctx, getAllConfis)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res []Config
	for rows.Next() {
		var config Config
		if err := rows.Scan(&config.Id, &config.Alias, &config.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, config)
	}

	return res, rows.Err()
}

const getConfigByIdOrAlias = `
	SELECT c.id, c.path, c.alias, c.created_at
	FROM config c
	WHERE c.id = ? OR c.alias = ?
`

func (r *Repository) GetConfigPathByAliasOrId(ctx context.Context, identifier string) (Config, error) {
	var config Config
	row := r.db.QueryRowContext(ctx, getConfigByIdOrAlias, identifier, identifier)
	err := row.Scan(&config.Id, &config.Path, &config.Alias, &config.CreatedAt)
	return config, err
}

const deleteByIdOrAlias = `
	DELETE FROM config WHERE id = ? OR alias = ?
`

func (r *Repository) DeleteByIdOrAlias(ctx context.Context, idOrAlias string) error {
	result, err := r.db.ExecContext(ctx, deleteByIdOrAlias, idOrAlias, idOrAlias)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &protheonErrors.NotFoundError{Identifier: idOrAlias}
	}

	return nil
}
