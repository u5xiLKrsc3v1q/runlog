package main

import (
	"fmt"
	"os"
	"time"

	"github.com/user/runlog/internal/db"
)

// pruneUsage prints help for the prune sub-command.
func pruneUsage() {
	fmt.Fprintf(os.Stderr, `Usage: runlog prune [flags]

Delete old run records from the log database.

Flags:
  --older-than <duration>   Remove runs older than this duration (required).
                            Examples: 24h, 7d (days), 30d.
  --exit-code <int>         Only remove runs with this exit code.
  --dry-run                 Print how many rows would be deleted without
                            actually deleting them.
`)
}

// runPruneCmd parses args and executes the prune sub-command.
func runPruneCmd(dbPath string, args []string) error {
	var (
		olderThanStr string
		exitCodeStr  string
		dryRun       bool
	)

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--older-than":
			if i+1 >= len(args) {
				return fmt.Errorf("--older-than requires a value")
			}
			i++
			olderThanStr = args[i]
		case "--exit-code":
			if i+1 >= len(args) {
				return fmt.Errorf("--exit-code requires a value")
			}
			i++
			exitCodeStr = args[i]
		case "--dry-run":
			dryRun = true
		case "--help", "-h":
			pruneUsage()
			return nil
		default:
			return fmt.Errorf("unknown flag: %s", args[i])
		}
	}

	if olderThanStr == "" {
		pruneUsage()
		return fmt.Errorf("--older-than is required")
	}

	olderThan, err := parseDuration(olderThanStr)
	if err != nil {
		return fmt.Errorf("invalid --older-than %q: %w", olderThanStr, err)
	}

	opts := db.PruneOptions{
		OlderThan: olderThan,
		DryRun:    dryRun,
	}

	if exitCodeStr != "" {
		var code int
		if _, err := fmt.Sscanf(exitCodeStr, "%d", &code); err != nil {
			return fmt.Errorf("invalid --exit-code %q", exitCodeStr)
		}
		opts.ExitCode = &code
	}

	conn, err := db.Open(dbPath)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer conn.Close()

	res, err := db.PruneRuns(conn, opts)
	if err != nil {
		return fmt.Errorf("prune: %w", err)
	}

	if dryRun {
		fmt.Printf("dry-run: %d run(s) would be deleted\n", res.Deleted)
	} else {
		fmt.Printf("deleted %d run(s)\n", res.Deleted)
	}
	return nil
}

// parseDuration extends time.ParseDuration with a simple "Nd" (days) suffix.
func parseDuration(s string) (time.Duration, error) {
	var days int
	if n, _ := fmt.Sscanf(s, "%dd", &days); n == 1 {
		return time.Duration(days) * 24 * time.Hour, nil
	}
	return time.ParseDuration(s)
}
