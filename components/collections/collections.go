package collections

import (
	"restman/app"
	"restman/components/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	normal = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(config.COLOR_SUBTLE).
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

type collections struct {
	focused bool
	mod     tea.Model
	smod    callModel
	state   app.App
}

func New() collections {
	return collections{
		mod:  NewModel(),
		smod: NewCallModel(),
	}
}

func (m collections) Init() tea.Cmd {
	return nil
}

func (m collections) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case config.WindowFocusedMsg:
		m.focused = msg.State

	case tea.WindowSizeMsg:
		normal.Width(msg.Width - 2)
		normal.Height(msg.Height - 2)
		focused.Width(msg.Width - 2)
		focused.Height(msg.Height - 2)

		newSModel, cmd2 := m.smod.Update(msg)
		m.smod = newSModel.(callModel)

		cmds = append(cmds, cmd2)
	}

	if app.Application.SelectedCollection != nil {
		cmd3 := m.smod.RefreshItems()
		newSModel, cmd2 := m.smod.Update(msg)
		m.smod = newSModel.(callModel)

		cmds = append(cmds, cmd2)
		cmds = append(cmds, cmd3)
	} else {
		newModel, cmd := m.mod.Update(msg)
		m.mod = newModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m collections) View() string {
	style := normal
	if m.focused {
		style = focused
	}

	// buttons := lipgloss.JoinHorizontal(lipgloss.Left,
	// )

	if app.Application.SelectedCollection != nil {
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			config.BoxHeader("󰅁 "+app.Application.SelectedCollection.Name),
			config.BoxDescription(app.Application.SelectedCollection.BaseUrl),
			m.smod.View(),
		)
		return style.Render(content)
	}

	return style.Render(m.mod.View())
}
