package collections

import (
	"restman/components/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Header struct {
	steps Steps
	mode  string //Create or Edit
}

func (h Header) Init() tea.Cmd {
	return nil
}

func (h Header) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return h, nil
}

func (h Header) View() string {
	icon := ""
	if h.mode == "edit" {
		icon = "󰷎"
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		config.BoxHeader.Render(icon+" "+h.mode+" collection"),
		h.steps.View(),
	)
}
