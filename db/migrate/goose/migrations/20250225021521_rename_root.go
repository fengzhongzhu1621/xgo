package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upRenameRoot, downRenameRoot)
}

func upRenameRoot(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.ExecContext(ctx, "UPDATE users SET username='admin' WHERE username='root';")
	return err
}

func downRenameRoot(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.ExecContext(ctx, "UPDATE users SET username='root' WHERE username='admin';")
	return err
}
