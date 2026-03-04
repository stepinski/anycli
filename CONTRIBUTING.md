# Contributing to anycli

Thanks for wanting to contribute. Here's how to get started fast.

## Setup
```bash
git clone https://github.com/stepinski/anycli
cd anycli
go mod tidy
make build
./anycli --help
```

## Development workflow
```bash
make test          # run tests (do this before every commit)
make test-verbose  # see individual test output
make lint          # run linters (install: brew install golangci-lint)
make build         # build binary
```

## Project structure
```
anycli/
├── main.go              # entry point — stays tiny
├── root.go              # cobra root command
├── cmd/                 # one file per subcommand
├── internal/
│   ├── api/             # AnythingLLM HTTP client
│   ├── config/          # config loading (viper)
│   └── tui/             # bubbletea TUI components
└── .github/
    └── workflows/       # CI + release automation
```

## Commit style

We use Conventional Commits:
```
feat: add --json output flag to chat command
fix: handle empty workspace slug gracefully
docs: add TUI usage examples to README
test: add mock server tests for streaming
chore: update bubbletea to v0.27
```

## Pull requests

- One PR per feature/fix
- Tests required for new functionality
- `make test` and `make lint` must pass
- Update README if you add a new command or flag
