package ui

import (
	"database/sql"

	"github.com/Joey-2134/SSHade/canvas"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
)

type FactionSelectionModel struct {
	renderer    *lipgloss.Renderer
	session     ssh.Session
	database    *sql.DB
	fingerprint string
	canvas      *canvas.Canvas
	broadcaster *canvas.Broadcaster
}

func FactionSelectionModelHandler(sess ssh.Session, db *sql.DB, fingerprint string, c *canvas.Canvas, bc *canvas.Broadcaster) tea.Model {
	renderer := bubbletea.MakeRenderer(sess)

	return FactionCreationModel{
		renderer:    renderer,
		session:     sess,
		database:    db,
		fingerprint: fingerprint,
		canvas:      c,
		broadcaster: bc,
	}
}

func (m FactionSelectionModel) Init() tea.Cmd {
	return nil
}

func (m FactionSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m FactionSelectionModel) View() string {
	return "faction selection"
}
