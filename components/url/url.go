package url

import (
	"restman/app"
	"restman/components/config"
	"restman/components/tabs"
	"restman/utils"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
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

	buttonStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(config.COLOR_FOREGROUND).
			Background(config.COLOR_HIGHLIGHT)

	promptStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(config.COLOR_HIGHLIGHT)

	title = lipgloss.NewStyle().
		Bold(true).
		PaddingLeft(1).
		Foreground(config.COLOR_HIGHLIGHT)
)

type MethodColor struct {
	Method string
	Color  string
}

type Url struct {
	placeholder string
	focused     bool
	width       int
	method      string
	t           textinput.Model
	defaultText string
	call        *app.Call
	collection  *app.Collection
}

func New() Url {
	t := textinput.New()
	t.PromptStyle = promptStyle
	t.Prompt = ""
	return Url{
		t:           t,
		method:      config.GET,
		placeholder: "Enter URL",
	}
}

func (m Url) Url() string {
	return m.t.Value()
}

func (m Url) Method() string {
	return m.method
}

func (m Url) Init() tea.Cmd {
	return nil
}

func (m Url) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tabs.TabFocusedMsg:
		if msg.Tab.Call != nil {
			m.call = msg.Tab.Call
			m.defaultText = m.call.Url
			m.t.SetValue(m.defaultText)
			m.method = m.call.Method
		} else {
			m.call = nil
			m.defaultText = ""
			m.t.SetValue(m.defaultText)
		}

	case tea.WindowSizeMsg:
		normal.Width(msg.Width - 2)
		focused.Width(msg.Width - 2)
		m.width = msg.Width

		method := config.Methods[m.method]
		send := buttonStyle.Render(" SEND ")

		m.t.Width = m.width - lipgloss.Width(method) - lipgloss.Width(send) - 10
		m.t.Placeholder = m.placeholder + strings.Repeat(" ", utils.MaxInt(0, m.t.Width-len(m.placeholder)))

	case config.WindowFocusedMsg:
		m.focused = msg.State
		if m.focused {
			m.t.Focus()
		} else {
			m.t.Blur()
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			call := m.call
			if m.call != nil {
				call.Url = m.t.Prompt + m.t.Value()
				call.Method = m.method
				return m, app.GetInstance().GetResponse(call)
			}

		case "ctrl+r":
			m.CycleOverMethods()
		}
	}

	newModel, cmd := m.t.Update(msg)
	m.t = newModel
	return m, cmd
}

func (m *Url) Value() string {
	return m.t.Value()
}

func (m *Url) CycleOverMethods() {
	switch m.method {
	case config.GET:
		m.method = config.POST
	case config.POST:
		m.method = config.PUT
	case config.PUT:
		m.method = config.DELETE
	case config.DELETE:
		m.method = config.GET
	}
}

func (m Url) View() string {
	style := normal
	if m.focused {
		style = focused
	}
	method := zone.Mark("method", config.Methods[m.method])
	send := zone.Mark("send", buttonStyle.Render(" SEND "))
	save := zone.Mark("save", " ")

	m.t.Width = m.width - lipgloss.Width(method) - lipgloss.Width(send) - 9
	m.t.Placeholder = m.placeholder + strings.Repeat(" ", utils.MaxInt(0, m.t.Width-len(m.placeholder)+1))

	v := zone.Mark("input", m.t.View())

	return style.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Center, method, " ", v, " ", send, " ", save,
		),
	)
}
