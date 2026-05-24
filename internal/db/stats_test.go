package db

import (
	"testing"
	"time"
)

func TestGetStatsEmpty(t *testing.T) {
	db := tempDB(t)
	s, err := GetStats(db, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TotalRuns != 0 {
		t.Errorf("expected 0 total runs, got %d", s.TotalRuns)
	}
}

func TestGetStatsAggregates(t *testing.T) {
	db := tempDB(t)
	now := time.Now().UTC()

	insertFinished := func(cmd string, exit int, startOffset, endOffset time.Duration) {
		id, err := InsertRun(db, cmd, []string{cmd})
		if err != nil {
			t.Fatal(err)
		}
		start := now.Add(startOffset)
		end := now.Add(endOffset)
		if err := FinishRun(db, id, exit, "", start, end); err != nil {
			t.Fatal(err)
		}
	}

	insertFinished("echo", 0, 0, 100*time.Millisecond)
	insertFinished("echo", 0, 0, 200*time.Millisecond)
	insertFinished("ls", 1, 0, 50*time.Millisecond)

	s, err := GetStats(db, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TotalRuns != 3 {
		t.Errorf("TotalRuns: want 3, got %d", s.TotalRuns)
	}
	if s.SuccessRuns != 2 {
		t.Errorf("SuccessRuns: want 2, got %d", s.SuccessRuns)
	}
	if s.FailedRuns != 1 {
		t.Errorf("FailedRuns: want 1, got %d", s.FailedRuns)
	}
	if s.UniqueCommands != 2 {
		t.Errorf("UniqueCommands: want 2, got %d", s.UniqueCommands)
	}
	if s.AvgDurationMs <= 0 {
		t.Errorf("AvgDurationMs: want > 0, got %f", s.AvgDurationMs)
	}
}

func TestGetStatsFilterByCommand(t *testing.T) {
	db := tempDB(t)
	now := time.Now().UTC()

	for _, cmd := range []string{"echo", "echo", "ls"} {
		id, err := InsertRun(db, cmd, []string{cmd})
		if err != nil {
			t.Fatal(err)
		}
		_ = FinishRun(db, id, 0, "", now, now.Add(10*time.Millisecond))
	}

	s, err := GetStats(db, "echo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TotalRuns != 2 {
		t.Errorf("filtered TotalRuns: want 2, got %d", s.TotalRuns)
	}
	if s.UniqueCommands != 1 {
		t.Errorf("filtered UniqueCommands: want 1, got %d", s.UniqueCommands)
	}
}
