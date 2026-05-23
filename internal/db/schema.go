package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const schemaSQL = `
CREATE TABLE IF NOT EXISTS runs (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	command     TEXT    NOT NULL,
	args        TEXT    NOT NULL DEFAULT '',
	started_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	finished_at DATETIME,
	exit_code   INTEGER,
	stdout      TEXT    NOT NULL DEFAULT '',
	stderr      TEXT    NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_runs_started_at ON runs(started_at);
CREATE INDEX IF NOT EXISTS idx_runs_command    ON runs(command);
CREATE INDEX IF NOT EXISTS idx_runs_exit_code  ON runs(exit_code);
`

// Open opens (or creates) the SQLite database at the given path and
// applies the schema migrations.
func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}

	if err := applySchema(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func applySchema(db *sql.DB) error {
	_, err := db.Exec(schemaSQL)
	if err != nil {
		return fmt.Errorf("apply schema: %w", err)
	}
	return nil
}
