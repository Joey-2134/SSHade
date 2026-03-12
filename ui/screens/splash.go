package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const splashDuration = 2 * time.Second

type SplashDoneMsg struct{}

type SplashModel struct {
	next     tea.Model
	width    int
	height   int
	renderer *lipgloss.Renderer
}

func NewSplashModel(next tea.Model, width, height int, renderer *lipgloss.Renderer) SplashModel {
	return SplashModel{
		next:     next,
		width:    width,
		height:   height,
		renderer: renderer,
	}
}

func (m SplashModel) Init() tea.Cmd {
	return tea.Tick(splashDuration, func(t time.Time) tea.Msg {
		return SplashDoneMsg{}
	})
}

func (m SplashModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SplashDoneMsg:
		return m.next, m.next.Init()
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	return m, nil
}

func (m SplashModel) View() string {
	titleStyle := m.renderer.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)
	taglineStyle := m.renderer.NewStyle().
		Foreground(lipgloss.Color("241"))

	title := titleStyle.Render("SSHade")
	tagline := taglineStyle.Render("place pixels • claim territory")
	content := lipgloss.JoinVertical(lipgloss.Center, title, tagline)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
