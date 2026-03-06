package ui

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
	gocrypto "golang.org/x/crypto/ssh"

	"github.com/Joey-2134/SSHade/canvas"
	"github.com/Joey-2134/SSHade/constants"
	"github.com/Joey-2134/SSHade/db"
	"github.com/Joey-2134/SSHade/ui/components"
)

const DefaultPlaceColour = "#ff6b6b"
const defaultCellColour = "#cccccc"
const gridSize = constants.GridSize

type CanvasUpdateMsg struct {
	Pixel canvas.Pixel
}

type Model struct {
	width          int
	height         int
	isTooSmall     bool
	renderer       *lipgloss.Renderer
	keyMap         constants.KeyMap
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
		keyMap:         constants.DefaultKeyMap,
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
		if m.width < constants.MinTerminalWidth || m.height < constants.MinTerminalHeight {
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
		m.cursor.X = ((m.cursor.X % gridSize) + gridSize) % gridSize
		m.cursor.Y = ((m.cursor.Y % gridSize) + gridSize) % gridSize
	default:
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	username := ""
	if m.user != nil {
		username = m.user.Username
	}

	header := components.Header(username)

	if m.isTooSmall {
		return fmt.Sprintf("Terminal too small.\nPlease resize to at least %dx%d.\n\nPress Q to quit.", constants.MinTerminalWidth, constants.MinTerminalHeight)
	}

	// Reserve space for header so grid scales to fit remaining area
	headerHeight := lipgloss.Height(header)
	availableHeight := m.height - headerHeight

	// Draw the canvas
	gridStr := components.Grid(
		m.width, availableHeight,
		m.renderer,
		m.canvasRef,
		m.cursor.X, m.cursor.Y,
		defaultCellColour,
	)

	fullView := lipgloss.JoinVertical(lipgloss.Center, header, gridStr)
	fullViewWidth := lipgloss.Width(fullView)
	fullViewHeight := lipgloss.Height(fullView)
	leftPad := (m.width - fullViewWidth) / 2
	topLines := (m.height - fullViewHeight) / 2

	styled := m.renderer.NewStyle().
		PaddingLeft(max(leftPad, 0)).
		PaddingTop(max(topLines, 0)).
		Render(fullView)

	return styled
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
