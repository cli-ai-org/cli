# Quick Usage Examples

Fast reference for common cli commands.

## Help & Documentation

```bash
cli help                    # Main help page
cli list --help             # Help for list command
cli debug --help            # Help for debug command
```

## List Commands

```bash
# Basic listing
cli list                    # Simple list of all CLI tools

# Detailed listing
cli list --all              # With full paths and details
cli list -a                 # Short flag version

# Verbose output
cli list --verbose          # Verbose output
cli list -v                 # Short flag version

# Combined flags
cli list --all --verbose    # Maximum information
cli list -a -v              # Short version
```

## Debug Commands

```bash
# Debug specific packages
cli debug npm               # Debug npm package
cli debug python            # Debug Python
cli debug docker            # Debug Docker
cli debug node              # Debug Node.js

# Debug all packages
cli debug --all             # All packages
cli debug -a                # Short flag version

# Debug with verbose
cli debug npm --verbose     # Verbose debug for npm
cli debug --all --verbose   # Verbose debug for all
cli debug npm -v            # Short flag version
cli debug -a -v             # All packages, verbose, short flags
```

## Real-world Examples

```bash
# Find all tools on your system
cli list

# Get detailed info about all tools including paths
cli list --all

# Debug your Node.js installation
cli debug node

# Debug your Python installation
cli debug python

# See everything cli can find (maximum detail)
cli list --all --verbose
cli debug --all --verbose

# Check if a specific tool exists
cli list | grep docker

# Count how many CLI tools you have
cli list | wc -l
```

## Flag Combinations

```bash
# Most common combinations
cli list -a                 # Detailed list
cli list -v                 # Verbose list
cli list -a -v              # Detailed + verbose list
cli debug npm -v            # Verbose debug for npm
cli debug -a                # Debug all packages
cli debug -a -v             # Debug all with verbose
```

## Output Examples

### Simple List
```
$ cli list
Found 142 CLI tools:

  awk
  bash
  cat
  docker
  git
  ...
```

### Detailed List (--all)
```
$ cli list --all
Found 142 CLI tools:

  awk                            /usr/bin/awk
  bash                           /bin/bash
  cat                            /bin/cat
  docker                         /usr/local/bin/docker
  git                            /usr/bin/git
  ...
```

### Debug Specific Package
```
$ cli debug npm
Debugging package: npm

Package: npm
Location: /usr/local/bin/npm
Version: 10.2.3
Dependencies: node, npx
...
```

### Debug All Packages
```
$ cli debug --all
Debugging all packages...

Package: npm
Location: /usr/local/bin/npm
...

Package: python
Location: /usr/bin/python3
...
```
