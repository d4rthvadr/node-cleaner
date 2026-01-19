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
	fmt.Fprintln(w, "SIZE\tLAST ACCESSED\tPATH")
	fmt.Fprintln(w, strings.Repeat("─", 80))

	for _, folder := range result.Folders {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			humanize.Bytes(uint64(folder.Size)),
			humanize.Time(folder.AccessTime),
			folder.Path,
		)

	}
	w.Flush()

	fmt.Println(strings.Repeat("─", 80))
	fmt.Printf("\n%s\n", headerStyle.Render("Summary:"))
	fmt.Printf(" Total folders: %d\n", result.TotalCount)
	fmt.Printf(" Total size: %s\n", humanize.Bytes(uint64(result.TotalSize)))
	fmt.Printf(" Scan duration: %s\n", result.Duration)
	if result.CacheHits > 0 {
		fmt.Printf("  Cache hits: %d (%.1f%%)\n",
			result.CacheHits,
			float64(result.CacheHits)/float64(result.CacheHits+result.CacheMisses)*100)
	}
}
