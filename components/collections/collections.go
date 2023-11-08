package collections

import (
	"restman/components/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	normal = lipgloss.NewStyle().
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(config.COLOR_SUBTLE).
		PaddingLeft(1)

	focused = lipgloss.NewStyle().
		Bold(true).
		Border(lipgloss.RoundedBorder()).
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
}

func New() collections {
	return collections{
		mod: NewModel(),
	}
}

func (m collections) Init() tea.Cmd {
	return nil
}

func (m collections) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case config.WindowFocusedMsg:
		m.focused = msg.State

	case tea.WindowSizeMsg:
		normal.Width(msg.Width - 2)
		normal.Height(msg.Height - 2)
		focused.Width(msg.Width - 2)
		focused.Height(msg.Height - 2)
	}

	var cmd tea.Cmd
	newModel, cmd := m.mod.Update(msg)
	m.mod = newModel
	return m, cmd
}

func (m collections) View() string {
	style := normal
	if m.focused {
		style = focused
	}

	// buttons := lipgloss.JoinHorizontal(lipgloss.Left,
	// 	config.BoxHeader("Collections"),
	// )

	// content := lipgloss.JoinVertical(lipgloss.Left,
	// 	buttons,
	// )

	return style.Render(m.mod.View())
}
