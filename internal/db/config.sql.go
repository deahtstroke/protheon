package db

import "context"

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
