package runner

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/user/runlog/internal/db"
)

// RunAndPersist executes the command, persists the result to the database,
// and returns the Result along with the assigned run ID.
//
// If the command itself fails (non-zero exit code), the result is still
// persisted and returned alongside a nil error. An error is only returned
// when execution cannot be attempted or persistence fails.
func RunAndPersist(ctx context.Context, database *sql.DB, command string, args ...string) (int64, *Result, error) {
	cmdLine := buildCmdLine(command, args)

	runID, err := db.InsertRun(database, cmdLine)
	if err != nil {
		return 0, nil, fmt.Errorf("insert run: %w", err)
	}

	res, runErr := Run(ctx, command, args...)

	if runErr != nil {
		// Best-effort finish with exit code -1 on unexpected error.
		_ = db.FinishRun(database, runID, "", runErr.Error(), -1)
		return runID, nil, fmt.Errorf("run command: %w", runErr)
	}

	if err := db.FinishRun(database, runID, res.Stdout, res.Stderr, res.ExitCode); err != nil {
		return runID, res, fmt.Errorf("finish run: %w", err)
	}

	return runID, res, nil
}

// buildCmdLine joins the command and its arguments into a single space-separated
// string suitable for display and storage.
func buildCmdLine(command string, args []string) string {
	if len(args) == 0 {
		return command
	}
	parts := append([]string{command}, args...)
	return strings.Join(parts, " ")
}
