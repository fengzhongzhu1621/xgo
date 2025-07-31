package goose

import "embed"

// EmbedMigrations within the openfga binary.
//
//go:embed migrations/*.sql
var EmbedMigrations embed.FS
