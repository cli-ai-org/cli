package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/cli-ai-org/cli/internal/models"
	"github.com/cli-ai-org/cli/internal/packages"
	"github.com/cli-ai-org/cli/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	auditOutput string
)

// auditCmd represents the audit command
var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit CLI environment and generate recommendations",
	Long: `Perform a comprehensive audit of your CLI environment and generate a report.

This command analyzes:
  - Installation clashes (tools from multiple package managers)
  - Shadowed installations (tools not being used)
  - Package manager coverage
  - System health recommendations

The audit generates a markdown report suitable for AI agents to analyze.`,
	Example: `  # Run audit and display to console
  cli-ai audit

  # Save audit report to file
  cli-ai audit --output cli-audit.md

  # Save with custom name
  cli-ai audit -o my-system-audit.md`,
	Run: func(cmd *cobra.Command, args []string) {
		s := scanner.New()

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

		// Perform audit
		report := performAudit(tools, pkgs)

		// Output report
		if auditOutput != "" {
			err := os.WriteFile(auditOutput, []byte(report), 0644)
			if err != nil {
				cmd.PrintErrf("Error writing audit report: %v\n", err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stdout, "✓ Audit report saved to: %s\n", auditOutput)
		} else {
			fmt.Fprint(os.Stdout, report)
		}
	},
}

type AuditResult struct {
	TotalTools        int
	PackageManagedTools int
	UnmanagedTools    int
	Clashes           []ToolClash
	ShadowedTools     []ShadowedTool
	PackageManagers   []PackageManagerInfo
	Recommendations   []Recommendation
}

type ToolClash struct {
	ToolName      string
	Installations []InstallationInfo
}

type InstallationInfo struct {
	Path           string
	PackageName    string
	PackageManager string
	Version        string
	IsActive       bool
}

type ShadowedTool struct {
	ToolName       string
	ActivePath     string
	ShadowedPath   string
	ActivePackage  string
	ShadowedPackage string
}

type PackageManagerInfo struct {
	Name         string
	PackageCount int
	ToolCount    int
}

type Recommendation struct {
	Severity string // "high", "medium", "low"
	Category string
	Issue    string
	Action   string
}

func performAudit(tools []models.Tool, pkgs []packages.Package) string {
	result := AuditResult{}

	// Count tools
	result.TotalTools = len(tools)
	for _, tool := range tools {
		if tool.PackageName != "" {
			result.PackageManagedTools++
		} else {
			result.UnmanagedTools++
		}
	}

	// Find clashes
	result.Clashes = findClashes(tools)

	// Find shadowed tools
	result.ShadowedTools = findShadowedTools(tools)

	// Analyze package managers
	result.PackageManagers = analyzePackageManagers(pkgs, tools)

	// Generate recommendations
	result.Recommendations = generateRecommendations(result, tools, pkgs)

	// Generate markdown report
	return generateMarkdownReport(result)
}

func findClashes(tools []models.Tool) []ToolClash {
	toolGroups := make(map[string][]models.Tool)
	for _, tool := range tools {
		if tool.PackageName != "" {
			toolGroups[tool.Name] = append(toolGroups[tool.Name], tool)
		}
	}

	var clashes []ToolClash
	for name, instances := range toolGroups {
		packageSeen := make(map[string]bool)
		for _, instance := range instances {
			packageSeen[instance.PackageName] = true
		}

		if len(packageSeen) > 1 {
			clash := ToolClash{ToolName: name}
			for i, instance := range instances {
				clash.Installations = append(clash.Installations, InstallationInfo{
					Path:           instance.Path,
					PackageName:    instance.PackageName,
					PackageManager: instance.PackageManager,
					Version:        instance.PackageVersion,
					IsActive:       i == 0,
				})
			}
			clashes = append(clashes, clash)
		}
	}

	return clashes
}

