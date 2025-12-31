# lnr - Linear CLI

A command-line interface for [Linear](https://linear.app), written in Go.

## Installation

### Using Go

```bash
go install github.com/stustirling/lnr@latest
```

### From Source

```bash
git clone https://github.com/stustirling/lnr.git
cd lnr
go build -o lnr .
```

## Authentication

Set your Linear API key as an environment variable:

```bash
export LINEAR_API_KEY=your_api_key
```

You can create an API key in Linear at:
**Settings > Account > Security & Access > Personal API keys**

Verify your authentication:

```bash
lnr auth status
```

## Usage

### Issues

```bash
# List issues
lnr issue list
lnr issue list --team <team-id>
lnr issue list --assignee <user-id>
lnr issue list --limit 100

# View an issue
lnr issue view ENG-123

# Search issues
lnr issue search "login bug"
```

### Projects

```bash
# List projects
lnr project list
lnr project list --state started

# View a project
lnr project view <project-id>
```

### Initiatives

```bash
# List initiatives
lnr initiative list

# View an initiative
lnr initiative view <initiative-id>
```

### Users

```bash
# Show current user
lnr user me

# List all users
lnr user list
```

### Teams

```bash
# List teams
lnr team list
```

### Cycles

```bash
# List cycles
lnr cycle list
lnr cycle list --team <team-id>

# Show active cycle
lnr cycle active <team-id>

# View a cycle
lnr cycle view <cycle-id>
```

### Labels & States

```bash
# List labels
lnr label list
lnr label list --team <team-id>

# List workflow states
lnr state list
lnr state list --team <team-id>
```

## Output Formats

By default, output is displayed as a table. Use `--json` for JSON output:

```bash
lnr issue list --json
lnr project view <id> --json
```

## Shell Completion

Generate shell completion scripts:

```bash
# Bash
lnr completion bash > /etc/bash_completion.d/lnr

# Zsh
lnr completion zsh > "${fpath[1]}/_lnr"

# Fish
lnr completion fish > ~/.config/fish/completions/lnr.fish
```

## Development

### Building

```bash
go build -o lnr .
```

### Testing

```bash
go test ./...
```

### Linting

```bash
golangci-lint run
```

## Licence

MIT - see [LICENSE](LICENSE)
