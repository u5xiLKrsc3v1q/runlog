package db

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// RunRecord is a fully-populated run suitable for export.
type RunRecord struct {
	ID        int64     `json:"id"`
	Command   string    `json:"command"`
	Args      string    `json:"args"`
	StartedAt time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	ExitCode  *int      `json:"exit_code,omitempty"`
	Stdout    string    `json:"stdout"`
	Stderr    string    `json:"stderr"`
	DurationMs *int64   `json:"duration_ms,omitempty"`
}

// ExportJSON writes runs as a JSON array to w.
func ExportJSON(runs []RunRow, w io.Writer) error {
	records := make([]RunRecord, 0, len(runs))
	for _, r := range runs {
		rec := RunRecord{
			ID:        r.ID,
			Command:   r.Command,
			Args:      r.Args,
			StartedAt: r.StartedAt,
			FinishedAt: r.FinishedAt,
			ExitCode:  r.ExitCode,
			Stdout:    r.Stdout,
			Stderr:    r.Stderr,
		}
		if r.FinishedAt != nil {
			ms := r.FinishedAt.Sub(r.StartedAt).Milliseconds()
			rec.DurationMs = &ms
		}
		records = append(records, rec)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

// ExportTable writes runs as a human-readable table to w.
func ExportTable(runs []RunRow, w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tCOMMAND\tARGS\tSTARTED\tEXIT\tDURATION")
	for _, r := range runs {
		exit := "-"
		if r.ExitCode != nil {
			exit = fmt.Sprintf("%d", *r.ExitCode)
		}
		dur := "-"
		if r.FinishedAt != nil {
			dur = r.FinishedAt.Sub(r.StartedAt).Round(time.Millisecond).String()
		}
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\n",
			r.ID,
			r.Command,
			r.Args,
			r.StartedAt.Format(time.RFC3339),
			exit,
			dur,
		)
	}
	return tw.Flush()
}
