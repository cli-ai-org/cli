package collector

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/cli-ai-org/cli/internal/models"
)

// Collector gathers detailed information about CLI tools
type Collector struct {
	timeoutSeconds int
}

// New creates a new Collector instance
func New() *Collector {
	return &Collector{
		timeoutSeconds: 3,
	}
}

// CollectToolInfo gathers detailed information about a specific tool
func (c *Collector) CollectToolInfo(toolName string, toolPath string) (*models.Tool, error) {
	tool := &models.Tool{
		Name: toolName,
		Path: toolPath,
	}

	// Get file info
	info, err := os.Lstat(toolPath)
	if err != nil {
		return tool, nil // Return basic info even if we can't stat
	}

	tool.Size = info.Size()

	// Check if symlink
	if info.Mode()&os.ModeSymlink != 0 {
		tool.IsSymlink = true
		if target, err := os.Readlink(toolPath); err == nil {
			tool.SymlinkTo = target
		}
	}

	// Try to get version
	tool.Version = c.getVersion(toolPath)

	// Try to get help text
	tool.HelpText = c.getHelpText(toolPath)

	return tool, nil
}

// getVersion attempts to extract version information from a tool
func (c *Collector) getVersion(toolPath string) string {
	versionFlags := []string{"--version", "-version", "version", "-v"}

	for _, flag := range versionFlags {
		cmd := exec.Command(toolPath, flag)
		output, err := cmd.CombinedOutput()
		if err == nil && len(output) > 0 {
			// Take first line of version output
			lines := strings.Split(string(output), "\n")
			if len(lines) > 0 && len(lines[0]) > 0 && len(lines[0]) < 200 {
				return strings.TrimSpace(lines[0])
			}
		}
	}

	return ""
}

// getHelpText attempts to extract help information from a tool
func (c *Collector) getHelpText(toolPath string) string {
	helpFlags := []string{"--help", "-help", "help", "-h"}

	for _, flag := range helpFlags {
		cmd := exec.Command(toolPath, flag)
		output, err := cmd.CombinedOutput()
		if err == nil && len(output) > 0 {
			// Limit help text size
			helpText := string(output)
			if len(helpText) > 5000 {
				helpText = helpText[:5000] + "\n... (truncated)"
			}
			return helpText
		}
	}

	return ""
}

// BuildCatalog creates a comprehensive catalog of all tools
func (c *Collector) BuildCatalog(tools []models.Tool, searchPaths []string) *models.ToolCatalog {
	return &models.ToolCatalog{
		TotalTools:  len(tools),
		Paths:       searchPaths,
		Tools:       tools,
		GeneratedAt: time.Now().Format(time.RFC3339),
	}
}

// GetToolByName finds a specific tool by name from a list
func GetToolByName(tools []models.Tool, name string) *models.Tool {
	for _, tool := range tools {
		if tool.Name == name {
			return &tool
		}
	}
	return nil
}

// FilterTools filters tools based on a predicate function
func FilterTools(tools []models.Tool, predicate func(models.Tool) bool) []models.Tool {
	filtered := []models.Tool{}
	for _, tool := range tools {
		if predicate(tool) {
			filtered = append(filtered, tool)
		}
	}
	return filtered
}

// GetToolPath finds the full path of a tool by name
func GetToolPath(toolName string) (string, error) {
	return exec.LookPath(toolName)
}

// ParseManPage attempts to extract information from a man page
func (c *Collector) ParseManPage(toolName string) string {
	cmd := exec.Command("man", toolName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	// Limit man page size
	manText := string(output)
	if len(manText) > 10000 {
		manText = manText[:10000] + "\n... (truncated)"
	}

	return manText
}
