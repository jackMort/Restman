package results

import (
	"restman/components"
	"restman/components/config"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	NONE = "None"
	TEXT = "Text"
)

var OPTIONS = []string{NONE, TEXT}

type BodyModel struct {
	Body      string
	width     int
	height    int
	textarea  textarea.Model
	toggle    components.ToggleModel
	body_type string
}

func NewBody(body string, width int, height int) BodyModel {

	focusedStyle, blurredStyle := textarea.DefaultStyles()
	focusedStyle.Base = lipgloss.
		NewStyle().
		MarginTop(1).
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

	// TODO: change logic to add param body_type to Call
	default_value := NONE
	if body != "" {
		default_value = TEXT
	}
	toggle := components.NewToggle("Body type", OPTIONS, default_value)

	return BodyModel{
		width:     width,
		height:    height,
		Body:      body,
		textarea:  ti,
		toggle:    toggle,
		body_type: default_value,
	}
}

func (m BodyModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m BodyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case components.OptionSelectedMsg:
		if msg.Id == m.toggle.Id {
			m.body_type = msg.Selected
		}
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
	content := m.textarea.View()
	if m.body_type == NONE {
		content = config.EmptyMessageStyle.Render("No body")
	}

	return lipgloss.
		NewStyle().
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.toggle.View(),
				content,
			),
		)
}
