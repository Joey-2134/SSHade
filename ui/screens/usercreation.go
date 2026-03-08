package ui

import (
	"database/sql"

	"github.com/Joey-2134/SSHade/canvas"
	"github.com/Joey-2134/SSHade/constants"
	"github.com/Joey-2134/SSHade/db"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
)

type UserCreationModel struct {
	renderer    *lipgloss.Renderer
	session     ssh.Session
	database    *sql.DB
	fingerprint string
	canvas      *canvas.Canvas
	broadcaster *canvas.Broadcaster
	textInput   textinput.Model
}

func UserCreationModelHandler(sess ssh.Session, db *sql.DB, fingerprint string, c *canvas.Canvas, bc *canvas.Broadcaster) tea.Model {
	renderer := bubbletea.MakeRenderer(sess)

	ti := textinput.New()
	ti.Placeholder = "Enter your username"
	ti.Focus()
	ti.CharLimit = 32
	ti.Width = 32
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	return UserCreationModel{
		renderer:    renderer,
		session:     sess,
		database:    db,
		fingerprint: fingerprint,
		canvas:      c,
		broadcaster: bc,
		textInput:   ti,
	}
}

func (m UserCreationModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m UserCreationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.DefaultKeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, constants.DefaultKeyMap.Enter):
			user, err := db.CreateUser(m.database, m.textInput.Value(), m.fingerprint)
			if err != nil {
				return m, tea.Batch(tea.Println("Error creating user"), tea.Quit)
			}
			// Transition to canvas model
			canvasUpdateCh, unsub := m.broadcaster.Subscribe()
			pty, _, _ := m.session.Pty()
			return Model{
				width:          pty.Window.Width,
				height:         pty.Window.Height,
				renderer:       m.renderer,
				keyMap:         constants.DefaultKeyMap,
				canvasRef:      m.canvas,
				db:             m.database,
				cursor:         DefaultCursor,
				canvasUpdateCh: canvasUpdateCh,
				unsub:          unsub,
				user:           user,
				session:        m.session,
				broadcaster:    m.broadcaster,
			}, waitForCanvasUpdate(canvasUpdateCh)
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m UserCreationModel) View() string {
	return m.headerView() + m.textInput.View() + m.footerView()
}

func (m UserCreationModel) headerView() string { return "Create your username\n" }
func (m UserCreationModel) footerView() string { return "\n(q to quit)" }
