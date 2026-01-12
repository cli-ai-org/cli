# Quick Start for AI Agents

The fastest way to integrate cli with AI agents.

## 3-Minute Quick Start

### 1. Install and Build

```bash
cd /Users/op/code/cli/cli
go mod download
go build -o cli .
```

### 2. Basic Usage

```bash
# Discover all CLI tools (JSON)
./cli list --json

# Export full catalog
./cli export --pretty --output tools.json

# Check if a tool exists
./cli list --json | jq -e '.[] | select(.name=="docker")'
```

### 3. Integrate with Your AI Agent

**Python:**
```python
import subprocess, json

def get_tools():
    result = subprocess.run(['cli', 'export'], capture_output=True, text=True)
    return json.loads(result.stdout)

tools = get_tools()
print(f"Found {tools['total_tools']} tools")
```

**Node.js:**
```javascript
const { execSync } = require('child_process');
const tools = JSON.parse(execSync('cli export', { encoding: 'utf-8' }));
console.log(`Found ${tools.total_tools} tools`);
```

**Shell:**
```bash
cli export | jq '.tools[] | select(.name=="git")'
```

## Common AI Agent Queries

```bash
# List all tools
cli list --json

# Get tool path
cli list --json | jq -r '.[] | select(.name=="git") | .path'

# Check tool exists
cli list --json | jq -e '.[] | select(.name=="docker")' && echo "exists"

# Find Python tools
cli list --json | jq '.[] | select(.name | contains("python"))'

# Get symlinked tools
cli list --json | jq '.[] | select(.is_symlink==true)'

# Export with versions (slower)
cli export --with-meta | jq '.tools[] | {name, version}'
```

## Best Practices

1. **Cache Results**: Generate catalog once, query many times
2. **Use Fast Commands**: `list --json` for quick checks
3. **Filter with jq**: Process JSON efficiently
4. **Background Generation**: Use `--with-meta` in background tasks
5. **Error Handling**: Check exit codes and parse JSON carefully

## Need More?

- Full guide: [docs/AI_AGENT_USAGE.md](docs/AI_AGENT_USAGE.md)
- Commands: [COMMANDS.md](COMMANDS.md)
- Examples: [USAGE_EXAMPLES.md](USAGE_EXAMPLES.md)
