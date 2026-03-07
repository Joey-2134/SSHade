package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func Footer(width int, cooldownTimer string) string {
	line := lipgloss.NewStyle().Render(strings.Repeat("─", width))
	quitMsg := "q quit"
	quitMsgWidth := lipgloss.Width(quitMsg)
	cooldownMsg := "cooldown: " + cooldownTimer

	footerInfoString := quitMsg + strings.Repeat(" ", width-quitMsgWidth-lipgloss.Width(cooldownMsg)) + cooldownMsg

	footerInfo := lipgloss.NewStyle().Render(footerInfoString)

	return lipgloss.JoinVertical(lipgloss.Center, line, footerInfo)
}
