package db

import "time"

// RunRow represents a single run record as returned from the database.
// It is shared across query, list, and export operations.
type RunRow struct {
	ID         int64
	Command    string
	Args       string
	StartedAt  time.Time
	FinishedAt *time.Time
	ExitCode   *int
	Stdout     string
	Stderr     string
}
