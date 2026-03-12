package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func Footer(width int, cooldownTimer string, factioncolour string) string {
	line := lipgloss.NewStyle().Foreground(lipgloss.Color(factioncolour)).Render(strings.Repeat("─", width))
	quitMsg := "q quit"
	quitMsgWidth := lipgloss.Width(quitMsg)
	cooldownMsg := "cooldown: " + cooldownTimer

	footerInfoString := quitMsg + strings.Repeat(" ", width-quitMsgWidth-lipgloss.Width(cooldownMsg)) + cooldownMsg
	footerInfo := lipgloss.NewStyle().Foreground(lipgloss.Color(factioncolour)).Render(footerInfoString)

	return lipgloss.JoinVertical(lipgloss.Center, line, footerInfo)
}
