-- +goose Up
-- +goose StatementBegin
-- Configs
CREATE TABLE IF NOT EXISTS config (
    id text PRIMARY KEY,
    path text NOT NULL,
    alias text UNIQUE,
    created_at integer NOT NULL DEFAULT (strftime ('%s', 'now')),
    updated_at integer
);

CREATE TRIGGER IF NOT EXISTS update_config_updated_at
    AFTER UPDATE ON config
BEGIN
    UPDATE config SET updated_at = strftime ('%s', 'now')
    WHERE
        id = new.id;

END;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_config_updated_at;

DROP TABLE IF EXISTS config;

-- +goose StatementEnd
