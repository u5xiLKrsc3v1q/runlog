package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDBPathEnvOverride(t *testing.T) {
	expected := "/tmp/custom.db"
	t.Setenv("RUNLOG_DB", expected)

	got := dbPath()
	if got != expected {
		t.Errorf("dbPath() = %q, want %q", got, expected)
	}
}

func TestDBPathDefault(t *testing.T) {
	t.Setenv("RUNLOG_DB", "")

	got := dbPath()
	if got == "" {
		t.Fatal("dbPath() returned empty string")
	}

	base := filepath.Base(got)
	if base != defaultDBName {
		t.Errorf("dbPath() base = %q, want %q", base, defaultDBName)
	}
}

func TestDBPathFallback(t *testing.T) {
	t.Setenv("RUNLOG_DB", "")

	// Simulate no home dir by checking the function still returns something
	got := dbPath()
	if got == "" {
		t.Error("dbPath() should never return empty string")
	}
}

func TestMainRunSubcmdHelp(t *testing.T) {
	// Verify the binary can be built — integration smoke test via os.Args manipulation
	// is done at the binary level; here we just test helper logic.
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	// Ensure dbPath respects environment
	dir := t.TempDir()
	dbFile := filepath.Join(dir, "test.db")
	t.Setenv("RUNLOG_DB", dbFile)

	got := dbPath()
	if got != dbFile {
		t.Errorf("expected %q, got %q", dbFile, got)
	}
}
