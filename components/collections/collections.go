package collections

import (
	"restman/app"
	"restman/components/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var (
	normal = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(config.COLOR_SUBTLE).
		PaddingLeft(1)

	minified = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(config.COLOR_SUBTLE).
			Foreground(config.COLOR_GRAY).
			PaddingLeft(1)

	focused = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(config.COLOR_HIGHLIGHT).
		PaddingLeft(1)

	methodStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(config.COLOR_FOREGROUND).
			Background(config.COLOR_HIGHLIGHT)

	buttonStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(config.COLOR_FOREGROUND).
			Background(config.COLOR_HIGHLIGHT)

	title = lipgloss.NewStyle().
		Bold(true).
		PaddingLeft(1).
		Foreground(config.COLOR_HIGHLIGHT)
)

type Collections struct {
	focused    bool
	minified   bool
	mod        tea.Model
	smod       callModel
	state      app.App
	collection *app.Collection
	width      int
	height     int
}

func New() Collections {
	return Collections{
		minified: true,
		mod:      NewModel(),
		smod:     NewCallModel(),
	}
}

func (m Collections) Init() tea.Cmd {
	return nil
}

func (m Collections) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {

	case app.FetchCollectionsSuccessMsg:
		newModel, cmd := m.mod.Update(msg)
		m.mod = newModel
		cmds = append(cmds, cmd)

	case app.CollectionSelectedMsg:
		m.collection = msg.Collection

	case config.WindowFocusedMsg:
		m.focused = msg.State

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		newSModel, cmd2 := m.smod.Update(msg)
		m.smod = newSModel.(callModel)

		cmds = append(cmds, cmd2)
	}

	if m.collection != nil {
		newSModel, cmd2 := m.smod.Update(msg)
		m.smod = newSModel.(callModel)
		cmds = append(cmds, cmd2)
	} else {
		newModel, cmd := m.mod.Update(msg)
		m.mod = newModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Collections) IsMinified() bool {
	return m.minified
}

func (m Collections) SetMinified(b bool) (tea.Model, tea.Cmd) {
	m.minified = b
	return m.Update(nil)
}

func (m Collections) View() string {
	style := normal
	if m.focused {
		style = focused
	} else if m.minified {
		style = minified
	}

	// set the width and height of the box
	style = style.Width(m.width - 2).Height(m.height - 2)

	// if the collection is minified, render the minified version
	if m.minified {
		return zone.Mark("collections_minified", style.Render(" \n\nC\nO\nL\nL\nE\nC\nT\nI\nO\nN\nS"))
	}

	if m.collection != nil {
		header := config.BoxHeader.MaxWidth(25).
			Render("󰅁 " + m.collection.Name)
		description := config.BoxDescription.MaxWidth(25).
			Render(m.collection.BaseUrl)

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			description,
			m.smod.View(),
		)
		return zone.Mark("collections", style.Render(content))
	}

	return zone.Mark("collections", style.Render(m.mod.View()))
}
