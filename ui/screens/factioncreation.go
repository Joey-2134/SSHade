package ui

import (
	"context"
	"database/sql"
	"strings"

	"github.com/Joey-2134/SSHade/canvas"
	"github.com/Joey-2134/SSHade/constants"
	"github.com/Joey-2134/SSHade/db"
	"github.com/Joey-2134/SSHade/ui/components"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
)

type FactionCreationModel struct {
	renderer    *lipgloss.Renderer
	session     ssh.Session
	database    *sql.DB
	fingerprint string
	canvas      *canvas.Canvas
	broadcaster *canvas.Broadcaster
	nameInput   textinput.Model
	colourInput textinput.Model
	focus       int // 0 = name, 1 = colour
	err         string
}

func FactionCreationModelHandler(sess ssh.Session, database *sql.DB, fingerprint string, c *canvas.Canvas, bc *canvas.Broadcaster) tea.Model {
	renderer := bubbletea.MakeRenderer(sess)

	style := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	nameInput := textinput.New()
	nameInput.Placeholder = "Faction name"
	nameInput.Focus()
	nameInput.CharLimit = 32
	nameInput.Width = 32
	nameInput.PromptStyle = style
	nameInput.TextStyle = style
	nameInput.PlaceholderStyle = style

	colourInput := textinput.New()
	colourInput.Placeholder = "#rrggbb (e.g. #ff0000)"
	colourInput.CharLimit = 7
	colourInput.Width = 32
	colourInput.PromptStyle = style
	colourInput.TextStyle = style
	colourInput.PlaceholderStyle = style

	return FactionCreationModel{
		renderer:    renderer,
		session:     sess,
		database:    database,
		fingerprint: fingerprint,
		canvas:      c,
		broadcaster: bc,
		nameInput:   nameInput,
		colourInput: colourInput,
		focus:       0,
	}
}

func (m FactionCreationModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m FactionCreationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.DefaultKeyMap.Quit):
			return FactionSelectionModelHandler(m.session, m.database, m.fingerprint, m.canvas, m.broadcaster), nil
		case key.Matches(msg, constants.DefaultKeyMap.Enter):
			if m.focus == 0 {
				m.nameInput.Blur()
				m.colourInput.Focus()
				m.focus = 1
				return m, textinput.Blink
			}
			// Submit
			name := strings.TrimSpace(m.nameInput.Value())
			colour := strings.TrimSpace(m.colourInput.Value())
			if name == "" {
				m.err = "Name is required"
				return m, nil
			}
			if colour == "" {
				m.err = "Colour is required"
				return m, nil
			}
			if colour[0] != '#' {
				colour = "#" + colour
			}
			_, err := db.CreateFaction(context.Background(), m.database, name, colour)
			if err != nil {
				m.err = err.Error()
				return m, nil
			}
			return FactionSelectionModelHandler(m.session, m.database, m.fingerprint, m.canvas, m.broadcaster), nil
		}
	}

	if m.focus == 0 {
		m.nameInput, cmd = m.nameInput.Update(msg)
	} else {
		m.colourInput, cmd = m.colourInput.Update(msg)
	}
	return m, cmd
}

func (m FactionCreationModel) View() string {
	return m.renderer.NewStyle().Render(
		components.FactionCreationForm(m.nameInput.View(), m.colourInput.View(), m.err),
	)
}
