package components

import (
	"github.com/charmbracelet/lipgloss"
	lipglosstable "github.com/charmbracelet/lipgloss/table"
)

var headers = []string{"SSHade", "f faction", "o options"}

func getHeaderCellWidth(username string) int {
	maxCellWidth := 0
	for _, h := range headers {
		if w := lipgloss.Width(h); w > maxCellWidth {
			maxCellWidth = w
		}
	}
	if w := lipgloss.Width(username); w > maxCellWidth {
		maxCellWidth = w
	}
	if maxCellWidth == 0 {
		maxCellWidth = lipgloss.Width("—")
	}
	return maxCellWidth + 2
}

func Header(username string, factionname string) string {
	userDisplay := username
	if userDisplay == "" {
		userDisplay = "—"
	}

	factionDisplay := ""
	if factionname != "" {
		factionDisplay = factionname
	} else {
		factionDisplay = headers[1]
	}

	headerTable := lipglosstable.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("15"))).
		BorderHeader(false).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Padding(0, 1).Width(getHeaderCellWidth(userDisplay)).Align(lipgloss.Center)
		}).
		Headers(headers[0], factionDisplay, headers[2], userDisplay)
	return headerTable.String()
}