func findShadowedTools(tools []models.Tool) []ShadowedTool {
	toolGroups := make(map[string][]models.Tool)
	for _, tool := range tools {
		toolGroups[tool.Name] = append(toolGroups[tool.Name], tool)
	}

	var shadowed []ShadowedTool
	for name, instances := range toolGroups {
		if len(instances) > 1 {
			for i := 1; i < len(instances); i++ {
				shadowed = append(shadowed, ShadowedTool{
					ToolName:        name,
					ActivePath:      instances[0].Path,
					ShadowedPath:    instances[i].Path,
					ActivePackage:   instances[0].PackageName,
					ShadowedPackage: instances[i].PackageName,
				})
			}
		}
	}

	return shadowed
}

func analyzePackageManagers(pkgs []packages.Package, tools []models.Tool) []PackageManagerInfo {
	managerStats := make(map[string]*PackageManagerInfo)

	for _, pkg := range pkgs {
		manager := string(pkg.Manager)
		if _, exists := managerStats[manager]; !exists {
			managerStats[manager] = &PackageManagerInfo{Name: manager}
		}
		managerStats[manager].PackageCount++
	}

	// Count tools per manager
	for _, tool := range tools {
		if tool.PackageManager != "" {
			if info, exists := managerStats[tool.PackageManager]; exists {
				info.ToolCount++
			}
		}
	}

	var result []PackageManagerInfo
	for _, info := range managerStats {
		result = append(result, *info)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ToolCount > result[j].ToolCount
	})

	return result
}

func generateRecommendations(result AuditResult, tools []models.Tool, pkgs []packages.Package) []Recommendation {
	var recs []Recommendation

	// Check for clashes
	if len(result.Clashes) > 0 {
		recs = append(recs, Recommendation{
			Severity: "high",
			Category: "Installation Conflicts",
			Issue:    fmt.Sprintf("Found %d tools with multiple installations from different package managers", len(result.Clashes)),
			Action:   "Review conflicting installations and uninstall duplicates to avoid version conflicts. Use `cli-ai debug --clashes` for details.",
		})
	}

	// Check for shadowed tools
	if len(result.ShadowedTools) > 0 {
		recs = append(recs, Recommendation{
			Severity: "medium",
			Category: "Shadowed Installations",
			Issue:    fmt.Sprintf("Found %d tools with shadowed installations that are not being used", len(result.ShadowedTools)),
			Action:   "Remove unused installations to free up disk space and reduce confusion. The shadowed installations are not in use.",
		})
	}

	// Check for unmanaged tools
	unmanagedPercent := float64(result.UnmanagedTools) / float64(result.TotalTools) * 100
	if unmanagedPercent > 20 {
		recs = append(recs, Recommendation{
			Severity: "low",
			Category: "Package Management",
			Issue:    fmt.Sprintf("%.1f%% of tools (%d/%d) are not managed by a package manager", unmanagedPercent, result.UnmanagedTools, result.TotalTools),
			Action:   "Consider installing tools via package managers (brew, npm, pip) for easier updates and management.",
		})
	}

	// Check package manager diversity
	if len(result.PackageManagers) == 1 {
		recs = append(recs, Recommendation{
			Severity: "low",
			Category: "Package Management",
			Issue:    "Only using one package manager on your system",
			Action:   "This is good for consistency! Continue managing all tools through " + result.PackageManagers[0].Name + ".",
		})
	}

	// If no issues found
	if len(recs) == 0 {
		recs = append(recs, Recommendation{
			Severity: "info",
			Category: "System Health",
			Issue:    "No issues detected",
			Action:   "Your CLI environment is well-maintained! All tools are properly managed and no conflicts detected.",
		})
	}

	return recs
}

