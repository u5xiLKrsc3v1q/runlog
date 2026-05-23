package db

import (
	"testing"
	"time"
)

func TestQueryRunsFilterByCommand(t *testing.T) {
	db := tempDB(t)

	id1, _ := InsertRun(db, "echo hello", time.Now())
	id2, _ := InsertRun(db, "ls -la", time.Now())
	FinishRun(db, id1, 0, "hello\n", "")
	FinishRun(db, id2, 0, "file1\n", "")

	runs, err := QueryRuns(db, RunFilter{Command: "echo"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(runs))
	}
	if runs[0].CmdLine != "echo hello" {
		t.Errorf("expected cmdline 'echo hello', got %q", runs[0].CmdLine)
	}
}

func TestQueryRunsFilterByExitCode(t *testing.T) {
	db := tempDB(t)

	id1, _ := InsertRun(db, "true", time.Now())
	id2, _ := InsertRun(db, "false", time.Now())
	FinishRun(db, id1, 0, "", "")
	FinishRun(db, id2, 1, "", "error")

	exitOne := 1
	runs, err := QueryRuns(db, RunFilter{ExitCode: &exitOne})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(runs))
	}
	if runs[0].CmdLine != "false" {
		t.Errorf("expected cmdline 'false', got %q", runs[0].CmdLine)
	}
}

func TestQueryRunsFilterBySince(t *testing.T) {
	db := tempDB(t)

	past := time.Now().Add(-2 * time.Hour)
	recent := time.Now()

	id1, _ := InsertRun(db, "old-cmd", past)
	id2, _ := InsertRun(db, "new-cmd", recent)
	FinishRun(db, id1, 0, "", "")
	FinishRun(db, id2, 0, "", "")

	cutoff := time.Now().Add(-1 * time.Hour)
	runs, err := QueryRuns(db, RunFilter{Since: &cutoff})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(runs))
	}
	if runs[0].CmdLine != "new-cmd" {
		t.Errorf("expected 'new-cmd', got %q", runs[0].CmdLine)
	}
}

func TestQueryRunsLimit(t *testing.T) {
	db := tempDB(t)

	for i := 0; i < 5; i++ {
		id, _ := InsertRun(db, "echo", time.Now())
		FinishRun(db, id, 0, "", "")
	}

	runs, err := QueryRuns(db, RunFilter{Limit: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(runs) != 3 {
		t.Fatalf("expected 3 runs, got %d", len(runs))
	}
}
