package cmd

import (
	"fmt"

	"github.com/Aayush9029/wiki/wikipedia"
	"github.com/spf13/cobra"
)

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Read a random Wikipedia article",
	Args:  cobra.NoArgs,
	RunE:  runRandom,
}

func init() {
	rootCmd.AddCommand(randomCmd)
}

func runRandom(cmd *cobra.Command, args []string) error {
	client := wikipedia.NewClient()
	summary, err := client.Random()
	if err != nil {
		return fmt.Errorf("failed to fetch random article: %w", err)
	}
	return displaySummary(summary)
}
