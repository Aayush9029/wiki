package cmd

import (
	"fmt"
	"strings"

	"github.com/Aayush9029/wiki/wikipedia"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var searchLimit int

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search Wikipedia articles",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSearch,
}

func init() {
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "n", 5, "number of results")
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := strings.Join(args, " ")
	client := wikipedia.NewClient()

	results, err := client.Search(query, searchLimit)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(results) == 0 {
		fmt.Println("No results found.")
		return nil
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	indexStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)

	for i, r := range results {
		fmt.Printf("%s %s\n", indexStyle.Render(fmt.Sprintf("[%d]", i+1)), titleStyle.Render(r.Title))
		if r.Snippet != "" {
			fmt.Printf("    %s\n", descStyle.Render(stripHTML(r.Snippet)))
		}
		fmt.Println()
	}

	fmt.Println(hintStyle.Render("Read an article: wiki read <title>"))
	return nil
}

// stripHTML removes HTML tags from a string.
func stripHTML(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			b.WriteRune(r)
		}
	}
	return b.String()
}
