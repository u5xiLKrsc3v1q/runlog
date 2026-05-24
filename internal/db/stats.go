package db

import (
	"database/sql"
	"fmt"
)

// Stats holds aggregate statistics over stored runs.
type Stats struct {
	TotalRuns    int
	SuccessRuns  int
	FailedRuns   int
	UniqueCommands int
	AvgDurationMs float64
}

// GetStats returns aggregate statistics for all runs, optionally filtered by
// command prefix. Pass an empty string to include all commands.
func GetStats(db *sql.DB, command string) (Stats, error) {
	var s Stats

	baseWhere := "WHERE 1=1"
	args := []any{}
	if command != "" {
		baseWhere += " AND command LIKE ?"
		args = append(args, command+"%")
	}

	row := db.QueryRow(fmt.Sprintf(`
		SELECT
			COUNT(*),
			COUNT(CASE WHEN exit_code = 0 THEN 1 END),
			COUNT(CASE WHEN exit_code != 0 THEN 1 END),
			COUNT(DISTINCT command),
			COALESCE(AVG(CASE WHEN finished_at IS NOT NULL
				THEN (julianday(finished_at) - julianday(started_at)) * 86400000
				END), 0)
		FROM runs %s`, baseWhere), args...)

	err := row.Scan(
		&s.TotalRuns,
		&s.SuccessRuns,
		&s.FailedRuns,
		&s.UniqueCommands,
		&s.AvgDurationMs,
	)
	if err != nil {
		return Stats{}, fmt.Errorf("stats query: %w", err)
	}
	return s, nil
}
