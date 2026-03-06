package ui

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
	gocrypto "golang.org/x/crypto/ssh"

	"github.com/Joey-2134/SSHade/canvas"
	"github.com/Joey-2134/SSHade/db"
	"github.com/Joey-2134/SSHade/ui/components"
)

const (
	CanvasWidth       = 20
	CanvasHeight      = 20
	MinTerminalWidth  = 50
	MinTerminalHeight = 25
)

const DefaultPlaceColour = "#ff6b6b"
const defaultCellColour = "#cccccc"

type CanvasUpdateMsg struct {
	Pixel canvas.Pixel
}

type Model struct {
	width          int
	height         int
	isTooSmall     bool
	renderer       *lipgloss.Renderer
	keyMap         KeyMap
	canvasRef      *canvas.Canvas
	db             *sql.DB
	cursor         Cursor
	canvasUpdateCh <-chan canvas.Pixel
	unsub          func()
	user           *db.User
}

func TeaHandler(s ssh.Session, c *canvas.Canvas, database *sql.DB, bc *canvas.Broadcaster) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	renderer := bubbletea.MakeRenderer(s)
	opts := []tea.ProgramOption{tea.WithAltScreen()}

	fingerprint := ""
	if pk := s.PublicKey(); pk != nil {
		fingerprint = gocrypto.FingerprintSHA256(pk)
	}

	user, err := db.GetUserByFingerprint(database, fingerprint)
	if err != nil || user == nil {
		// New user → show username creation screen
		return UserCreationModelHandler(s, database, fingerprint, c, bc), opts
	}

	canvasUpdateCh, unsub := bc.Subscribe()
	m := Model{
		width:          pty.Window.Width,
		height:         pty.Window.Height,
		renderer:       renderer,
		keyMap:         DefaultKeyMap,
		canvasRef:      c,
		db:             database,
		cursor:         DefaultCursor,
		canvasUpdateCh: canvasUpdateCh,
		unsub:          unsub,
		user:           user,
	}
	return m, opts
}

func waitForCanvasUpdate(ch <-chan canvas.Pixel) tea.Cmd {
	return func() tea.Msg {
		p, ok := <-ch
		if !ok {
			return nil
		}
		return CanvasUpdateMsg{Pixel: p}
	}
}

func (m Model) Init() tea.Cmd {
	if m.canvasUpdateCh == nil {
		return nil
	}
	return waitForCanvasUpdate(m.canvasUpdateCh)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.width < MinTerminalWidth || m.height < MinTerminalHeight {
			m.isTooSmall = true
		} else {
			m.isTooSmall = false
		}
	case CanvasUpdateMsg:
		return m, waitForCanvasUpdate(m.canvasUpdateCh)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Up):
			m.cursor.Y--
		case key.Matches(msg, m.keyMap.Down):
			m.cursor.Y++
		case key.Matches(msg, m.keyMap.Left):
			m.cursor.X--
		case key.Matches(msg, m.keyMap.Right):
			m.cursor.X++
		case key.Matches(msg, m.keyMap.Place):
			if m.canvasRef != nil && m.db != nil {
				_ = m.canvasRef.Set(context.Background(), m.db, m.cursor.X, m.cursor.Y, DefaultPlaceColour)
			}
		case key.Matches(msg, m.keyMap.Quit):
			if m.unsub != nil {
				m.unsub()
			}
			return m, tea.Quit
		}

		// Wrap cursor around canvas edges
		w, h := CanvasWidth, CanvasHeight
		if m.canvasRef != nil {
			w, h = m.canvasRef.Width(), m.canvasRef.Height()
		}
		m.cursor.X = ((m.cursor.X % w) + w) % w
		m.cursor.Y = ((m.cursor.Y % h) + h) % h
	default:
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	if m.isTooSmall {
		return fmt.Sprintf("Terminal too small.\nPlease resize to at least %dx%d.\n\nPress Q to quit.", MinTerminalWidth, MinTerminalHeight)
	}

	// Scale to fit: limit cell size by both width and height so grid fits and cells stay square (2:1 line-to-column aspect).
	scaleByWidth := m.width / CanvasWidth
	scaleByHeight := 2 * (m.height / CanvasHeight) // 2 cols per line for square cells
	cellWidth := max(min(scaleByWidth, scaleByHeight), 1)
	linesPerRow := max(cellWidth/2, 1)

	leftPad := (m.width - CanvasWidth*cellWidth) / 2
	topLines := (m.height - CanvasHeight*linesPerRow) / 2

	for y := range CanvasHeight {
		for range linesPerRow {
			for x := range CanvasWidth {
				colour := defaultCellColour
				if m.canvasRef != nil {
					if p, ok := m.canvasRef.PixelAt(x, y); ok {
						colour = p.ColourHex
					}
				}
				if m.cursor.X == x && m.cursor.Y == y {
					style := m.renderer.NewStyle().Background(lipgloss.Color("241")).SetString(strings.Repeat(" ", cellWidth))
					b.WriteString(style.String())
				} else {
					style := m.renderer.NewStyle().Background(lipgloss.Color(colour)).SetString(strings.Repeat(" ", cellWidth))
					b.WriteString(style.String())
				}
			}
			b.WriteString("\n")
		}
	}

	username := ""
	if m.user != nil {
		username = m.user.Username
	}

	header := components.Header(username)
	gridStr := b.String()
	fullView := lipgloss.JoinVertical(lipgloss.Center, header, gridStr)
	styled := m.renderer.NewStyle().
		PaddingLeft(max(leftPad, 0)).
		PaddingTop(max(topLines, 0)).
		Render(fullView)
	return styled
}

type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Place key.Binding
	Quit  key.Binding
	Enter key.Binding
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
}

type Cursor struct {
	X      int
	Y      int
	Colour string
}

var DefaultCursor = Cursor{
	X:      0,
	Y:      0,
	Colour: "241",
}
