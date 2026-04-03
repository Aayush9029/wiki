package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Version is set via ldflags at build time.
var Version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "wiki",
	Short: "Search and read Wikipedia from your terminal",
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = Version
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
