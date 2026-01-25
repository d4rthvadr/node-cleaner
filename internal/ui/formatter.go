package ui

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
	"github.com/d4rthvadr/node-cleaner/pkg/models"
	"github.com/dustin/go-humanize"
)

var (
	headerStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11"))
)

func DisplayScanResults(result *models.ScanResult) {

	fmt.Println(headerStyle.Render("Scan Results:"))
	fmt.Println(strings.Repeat("-", 80))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Colorful header
	fmt.Fprintln(w, headerStyle.Render("SIZE")+"\t"+
		headerStyle.Render("LAST ACCESSED")+"\t"+
		headerStyle.Render("PATH"))
	fmt.Fprintln(w, strings.Repeat("─", 80))

	pathStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14"))

	for _, folder := range result.Folders {
		// Color code by size
		sizeStr := humanize.Bytes(uint64(folder.Size))
		if folder.Size > 500*1024*1024 { // > 500MB
			sizeStr = errorStyle.Render(sizeStr) // red for large
		} else if folder.Size > 100*1024*1024 { // > 100MB
			sizeStr = warningStyle.Render(sizeStr) // yellow for medium
		} else {
			sizeStr = successStyle.Render(sizeStr) // green for small
		}

		fmt.Fprintf(w, "%s\t%s\t%s\n",
			sizeStr,
			humanize.Time(folder.AccessTime),
			pathStyle.Render(folder.Path),
		)

	}
	w.Flush()

	fmt.Println(strings.Repeat("─", 80))
	fmt.Printf("\n%s\n", headerStyle.Render("✨ Summary:"))
	fmt.Printf(" Total folders: %s\n", successStyle.Render(fmt.Sprintf("%d", result.TotalCount)))
	fmt.Printf(" Total size: %s\n", errorStyle.Render(humanize.Bytes(uint64(result.TotalSize))))
	fmt.Printf(" Scan duration: %s\n", result.Duration)
	if result.CacheHits > 0 {
		fmt.Printf(" Cache hits: %s (%.1f%%)\n",
			successStyle.Render(fmt.Sprintf("%d", result.CacheHits)),
			float64(result.CacheHits)/float64(result.CacheHits+result.CacheMisses)*100)
	}
}

func DisplayCleanResults(result *models.CleanResult) {

	fmt.Println()

	if result.DryRun {
		fmt.Println(warningStyle.Render("Dry Run: No folders were actually deleted."))
		fmt.Println()
	}
}
