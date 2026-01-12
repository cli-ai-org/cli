package models

// Tool represents a CLI tool discovered on the system
type Tool struct {
	Name           string   `json:"name"`
	Path           string   `json:"path"`
	Description    string   `json:"description,omitempty"`
	Version        string   `json:"version,omitempty"`
	HelpText       string   `json:"help_text,omitempty"`
	IsSymlink      bool     `json:"is_symlink"`
	SymlinkTo      string   `json:"symlink_to,omitempty"`
	Size           int64    `json:"size"`
	Aliases        []string `json:"aliases,omitempty"`
	PackageName    string   `json:"package_name,omitempty"`
	PackageManager string   `json:"package_manager,omitempty"`
	PackageVersion string   `json:"package_version,omitempty"`
}

// ToolCatalog represents a collection of tools for AI agent consumption
type ToolCatalog struct {
	TotalTools    int              `json:"total_tools"`
	TotalPackages int              `json:"total_packages,omitempty"`
	Paths         []string         `json:"search_paths"`
	Tools         []Tool           `json:"tools"`
	Packages      []PackageInfo    `json:"packages,omitempty"`
	GeneratedAt   string           `json:"generated_at"`
}

// PackageInfo represents a package that provides CLI tools
type PackageInfo struct {
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Manager  string   `json:"manager"`
	Binaries []string `json:"binaries,omitempty"`
	Location string   `json:"location,omitempty"`
	Global   bool     `json:"global"`
}

// ToolInfo provides structured information about a tool for AI agents
type ToolInfo struct {
	Name         string            `json:"name"`
	Location     string            `json:"location"`
	Version      string            `json:"version,omitempty"`
	Description  string            `json:"description,omitempty"`
	Usage        string            `json:"usage,omitempty"`
	CommonFlags  []Flag            `json:"common_flags,omitempty"`
	Examples     []string          `json:"examples,omitempty"`
	Dependencies []string          `json:"dependencies,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// Flag represents a command-line flag
type Flag struct {
	Name        string `json:"name"`
	Short       string `json:"short,omitempty"`
	Description string `json:"description"`
	Default     string `json:"default,omitempty"`
}
