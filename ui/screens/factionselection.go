package ui

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/Joey-2134/SSHade/canvas"
	"github.com/Joey-2134/SSHade/constants"
	"github.com/Joey-2134/SSHade/db"
	"github.com/Joey-2134/SSHade/ui/components"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
)

type errModel struct{ msg string }

func (e errModel) Init() tea.Cmd { return nil }
func (e errModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok {
		return e, tea.Quit
	}
	return e, nil
}
func (e errModel) View() string { return e.msg }

func normalizeHex(hex string) string {
	if hex == "" {
		return "#240"
	}
	if hex[0] != '#' {
		return "#" + hex
	}
	return hex
}

func tableStylesWithHeaderColor(r *lipgloss.Renderer, hex string) table.Styles {
	colour := lipgloss.Color(normalizeHex(hex))
	style := table.DefaultStyles()
	style.Header = r.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(colour).
		BorderBottom(true).
		Bold(false).
		Foreground(colour)
	style.Selected = r.NewStyle().Bold(true)
	return style
}

type FactionSelectionModel struct {
	renderer        *lipgloss.Renderer
	session         ssh.Session
	database        *sql.DB
	fingerprint     string
	canvas          *canvas.Canvas
	broadcaster     *canvas.Broadcaster
	factions        []db.Faction
	selectedFaction int
	table           table.Model
	width           int
	height          int
	user            *db.User
}

func FactionSelectionModelHandler(sess ssh.Session, database *sql.DB, user *db.User, fingerprint string, c *canvas.Canvas, bc *canvas.Broadcaster, width, height int) tea.Model {
	if sess == nil {
		return errModel{msg: "session is nil"}
	}
	renderer := bubbletea.MakeRenderer(sess)
	factions, _ := db.GetAllFactions(context.Background(), database)
	rows := make([]table.Row, len(factions))
	for idx, f := range factions {
		rows[idx] = table.Row{strconv.Itoa(idx), f.Name, f.ColourHex}
	}

	columns := []table.Column{
		{Title: "Idx", Width: 5},
		{Title: "Name", Width: 10},
		{Title: "Colour", Width: 10},
	}

	factionTable := table.New(
		table.WithRows(rows),
		table.WithColumns(columns),
		table.WithFocused(true),
	)

	headerHex := "240"
	if len(factions) > 0 && factions[0].ColourHex != "" {
		headerHex = normalizeHex(factions[0].ColourHex)
	}
	factionTable.SetStyles(tableStylesWithHeaderColor(renderer, headerHex))

	return FactionSelectionModel{
		renderer:        renderer,
		session:         sess,
		database:        database,
		fingerprint:     fingerprint,
		canvas:          c,
		broadcaster:     bc,
		factions:        factions,
		selectedFaction: 0,
		table:           factionTable,
		width:           width,
		height:          height,
		user:            user,
	}
}

func (m FactionSelectionModel) Init() tea.Cmd {
	return nil
}

func (m FactionSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.DefaultKeyMap.FactionCreation):
			return FactionCreationModelHandler(m.session, m.database, m.user, m.fingerprint, m.canvas, m.broadcaster, m.width, m.height), nil
		case key.Matches(msg, constants.DefaultKeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, constants.DefaultKeyMap.Enter):
			if len(m.factions) > 0 {
				db.UpdateUserFaction(m.database, m.user.ID, m.factions[m.selectedFaction].ID)
			}
			mainModel, _ := TeaHandler(m.session, m.canvas, m.database, m.broadcaster)
			return mainModel, nil
		default:
			var cmd tea.Cmd
			m.table, cmd = m.table.Update(msg)
			m.selectedFaction = m.table.Cursor()
			if len(m.factions) > 0 && m.selectedFaction < len(m.factions) {
				m.table.SetStyles(tableStylesWithHeaderColor(m.renderer, m.factions[m.selectedFaction].ColourHex))
			}
			return m, cmd
		}
	}
	return m, nil
}

func (m FactionSelectionModel) View() string {

	if len(m.factions) == 0 {
		return components.EmptyFactionsView(m.renderer, m.width, m.height)
	}

	hex := normalizeHex(m.factions[m.selectedFaction].ColourHex)
	borderStyle := m.renderer.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(hex))
	tableContent := borderStyle.Render(m.table.View())

	// Center the table in the terminal
	tableWidth := lipgloss.Width(tableContent)
	tableHeight := lipgloss.Height(tableContent)
	leftPad := (m.width - tableWidth) / 2
	topPad := (m.height - tableHeight) / 2
	if leftPad < 0 {
		leftPad = 0
	}
	if topPad < 0 {
		topPad = 0
	}

	return m.renderer.NewStyle().
		PaddingLeft(leftPad).
		PaddingTop(topPad).
		Render(tableContent)
}
