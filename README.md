# cli-ai

A CLI tool to discover and explore all command-line tools installed on your system.

## Features

- **Discover CLI Tools**: Scan your PATH and find all available command-line tools
- **Detailed Information**: View paths, versions, and metadata for each tool
- **AI Agent Integration**: Export tool catalogs in JSON format for AI agent consumption
- **Multiple Output Formats**: Human-readable lists or machine-readable JSON
- **Debug Capabilities**: Investigate specific tools and packages
- **Fast Scanning**: Efficiently discover hundreds of tools in seconds

## Prerequisites

- Go 1.21 or higher

## Installation

### Install Go

If you don't have Go installed, visit [https://golang.org/dl/](https://golang.org/dl/) to download and install it.

### Build the project

```bash
# Install dependencies
go mod download

# Build the binary
go build -o cli-ai .

# Or install it to your $GOPATH/bin
go install .
```

## Usage

### Getting Help

```bash
# Show main help page with all commands
cli-ai help

# Show help for a specific command
cli-ai list --help
cli-ai debug --help
```

### Commands

#### `cli-ai list` - List CLI Tools

List all available CLI tools discovered on your system.

```bash
# List all CLI tools (simple output)
cli-ai list

# List with detailed information and full paths
cli-ai list --all
cli-ai list -a

# List with verbose output
cli-ai list --verbose
cli-ai list -v

# Combine flags
cli-ai list --all --verbose

# List in JSON format (for AI agents)
cli-ai list --json
```

#### `cli-ai export` - Export Tools Catalog for AI Agents

Export a comprehensive catalog of all CLI tools in JSON format, optimized for AI agent consumption.

```bash
# Export catalog to stdout
cli-ai export

# Export to file with pretty formatting
cli-ai export --pretty --output tools.json

# Export with metadata (version, help text) - slower but more detailed
cli-ai export --with-meta --output tools-detailed.json

# Pipe to other tools
cli-ai export | jq '.tools[] | .name'
```

**AI Agent Usage**: See [docs/AI_AGENT_USAGE.md](docs/AI_AGENT_USAGE.md) for comprehensive AI integration guide.

#### `cli-ai debug` - Debug Package Information

Show detailed debug information for CLI tools and packages.

```bash
# Debug a specific package
cli-ai debug npm
cli-ai debug python
cli-ai debug docker

# Debug all packages
cli-ai debug --all
cli-ai debug -a

# Debug with verbose output
cli-ai debug npm --verbose
cli-ai debug --all -v
```

### Global Flags

These flags work with any command:

```bash
-v, --verbose          Enable verbose output
--config <file>        Specify config file (default: $HOME/.cli-ai.yaml)
```

### Examples

```bash
# Quick list of all tools
cli-ai list

# Detailed view with paths
cli-ai list --all

# Export for AI agents
cli-ai export --pretty --output tools.json

# Debug specific package
cli-ai debug nodejs

# Debug all packages with verbose output
cli-ai debug --all --verbose
```

### AI Agent Integration

cli-ai is designed to make CLI tools discoverable to AI agents:

```bash
# Export all tools in JSON format
cli-ai export --output tools.json

# Quick JSON output
cli-ai list --json

# Check if a tool exists (for AI agents)
cli-ai list --json | jq -e '.[] | select(.name=="docker")'

# Get tool path programmatically
cli-ai list --json | jq -r '.[] | select(.name=="git") | .path'
```

For comprehensive AI agent integration examples, see [docs/AI_AGENT_USAGE.md](docs/AI_AGENT_USAGE.md).

## Project Structure

```
.
├── main.go                 # Entry point
├── cmd/                    # Cobra commands
│   ├── root.go            # Root command setup
│   ├── list.go            # List command
│   └── debug.go           # Debug command
├── internal/               # Internal packages
│   ├── scanner/           # Tool discovery logic
│   │   └── scanner.go
│   └── display/           # Output formatting
│       └── display.go
├── go.mod                 # Go module definition
├── Makefile               # Build automation
└── README.md             # This file
```

## Development

### Adding New Commands

Commands are built using [Cobra](https://github.com/spf13/cobra). To add a new command:

1. Create a new file in the `cmd/` directory
2. Define your command using `cobra.Command`
3. Add it to the root command in the `init()` function

Example:
```go
var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Description of my command",
    Run: func(cmd *cobra.Command, args []string) {
        // Your logic here
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

### Adding Internal Packages

Internal packages go in the `internal/` directory and contain the core business logic:
- `scanner/`: Tool discovery and system scanning
- `display/`: Output formatting and presentation
- Add more packages as needed for your features

## Contributing

This is the initial scaffolding. Features and contribution guidelines coming soon.

## License

TBD
# cli
