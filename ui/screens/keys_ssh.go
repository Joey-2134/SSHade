package ui

import (
	"unicode"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// matchesEnter returns true if msg is the Enter key, using either the binding
// or rune fallback (\r, \n). Over SSH/PTY, Enter sometimes arrives as KeyRunes
// instead of KeyEnter, so the first keypress may not match the binding.
func matchesEnter(msg tea.KeyMsg, enter key.Binding) bool {
	if key.Matches(msg, enter) {
		return true
	}
	return msg.Type == tea.KeyRunes && len(msg.Runes) == 1 &&
		(msg.Runes[0] == '\r' || msg.Runes[0] == '\n')
}

// matchesRuneIgnoreCase returns true if msg is a single rune matching r (case-insensitive).
// Used so letter keys (f, c, q) register reliably when sent as KeyRunes over SSH.
func matchesRuneIgnoreCase(msg tea.KeyMsg, r rune) bool {
	return msg.Type == tea.KeyRunes && len(msg.Runes) == 1 &&
		unicode.ToLower(msg.Runes[0]) == unicode.ToLower(r)
}
