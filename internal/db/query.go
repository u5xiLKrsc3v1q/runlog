package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// RunFilter holds optional filter criteria for querying runs.
type RunFilter struct {
	Command    string
	ExitCode   *int
	Since      *time.Time
	Limit      int
}

// QueryRuns returns runs matching the given filter criteria.
func QueryRuns(db *sql.DB, f RunFilter) ([]Run, error) {
	var conditions []string
	var args []interface{}

	if f.Command != "" {
		conditions = append(conditions, "cmdline LIKE ?")
		args = append(args, "%"+f.Command+"%")
	}

	if f.ExitCode != nil {
		conditions = append(conditions, "exit_code = ?")
		args = append(args, *f.ExitCode)
	}

	if f.Since != nil {
		conditions = append(conditions, "started_at >= ?")
		args = append(args, f.Since.UTC().Format(time.RFC3339))
	}

	query := "SELECT id, cmdline, started_at, finished_at, exit_code, stdout, stderr FROM runs"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY started_at DESC"

	if f.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", f.Limit)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query runs: %w", err)
	}
	defer rows.Close()

	var runs []Run
	for rows.Next() {
		var r Run
		if err := rows.Scan(&r.ID, &r.CmdLine, &r.StartedAt, &r.FinishedAt, &r.ExitCode, &r.Stdout, &r.Stderr); err != nil {
			return nil, fmt.Errorf("scan run: %w", err)
		}
		runs = append(runs, r)
	}
	return runs, rows.Err()
}
