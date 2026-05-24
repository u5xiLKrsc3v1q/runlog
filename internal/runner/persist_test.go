package runner

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/user/runlog/internal/db"
)

func tempDB(t *testing.T) *sql.DB {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")
	database, err := db.Open(path)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { database.Close() })
	return database
}

func TestRunAndPersistSuccess(t *testing.T) {
	database := tempDB(t)
	ctx := context.Background()

	runID, res, err := RunAndPersist(ctx, database, "echo", "persisted")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if runID <= 0 {
		t.Errorf("expected positive run ID, got %d", runID)
	}
	if res.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", res.ExitCode)
	}

	runs, err := db.ListRuns(database, 10)
	if err != nil {
		t.Fatalf("list runs: %v", err)
	}
	if len(runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(runs))
	}
	if runs[0].ExitCode != 0 {
		t.Errorf("persisted exit code mismatch: got %d", runs[0].ExitCode)
	}
}

func TestRunAndPersistNonZeroExit(t *testing.T) {
	database := tempDB(t)
	ctx := context.Background()

	_, res, err := RunAndPersist(ctx, database, "sh", "-c", "exit 7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ExitCode != 7 {
		t.Errorf("expected exit code 7, got %d", res.ExitCode)
	}

	runs, err := db.ListRuns(database, 10)
	if err != nil {
		t.Fatalf("list runs: %v", err)
	}
	if runs[0].ExitCode != 7 {
		t.Errorf("expected persisted exit code 7, got %d", runs[0].ExitCode)
	}
}

// TestRunAndPersistStoresCmdLine verifies that the command line string is
// correctly recorded in the database alongside the run result.
func TestRunAndPersistStoresCmdLine(t *testing.T) {
	database := tempDB(t)
	ctx := context.Background()

	_, _, err := RunAndPersist(ctx, database, "echo", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	runs, err := db.ListRuns(database, 10)
	if err != nil {
		t.Fatalf("list runs: %v", err)
	}
	if len(runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(runs))
	}
	want := "echo hello"
	if runs[0].CmdLine != want {
		t.Errorf("expected cmd line %q, got %q", want, runs[0].CmdLine)
	}
}

func TestBuildCmdLine(t *testing.T) {
	if got := buildCmdLine("echo", nil); got != "echo" {
		t.Errorf("got %q", got)
	}
	if got := buildCmdLine("echo", []string{"a", "b"}); got != "echo a b" {
		t.Errorf("got %q", got)
	}
}
