package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/runlog/internal/db"
)

func queryUsage() {
	fmt.Fprintf(os.Stderr, `Usage: runlog query [options]

Query stored run logs.

Options:
  -cmd string      Filter by command substring
  -exit int        Filter by exit code (default: -1 = any)
  -since duration  Only show runs newer than this duration (e.g. 24h)
  -limit int       Maximum number of results (default 50)
  -format string   Output format: table (default) or json
  -db string       Path to SQLite database (default: runlog.db)

`)
}

func runQueryCmd(args []string) error {
	fs := flag.NewFlagSet("query", flag.ContinueOnError)
	fs.Usage = queryUsage

	cmdFilter := fs.String("cmd", "", "filter by command substring")
	exitFilter := fs.Int("exit", -1, "filter by exit code (-1 = any)")
	sinceFlag := fs.Duration("since", 0, "only show runs newer than duration")
	limitFlag := fs.Int("limit", 50, "max results")
	formatFlag := fs.String("format", "table", "output format: table or json")
	dbPath := fs.String("db", "runlog.db", "path to SQLite database")

	if err := fs.Parse(args); err != nil {
		return err
	}

	conn, err := db.Open(*dbPath)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer conn.Close()

	filter := db.RunFilter{
		CommandLike: *cmdFilter,
		Limit:       *limitFlag,
	}
	if *exitFilter >= 0 {
		code := *exitFilter
		filter.ExitCode = &code
	}
	if *sinceFlag > 0 {
		t := time.Now().Add(-*sinceFlag)
		filter.Since = &t
	}

	runs, err := db.QueryRuns(conn, filter)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	switch *formatFlag {
	case "json":
		return db.ExportJSON(os.Stdout, runs)
	default:
		return db.ExportTable(os.Stdout, runs)
	}
}