func generateMarkdownReport(result AuditResult) string {
	var sb strings.Builder

	// Header
	sb.WriteString("# CLI Environment Audit Report\n\n")
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	// Executive Summary
	sb.WriteString("## Executive Summary\n\n")
	sb.WriteString(fmt.Sprintf("- **Total CLI Tools:** %d\n", result.TotalTools))
	sb.WriteString(fmt.Sprintf("- **Package-Managed:** %d (%.1f%%)\n",
		result.PackageManagedTools,
		float64(result.PackageManagedTools)/float64(result.TotalTools)*100))
	sb.WriteString(fmt.Sprintf("- **Unmanaged:** %d (%.1f%%)\n",
		result.UnmanagedTools,
		float64(result.UnmanagedTools)/float64(result.TotalTools)*100))
	sb.WriteString(fmt.Sprintf("- **Installation Conflicts:** %d\n", len(result.Clashes)))
	sb.WriteString(fmt.Sprintf("- **Shadowed Installations:** %d\n\n", len(result.ShadowedTools)))

	// Package Managers
	sb.WriteString("## Package Managers\n\n")
	sb.WriteString("| Manager | Packages | Tools Provided |\n")
	sb.WriteString("|---------|----------|----------------|\n")
	for _, pm := range result.PackageManagers {
		sb.WriteString(fmt.Sprintf("| %s | %d | %d |\n", pm.Name, pm.PackageCount, pm.ToolCount))
	}
	sb.WriteString("\n")

	// Recommendations
	sb.WriteString("## Recommendations\n\n")
	if len(result.Recommendations) > 0 {
		for i, rec := range result.Recommendations {
			icon := "ℹ️"
			switch rec.Severity {
			case "high":
				icon = "🔴"
			case "medium":
				icon = "🟡"
			case "low":
				icon = "🟢"
			}

			sb.WriteString(fmt.Sprintf("### %d. %s %s - %s\n\n", i+1, icon, strings.ToUpper(rec.Severity), rec.Category))
			sb.WriteString(fmt.Sprintf("**Issue:** %s\n\n", rec.Issue))
			sb.WriteString(fmt.Sprintf("**Action:** %s\n\n", rec.Action))
		}
	}

	// Installation Conflicts Details
	if len(result.Clashes) > 0 {
		sb.WriteString("## Installation Conflicts (Detailed)\n\n")
		sb.WriteString("The following tools have multiple installations from different package managers:\n\n")

		for _, clash := range result.Clashes {
			sb.WriteString(fmt.Sprintf("### `%s`\n\n", clash.ToolName))
			for _, inst := range clash.Installations {
				status := ""
				if inst.IsActive {
					status = " ✓ **ACTIVE**"
				} else {
					status = " (shadowed)"
				}
				sb.WriteString(fmt.Sprintf("- `%s` via **%s** (v%s)%s\n",
					inst.Path, inst.PackageManager, inst.Version, status))
			}
			sb.WriteString("\n")
		}
	}

	// Shadowed Tools Details
	if len(result.ShadowedTools) > 0 {
		sb.WriteString("## Shadowed Installations (Detailed)\n\n")
		sb.WriteString("These tool installations exist but are not being used:\n\n")
		sb.WriteString("| Tool | Active | Shadowed |\n")
		sb.WriteString("|------|--------|----------|\n")

		for _, shadow := range result.ShadowedTools {
			sb.WriteString(fmt.Sprintf("| `%s` | %s (%s) | %s (%s) |\n",
				shadow.ToolName,
				shadow.ActivePath,
				shadow.ActivePackage,
				shadow.ShadowedPath,
				shadow.ShadowedPackage))
		}
		sb.WriteString("\n")
	}

	// AI Agent Notes
	sb.WriteString("## Notes for AI Agents\n\n")
	sb.WriteString("This audit report can be used to:\n")
	sb.WriteString("1. Identify package manager conflicts before installing new tools\n")
	sb.WriteString("2. Recommend cleanup actions to users\n")
	sb.WriteString("3. Understand which package managers are available on the system\n")
	sb.WriteString("4. Detect potential PATH issues or version conflicts\n")
	sb.WriteString("5. Provide context when troubleshooting tool-related issues\n\n")

	sb.WriteString("**Command to re-run audit:**\n")
	sb.WriteString("```bash\n")
	sb.WriteString("cli-ai audit --output cli-audit.md\n")
	sb.WriteString("```\n")

	return sb.String()
}

func init() {
	rootCmd.AddCommand(auditCmd)
	auditCmd.Flags().StringVarP(&auditOutput, "output", "o", "", "save audit report to file (default: display to console)")
}
