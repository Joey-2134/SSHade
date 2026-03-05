package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
)

const (
	CanvasWidth  = 20
	CanvasHeight = 10
)

type Model struct {
	width    int
	height   int
	renderer *lipgloss.Renderer
	keyMap   KeyMap
	canvas   [CanvasHeight][CanvasWidth]string
	cursor   [][]int
}

func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	renderer := bubbletea.MakeRenderer(s)
	m := Model{
		width:    pty.Window.Width,
		height:   pty.Window.Height,
		renderer: renderer,
		keyMap:   DefaultKeyMap,
		cursor:   [][]int{{0, 0}},
	}

	colours := []string{"#ff6b6b", "#4ecdc4", "#45b7d1", "#96ceb4", "#ffeaa7"}
	for y := 0; y < CanvasHeight; y++ {
		for x := 0; x < CanvasWidth; x++ {
			m.canvas[y][x] = colours[(x+y)%len(colours)]
		}
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Up):
			m.cursor[0][1]--
		case key.Matches(msg, m.keyMap.Down):
			m.cursor[0][1]++
		case key.Matches(msg, m.keyMap.Left):
			m.cursor[0][0]--
		case key.Matches(msg, m.keyMap.Right):
			m.cursor[0][0]++
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		}
	default:
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	var b strings.Builder
	for y := 0; y < CanvasHeight; y++ {
		for x := 0; x < CanvasWidth; x++ {
			if m.cursor[0][0] == x && m.cursor[0][1] == y {
				style := m.renderer.NewStyle().Background(lipgloss.Color("241")).SetString("  ")
				b.WriteString(style.String())
			} else {
				style := m.renderer.NewStyle().Background(lipgloss.Color(m.canvas[y][x])).SetString("  ")
				b.WriteString(style.String())
			}
		}
		b.WriteString("\n")
	}
	b.WriteString(m.renderer.NewStyle().Foreground(lipgloss.Color("241")).SetString("Press 'q' to quit").String())
	return b.String()
}

type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Quit  key.Binding
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
	Quit: key.NewBinding(
		key.WithKeys("q", "quit"),
		key.WithHelp("q", "quit the game"),
	),
}
