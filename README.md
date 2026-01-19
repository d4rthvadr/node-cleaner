# NodeCleaner 🧹

Reclaim disk space by safely identifying and reviewing large dependency directories across common ecosystems.

## Supported Ecosystems

NodeCleaner scans for dependency folders in these stacks:

- Node.js: `node_modules`
- Go: `vendor`
- PHP: `vendor`
- Python: `venv`, `.venv`
- Rust: `target`

## Features

- Smart scanning across supported ecosystems
- Fast, concurrent worker pool
- Caching to skip unchanged directories
- Clear, colorful terminal output
- Safe by design: preview, confirm, then act

## Installation

This is a Go CLI.

```bash
# Build locally
go build

# Run (defaults use config at $HOME/.nodecleaner/config.yaml)
./node-cleaner --workers 4 scan [path]
```

Alternatively, add it to your PATH or package per your workflow.

## Usage

### Scan

```bash
# Scan your home (default from config)
./node-cleaner scan

# Scan a specific path
./node-cleaner scan /path/to/projects

# Control concurrency
./node-cleaner --workers 8 scan /path/to/projects

# Disable cache for a fresh run
./node-cleaner --workers 4 scan --no-cache /path/to/projects
```

### Config

NodeCleaner uses Viper. Configuration priority is: flags > env vars > config file > defaults.

- Config file: `$HOME/.nodecleaner/config.yaml`
- Env vars prefix: `NODECLEANER_` (e.g., `NODECLEANER_WORKERS=8`)

Display current config:

```bash
./node-cleaner config show
```

### Cache

Inspect or reset the cache:

```bash
./node-cleaner cache clear
```

## How It Works

1. Walks directories and detects dependency folders
2. Checks cache validity by modification time
3. Analyzes size and metadata (concurrently)
4. Streams results to the UI formatter
5. Enables optional interactive selection (planned)

## Safety First

- Preview and confirm before destructive actions
- Clear reporting of sizes and paths
- Conservative defaults for traversal depth and symlink handling

## Contributing

Pull requests are welcome. Please keep changes focused and tested.

## License

MIT

---

Made for developers who keep many projects and want clarity before cleanup.
