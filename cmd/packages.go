package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/cli-ai-org/cli/internal/packages"
	"github.com/cli-ai-org/cli/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	packagesJSON    bool
	packagesManager string
)

// packagesCmd represents the packages command
var packagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "List packages that provide CLI tools",
	Long: `List all packages from various package managers (npm, pip, brew, cargo, gem)
that provide command-line tools.

This helps identify which package a CLI tool comes from, useful for tools
like vercel, supabase, aws-cli, etc.`,
	Example: `  # List all packages with CLI tools
  cli packages

  # List npm packages only
  cli packages --manager npm

  # List in JSON format
  cli packages --json

  # Find which package provides a tool
  cli packages | grep vercel`,
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			fmt.Fprintln(os.Stderr, "Detecting packages from package managers...")
		}

		// Detect packages
		detector := packages.NewDetector()
		pkgs, err := detector.DetectAll()
		if err != nil {
			cmd.PrintErrf("Error detecting packages: %v\n", err)
			os.Exit(1)
		}

		// Filter by manager if specified
		if packagesManager != "" {
			filtered := []packages.Package{}
			for _, pkg := range pkgs {
				if string(pkg.Manager) == packagesManager {
					filtered = append(filtered, pkg)
				}
			}
			pkgs = filtered
		}

		// Link packages to tools to find which packages provide CLIs
		s := scanner.New()
		tools, err := s.ScanAllDetailed()
		if err != nil {
			cmd.PrintErrf("Error scanning tools: %v\n", err)
			os.Exit(1)
		}

		linker := packages.NewLinker(pkgs)
		enrichedTools := linker.LinkTools(tools)

		// Get packages that have binaries
		pkgsWithBinaries := packages.GetPackagesWithBinaries(pkgs, enrichedTools)

		if packagesJSON {
			// JSON output
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(pkgsWithBinaries); err != nil {
				cmd.PrintErrf("Error encoding JSON: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Human-readable output
			if len(pkgsWithBinaries) == 0 {
				fmt.Fprintln(os.Stdout, "No packages with CLI tools found.")
				return
			}

			// Sort by name
			sort.Slice(pkgsWithBinaries, func(i, j int) bool {
				return pkgsWithBinaries[i].Name < pkgsWithBinaries[j].Name
			})

			fmt.Fprintf(os.Stdout, "Found %d packages with CLI tools:\n\n", len(pkgsWithBinaries))
			fmt.Fprintf(os.Stdout, "%-30s %-10s %-15s %s\n", "PACKAGE", "MANAGER", "VERSION", "CLIs")
			fmt.Fprintf(os.Stdout, "%-30s %-10s %-15s %s\n", "-------", "-------", "-------", "----")

			for _, pkg := range pkgsWithBinaries {
				binaries := "none"
				if len(pkg.Binaries) > 0 {
					if len(pkg.Binaries) == 1 {
						binaries = pkg.Binaries[0]
					} else if len(pkg.Binaries) <= 3 {
						binaries = fmt.Sprintf("%v", pkg.Binaries)
					} else {
						binaries = fmt.Sprintf("%d binaries", len(pkg.Binaries))
					}
				}
				fmt.Fprintf(os.Stdout, "%-30s %-10s %-15s %s\n",
					pkg.Name,
					pkg.Manager,
					pkg.Version,
					binaries,
				)
			}

			if verbose {
				fmt.Fprintf(os.Stderr, "\nTotal packages scanned: %d\n", len(pkgs))
				fmt.Fprintf(os.Stderr, "Packages with CLIs: %d\n", len(pkgsWithBinaries))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(packagesCmd)
	packagesCmd.Flags().BoolVarP(&packagesJSON, "json", "j", false, "output in JSON format")
	packagesCmd.Flags().StringVarP(&packagesManager, "manager", "m", "", "filter by package manager (npm, pip, brew, cargo, gem)")
}
