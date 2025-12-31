# lnr - Linear CLI

## Project Overview
A read-only CLI tool for Linear written in Go. Uses Cobra for commands and GraphQL for the Linear API.

## Build & Test Commands
- Build: `go build -o lnr .`
- Test: `go test ./...`
- Test with coverage: `go test -cover ./...`
- Lint: `golangci-lint run`
- Run locally: `go run . <command>`

## Architecture
- `cmd/` - Root command setup
- `internal/api/` - Linear GraphQL client
- `internal/cmd/` - Command implementations (one package per noun)
- `internal/config/` - Configuration and environment
- `internal/output/` - Table/JSON formatting

## Conventions
- Command pattern: `lnr <noun> <verb>` (e.g., `lnr issue list`)
- Each command in its own file with `NewCmd<Name>` function
- Use interfaces for API client to enable mocking
- Tests alongside source files with `_test.go` suffix
- British English spellings in user-facing text

## Environment Variables
- `LINEAR_API_KEY` - Required for authentication

## Adding a New Command
1. Create package under `internal/cmd/<noun>/`
2. Add `<noun>.go` with parent command
3. Add `<verb>.go` with subcommand implementation
4. Add `<verb>_test.go` with tests
5. Register in `cmd/root.go`
