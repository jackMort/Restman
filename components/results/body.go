package results

import (
	"restman/components"
	"restman/components/config"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BodyModel struct {
	Body     string
	width    int
	height   int
	textarea textarea.Model
	toggle   components.ToggleModel
}

func NewBody(body string, width int, height int) BodyModel {

	focusedStyle, blurredStyle := textarea.DefaultStyles()
	focusedStyle.Base = lipgloss.
		NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(config.COLOR_SUBTLE).
		BorderRight(false).
		BorderBottom(false).
		BorderTop(false)

	ti := textarea.New()
	ti.SetValue(body)
	ti.SetWidth(width - 4)
	ti.SetHeight(height - 4)
	ti.Focus()
	ti.Prompt = ""
	ti.FocusedStyle = focusedStyle
	ti.BlurredStyle = blurredStyle

	default_value := "None"
	if body != "" {
		default_value = "Text"
	}
	toggle := components.NewToggle("Body type", []string{"None", "Text"}, default_value)

	return BodyModel{
		width:    width,
		height:   height,
		Body:     body,
		textarea: ti,
		toggle:   toggle,
	}
}

func (m BodyModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m BodyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	}

	m.toggle, cmd = m.toggle.Update(msg)
	cmds = append(cmds, cmd)

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m BodyModel) View() string {
	return lipgloss.
		NewStyle().
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.toggle.View(),
				"",
				m.textarea.View(),
			),
		)
}
