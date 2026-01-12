# cli

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
go build -o cli .

# Or install it to your $GOPATH/bin
go install .
```

## Usage

### Getting Help

```bash
# Show main help page with all commands
cli help

# Show help for a specific command
cli list --help
cli debug --help
```

### Commands

#### `cli list` - List CLI Tools

List all available CLI tools discovered on your system.

```bash
# List all CLI tools (simple output)
cli list

# List with detailed information and full paths
cli list --all
cli list -a

# List with verbose output
cli list --verbose
cli list -v

# Combine flags
cli list --all --verbose

# List in JSON format (for AI agents)
cli list --json
```

#### `cli export` - Export Tools Catalog for AI Agents

Export a comprehensive catalog of all CLI tools in JSON format, optimized for AI agent consumption.

```bash
# Export catalog to stdout
cli export

# Export to file with pretty formatting
cli export --pretty --output tools.json

# Export with metadata (version, help text) - slower but more detailed
cli export --with-meta --output tools-detailed.json

# Pipe to other tools
cli export | jq '.tools[] | .name'
```

**AI Agent Usage**: See [docs/AI_AGENT_USAGE.md](docs/AI_AGENT_USAGE.md) for comprehensive AI integration guide.

#### `cli debug` - Debug Package Information

Show detailed debug information for CLI tools and packages.

```bash
# Debug a specific package
cli debug npm
cli debug python
cli debug docker

# Debug all packages
cli debug --all
cli debug -a

# Debug with verbose output
cli debug npm --verbose
cli debug --all -v
```

### Global Flags

These flags work with any command:

```bash
-v, --verbose          Enable verbose output
--config <file>        Specify config file (default: $HOME/.cli.yaml)
```

### Examples

```bash
# Quick list of all tools
cli list

# Detailed view with paths
cli list --all

# Export for AI agents
cli export --pretty --output tools.json

# Debug specific package
cli debug nodejs

# Debug all packages with verbose output
cli debug --all --verbose
```

### AI Agent Integration

cli is designed to make CLI tools discoverable to AI agents:

```bash
# Export all tools in JSON format
cli export --output tools.json

# Quick JSON output
cli list --json

# Check if a tool exists (for AI agents)
cli list --json | jq -e '.[] | select(.name=="docker")'

# Get tool path programmatically
cli list --json | jq -r '.[] | select(.name=="git") | .path'
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

We welcome contributions! Here's how you can help:

### Reporting Issues

- **Bug Reports**: Open an issue with details about the bug, including steps to reproduce
- **Feature Requests**: Open an issue describing the feature and why it would be useful
- **Questions**: Feel free to open an issue for questions about usage or development

### Contributing Code

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/cli.git
   cd cli
   ```
3. **Create a branch** for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```
4. **Make your changes** and test them:
   ```bash
   go build -o cli .
   ./cli list  # Test your changes
   ```
5. **Commit your changes** with clear commit messages:
   ```bash
   git add .
   git commit -m "Add feature: description of your changes"
   ```
6. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```
7. **Open a Pull Request** on GitHub with:
   - Clear description of what you changed and why
   - Any related issue numbers (e.g., "Fixes #123")
   - Test results showing your changes work

### Development Guidelines

- **Code Style**: Follow standard Go conventions (run `go fmt`)
- **Testing**: Test your changes thoroughly before submitting
- **Documentation**: Update README.md if adding new features
- **Commit Messages**: Use clear, descriptive commit messages
- **Small PRs**: Keep pull requests focused on a single feature or fix

### Areas to Contribute

- **Package Manager Support**: Add detection for more package managers (Poetry, Conda, etc.)
- **Filtering Improvements**: Better heuristics for identifying actual CLI tools
- **Performance**: Optimize scanning for large systems
- **Documentation**: Improve docs, add examples, write guides
- **Testing**: Add unit tests and integration tests
- **Bug Fixes**: Fix reported issues

### Getting Help

- Open an issue if you need help or have questions
- Check existing issues and PRs to avoid duplicates
- Feel free to ask questions in your PR

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
