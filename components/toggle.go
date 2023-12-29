package components

import (
	"restman/components/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var (
	toggleStyle = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 2).
		Foreground(config.COLOR_FOREGROUND).
		Background(config.COLOR_HIGHLIGHT)
)

type ToggleModel struct {
	Id            string
	default_value string
	options       []string
	label         string
	selected      int
}

type OptionSelectedMsg struct {
	Id       string
	Selected string
}

func NewToggle(label string, options []string, default_value string) ToggleModel {

	selected := 0
	for i, option := range options {
		if option == default_value {
			selected = i
			break
		}
	}

	return ToggleModel{
		Id:            zone.NewPrefix(),
		default_value: default_value,
		options:       options,
		selected:      selected,
		label:         label,
	}
}

func (c ToggleModel) Init() tea.Cmd {
	return nil
}

func (c *ToggleModel) Next() tea.Cmd {
	c.selected = (c.selected + 1) % len(c.options)
	return func() tea.Msg {
		return OptionSelectedMsg{c.Id, c.options[c.selected]}
	}
}

func (c ToggleModel) Update(msg tea.Msg) (ToggleModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			if zone.Get(c.Id + "_toggle").InBounds(msg) {
				return c, c.Next()
			}
		}
	}
	return c, nil
}

func (c ToggleModel) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		c.label+": ",
		zone.Mark(c.Id+"_toggle", toggleStyle.Render(c.options[c.selected]+" ")),
	)
}
