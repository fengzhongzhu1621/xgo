package migrations

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationNoTxContext(upAddUserNotTx, downAddUserNotTx)
}

func getUserID(db *sql.DB, username string) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	return id, nil
}

func upAddUserNotTx(ctx context.Context, db *sql.DB) error {
	// This code is executed when the migration is applied.
	id, err := getUserID(db, "jamesbond")
	if err != nil {
		return err
	}
	if id == 0 {
		query := "INSERT INTO users (username, name, surname) VALUES ($1, $2, $3)"
		if _, err := db.ExecContext(ctx, query, "jamesbond", "James", "Bond"); err != nil {
			return err
		}
	}
	return nil
}

func downAddUserNotTx(ctx context.Context, db *sql.DB) error {
	// This code is executed when the migration is rolled back.
	query := "DELETE FROM users WHERE username = $1"
	if _, err := db.ExecContext(ctx, query, "jamesbond"); err != nil {
		return err
	}
	return nil
}
