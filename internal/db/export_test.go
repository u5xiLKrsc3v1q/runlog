package db

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func makeRuns() []RunRow {
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	finished := now.Add(2 * time.Second)
	exit0 := 0
	exit1 := 1
	return []RunRow{
		{
			ID: 1, Command: "echo", Args: "hello",
			StartedAt: now, FinishedAt: &finished,
			ExitCode: &exit0, Stdout: "hello\n", Stderr: "",
		},
		{
			ID: 2, Command: "ls", Args: "-la",
			StartedAt: now, FinishedAt: &finished,
			ExitCode: &exit1, Stdout: "", Stderr: "error\n",
		},
	}
}

func TestExportJSON(t *testing.T) {
	runs := makeRuns()
	var buf bytes.Buffer
	if err := ExportJSON(runs, &buf); err != nil {
		t.Fatalf("ExportJSON error: %v", err)
	}
	var records []RunRecord
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].Command != "echo" {
		t.Errorf("expected command 'echo', got %q", records[0].Command)
	}
	if records[0].DurationMs == nil || *records[0].DurationMs != 2000 {
		t.Errorf("expected duration_ms 2000, got %v", records[0].DurationMs)
	}
	if records[1].ExitCode == nil || *records[1].ExitCode != 1 {
		t.Errorf("expected exit_code 1, got %v", records[1].ExitCode)
	}
}

func TestExportTable(t *testing.T) {
	runs := makeRuns()
	var buf bytes.Buffer
	if err := ExportTable(runs, &buf); err != nil {
		t.Fatalf("ExportTable error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "COMMAND") {
		t.Error("expected header row with COMMAND")
	}
	if !strings.Contains(out, "echo") {
		t.Error("expected 'echo' in table output")
	}
	if !strings.Contains(out, "2s") {
		t.Error("expected duration '2s' in table output")
	}
}

func TestExportJSONEmpty(t *testing.T) {
	var buf bytes.Buffer
	if err := ExportJSON([]RunRow{}, &buf); err != nil {
		t.Fatalf("ExportJSON empty error: %v", err)
	}
	if strings.TrimSpace(buf.String()) != "[]" {
		t.Errorf("expected '[]', got %q", buf.String())
	}
}
