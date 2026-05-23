# runlog

Structured process runner that captures stdout/stderr per-invocation to a local SQLite log with query support.

---

## Installation

```bash
go install github.com/yourname/runlog@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/runlog.git && cd runlog && go build -o runlog .
```

---

## Usage

Run any command through `runlog` and it will execute normally while logging all output to a local SQLite database (`~/.runlog.db` by default).

```bash
# Run a command and capture its output
runlog run -- make build

# Run with a custom label
runlog run --label "nightly-build" -- ./scripts/deploy.sh

# List recent invocations
runlog list

# Query logs for a specific run by ID
runlog show 42

# Search output across all runs
runlog search "error"

# Tail the output of the last run
runlog tail
```

Output is always passed through to your terminal in real time — `runlog` captures silently in the background.

---

## Database

Logs are stored in `~/.runlog.db` by default. Override with the `--db` flag or the `RUNLOG_DB` environment variable.

```bash
RUNLOG_DB=/var/log/myapp.db runlog run -- ./server
```

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss significant changes.

---

## License

MIT © yourname