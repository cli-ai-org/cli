# cli Command Reference

Complete reference for all `cli` commands and their usage.

## Command Syntax

All commands follow the pattern:
```
cli <command> [arguments] [flags]
```

## Available Commands

### `cli help`

Display help information for cli.

**Usage:**
```bash
cli help                    # Show main help
cli help <command>          # Show help for specific command
cli <command> --help        # Alternative syntax
```

**Examples:**
```bash
cli help
cli help list
cli list --help
```

---

### `cli list`

List all CLI tools available on your system by scanning PATH directories.

**Usage:**
```bash
cli list [flags]
```

**Flags:**
- `-a, --all` - Show detailed information including full paths
- `-v, --verbose` - Enable verbose output
- `--config <file>` - Specify config file

**Examples:**
```bash
# Simple list
cli list

# Detailed list with paths
cli list --all

# Verbose output
cli list --verbose

# Combine flags
cli list --all --verbose
```

**Output:**
- Default: Simple list of tool names
- With `--all`: Tool names with full paths and metadata
- With `--verbose`: Additional debugging information

---

### `cli export`

Export a comprehensive catalog of CLI tools in JSON format for AI agents and programmatic access.

**Usage:**
```bash
cli export [flags]
```

**Flags:**
- `-j, --json` - Output in JSON format (default: true)
- `-p, --pretty` - Pretty-print JSON output
- `-o, --output <file>` - Write to file instead of stdout
- `-m, --with-meta` - Include version and help text (slower)
- `-v, --verbose` - Enable verbose output

**Examples:**
```bash
# Export to stdout
cli export

# Export to file with formatting
cli export --pretty --output tools.json

# Export with full metadata
cli export --with-meta --pretty --output tools-detailed.json

# Pipe to jq for processing
cli export | jq '.tools[] | select(.name=="docker")'

# Export verbose mode
cli export --verbose --output tools.json
```

**Output:**
The export command generates a JSON catalog containing:
- Total number of tools discovered
- List of search paths (PATH directories)
- Array of tool objects with detailed information
- Timestamp of catalog generation

**JSON Structure:**
```json
{
  "total_tools": 142,
  "search_paths": ["/usr/local/bin", "/usr/bin", "/bin"],
  "tools": [
    {
      "name": "git",
      "path": "/usr/bin/git",
      "size": 2847216,
      "is_symlink": false,
      "version": "git version 2.39.2",
      "help_text": "usage: git [--version]..."
    }
  ],
  "generated_at": "2025-01-10T15:30:00Z"
}
```

**Performance Notes:**
- Basic export: Fast (< 1 second)
- With `--with-meta`: Slower (may take 10-30 seconds for 100+ tools)
- Use `--verbose` to see progress during metadata collection

**AI Agent Usage:**
This command is optimized for AI agents. See [docs/AI_AGENT_USAGE.md](docs/AI_AGENT_USAGE.md) for integration examples.

---

### `cli debug`

Show detailed debug information for CLI tools and packages.

**Usage:**
```bash
cli debug <package_name>    # Debug specific package
cli debug --all             # Debug all packages
```

**Arguments:**
- `<package_name>` - Name of the package to debug (optional if using --all)

**Flags:**
- `-a, --all` - Show debug info for all packages
- `-v, --verbose` - Enable verbose output
- `--config <file>` - Specify config file

**Examples:**
```bash
# Debug specific package
cli debug npm
cli debug python
cli debug docker

# Debug all packages
cli debug --all

# Debug with verbose output
cli debug npm --verbose
cli debug --all --verbose
```

**Output:**
Displays:
- Package location
- Installation path
- Binary details
- Dependencies (if available)
- Version information (if available)

---

## Global Flags

These flags work with any command:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--verbose` | `-v` | Enable verbose output | `false` |
| `--config` | - | Config file path | `$HOME/.cli.yaml` |
| `--help` | `-h` | Show help for command | - |

---

## Command Quick Reference

| Command | Purpose | Common Usage |
|---------|---------|--------------|
| `cli help` | Show help | `cli help` |
| `cli list` | List all CLI tools | `cli list --all` |
| `cli list --json` | List in JSON format | `cli list --json` |
| `cli export` | Export catalog for AI | `cli export --pretty -o tools.json` |
| `cli debug <pkg>` | Debug package | `cli debug npm` |
| `cli debug --all` | Debug all packages | `cli debug --all` |

---

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | General error |
| `2` | Invalid command or arguments |

---

## Environment Variables

| Variable | Description |
|----------|-------------|
| `PATH` | Directories scanned for CLI tools |
| `CLIIL_CONFIG` | Override default config file location |

---

## Configuration File

Default location: `$HOME/.cli.yaml`

**Example configuration:**
```yaml
# TODO: Configuration schema to be defined
```

---

## Common Workflows

### Discovering all tools on your system
```bash
cli list --all
```

### Finding information about a specific tool
```bash
cli debug <tool_name>
```

### Getting detailed output for everything
```bash
cli list --all --verbose
cli debug --all --verbose
```

---

## Tips

1. Use `--verbose` for troubleshooting
2. Use `--all` to see full paths and details
3. Combine flags for maximum information: `--all --verbose`
4. Use `cli help <command>` for detailed command help
5. Use tab completion (if installed) for faster command entry
