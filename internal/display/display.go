package display

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/cli-ai-org/cli/internal/models"
)

// Display handles the output formatting for CLI tools
type Display struct {
	writer io.Writer
}

// New creates a new Display instance
func New(w io.Writer) *Display {
	return &Display{writer: w}
}

// ShowTools displays a list of tools
func (d *Display) ShowTools(tools []string) {
	if len(tools) == 0 {
		fmt.Fprintln(d.writer, "No CLI tools found.")
		return
	}

	// Sort alphabetically
	sorted := make([]string, len(tools))
	copy(sorted, tools)
	sort.Strings(sorted)

	fmt.Fprintf(d.writer, "Found %d CLI tools:\n\n", len(sorted))
	for _, tool := range sorted {
		fmt.Fprintf(d.writer, "  %s\n", tool)
	}
}

// ShowToolsVerbose displays tools with additional information
func (d *Display) ShowToolsVerbose(tools []string, paths map[string]string) {
	if len(tools) == 0 {
		fmt.Fprintln(d.writer, "No CLI tools found.")
		return
	}

	// Sort alphabetically
	sorted := make([]string, len(tools))
	copy(sorted, tools)
	sort.Strings(sorted)

	fmt.Fprintf(d.writer, "Found %d CLI tools:\n\n", len(sorted))
	for _, tool := range sorted {
		path := paths[tool]
		if path != "" {
			fmt.Fprintf(d.writer, "  %-30s %s\n", tool, path)
		} else {
			fmt.Fprintf(d.writer, "  %s\n", tool)
		}
	}
}

// ShowToolsDetailed displays detailed tool information
func (d *Display) ShowToolsDetailed(tools []models.Tool) {
	if len(tools) == 0 {
		fmt.Fprintln(d.writer, "No CLI tools found.")
		return
	}

	// Sort alphabetically by name
	sorted := make([]models.Tool, len(tools))
	copy(sorted, tools)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})

	fmt.Fprintf(d.writer, "Found %d CLI tools:\n\n", len(sorted))
	for _, tool := range sorted {
		fmt.Fprintf(d.writer, "  %-30s %s", tool.Name, tool.Path)
		if tool.IsSymlink {
			fmt.Fprintf(d.writer, " -> %s", tool.SymlinkTo)
		}
		fmt.Fprintln(d.writer)
	}
}

// ShowToolsJSON outputs tools in JSON format for AI agents
func (d *Display) ShowToolsJSON(tools []models.Tool, pretty bool) error {
	if pretty {
		encoder := json.NewEncoder(d.writer)
		encoder.SetIndent("", "  ")
		return encoder.Encode(tools)
	}

	encoder := json.NewEncoder(d.writer)
	return encoder.Encode(tools)
}

// ShowCatalogJSON outputs a complete tool catalog in JSON format
func (d *Display) ShowCatalogJSON(catalog *models.ToolCatalog, pretty bool) error {
	if pretty {
		encoder := json.NewEncoder(d.writer)
		encoder.SetIndent("", "  ")
		return encoder.Encode(catalog)
	}

	encoder := json.NewEncoder(d.writer)
	return encoder.Encode(catalog)
}

// ShowToolInfo displays detailed information about a single tool
func (d *Display) ShowToolInfo(tool *models.Tool, detailed bool) {
	fmt.Fprintf(d.writer, "Tool: %s\n", tool.Name)
	fmt.Fprintf(d.writer, "Path: %s\n", tool.Path)

	if tool.IsSymlink {
		fmt.Fprintf(d.writer, "Symlink: -> %s\n", tool.SymlinkTo)
	}

	if tool.Size > 0 {
		fmt.Fprintf(d.writer, "Size: %d bytes\n", tool.Size)
	}

	if tool.Version != "" {
		fmt.Fprintf(d.writer, "Version: %s\n", tool.Version)
	}

	if detailed && tool.Description != "" {
		fmt.Fprintf(d.writer, "Description: %s\n", tool.Description)
	}

	if detailed && tool.HelpText != "" {
		fmt.Fprintf(d.writer, "\nHelp Text:\n%s\n", tool.HelpText)
	}
}
