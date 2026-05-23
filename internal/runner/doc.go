// Package runner provides utilities for executing external processes and
// capturing their output.
//
// The primary entry point for most callers is RunAndPersist, which executes a
// command and atomically records the invocation — including stdout, stderr, exit
// code, and timing — into the runlog SQLite database via the db package.
//
// For cases where persistence is not required, Run and RunWithWriter offer
// lightweight process execution with output capture.
//
// Example usage:
//
//	runID, result, err := runner.RunAndPersist(ctx, database, "go", "test", "./...")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("run %d finished with exit code %d\n", runID, result.ExitCode)
package runner
