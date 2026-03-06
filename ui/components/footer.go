package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func Footer(width int) string {
	line := lipgloss.NewStyle().Render(strings.Repeat("─", width))
	footerInfo := lipgloss.NewStyle().Render("q quit")
	return lipgloss.JoinVertical(lipgloss.Center, line, footerInfo)
}
