package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/example/runlog/internal/db"
	"github.com/example/runlog/internal/runner"
)

const defaultDBName = ".runlog.db"

func usage() {
	fmt.Fprintf(os.Stderr, `runlog - structured process runner with SQLite logging

Usage:
  runlog run <command> [args...]   Run a command and log it
  runlog query [options]           Query past runs
  runlog prune [options]           Remove old run records
  runlog help                      Show this help

Options:
  --db <path>   Path to SQLite database (default: $HOME/.runlog.db)

`)
}

func dbPath() string {
	if v := os.Getenv("RUNLOG_DB"); v != "" {
		return v
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return defaultDBName
	}
	return filepath.Join(home, defaultDBName)
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		usage()
		os.Exit(0)
	}

	path := dbPath()

	// Allow global --db flag before subcommand
	if len(args) >= 2 && args[0] == "--db" {
		path = args[1]
		args = args[2:]
	}

	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	conn, err := db.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runlog: failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	subcmd := args[0]
	rest := args[1:]

	switch subcmd {
	case "run":
		if len(rest) == 0 {
			fmt.Fprintln(os.Stderr, "runlog run: command required")
			os.Exit(1)
		}
		result, err := runner.RunAndPersist(conn, rest[0], rest[1:]...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runlog: %v\n", err)
			os.Exit(1)
		}
		os.Exit(result.ExitCode)
	case "query":
		runQueryCmd(conn, rest)
	case "prune":
		runPruneCmd(conn, rest)
	default:
		fmt.Fprintf(os.Stderr, "runlog: unknown subcommand %q\n", subcmd)
		usage()
		os.Exit(1)
	}
}
