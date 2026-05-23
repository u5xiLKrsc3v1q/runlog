package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Run represents a single recorded process invocation.
type Run struct {
	ID         int64
	Command    string
	Args       string
	StartedAt  time.Time
	FinishedAt *time.Time
	ExitCode   *int
	Stdout     string
	Stderr     string
}

// InsertRun writes a new run record and returns the assigned ID.
func InsertRun(db *sql.DB, command string, args []string) (int64, error) {
	result, err := db.Exec(
		`INSERT INTO runs (command, args, started_at) VALUES (?, ?, ?)`,
		command, strings.Join(args, " "), time.Now().UTC(),
	)
	if err != nil {
		return 0, fmt.Errorf("insert run: %w", err)
	}
	return result.LastInsertId()
}

// FinishRun updates a run record with captured output and exit code.
func FinishRun(db *sql.DB, id int64, stdout, stderr string, exitCode int) error {
	_, err := db.Exec(
		`UPDATE runs SET finished_at=?, exit_code=?, stdout=?, stderr=? WHERE id=?`,
		time.Now().UTC(), exitCode, stdout, stderr, id,
	)
	if err != nil {
		return fmt.Errorf("finish run: %w", err)
	}
	return nil
}

// ListRuns returns the most recent runs, optionally filtered by command.
func ListRuns(db *sql.DB, command string, limit int) ([]Run, error) {
	query := `SELECT id, command, args, started_at, finished_at, exit_code, stdout, stderr
	          FROM runs`
	args := []any{}
	if command != "" {
		query += " WHERE command = ?"
		args = append(args, command)
	}
	query += " ORDER BY started_at DESC LIMIT ?"
	args = append(args, limit)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list runs: %w", err)
	}
	defer rows.Close()

	var runs []Run
	for rows.Next() {
		var r Run
		if err := rows.Scan(&r.ID, &r.Command, &r.Args, &r.StartedAt,
			&r.FinishedAt, &r.ExitCode, &r.Stdout, &r.Stderr); err != nil {
			return nil, fmt.Errorf("scan run: %w", err)
		}
		runs = append(runs, r)
	}
	return runs, rows.Err()
}
