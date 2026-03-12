package components

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	emptyFactionsBoxWidth = 44
)

// EmptyFactionsView renders a styled empty state when no factions exist.
// It is centered within the given width and height.
func EmptyFactionsView(r *lipgloss.Renderer, width, height int) string {
	accent := lipgloss.Color("99") // cyan
	muted := lipgloss.Color("241") // dim
	keyFg := lipgloss.Color("205") // pink, matches splash

	titleStyle := r.NewStyle().
		Bold(true).
		Foreground(accent).
		MarginBottom(1)

	bodyStyle := r.NewStyle().
		Foreground(muted).
		MarginBottom(2).
		Width(emptyFactionsBoxWidth - 4).
		Align(lipgloss.Center)

	keyStyle := r.NewStyle().
		Foreground(keyFg).
		Bold(true)

	hintStyle := r.NewStyle().Foreground(muted)

	title := titleStyle.Render("No factions yet")
	body := bodyStyle.Render("The canvas has no factions.\nCreate one to claim territory\nand compete.")
	key := keyStyle.Render("c")
	hint := hintStyle.Render("Press ") + key + hintStyle.Render(" to create a new faction")

	content := lipgloss.JoinVertical(lipgloss.Center, title, body, hint)

	boxStyle := r.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(accent).
		Padding(1, 2).
		Align(lipgloss.Center)

	box := boxStyle.Render(content)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box)
}
