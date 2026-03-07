package components

import (
	"github.com/charmbracelet/lipgloss"
)

// FactionCreationForm renders the faction creation form layout.
// nameInputView and colourInputView are the result of calling .View() on the respective textinput models.
func FactionCreationForm(nameInputView, colourInputView, err string) string {
	s := "Create a faction\n\n"
	s += "Name:\n" + nameInputView + "\n\n"
	s += "Colour (hex):\n" + colourInputView + "\n\n"
	if err != "" {
		s += lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(err) + "\n\n"
	}
	s += "\n(enter: next / submit, q: back to faction list)"
	return lipgloss.NewStyle().Render(s)
}
