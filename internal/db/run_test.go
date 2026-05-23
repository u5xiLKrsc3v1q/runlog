package db_test

import (
	"os"
	"testing"

	"github.com/user/runlog/internal/db"
)

func tempDB(t *testing.T) *sqlDB {
	t.Helper()
	f, err := os.CreateTemp("", "runlog-test-*.db")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	conn, err := db.Open(f.Name())
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

func TestInsertAndFinishRun(t *testing.T) {
	conn := tempDB(t)

	id, err := db.InsertRun(conn, "echo", []string{"hello", "world"})
	if err != nil {
		t.Fatalf("InsertRun: %v", err)
	}
	if id <= 0 {
		t.Fatalf("expected positive id, got %d", id)
	}

	if err := db.FinishRun(conn, id, "hello world\n", "", 0); err != nil {
		t.Fatalf("FinishRun: %v", err)
	}
}

func TestListRuns(t *testing.T) {
	conn := tempDB(t)

	for i := 0; i < 5; i++ {
		id, err := db.InsertRun(conn, "ls", []string{"-la"})
		if err != nil {
			t.Fatal(err)
		}
		db.FinishRun(conn, id, "output", "", 0)
	}

	runs, err := db.ListRuns(conn, "", 10)
	if err != nil {
		t.Fatalf("ListRuns: %v", err)
	}
	if len(runs) != 5 {
		t.Fatalf("expected 5 runs, got %d", len(runs))
	}

	filtered, err := db.ListRuns(conn, "ls", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(filtered) != 5 {
		t.Fatalf("expected 5 filtered runs, got %d", len(filtered))
	}

	none, err := db.ListRuns(conn, "nonexistent", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(none) != 0 {
		t.Fatalf("expected 0 runs for unknown command, got %d", len(none))
	}
}
