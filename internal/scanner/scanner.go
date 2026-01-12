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

			name := entry.Name()

			// Filter out non-CLI tools
			if !shouldIncludeTool(name) {
				continue
			}

			// Check if file is executable
			info, err := entry.Info()
			if err != nil {
				continue
			}

			if isExecutable(info) {
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

// shouldIncludeTool filters out system daemons, test utilities, and internal tools
func shouldIncludeTool(name string) bool {
	lower := strings.ToLower(name)

	// Skip Python cache and obvious non-tools
	if name == "__pycache__" || name == "." || name == ".." {
		return false
	}

	// Skip DTrace scripts (end with .d)
	if strings.HasSuffix(name, ".d") {
		return false
	}

	// Skip obvious test utilities and demos
	excludePatterns := []string{
		"test", "demo", "bench", "example", "sample",
		"_test", "_demo", "_bench", "_example",
	}
	for _, pattern := range excludePatterns {
		if strings.Contains(lower, pattern) {
			return false
		}
	}

	// Skip server/daemon/agent patterns
	if strings.HasSuffix(lower, "server") || strings.HasSuffix(lower, "agent") ||
	   strings.HasSuffix(lower, "daemon") || strings.HasSuffix(lower, "serverd") {
		// Allow some legitimate tools
		allowed := []string{"transmission-daemon", "jupyter-server"}
		isAllowed := false
		for _, allow := range allowed {
			if lower == allow {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			return false
		}
	}

	// Skip common daemon patterns (but allow some legitimate tools)
	daemonExclusions := []string{
		"bluetoothd", "coreaudiod", "cfprefsd", "distnoted",
		"launchd", "notifyd", "securityd", "syslogd", "configd",
		"kerneleventd", "powerd", "cupsd", "httpd", "sshd",
		"snmpd", "named", "ntpd", "syslogd",
		"btleserver", "btleserveragent",
	}
	for _, daemon := range daemonExclusions {
		if lower == daemon {
			return false
		}
	}

	// Skip Apple internal tools (specific patterns)
	appleInternalPrefixes := []string{
		"appleh", "assetcache", "bluetool", "bootcache",
		"createdom", "domcount", "domprint", "derez",
		"devtools", "directory", "enumval", "getfileinfo",
		"ioaccel", "iomfb", "iosdebug", "kernel",
		"pparse", "psviwriter", "password", "protocol",
		"redirect", "resmerger", "rez", "sax", "scmprint",
		"senumval", "safeeject", "setfile", "splitforks",
		"stdin", "svtav1", "wireless", "xinclude",
	}
	for _, prefix := range appleInternalPrefixes {
		if strings.HasPrefix(lower, prefix) {
			return false
		}
	}

	// Skip more system internals
	systemInternals := []string{
		"mDNSResponder", "mDNSResponderHelper",
	}
	for _, internal := range systemInternals {
		if name == internal {
			return false
		}
	}

	return true
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

			name := entry.Name()

			// Filter out non-CLI tools
			if !shouldIncludeTool(name) {
				continue
			}

			// Check if file is executable
			info, err := entry.Info()
			if err != nil {
				continue
			}

			if isExecutable(info) {
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
