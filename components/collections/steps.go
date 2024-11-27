package collections

import (
	"restman/components/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	stepStyle = lipgloss.NewStyle().
			Foreground(config.COLOR_SUBTLE).Italic(true)
	currentStepStyle = lipgloss.NewStyle().
				Foreground(config.COLOR_HIGHLIGHT).Italic(true)
	spacer = lipgloss.NewStyle().
		Foreground(config.COLOR_SUBTLE)
)

var STEPS = []string{"󰲠 Basic Info", "󰲢 Authentication"} // "󰲤 Variables"}

type Steps struct {
	Current int
}

func (s Steps) Init() tea.Cmd {
	return nil
}

func (s Steps) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, nil
}

func (s Steps) View() string {
	toRender := []string{}
	for i, step := range STEPS {
		style := stepStyle
		if i == s.Current {
			style = currentStepStyle
		}
		toRender = append(toRender, style.Render(step))
		if i < len(STEPS)-1 {
			toRender = append(toRender, spacer.Render(" ─── "))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, toRender...)
}
