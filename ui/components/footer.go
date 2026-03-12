package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const githubLink = "https://github.com/Joey-2134"
const quitMsg = "q quit"

func Footer(width int, factioncolour string) string {
	colour := effectiveFactionColour(factioncolour)
	line := lipgloss.NewStyle().Foreground(lipgloss.Color(colour)).Render(strings.Repeat("─", width))
	quitMsgWidth := lipgloss.Width(quitMsg)

	footerInfoString := quitMsg + strings.Repeat(" ", width-quitMsgWidth-lipgloss.Width(githubLink)) + githubLink
	footerInfo := lipgloss.NewStyle().Foreground(lipgloss.Color(colour)).Render(footerInfoString)
	return lipgloss.JoinVertical(lipgloss.Center, line, footerInfo)
}
