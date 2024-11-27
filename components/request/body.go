package request

import (
	"os"
	"restman/app"
	"restman/components"
	"restman/components/config"
	"restman/utils"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	NONE = "None"
	TEXT = "Text"
	JSON = "JSON"
)

type editorFinishedMsg struct {
	file *os.File
	err  error
}

var OPTIONS = []string{NONE, TEXT, JSON}

type BodyModel struct {
	call     *app.Call
	width    int
	height   int
	textarea textarea.Model
	toggle   components.ToggleModel
}

func NewBody(call *app.Call, width int, height int) BodyModel {

	focusedStyle, blurredStyle := textarea.DefaultStyles()
	focusedStyle.Base = lipgloss.
		NewStyle().
		MarginTop(1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(config.COLOR_SUBTLE).
		BorderRight(false).
		BorderBottom(false).
		BorderTop(false)

	data := ""
	if call != nil {
		data = call.Data
	}

	ti := textarea.New()
	ti.CharLimit = 0
	ti.SetValue(data)
	ti.SetWidth(width - 4)
	ti.SetHeight(height - 4)
	ti.Focus()
	ti.Prompt = ""
	ti.FocusedStyle = focusedStyle
	ti.BlurredStyle = blurredStyle

	defaultValue := NONE
	if call != nil {
		if call.DataType != "" {
			defaultValue = call.DataType
		}
	}
	toggle := components.NewToggle("Body type", OPTIONS, defaultValue)

	return BodyModel{
		width:    width,
		height:   height,
		call:     call,
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

	case editorFinishedMsg:
		// TODO: handle error
		content, _ := os.ReadFile(msg.file.Name())
		m.textarea.SetValue(string(content))
		utils.RemoveTempFile(msg.file)

	case components.OptionSelectedMsg:
		if msg.Id == m.toggle.Id {
			if m.call != nil {
				m.call.DataType = msg.Selected
				m.call.Data = ""
				m.textarea.SetValue("")
			}
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlT:
			return m, m.toggle.Next()

		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlE:
			extension := "txt"
			if m.call != nil && m.call.DataType == JSON {
				extension = "json"
			}
			// TODO: handle error
			tmpFile, _ := utils.CreateTempFile(m.textarea.Value(), extension)
			return m, tea.ExecProcess(utils.OpenInEditorCommand(tmpFile), func(err error) tea.Msg {
				return editorFinishedMsg{tmpFile, err}
			})

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
	if m.call == nil || m.call.DataType == NONE {
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
