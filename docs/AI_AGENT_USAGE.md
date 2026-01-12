# Using cli with AI Agents

This guide explains how to use `cli` to make CLI tools discoverable and accessible to AI agents.

## Overview

`cli` provides machine-readable JSON output that AI agents can consume to:
- Discover all available CLI tools on a system
- Understand tool locations and paths
- Access version information and help text
- Make informed decisions about which tools to use

## Quick Start for AI Agents

### Basic Tool Discovery

```bash
# Get a simple JSON list of all tools
cli list --json

# Get detailed tool information with paths
cli export --pretty

# Export to a file for persistent access
cli export --output /path/to/tools-catalog.json
```

## Output Formats

### Simple List (cli list --json)

Returns an array of tool objects with basic information:

```json
[
  {
    "name": "git",
    "path": "/usr/bin/git",
    "size": 2847216,
    "is_symlink": false
  },
  {
    "name": "node",
    "path": "/usr/local/bin/node",
    "size": 98234896,
    "is_symlink": true,
    "symlink_to": "../lib/node_modules/nodejs/bin/node"
  }
]
```

### Full Catalog (cli export)

Returns a comprehensive catalog with metadata:

```json
{
  "total_tools": 142,
  "search_paths": [
    "/usr/local/bin",
    "/usr/bin",
    "/bin",
    "/usr/sbin",
    "/sbin"
  ],
  "tools": [
    {
      "name": "git",
      "path": "/usr/bin/git",
      "size": 2847216,
      "is_symlink": false
    },
    {
      "name": "npm",
      "path": "/usr/local/bin/npm",
      "version": "npm/10.2.3 node/v20.10.0",
      "size": 2145,
      "is_symlink": true,
      "symlink_to": "../lib/node_modules/npm/bin/npm-cli.js"
    }
  ],
  "generated_at": "2025-01-10T15:30:00Z"
}
```

### Detailed Catalog with Metadata (cli export --with-meta)

Includes version information and help text (slower to generate):

```json
{
  "total_tools": 142,
  "search_paths": ["/usr/local/bin", "/usr/bin", "/bin"],
  "tools": [
    {
      "name": "git",
      "path": "/usr/bin/git",
      "version": "git version 2.39.2",
      "help_text": "usage: git [--version] [--help] [-C <path>]...",
      "size": 2847216,
      "is_symlink": false
    }
  ],
  "generated_at": "2025-01-10T15:30:00Z"
}
```

## AI Agent Integration Patterns

### Pattern 1: Discover Available Tools

AI agents can query the system to see what tools are available:

```bash
# Get all tools in JSON format
cli list --json | jq -r '.[].name'

# Check if a specific tool exists
cli list --json | jq -r '.[] | select(.name=="docker") | .path'

# Find all Python-related tools
cli list --json | jq -r '.[] | select(.name | contains("python")) | .name'
```

### Pattern 2: Build Tool Knowledge Base

Create a persistent catalog for repeated querying:

```bash
# Generate catalog once
cli export --pretty --output ~/tools-catalog.json

# AI agent can then read from file
cat ~/tools-catalog.json | jq '.tools[] | select(.name=="npm")'
```

### Pattern 3: Enhanced Discovery with Metadata

Get version and help information for tool selection:

```bash
# Generate detailed catalog (takes longer)
cli export --with-meta --pretty --output ~/tools-detailed.json

# Query for tools with specific versions
cat ~/tools-detailed.json | jq '.tools[] | select(.version | contains("3.9"))'
```

## Use Cases

### 1. Tool Availability Check

Before executing a command, check if the tool exists:

```bash
# Check if docker is available
if cli list --json | jq -e '.[] | select(.name=="docker")' > /dev/null; then
  echo "Docker is available"
fi
```

### 2. Path Resolution

Get the exact path of a tool:

```bash
# Get git path
GIT_PATH=$(cli list --json | jq -r '.[] | select(.name=="git") | .path')
echo "Git is located at: $GIT_PATH"
```

### 3. Version Detection

Discover tool versions for compatibility checking:

```bash
# Export with metadata and check node version
cli export --with-meta --json | jq -r '.tools[] | select(.name=="node") | .version'
```

### 4. Symlink Resolution

Understand tool installation structure:

```bash
# Find all symlinked tools
cli list --json | jq '.[] | select(.is_symlink==true) | {name, symlink_to}'
```

## API Reference for AI Agents

### Commands

