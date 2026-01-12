package cmd

import (
	"os"

	"github.com/cli-ai-org/cli/internal/display"
	"github.com/cli-ai-org/cli/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	listAll  bool
	listJSON bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available CLI tools",
	Long: `Scans your system's PATH and lists all available CLI tools.
This command discovers all executable files in your PATH directories.

By default, shows a simple list of tool names. Use --all flag for detailed information
including full paths and additional metadata.

Use --json flag to output in JSON format for programmatic access or AI agent consumption.`,
	Example: `  # List all CLI tools
  cli list

  # List with full paths and details
  cli list --all

  # List in JSON format for AI agents
  cli list --json

  # List with verbose output
  cli list --verbose`,
	Run: func(cmd *cobra.Command, args []string) {
		s := scanner.New()
		d := display.New(os.Stdout)

		// Use detailed scan if JSON or --all is requested
		if listJSON || listAll {
			tools, err := s.ScanAllDetailed()
			if err != nil {
				cmd.PrintErrf("Error scanning for tools: %v\n", err)
				os.Exit(1)
			}

			if listJSON {
				if err := d.ShowToolsJSON(tools, true); err != nil {
					cmd.PrintErrf("Error encoding JSON: %v\n", err)
					os.Exit(1)
				}
			} else {
				d.ShowToolsDetailed(tools)
			}
		} else {
			// Simple list
			tools, err := s.ScanAll()
			if err != nil {
				cmd.PrintErrf("Error scanning for tools: %v\n", err)
				os.Exit(1)
			}
			d.ShowTools(tools)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "show detailed information including paths")
	listCmd.Flags().BoolVarP(&listJSON, "json", "j", false, "output in JSON format for AI agents")
}
