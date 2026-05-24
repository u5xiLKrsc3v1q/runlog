package db

import (
	"testing"
	"time"
)

func insertOldRun(t *testing.T, db_ interface{ Exec(string, ...any) (interface{}, error) }) {}

// seedRuns inserts runs with controlled timestamps via direct SQL.
func seedForPrune(t *testing.T) *tempDBHandle {
	t.Helper()
	h := tempDB(t)

	old := time.Now().Add(-48 * time.Hour).UTC().Format(time.RFC3339)
	recent := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)

	for _, row := range []struct {
		started string
		exit    int
	}{
		{old, 0},
		{old, 1},
		{recent, 0},
	} {
		_, err := h.db.Exec(
			`INSERT INTO runs (command, args, started_at, finished_at, exit_code, output)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			"echo", "[]", row.started, row.started, row.exit, "hi",
		)
		if err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
	return h
}

type tempDBHandle struct {
	db interface {
		Exec(string, ...any) (interface{ LastInsertId() (int64, error); RowsAffected() (int64, error) }, error)
		Query(string, ...any) (interface{}, error)
	}
}

func TestPruneDryRun(t *testing.T) {
	db := tempDB(t)

	old := time.Now().Add(-48 * time.Hour).UTC().Format(time.RFC3339)
	recent := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	for _, ts := range []string{old, old, recent} {
		_, err := db.Exec(`INSERT INTO runs (command,args,started_at,finished_at,exit_code,output) VALUES (?,?,?,?,?,?)`,
			"echo", "[]", ts, ts, 0, "")
		if err != nil {
			t.Fatal(err)
		}
	}

	res, err := PruneRuns(db, PruneOptions{OlderThan: 24 * time.Hour, DryRun: true})
	if err != nil {
		t.Fatalf("dry run error: %v", err)
	}
	if res.Deleted != 2 {
		t.Errorf("dry run: want 2 would-delete, got %d", res.Deleted)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM runs").Scan(&count)
	if count != 3 {
		t.Errorf("dry run must not delete rows; want 3, got %d", count)
	}
}

func TestPruneDeletes(t *testing.T) {
	db := tempDB(t)

	old := time.Now().Add(-48 * time.Hour).UTC().Format(time.RFC3339)
	recent := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	for _, ts := range []string{old, old, recent} {
		_, err := db.Exec(`INSERT INTO runs (command,args,started_at,finished_at,exit_code,output) VALUES (?,?,?,?,?,?)`,
			"echo", "[]", ts, ts, 0, "")
		if err != nil {
			t.Fatal(err)
		}
	}

	res, err := PruneRuns(db, PruneOptions{OlderThan: 24 * time.Hour})
	if err != nil {
		t.Fatalf("prune error: %v", err)
	}
	if res.Deleted != 2 {
		t.Errorf("want 2 deleted, got %d", res.Deleted)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM runs").Scan(&count)
	if count != 1 {
		t.Errorf("want 1 remaining, got %d", count)
	}
}

func TestPruneByExitCode(t *testing.T) {
	db := tempDB(t)

	old := time.Now().Add(-48 * time.Hour).UTC().Format(time.RFC3339)
	for _, exit := range []int{0, 1, 1} {
		_, err := db.Exec(`INSERT INTO runs (command,args,started_at,finished_at,exit_code,output) VALUES (?,?,?,?,?,?)`,
			"cmd", "[]", old, old, exit, "")
		if err != nil {
			t.Fatal(err)
		}
	}

	exitOne := 1
	res, err := PruneRuns(db, PruneOptions{OlderThan: 24 * time.Hour, ExitCode: &exitOne})
	if err != nil {
		t.Fatalf("prune error: %v", err)
	}
	if res.Deleted != 2 {
		t.Errorf("want 2 deleted, got %d", res.Deleted)
	}
}

func TestPruneInvalidDuration(t *testing.T) {
	db := tempDB(t)
	_, err := PruneRuns(db, PruneOptions{OlderThan: 0})
	if err == nil {
		t.Error("expected error for zero OlderThan")
	}
}
