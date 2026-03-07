package ui

import (
	"context"
	"database/sql"
	"strings"

	"github.com/Joey-2134/SSHade/canvas"
	"github.com/Joey-2134/SSHade/constants"
	"github.com/Joey-2134/SSHade/db"
	"github.com/charmbracelet/bubbles/key"
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
	factions    []db.Faction
}

func FactionSelectionModelHandler(sess ssh.Session, database *sql.DB, fingerprint string, c *canvas.Canvas, bc *canvas.Broadcaster) tea.Model {
	renderer := bubbletea.MakeRenderer(sess)
	factions, _ := db.GetAllFactions(context.Background(), database)

	return FactionSelectionModel{
		renderer:    renderer,
		session:     sess,
		database:    database,
		fingerprint: fingerprint,
		canvas:      c,
		broadcaster: bc,
		factions:    factions,
	}
}

func (m FactionSelectionModel) Init() tea.Cmd {
	return nil
}

func (m FactionSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.DefaultKeyMap.FactionCreation):
			return FactionCreationModelHandler(m.session, m.database, m.fingerprint, m.canvas, m.broadcaster), nil
		case key.Matches(msg, constants.DefaultKeyMap.Quit):
			return m, tea.Quit
		}
	default:
		return m, nil
	}
	return m, nil
}

func (m FactionSelectionModel) View() string {

	if len(m.factions) == 0 {
		return "No factions found, press c to create a new faction" //TODO replace with styled component
	}

	//TODO replace with styled table component
	names := make([]string, 0, len(m.factions))
	for _, f := range m.factions {
		names = append(names, f.Name)
	}
	return lipgloss.NewStyle().Render(
		"Faction Selection\n",
		"\n",
		"Factions:\n",
		"\n",
		strings.Join(names, "\n"),
		"\n",
		"\n",
		"(q to quit)",
	)
}
