# NodeCleaner 🧹

**Reclaim gigabytes of disk space by safely removing stale dependency directories**

NodeCleaner is a CLI tool that helps developers identify and remove unused `node_modules` directories from old projects, learning repos, and experiments—without breaking active projects.

## Why NodeCleaner?

Working across multiple projects means dependency directories pile up fast. A single `node_modules` folder can be hundreds of megabytes, and they add up quickly across:

- 🎓 Tutorial and learning projects
- 🏗️ Hackathon experiments
- 📦 Archived repositories
- 🔄 Active projects you haven't touched in months

NodeCleaner scans your system, finds these directories, and helps you safely reclaim disk space with full transparency and control.

## Features

- 🔍 **Smart Scanning** - Quickly identifies dependency directories across your system
- 🛡️ **Safe by Default** - Dry-run mode and confirmations before any deletion
- 📊 **Clear Reporting** - See exactly what will be removed and how much space you'll reclaim
- 🎯 **Developer-Focused** - Understands your workflow, unlike generic system cleaners
- ⚡ **Fast Performance** - Intelligent caching minimizes scan time

## Quick Start

```bash
# Install
npm install -g nodecleaner

# Scan for stale node_modules (dry-run)
nodecleaner scan

# Clean with interactive confirmation
nodecleaner clean

# Clean with size threshold (only projects > 100MB)
nodecleaner clean --min-size 100
```

## Installation

```bash
npm install -g nodecleaner
```

## Usage

### Scan Mode

Preview what would be cleaned without making changes:

```bash
nodecleaner scan
```

### Clean Mode

Remove stale dependency directories:

```bash
nodecleaner clean
```

### Options

- `--min-size <MB>` - Only clean directories larger than specified size
- `--dry-run` - Preview actions without deleting
- `--help` - Show all available commands and options

## How It Works

1. **Scans** your specified directories for dependency folders
2. **Analyzes** project metadata to determine staleness
3. **Reports** findings with clear size information
4. **Prompts** for confirmation before any deletion
5. **Cleans** safely, one directory at a time

## Safety First

NodeCleaner is designed with safety in mind:

- ✅ Dry-run mode by default
- ✅ Interactive confirmation prompts
- ✅ Detailed reporting before actions
- ✅ Smart detection to avoid breaking active projects

## Roadmap

Future enhancements may include:

- Support for other package managers (Python venv, Ruby gems, etc.)
- Automatic scheduling options
- Advanced filtering and exclusion rules
- GUI interface

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT

---

**Made with ❤️ for developers tired of "Disk Almost Full" notifications**