| Command | Output Format | Speed | Use Case |
|---------|--------------|-------|----------|
| `cli list --json` | JSON array | Fast | Quick tool discovery |
| `cli export` | JSON catalog | Fast | Comprehensive catalog |
| `cli export --with-meta` | JSON catalog + metadata | Slow | Detailed analysis |
| `cli export --pretty` | Pretty JSON | Fast | Human-readable |

### JSON Schema

#### Tool Object

```typescript
interface Tool {
  name: string;           // Tool executable name
  path: string;           // Full filesystem path
  size: number;           // File size in bytes
  is_symlink: boolean;    // Whether file is a symlink
  symlink_to?: string;    // Target of symlink (if applicable)
  version?: string;       // Version string (if --with-meta)
  help_text?: string;     // Help output (if --with-meta)
  description?: string;   // Tool description (if available)
}
```

#### Catalog Object

```typescript
interface ToolCatalog {
  total_tools: number;    // Count of discovered tools
  search_paths: string[]; // Directories searched
  tools: Tool[];          // Array of tool objects
  generated_at: string;   // ISO 8601 timestamp
}
```

## Performance Considerations

- **Fast**: `cli list --json` - Quick scan, basic info
- **Medium**: `cli export` - Full scan with paths and symlinks
- **Slow**: `cli export --with-meta` - Executes each tool to get version/help

For AI agents:
- Use `list --json` for quick checks
- Use `export` for building initial knowledge base
- Use `export --with-meta` sparingly, cache results
- Consider generating catalog once and refreshing periodically

## Examples for Common AI Tasks

### Task: Find Python Interpreter

```bash
cli list --json | jq -r '.[] | select(.name | test("^python[0-9.]*$")) | {name, path}'
```

### Task: Get All Version Control Tools

```bash
cli export | jq '.tools[] | select(.name | test("git|svn|hg|bzr")) | .name'
```

### Task: Check Docker Availability and Version

```bash
cli export --with-meta | jq '.tools[] | select(.name=="docker") | {name, path, version}'
```

### Task: List All Node.js Related Tools

```bash
cli list --json | jq -r '.[] | select(.name | contains("node") or .name | contains("npm") or .name | contains("npx"))'
```

### Task: Generate Markdown Tool List

```bash
cli export --pretty | jq -r '.tools[] | "- [\(.name)](\(.path))"'
```

## Integration Examples

### Python AI Agent

```python
import subprocess
import json

def get_available_tools():
    result = subprocess.run(
        ['cli', 'list', '--json'],
        capture_output=True,
        text=True
    )
    return json.loads(result.stdout)

def check_tool_exists(tool_name):
    tools = get_available_tools()
    return any(tool['name'] == tool_name for tool in tools)

# Usage
if check_tool_exists('docker'):
    print("Docker is available on this system")
```

### Node.js AI Agent

```javascript
const { execSync } = require('child_process');

function getAvailableTools() {
  const output = execSync('cli list --json', { encoding: 'utf-8' });
  return JSON.parse(output);
}

function getToolPath(toolName) {
  const tools = getAvailableTools();
  const tool = tools.find(t => t.name === toolName);
  return tool ? tool.path : null;
}

// Usage
const gitPath = getToolPath('git');
console.log(`Git is at: ${gitPath}`);
```

### Shell Script

```bash
#!/bin/bash

# Get all available tools as JSON
TOOLS=$(cli list --json)

# Check if required tools exist
REQUIRED_TOOLS=("git" "docker" "npm")

for tool in "${REQUIRED_TOOLS[@]}"; do
  if echo "$TOOLS" | jq -e ".[] | select(.name==\"$tool\")" > /dev/null; then
    echo "✓ $tool is available"
  else
    echo "✗ $tool is missing"
  fi
done
```

## Best Practices

1. **Cache Results**: Generate catalog once, refresh periodically
2. **Use Fast Commands**: Prefer `list --json` for real-time queries
3. **Filter Efficiently**: Use `jq` or similar tools to filter JSON
4. **Check Availability**: Always verify tool existence before execution
5. **Handle Errors**: Tool lists may change between invocations
6. **Pretty Print for Debugging**: Use `--pretty` when inspecting output manually

## Troubleshooting

### Empty Results

```bash
# Check if PATH is set
echo $PATH

# Verify cli can see PATH
cli export | jq '.search_paths'
```

### Missing Tools

```bash
# Tool might be in a directory not in PATH
# Check specific directory
ls -la /usr/local/bin | grep <tool-name>
```

### Slow Performance

```bash
# Avoid --with-meta for large tool sets
# Use basic export instead
cli export --output tools.json

# Then query specific tools if needed
cli debug <tool-name>
```
