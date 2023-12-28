package results

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type DataModel struct {
	Body     string
	width    int
	height   int
	textarea textarea.Model
}

func NewData(body string, width int, height int) DataModel {
	ti := textarea.New()
	ti.SetValue(body)
	ti.SetWidth(width)
	ti.SetHeight(height)
	ti.Focus()

	return DataModel{
		width:    width,
		height:   height,
		Body:     body,
		textarea: ti,
	}
}

func (m DataModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m DataModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m DataModel) View() string {
	return m.textarea.View()
}
