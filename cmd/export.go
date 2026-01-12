package cmd

import (
	"fmt"
	"os"

	"github.com/cli-ai-org/cli/internal/collector"
	"github.com/cli-ai-org/cli/internal/display"
	"github.com/cli-ai-org/cli/internal/packages"
	"github.com/cli-ai-org/cli/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	exportJSON        bool
	exportPretty      bool
	exportOutput      string
	exportWithMeta    bool
	exportWithPackages bool
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export CLI tools catalog for AI agents",
	Long: `Export a comprehensive catalog of all CLI tools in a format optimized for AI agents.

This command generates a machine-readable JSON catalog containing:
  - Complete list of all CLI tools
  - Full paths and locations
  - Tool metadata (size, symlinks, etc.)
  - Optional: Version information (slower, requires running tools)
  - Optional: Help text extraction (slower, requires running tools)
  - Optional: Package information (which package each tool comes from)

The exported catalog can be used by AI agents to discover and understand
available CLI tools on the system.`,
	Example: `  # Export basic catalog to stdout
  cli export

  # Export to file
  cli export --output tools.json

  # Export with pretty formatting
  cli export --pretty --output tools.json

  # Export with metadata (version, help text) - slower
  cli export --with-meta --output tools-detailed.json

  # Export with package information
  cli export --with-packages --pretty --output tools-with-packages.json

  # Pipe to AI agent or other tool
  cli export | jq '.tools[] | .name'`,
	Run: func(cmd *cobra.Command, args []string) {
		s := scanner.New()

		if verbose {
			fmt.Fprintln(os.Stderr, "Scanning for CLI tools...")
		}

		tools, err := s.ScanAllDetailed()
		if err != nil {
			cmd.PrintErrf("Error scanning for tools: %v\n", err)
			os.Exit(1)
		}

		if verbose {
			fmt.Fprintf(os.Stderr, "Found %d tools\n", len(tools))
		}

		// Detect packages if requested
		var pkgs []packages.Package
		if exportWithPackages {
			if verbose {
				fmt.Fprintln(os.Stderr, "Detecting packages...")
			}

			detector := packages.NewDetector()
			var err error
			pkgs, err = detector.DetectAll()
			if err != nil && verbose {
				fmt.Fprintf(os.Stderr, "Warning: some package managers failed: %v\n", err)
			}

			if verbose {
				fmt.Fprintf(os.Stderr, "Found %d packages\n", len(pkgs))
				fmt.Fprintln(os.Stderr, "Linking tools to packages...")
			}

			// Link tools to packages
			linker := packages.NewLinker(pkgs)
			tools = linker.LinkTools(tools)
		}

		// Collect additional metadata if requested
		if exportWithMeta {
			if verbose {
				fmt.Fprintln(os.Stderr, "Collecting metadata (this may take a while)...")
			}

			c := collector.New()
			for i := range tools {
				if verbose && i%50 == 0 {
					fmt.Fprintf(os.Stderr, "Processing tool %d/%d...\n", i+1, len(tools))
				}

				enriched, err := c.CollectToolInfo(tools[i].Name, tools[i].Path)
				if err == nil && enriched != nil {
					tools[i].Version = enriched.Version
					tools[i].HelpText = enriched.HelpText
				}
			}
		}

		// Build catalog
		c := collector.New()
		catalog := c.BuildCatalog(tools, s.GetPaths())

		// Add package information to catalog if available
		if exportWithPackages && len(pkgs) > 0 {
			pkgsWithBinaries := packages.GetPackagesWithBinaries(pkgs, tools)
			catalog.Packages = pkgsWithBinaries
			catalog.TotalPackages = len(pkgsWithBinaries)
		}

		// Determine output writer
		writer := os.Stdout
		if exportOutput != "" {
			file, err := os.Create(exportOutput)
			if err != nil {
				cmd.PrintErrf("Error creating output file: %v\n", err)
				os.Exit(1)
			}
			defer file.Close()
			writer = file
		}

		// Output catalog
		d := display.New(writer)
		if err := d.ShowCatalogJSON(catalog, exportPretty); err != nil {
			cmd.PrintErrf("Error encoding JSON: %v\n", err)
			os.Exit(1)
		}

		if verbose && exportOutput != "" {
			fmt.Fprintf(os.Stderr, "Catalog exported to %s\n", exportOutput)
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().BoolVarP(&exportJSON, "json", "j", true, "output in JSON format (default)")
	exportCmd.Flags().BoolVarP(&exportPretty, "pretty", "p", false, "pretty-print JSON output")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "output file (default: stdout)")
	exportCmd.Flags().BoolVarP(&exportWithMeta, "with-meta", "m", false, "include version and help text (slower)")
	exportCmd.Flags().BoolVarP(&exportWithPackages, "with-packages", "P", false, "include package information (npm, pip, brew, etc.)")
}
