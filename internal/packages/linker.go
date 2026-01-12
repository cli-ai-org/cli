package packages

import (
	"path/filepath"
	"strings"

	"github.com/cli-ai-org/cli/internal/models"
)

// Linker links CLI tools to their source packages
type Linker struct {
	packages map[string]Package
}

// NewLinker creates a new package linker
func NewLinker(packages []Package) *Linker {
	pkgMap := make(map[string]Package)
	for _, pkg := range packages {
		pkgMap[pkg.Name] = pkg
	}
	return &Linker{packages: pkgMap}
}

// LinkTools links tools to their source packages using various heuristics
func (l *Linker) LinkTools(tools []models.Tool) []models.Tool {
	enriched := make([]models.Tool, len(tools))
	copy(enriched, tools)

	for i := range enriched {
		l.linkTool(&enriched[i])
	}

	return enriched
}

// linkTool attempts to link a single tool to its package
func (l *Linker) linkTool(tool *models.Tool) {
	// Strategy 1: Direct name match (e.g., "vercel" package -> "vercel" cli)
	if pkg, ok := l.packages[tool.Name]; ok {
		tool.PackageName = pkg.Name
		tool.PackageManager = string(pkg.Manager)
		tool.PackageVersion = pkg.Version
		return
	}

	// Strategy 2: Path-based detection
	l.detectFromPath(tool)

	// Strategy 3: Common patterns (e.g., @supabase/cli -> supabase)
	if tool.PackageName == "" {
		l.detectFromPatterns(tool)
	}
}

// detectFromPath attempts to detect package from the tool's path
func (l *Linker) detectFromPath(tool *models.Tool) {
	// Check both the path and symlink target
	paths := []string{tool.Path}
	if tool.IsSymlink && tool.SymlinkTo != "" {
		paths = append(paths, tool.SymlinkTo)
	}

	for _, path := range paths {
		if l.checkPath(tool, path) {
			return
		}
	}
}

// checkPath checks a single path for package information
func (l *Linker) checkPath(tool *models.Tool, path string) bool {
	// NPM global modules (.nvm, node_modules)
	if strings.Contains(path, "node_modules") {
		parts := strings.Split(path, "node_modules")
		if len(parts) > 1 {
			// Extract package name from path after node_modules
			remaining := strings.Trim(parts[1], "/")
			pkgParts := strings.Split(remaining, "/")
			if len(pkgParts) > 0 {
				pkgName := pkgParts[0]
				// Handle scoped packages (@org/package)
				if strings.HasPrefix(pkgName, "@") && len(pkgParts) > 1 {
					pkgName = pkgName + "/" + pkgParts[1]
				}
				if pkg, ok := l.packages[pkgName]; ok {
					tool.PackageName = pkg.Name
					tool.PackageManager = string(pkg.Manager)
					tool.PackageVersion = pkg.Version
					return true
				}
			}
		}
	}

	// Homebrew packages
	if strings.Contains(path, "/opt/homebrew/") || strings.Contains(path, "/usr/local/Cellar/") || strings.Contains(path, "Cellar/") {
		// Extract from /opt/homebrew/Cellar/package/version/bin/tool or ../Cellar/package/version/bin/tool
		if strings.Contains(path, "Cellar/") {
			parts := strings.Split(path, "Cellar/")
			if len(parts) > 1 {
				remaining := parts[1]
				pkgName := strings.Split(remaining, "/")[0]
				if pkg, ok := l.packages[pkgName]; ok {
					tool.PackageName = pkg.Name
					tool.PackageManager = string(pkg.Manager)
					tool.PackageVersion = pkg.Version
					return true
				}
			}
		}

		// Try extracting from /opt/homebrew/opt/package
		if strings.Contains(path, "/opt/") {
			parts := strings.Split(path, "/opt/")
			if len(parts) > 1 {
				remaining := parts[1]
				pkgName := strings.Split(remaining, "/")[0]
				if pkg, ok := l.packages[pkgName]; ok {
					tool.PackageName = pkg.Name
					tool.PackageManager = string(pkg.Manager)
					tool.PackageVersion = pkg.Version
					return true
				}
			}
		}
	}

	// Python packages (.pyenv, site-packages)
	if strings.Contains(path, "site-packages") || strings.Contains(path, ".pyenv") {
		// Python CLIs are harder to detect, skip for now
		return false
	}

	// Cargo packages (.cargo/bin)
	if strings.Contains(path, ".cargo/bin") {
		toolName := filepath.Base(path)
		if pkg, ok := l.packages[toolName]; ok && pkg.Manager == Cargo {
			tool.PackageName = pkg.Name
			tool.PackageManager = string(pkg.Manager)
			tool.PackageVersion = pkg.Version
			return true
		}
	}

	return false
}

// detectFromPatterns uses common naming patterns to detect packages
func (l *Linker) detectFromPatterns(tool *models.Tool) {
	name := tool.Name

	// Common patterns:
	// - @scope/cli -> scope or cli
	// - package-cli -> package
	// - cli-package -> package

	// Try removing common suffixes/prefixes
	patterns := []string{
		strings.TrimSuffix(name, "-cli"),
		strings.TrimPrefix(name, "cli-"),
		strings.TrimSuffix(name, "cli"),
	}

	for _, pattern := range patterns {
		if pattern != name {
			if pkg, ok := l.packages[pattern]; ok {
				tool.PackageName = pkg.Name
				tool.PackageManager = string(pkg.Manager)
				tool.PackageVersion = pkg.Version
				return
			}
		}
	}

	// Handle scoped packages
	if strings.HasPrefix(name, "@") {
		parts := strings.Split(name, "/")
		if len(parts) == 2 {
			// Try @scope/package
			if pkg, ok := l.packages[name]; ok {
				tool.PackageName = pkg.Name
				tool.PackageManager = string(pkg.Manager)
				tool.PackageVersion = pkg.Version
				return
			}
		}
	}
}

// GetPackagesWithBinaries enriches packages with their binary information
func GetPackagesWithBinaries(packages []Package, tools []models.Tool) []models.PackageInfo {
	pkgBinaries := make(map[string][]string)

	for _, tool := range tools {
		if tool.PackageName != "" {
			pkgBinaries[tool.PackageName] = append(pkgBinaries[tool.PackageName], tool.Name)
		}
	}

	var result []models.PackageInfo
	for _, pkg := range packages {
		binaries := pkgBinaries[pkg.Name]
		if len(binaries) > 0 {
			result = append(result, models.PackageInfo{
				Name:     pkg.Name,
				Version:  pkg.Version,
				Manager:  string(pkg.Manager),
				Binaries: binaries,
				Location: pkg.Location,
				Global:   pkg.Global,
			})
		}
	}

	return result
}
