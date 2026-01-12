# Package Detection Feature

cli can now detect which packages (npm, pip, brew, cargo, gem) provide CLI tools and link CLI tools back to their source packages.

## Overview

When you install a package like `vercel` or `supabase`, you get CLI tools. cli now:
1. Detects packages from various package managers
2. Links CLI tools to their source packages
3. Shows which tools come from which packages
4. Exports this information in JSON for AI agents

## Commands

### `cli packages`

List all packages that provide CLI tools.

```bash
# List all packages with CLIs
cli packages

# Show output
Found 158 packages with CLI tools:

PACKAGE                        MANAGER    VERSION         CLIs
-------                        -------    -------         ----
supabase                       brew       2.65.5          supabase
vercel-cli                     brew       41.6.2          [vc vercel]
@openai/codex                  npm        0.65.0          codex
pytest                         brew       8.3.4           pytest
docker                         brew       27.5.1          docker
```

### Filter by Package Manager

```bash
# Show only npm packages
cli packages --manager npm

# Show only brew packages
cli packages --manager brew

# Show only pip packages
cli packages --manager pip
```

### JSON Output

```bash
# Export packages in JSON format
cli packages --json

# Example output
[
  {
    "name": "vercel-cli",
    "version": "41.6.2",
    "manager": "brew",
    "binaries": ["vc", "vercel"],
    "global": true
  },
  {
    "name": "supabase",
    "version": "2.65.5",
    "manager": "brew",
    "binaries": ["supabase"],
    "global": true
  }
]
```

## Export with Package Information

The `export` command can now include package information:

```bash
# Export catalog with package info
cli export --with-packages --pretty --output tools-with-packages.json

# Result includes package_name, package_manager, package_version
{
  "total_tools": 2960,
  "total_packages": 158,
  "tools": [
    {
      "name": "vercel",
      "path": "/opt/homebrew/bin/vercel",
      "package_name": "vercel-cli",
      "package_manager": "brew",
      "package_version": "41.6.2",
      ...
    }
  ],
  "packages": [
    {
      "name": "vercel-cli",
      "version": "41.6.2",
      "manager": "brew",
      "binaries": ["vc", "vercel"]
    }
  ]
}
```

## Use Cases

### 1. Find Which Package Provides a CLI

```bash
# List packages and grep for the one you want
cli packages | grep vercel

# Output:
vercel-cli                     brew       41.6.2          [vc vercel]
```

### 2. See All CLIs from a Package

```bash
cli packages --json | jq '.[] | select(.name=="vercel-cli")'

# Shows all binaries provided by vercel-cli
{
  "name": "vercel-cli",
  "version": "41.6.2",
  "manager": "brew",
  "binaries": ["vc", "vercel"],
  "global": true
}
```

### 3. List All npm Package CLIs

```bash
cli packages --manager npm

# Shows only npm global packages with CLIs
Found 3 packages with CLI tools:

PACKAGE                        MANAGER    VERSION         CLIs
-------                        -------    -------         ----
@openai/codex                  npm        0.65.0          codex
corepack                       npm        0.34.0          corepack
npm                            npm        11.6.1          [npm npx]
```

### 4. AI Agent Query: Check Tool Source

```bash
# Get package info for a specific tool
cli export --with-packages | jq '.tools[] | select(.name=="vercel") | {name, package_name, package_manager, package_version}'

# Output:
{
  "name": "vercel",
  "package_name": "vercel-cli",
  "package_manager": "brew",
  "package_version": "41.6.2"
}
```

### 5. Find All Tools from a Specific Manager

```bash
# Find all Homebrew tools
cli export --with-packages | jq '.tools[] | select(.package_manager=="brew") | .name' | head -20
```

## Package Managers Supported

| Manager | Detection | Linking |
|---------|-----------|---------|
| npm | Global packages via `npm list -g` | ✓ Path-based + node_modules |
| pip | All packages via `pip list` | ✓ Path-based |
| Homebrew | All packages via `brew list` | ✓ Cellar path + symlinks |
| cargo | Installed packages via `cargo install --list` | ✓ .cargo/bin path |
| gem | Local gems via `gem list` | ✓ Path-based |

## How Linking Works

cli uses multiple strategies to link CLIs to packages:

1. **Direct Name Match**: Tool name matches package name (e.g., `supabase` → `supabase`)
2. **Path Detection**: Extracts package from installation path:
   - npm: `/path/node_modules/package/bin/tool`
   - Homebrew: `/opt/homebrew/Cellar/package/version/bin/tool`
   - pip: Detected via package manager
3. **Symlink Following**: Checks symlink targets for package information
4. **Pattern Matching**: Handles common patterns like `package-cli` → `package`

## Examples

### Before Package Detection

```bash
cli list --json | jq '.[] | select(.name=="vercel")'

{
  "name": "vercel",
  "path": "/opt/homebrew/bin/vercel",
  "size": 38
}
```

### After Package Detection

```bash
cli export --with-packages | jq '.tools[] | select(.name=="vercel")'

{
  "name": "vercel",
  "path": "/opt/homebrew/bin/vercel",
  "package_name": "vercel-cli",
  "package_manager": "brew",
  "package_version": "41.6.2",
  "is_symlink": true,
  "symlink_to": "../Cellar/vercel-cli/41.6.2/bin/vercel",
  "size": 38
}
```

## Benefits for AI Agents

1. **Dependency Management**: Know which package to install to get a specific CLI
2. **Version Tracking**: See which version of a package provides a CLI
3. **Package Manager Context**: Understand how tools were installed
4. **Multiple Binaries**: Discover all CLIs provided by a single package
5. **Installation Source**: Differentiate between system tools and package-installed tools

## Performance

- Package detection adds ~2-3 seconds to the export command
- Uses existing package manager commands (npm, pip, brew, etc.)
- Caches package list for linking all tools
- No performance impact if `--with-packages` flag is not used

## Future Enhancements

- yarn/pnpm support
- apt/dnf/pacman support for Linux
- Package dependency graphs
- Installation command suggestions
- Package update notifications
