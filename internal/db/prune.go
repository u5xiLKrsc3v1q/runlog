package db

import (
	"database/sql"
	"fmt"
	"time"
)

// PruneOptions controls which runs are deleted.
type PruneOptions struct {
	// OlderThan removes runs whose started_at is before Now-OlderThan.
	OlderThan time.Duration
	// ExitCode, when non-nil, restricts deletion to runs with that exit code.
	ExitCode *int
	// DryRun reports how many rows would be deleted without deleting them.
	DryRun bool
}

// PruneResult summarises the outcome of a prune operation.
type PruneResult struct {
	Deleted int64
}

// PruneRuns deletes run records (and their output) that match opts.
// It returns the number of rows affected (or that would be affected in dry-run mode).
func PruneRuns(db *sql.DB, opts PruneOptions) (PruneResult, error) {
	if opts.OlderThan <= 0 {
		return PruneResult{}, fmt.Errorf("prune: OlderThan must be positive")
	}

	cutoff := time.Now().Add(-opts.OlderThan).UTC().Format(time.RFC3339)

	query := "SELECT id FROM runs WHERE started_at < ?"
	args := []any{cutoff}

	if opts.ExitCode != nil {
		query += " AND exit_code = ?"
		args = append(args, *opts.ExitCode)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return PruneResult{}, fmt.Errorf("prune query: %w", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return PruneResult{}, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return PruneResult{}, err
	}

	if opts.DryRun || len(ids) == 0 {
		return PruneResult{Deleted: int64(len(ids))}, nil
	}

	tx, err := db.Begin()
	if err != nil {
		return PruneResult{}, err
	}
	defer tx.Rollback() //nolint:errcheck

	var total int64
	for _, id := range ids {
		if _, err := tx.Exec("DELETE FROM runs WHERE id = ?", id); err != nil {
			return PruneResult{}, fmt.Errorf("prune delete id=%d: %w", id, err)
		}
		total++
	}

	if err := tx.Commit(); err != nil {
		return PruneResult{}, err
	}
	return PruneResult{Deleted: total}, nil
}
