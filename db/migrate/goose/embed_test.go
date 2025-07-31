package goose

import (
	"database/sql"
	"embed"
	"testing"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func TestEmbed(t *testing.T) {
	var db *sql.DB
	// setup database

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	// run app
}
