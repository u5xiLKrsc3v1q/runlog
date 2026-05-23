package runner

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"time"
)

// Result holds the outcome of a single process invocation.
type Result struct {
	Command  string
	Args     []string
	Stdout   string
	Stderr   string
	ExitCode int
	Started  time.Time
	Finished time.Time
}

// Run executes the given command with args, capturing stdout and stderr.
// It respects context cancellation.
func Run(ctx context.Context, command string, args ...string) (*Result, error) {
	var stdoutBuf, stderrBuf bytes.Buffer

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	started := time.Now()
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	waitErr := cmd.Wait()
	finished := time.Now()

	exitCode := 0
	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return nil, waitErr
		}
	}

	return &Result{
		Command:  command,
		Args:     args,
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		ExitCode: exitCode,
		Started:  started,
		Finished: finished,
	}, nil
}

// RunWithWriter executes the command and also streams stdout/stderr to the
// provided writer in addition to capturing them.
func RunWithWriter(ctx context.Context, w io.Writer, command string, args ...string) (*Result, error) {
	var stdoutBuf, stderrBuf bytes.Buffer

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = io.MultiWriter(&stdoutBuf, w)
	cmd.Stderr = io.MultiWriter(&stderrBuf, w)

	started := time.Now()
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	waitErr := cmd.Wait()
	finished := time.Now()

	exitCode := 0
	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return nil, waitErr
		}
	}

	return &Result{
		Command:  command,
		Args:     args,
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		ExitCode: exitCode,
		Started:  started,
		Finished: finished,
	}, nil
}
