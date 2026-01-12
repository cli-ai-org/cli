package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/cli-ai-org/cli/internal/models"
)

// Scanner handles the discovery of CLI tools on the system
type Scanner struct {
	paths []string
}

// New creates a new Scanner instance
func New() *Scanner {
	return &Scanner{
		paths: getPathDirectories(),
	}
}

// getPathDirectories returns all directories in the system PATH
func getPathDirectories() []string {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return []string{}
	}
	return strings.Split(pathEnv, string(os.PathListSeparator))
}

// ScanAll scans all PATH directories for CLI tools
func (s *Scanner) ScanAll() ([]string, error) {
	var tools []string
	seen := make(map[string]bool)

	for _, dir := range s.paths {
		entries, err := os.ReadDir(dir)
		if err != nil {
			// Skip directories we can't read
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			// Check if file is executable
			info, err := entry.Info()
			if err != nil {
				continue
			}

			if isExecutable(info) {
				name := entry.Name()
				if !seen[name] {
					seen[name] = true
					tools = append(tools, name)
				}
			}
		}
	}

	return tools, nil
}

// isExecutable checks if a file has executable permissions
func isExecutable(info os.FileInfo) bool {
	mode := info.Mode()
	return mode&0111 != 0
}

// GetPaths returns the list of PATH directories
func (s *Scanner) GetPaths() []string {
	return s.paths
}

// ScanAllDetailed scans all PATH directories and returns detailed Tool information
func (s *Scanner) ScanAllDetailed() ([]models.Tool, error) {
	var tools []models.Tool
	seen := make(map[string]bool)

	for _, dir := range s.paths {
		entries, err := os.ReadDir(dir)
		if err != nil {
			// Skip directories we can't read
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			// Check if file is executable
			info, err := entry.Info()
			if err != nil {
				continue
			}

			if isExecutable(info) {
				name := entry.Name()
				if !seen[name] {
					seen[name] = true
					fullPath := filepath.Join(dir, name)

					tool := models.Tool{
						Name: name,
						Path: fullPath,
						Size: info.Size(),
					}

					// Check if symlink
					if info.Mode()&os.ModeSymlink != 0 {
						tool.IsSymlink = true
						if target, err := os.Readlink(fullPath); err == nil {
							tool.SymlinkTo = target
						}
					}

					tools = append(tools, tool)
				}
			}
		}
	}

	return tools, nil
}

// FindTool finds a specific tool by name and returns detailed information
func (s *Scanner) FindTool(name string) (*models.Tool, error) {
	for _, dir := range s.paths {
		fullPath := filepath.Join(dir, name)
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}

		if isExecutable(info) {
			tool := &models.Tool{
				Name: name,
				Path: fullPath,
				Size: info.Size(),
			}

			// Check if symlink
			linkInfo, err := os.Lstat(fullPath)
			if err == nil && linkInfo.Mode()&os.ModeSymlink != 0 {
				tool.IsSymlink = true
				if target, err := os.Readlink(fullPath); err == nil {
					tool.SymlinkTo = target
				}
			}

			return tool, nil
		}
	}

	return nil, os.ErrNotExist
}
