package constants

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Up               key.Binding
	Down             key.Binding
	Left             key.Binding
	Right            key.Binding
	Place            key.Binding
	Quit             key.Binding
	Enter            key.Binding
	FactionSelection key.Binding
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("w", "up"),
		key.WithHelp("↑/w", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("s", "down"),
		key.WithHelp("↓/s", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("a", "left"),
		key.WithHelp("←/a", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("d", "right"),
		key.WithHelp("→/d", "move right"),
	),
	Place: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "place pixel"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "quit"),
		key.WithHelp("q", "quit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "enter"),
	),
	FactionSelection: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "faction selection"),
	),
}
