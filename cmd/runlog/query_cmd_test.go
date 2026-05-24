package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/runlog/internal/db"
)

func tempQueryDB(t *testing.T) (string, func()) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	conn, err := db.Open(path)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	now := time.Now()
	for i, cmd := range []string{"echo hello", "ls -la", "echo world"} {
		id, _ := db.InsertRun(conn, cmd, []string{})
		exitCode := 0
		if i == 1 {
			exitCode = 1
		}
		_ = db.FinishRun(conn, id, exitCode, "output", now.Add(time.Duration(i)*time.Second))
	}
	conn.Close()
	return path, func() { os.RemoveAll(dir) }
}

func TestRunQueryCmdTable(t *testing.T) {
	dbPath, cleanup := tempQueryDB(t)
	defer cleanup()

	err := runQueryCmd([]string{"-db", dbPath, "-format", "table"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunQueryCmdJSON(t *testing.T) {
	dbPath, cleanup := tempQueryDB(t)
	defer cleanup()

	err := runQueryCmd([]string{"-db", dbPath, "-format", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunQueryCmdFilterCmd(t *testing.T) {
	dbPath, cleanup := tempQueryDB(t)
	defer cleanup()

	err := runQueryCmd([]string{"-db", dbPath, "-cmd", "echo", "-format", "table"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunQueryCmdFilterExit(t *testing.T) {
	dbPath, cleanup := tempQueryDB(t)
	defer cleanup()

	err := runQueryCmd([]string{"-db", dbPath, "-exit", "1", "-format", "table"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunQueryCmdBadDB(t *testing.T) {
	err := runQueryCmd([]string{"-db", "/nonexistent/path/runlog.db"})
	if err == nil {
		t.Fatal("expected error for bad db path")
	}
}
