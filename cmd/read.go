package cmd

import (
	"fmt"
	"strings"

	"github.com/Aayush9029/wiki/wikipedia"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var summaryOnly bool

var readCmd = &cobra.Command{
	Use:   "read <title>",
	Short: "Read a Wikipedia article",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runRead,
}

func init() {
	readCmd.Flags().BoolVarP(&summaryOnly, "summary", "s", false, "show summary only")
	rootCmd.AddCommand(readCmd)
}

func runRead(cmd *cobra.Command, args []string) error {
	title := strings.Join(args, " ")
	client := wikipedia.NewClient()

	if summaryOnly {
		summary, err := client.Summary(title)
		if err != nil {
			return fmt.Errorf("failed to fetch article: %w", err)
		}
		return displaySummary(summary)
	}
	return showFullArticle(client, title)
}

func displaySummary(s *wikipedia.Summary) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
	urlStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Underline(true)

	fmt.Println()
	fmt.Println(titleStyle.Render(s.Title))
	if s.Description != "" {
		fmt.Println(descStyle.Render(s.Description))
	}
	fmt.Println()

	rendered, err := renderMarkdown(s.Extract)
	if err != nil {
		fmt.Println(s.Extract)
	} else {
		fmt.Print(rendered)
	}

	if s.ContentURLs.Desktop.Page != "" {
		fmt.Println(urlStyle.Render(s.ContentURLs.Desktop.Page))
	}

	return nil
}

func showFullArticle(client *wikipedia.Client, title string) error {
	article, err := client.FullArticle(title)
	if err != nil {
		return fmt.Errorf("failed to fetch article: %w", err)
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	urlStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Underline(true)

	fmt.Println()
	fmt.Println(titleStyle.Render(article.Title))
	fmt.Println()

	rendered, err := renderMarkdown(formatSections(article.Extract))
	if err != nil {
		fmt.Println(article.Extract)
	} else {
		fmt.Print(rendered)
	}

	fmt.Println(urlStyle.Render(fmt.Sprintf("https://en.wikipedia.org/wiki/%s", strings.ReplaceAll(article.Title, " ", "_"))))
	return nil
}

func renderMarkdown(text string) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return "", err
	}
	return renderer.Render(text)
}

// formatSections converts Wikipedia's plain-text section headers (== Title ==)
// into markdown headings so glamour can render them.
func formatSections(text string) string {
	var b strings.Builder
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "==") && strings.HasSuffix(trimmed, "==") {
			level := 0
			for _, r := range trimmed {
				if r == '=' {
					level++
				} else {
					break
				}
			}
			heading := strings.Trim(trimmed, "= ")
			prefix := strings.Repeat("#", level-1)
			if prefix == "" {
				prefix = "##"
			}
			b.WriteString(prefix + " " + heading + "\n")
		} else {
			b.WriteString(line + "\n")
		}
	}
	return b.String()
}
