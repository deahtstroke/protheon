package db

import "embed"

//go:embed migrations/*.sql
var dir embed.FS
