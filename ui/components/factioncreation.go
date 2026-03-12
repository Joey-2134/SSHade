package components

import (
	"github.com/charmbracelet/lipgloss"
)

func normalizeHexForStyle(hex string) string {
	if hex == "" {
		return "99"
	}
	// Preserve 256-color codes (e.g. "99", "240"); only add # for 6-char hex
	if len(hex) == 6 && hex[0] != '#' {
		return "#" + hex
	}
	return hex
}

// FactionCreationForm renders the faction creation form layout with styling
// consistent with the faction selection screen (bordered box, accent colour, centered).
// nameInputView and colourInputView are the result of calling .View() on the respective textinput models.
func FactionCreationForm(r *lipgloss.Renderer, width, height int, accentHex, nameInputView, colourInputView, err string) string {
	accent := lipgloss.Color(normalizeHexForStyle(accentHex))
	muted := lipgloss.Color("241")
	keyFg := lipgloss.Color("205")
	errFg := lipgloss.Color("9")

	titleStyle := r.NewStyle().
		Bold(true).
		Foreground(accent).
		MarginBottom(1)

	labelStyle := r.NewStyle().
		Foreground(muted).
		MarginBottom(0)

	keyStyle := r.NewStyle().
		Foreground(keyFg).
		Bold(true)

	hintStyle := r.NewStyle().Foreground(muted)
	errStyle := r.NewStyle().Foreground(errFg)

	title := titleStyle.Render("Create a faction")
	nameLabel := labelStyle.Render("Name:")
	colourLabel := labelStyle.Render("Colour (hex):")

	enterKey := keyStyle.Render("enter")
	qKey := keyStyle.Render("q")
	hint := hintStyle.Render("Press ") + enterKey + hintStyle.Render(" next / submit · ") + qKey + hintStyle.Render(" back to list")

	inner := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		nameLabel,
		nameInputView,
		"",
		colourLabel,
		colourInputView,
		"",
	)
	if err != "" {
		inner += "\n" + errStyle.Render(err) + "\n"
	}
	inner += "\n" + hint

	boxStyle := r.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(accent).
		Padding(1, 2).
		Width(44).
		Align(lipgloss.Left)

	box := boxStyle.Render(inner)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box)
}
