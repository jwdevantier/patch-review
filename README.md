# patch-review

A CLI tool for managing git worktrees during patch review workflows.

## Installation

Build from source:

```bash
CGO_ENABLED=0 go build
```

Or use `go install` to fetch and build the binary, installing it into `$HOME/go/bin/patch-review`:
```bash
CGO_ENABLED=0 GOPROXY=direct go install github.com/jwdevantier/patch-review@latest
```

## Quick Start

1. Create a configuration file at `~/.config/patch-review/patch-review.config.toml`:

```toml
[settings]
branch_prefix = "review"
default_source = "qemu"

[sources.qemu]
path = "~/repos/qemu"
branch = "nvme.next"
remote = "staging"

[sources.linux]
path = "/path/to/linux"
branch = "main"
remote = "origin"
```

2. Create a review directory:

```bash
patch-review reset ~/reviews/my-patch
```

This creates a new git worktree with a fresh branch based on the latest upstream code.

3. Apply a patch:

```bash
patch-review apply ~/reviews/my-patch.patch ~/reviews/my-patch
```

4. Clean up when done:

```bash
patch-review rm ~/reviews/my-patch
```

## Commands

### reset

Creates or refreshes a git worktree directory.

```bash
patch-review reset [--source <name>] <review-dir>
```

- Fetches latest updates from the configured remote
- Creates a new branch from the upstream branch
- Creates a new worktree at the specified path
- If the worktree already exists, it removes and recreates it

Use `--source` to specify a source (defaults to `default_source` in config).

### rm

Removes a worktree and cleans up git state.

```bash
patch-review rm <review-dir>
```

- Validates the path is a known worktree
- Removes the worktree directory
- Deletes the associated git branch
- Updates state tracking

### apply

Applies a patch file to a worktree.

```bash
patch-review apply <review-dir> <patch-file>
```

- Validates the worktree path is known
- Applies the patch using `git am`

## Configuration

Config files are stored in `~/.config/patch-review/` by default. Override with `--dir`.

### patch-review.config.toml

```toml
[settings]
branch_prefix = "review"       # Prefix for created branches
default_source = "qemu"        # Default source when --source not specified

[sources.<name>]
path = "~/repos/qemu"          # Path to repository
branch = "nvme.next"           # Branch to use as base
remote = "staging"             # Git remote to fetch from
```

### State File

`patch-review.state.json` tracks active worktrees. Do not edit manually.

## Global Flags

- `--dir <path>`: Override config directory location

---

## Developer Notes

### Project Structure

```
.
├── main.go              # Entry point, CLI setup
├── cmds/
│   ├── reset.go         # reset command
│   ├── rm.go            # rm command
│   └── apply.go         # apply command
├── internal/
│   ├── config.go        # Configuration loading (TOML)
│   ├── state.go         # State persistence (JSON)
│   └── git.go           # Git operations
├── go.mod / go.sum      # Dependencies
└── DESIGN.md           # Design document
```

### Dependencies

- [spf13/cobra](https://github.com/spf13/cobra) - CLI framework
- [BurntSushi/toml](https://github.com/BurntSushi/toml) - TOML parsing

### Building

```bash
go build .
```

### Running Tests

```bash
go test ./...
```
