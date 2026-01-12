package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Used for flags
	cfgFile string
	verbose bool

	// Version information (set by main.go)
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "Discover and explore CLI tools installed on your system",
	Long: `cli is a CLI tool that helps you discover, explore, and manage
all the command-line tools available on your system. It scans your PATH
and installed packages to find all accessible CLI tools.

Available Commands:
  cli help              Show this help message
  cli list              List all available CLI tools
  cli list --all        List all CLI tools with detailed information
  cli packages          List packages that provide CLI tools (npm, pip, brew, etc.)
  cli export            Export tools catalog in JSON format for AI agents
  cli export --output   Export catalog to a file
  cli debug <package>   Show debug information for a specific package
  cli debug --all       Show debug information for all packages

Global Flags:
  -v, --verbose           Enable verbose output
  --config <file>         Specify config file (default: $HOME/.cli.yaml)

Use "cli [command] --help" for more information about a command.`,
	Example: `  # Show help
  cli help

  # List all CLI tools
  cli list

  # List with detailed information
  cli list --all

  # List packages with CLI tools
  cli packages

  # Export catalog for AI agents
  cli export --pretty --output tools.json

  # Export with package information
  cli export --with-packages --pretty -o tools.json

  # Debug a specific package
  cli debug npm

  # Debug all packages
  cli debug --all`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// SetVersion sets the version information for the CLI
func SetVersion(v, c, d string) {
	version = v
	commit = c
	date = d
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// TODO: Implement config file reading if needed
}
