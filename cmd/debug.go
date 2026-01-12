package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	debugAll bool
)

// debugCmd represents the debug command
var debugCmd = &cobra.Command{
	Use:   "debug [package_name]",
	Short: "Show debug information for CLI tools",
	Long: `Display detailed debug information about CLI tools and packages.

You can either debug a specific package by name, or use the --all flag
to show debug information for all discovered packages.

This command will show:
  - Package location
  - Installation path
  - Binary details
  - Dependencies (if available)
  - Version information (if available)`,
	Example: `  # Debug a specific package
  cli debug npm
  cli debug python

  # Debug all packages
  cli debug --all

  # Debug with verbose output
  cli debug npm --verbose`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if debugAll {
			// Debug all packages
			fmt.Fprintln(os.Stdout, "Debugging all packages...")
			// TODO: Implement debug all logic
			fmt.Fprintln(os.Stdout, "This feature will be implemented soon.")
		} else if len(args) == 0 {
			// No package specified and --all not set
			cmd.PrintErr("Error: must specify a package name or use --all flag\n\n")
			cmd.Usage()
			os.Exit(1)
		} else {
			// Debug specific package
			packageName := args[0]
			fmt.Fprintf(os.Stdout, "Debugging package: %s\n", packageName)
			// TODO: Implement package-specific debug logic
			fmt.Fprintln(os.Stdout, "This feature will be implemented soon.")
		}
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
	debugCmd.Flags().BoolVarP(&debugAll, "all", "a", false, "show debug information for all packages")
}
