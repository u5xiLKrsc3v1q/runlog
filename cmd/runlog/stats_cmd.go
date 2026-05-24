package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/nicholasgasior/runlog/internal/db"
)

const statsUsage = `Usage: runlog stats [options]

Print aggregate statistics about stored runs.

Options:
  -cmd string    Filter statistics to runs whose command starts with this prefix
  -h             Show this help
`

func runStatsCmd(args []string, out io.Writer) int {
	fs := flag.NewFlagSet("stats", flag.ContinueOnError)
	fs.SetOutput(out)

	var cmd string
	fs.StringVar(&cmd, "cmd", "", "filter by command prefix")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(out, statsUsage)
		return 1
	}

	if fs.NArg() > 0 && fs.Arg(0) == "-h" {
		fmt.Fprintln(out, statsUsage)
		return 0
	}

	path := dbPath()
	database, err := db.Open(path)
	if err != nil {
		fmt.Fprintf(out, "error opening db: %v\n", err)
		return 1
	}
	defer database.Close()

	s, err := db.GetStats(database, cmd)
	if err != nil {
		fmt.Fprintf(out, "error fetching stats: %v\n", err)
		return 1
	}

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Metric\tValue")
	fmt.Fprintln(w, "------\t-----")
	fmt.Fprintf(w, "Total runs\t%d\n", s.TotalRuns)
	fmt.Fprintf(w, "Successful\t%d\n", s.SuccessRuns)
	fmt.Fprintf(w, "Failed\t%d\n", s.FailedRuns)
	fmt.Fprintf(w, "Unique commands\t%d\n", s.UniqueCommands)
	fmt.Fprintf(w, "Avg duration (ms)\t%.1f\n", s.AvgDurationMs)
	w.Flush()

	return 0
}

func init() {
	_ = os.Stderr // ensure os import used
}
