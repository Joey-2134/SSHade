package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const githubLink = "https://github.com/Joey-2134"
const quitMsg = "q quit"

func Footer(width int, factioncolour string) string {
	line := lipgloss.NewStyle().Foreground(lipgloss.Color(factioncolour)).Render(strings.Repeat("─", width))
	quitMsgWidth := lipgloss.Width(quitMsg)

	footerInfoString := quitMsg + strings.Repeat(" ", width-quitMsgWidth-lipgloss.Width(githubLink)) + githubLink
	footerInfo := lipgloss.NewStyle().Foreground(lipgloss.Color(factioncolour)).Render(footerInfoString)
	return lipgloss.JoinVertical(lipgloss.Center, line, footerInfo)
}
