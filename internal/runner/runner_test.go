package runner

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunSuccess(t *testing.T) {
	ctx := context.Background()
	res, err := Run(ctx, "echo", "hello", "world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", res.ExitCode)
	}
	if !strings.Contains(res.Stdout, "hello world") {
		t.Errorf("expected stdout to contain 'hello world', got %q", res.Stdout)
	}
	if res.Finished.Before(res.Started) {
		t.Error("finished time should be after started time")
	}
}

func TestRunNonZeroExit(t *testing.T) {
	ctx := context.Background()
	res, err := Run(ctx, "sh", "-c", "exit 42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ExitCode != 42 {
		t.Errorf("expected exit code 42, got %d", res.ExitCode)
	}
}

func TestRunCapturesStderr(t *testing.T) {
	ctx := context.Background()
	res, err := Run(ctx, "sh", "-c", "echo errout >&2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Stderr, "errout") {
		t.Errorf("expected stderr to contain 'errout', got %q", res.Stderr)
	}
}

func TestRunContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	_, err := Run(ctx, "sleep", "10")
	if err == nil {
		t.Error("expected an error for cancelled context")
	}
}

func TestRunWithWriter(t *testing.T) {
	var buf bytes.Buffer
	ctx := context.Background()
	res, err := RunWithWriter(ctx, &buf, "echo", "streamed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "streamed") {
		t.Errorf("expected writer to contain 'streamed', got %q", buf.String())
	}
	if !strings.Contains(res.Stdout, "streamed") {
		t.Errorf("expected captured stdout to contain 'streamed', got %q", res.Stdout)
	}
}
