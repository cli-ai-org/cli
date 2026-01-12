package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/cli-ai-org/cli/internal/display"
	"github.com/cli-ai-org/cli/internal/models"
	"github.com/cli-ai-org/cli/internal/packages"
	"github.com/cli-ai-org/cli/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	debugAll     bool
	debugClashes bool
)

// debugCmd represents the debug command
var debugCmd = &cobra.Command{
	Use:   "debug [tool_name]",
	Short: "Debug CLI tool installations and detect clashes",
	Long: `Display detailed debug information about CLI tools and detect installation conflicts.

This command helps you identify when the same tool is installed by multiple package
managers (brew, pip, npm, etc.) and shows which installation is active in your PATH.

Modes:
  - debug TOOL_NAME: Show all installations of a specific tool
  - debug --clashes: Show all tools with conflicting installations
  - debug --all: Show debug info for all tools`,
	Example: `  # Debug a specific tool
  cli-ai debug python
  cli-ai debug docker

  # Show all installation clashes
  cli-ai debug --clashes

  # Debug all tools
  cli-ai debug --all`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s := scanner.New()
		d := display.New(os.Stdout)

		// Scan all tools
		tools, err := s.ScanAllDetailed()
		if err != nil {
			cmd.PrintErrf("Error scanning tools: %v\n", err)
			os.Exit(1)
		}

		// Detect packages
		detector := packages.NewDetector()
		pkgs, err := detector.DetectAll()
		if err != nil {
			cmd.PrintErrf("Error detecting packages: %v\n", err)
			os.Exit(1)
		}

		// Link tools to packages
		linker := packages.NewLinker(pkgs)
		tools = linker.LinkTools(tools)

		if debugClashes {
			showClashes(tools, d)
		} else if debugAll {
			showAllDebug(tools, d)
		} else if len(args) == 0 {
			cmd.PrintErr("Error: must specify a tool name or use --clashes or --all flag\n\n")
			cmd.Usage()
			os.Exit(1)
		} else {
			showToolDebug(args[0], tools, d)
		}
	},
}

func showClashes(tools []models.Tool, d *display.Display) {
	// Group tools by name
	toolGroups := make(map[string][]models.Tool)
	for _, tool := range tools {
		if tool.PackageName != "" {
			toolGroups[tool.Name] = append(toolGroups[tool.Name], tool)
		}
	}

	// Find clashes (tools installed by multiple packages)
	var clashes []string
	for name, instances := range toolGroups {
		// Check if multiple different packages provide this tool
		packageSeen := make(map[string]bool)
		for _, instance := range instances {
			packageSeen[instance.PackageName] = true
		}
		if len(packageSeen) > 1 {
			clashes = append(clashes, name)
		}
	}

	if len(clashes) == 0 {
		fmt.Fprintln(os.Stdout, "No installation clashes found!")
		return
	}

	sort.Strings(clashes)

	fmt.Fprintf(os.Stdout, "Found %d tools with multiple installations:\n\n", len(clashes))

	for _, name := range clashes {
		instances := toolGroups[name]
		fmt.Fprintf(os.Stdout, "🔴 %s (%d installations)\n", name, len(instances))

		// Sort by PATH order (first is active)
		for i, instance := range instances {
			active := ""
			if i == 0 {
				active = " ✓ ACTIVE"
			}
			fmt.Fprintf(os.Stdout, "   %s via %s%s\n", instance.Path, instance.PackageManager, active)
			if instance.PackageVersion != "" {
				fmt.Fprintf(os.Stdout, "      Version: %s\n", instance.PackageVersion)
			}
		}
		fmt.Fprintln(os.Stdout)
	}
}

func showToolDebug(toolName string, tools []models.Tool, d *display.Display) {
	var matches []models.Tool
	for _, tool := range tools {
		if tool.Name == toolName {
			matches = append(matches, tool)
		}
	}

	if len(matches) == 0 {
		fmt.Fprintf(os.Stdout, "Tool '%s' not found in PATH\n", toolName)
		return
	}

	fmt.Fprintf(os.Stdout, "Debug information for: %s\n", toolName)
	fmt.Fprintf(os.Stdout, "Total installations: %d\n\n", len(matches))

	for i, tool := range matches {
		fmt.Fprintf(os.Stdout, "Installation #%d:\n", i+1)
		if i == 0 {
			fmt.Fprintln(os.Stdout, "  Status: ✓ ACTIVE (first in PATH)")
		} else {
			fmt.Fprintln(os.Stdout, "  Status: ⚠ SHADOWED (not used)")
		}
		fmt.Fprintf(os.Stdout, "  Path: %s\n", tool.Path)

		if tool.IsSymlink {
			fmt.Fprintf(os.Stdout, "  Symlink to: %s\n", tool.SymlinkTo)
		}

		if tool.PackageName != "" {
			fmt.Fprintf(os.Stdout, "  Package: %s\n", tool.PackageName)
			fmt.Fprintf(os.Stdout, "  Manager: %s\n", tool.PackageManager)
			if tool.PackageVersion != "" {
				fmt.Fprintf(os.Stdout, "  Version: %s\n", tool.PackageVersion)
			}
		} else {
			fmt.Fprintln(os.Stdout, "  Package: (not detected)")
		}

		if tool.Size > 0 {
			fmt.Fprintf(os.Stdout, "  Size: %d bytes\n", tool.Size)
		}

		fmt.Fprintln(os.Stdout)
	}

	// Show recommendation if multiple installations
	if len(matches) > 1 {
		fmt.Fprintln(os.Stdout, "⚠️  RECOMMENDATION:")
		fmt.Fprintln(os.Stdout, "Multiple installations detected. Consider:")
		fmt.Fprintf(os.Stdout, "  - Using the active installation via %s\n", matches[0].PackageManager)
		fmt.Fprintln(os.Stdout, "  - Uninstalling unused versions to avoid conflicts")
	}
}

func showAllDebug(tools []models.Tool, d *display.Display) {
	// Group by package
	packageTools := make(map[string][]models.Tool)
	for _, tool := range tools {
		if tool.PackageName != "" {
			key := fmt.Sprintf("%s:%s", tool.PackageManager, tool.PackageName)
			packageTools[key] = append(packageTools[key], tool)
		}
	}

	var keys []string
	for k := range packageTools {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(os.Stdout, "Showing debug info for %d packages:\n\n", len(keys))

	for _, key := range keys {
		parts := strings.Split(key, ":")
		manager := parts[0]
		pkg := parts[1]
		instances := packageTools[key]

		fmt.Fprintf(os.Stdout, "📦 %s (via %s)\n", pkg, manager)
		fmt.Fprintf(os.Stdout, "   Provides %d tool(s):", len(instances))

		var toolNames []string
		for _, t := range instances {
			toolNames = append(toolNames, t.Name)
		}
		fmt.Fprintf(os.Stdout, " %s\n", strings.Join(toolNames, ", "))

		if instances[0].PackageVersion != "" {
			fmt.Fprintf(os.Stdout, "   Version: %s\n", instances[0].PackageVersion)
		}
		fmt.Fprintln(os.Stdout)
	}
}

func init() {
	rootCmd.AddCommand(debugCmd)
	debugCmd.Flags().BoolVarP(&debugAll, "all", "a", false, "show debug information for all packages")
	debugCmd.Flags().BoolVarP(&debugClashes, "clashes", "c", false, "show only tools with conflicting installations")
}
