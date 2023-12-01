package url

import (
	"restman/app"
	"restman/components/config"
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

	methodStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(config.COLOR_FOREGROUND).
			Background(config.COLOR_HIGHLIGHT)

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

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type MethodColor struct {
	Method string
	Color  string
}

var methodColors = map[string]string{
	GET:    "#43BF6D",
	POST:   "#FFB454",
	PUT:    "#F2C94C",
	DELETE: "#F25C54",
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
		method:      GET,
		placeholder: "Enter URL",
	}
}

func (m Url) Init() tea.Cmd {
	return nil
}

func (m Url) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case app.CallSelectedMsg:
		m.call = msg.Call
		m.defaultText = m.call.Endpoint
		m.t.SetValue(m.defaultText)

	case tea.WindowSizeMsg:
		normal.Width(msg.Width - 2)
		focused.Width(msg.Width - 2)
		m.width = msg.Width

		method := methodStyle.Render(" " + m.method + " ")
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
			url := m.t.Prompt + m.t.Value()
			return m, app.GetInstance().GetResponse(url)

		case "ctrl+r":
			m.CycleOverMethods()
		}
	}

	newModel, cmd := m.t.Update(msg)
	m.t = newModel
	return m, cmd
}

func (m *Url) SaveToCollection() (tea.Model, tea.Cmd) {
	return m, app.GetInstance().GetAndSaveEndpoint(m.t.Value())
}

func (m *Url) Value() string {
	return m.t.Value()
}

func (m *Url) CycleOverMethods() {
	switch m.method {
	case GET:
		m.method = POST
	case POST:
		m.method = PUT
	case PUT:
		m.method = DELETE
	case DELETE:
		m.method = GET
	}
}

func (m Url) View() string {
	style := normal
	if m.focused {
		style = focused
	}
	methodStyle.Background(lipgloss.Color(methodColors[m.method]))
	method := zone.Mark("method", methodStyle.Render(" "+m.method+" "))
	send := zone.Mark("send", buttonStyle.Render(" SEND "))
	save := zone.Mark("save", "󰐒 ")

	m.t.Width = m.width - lipgloss.Width(method) - lipgloss.Width(send) - 9
	m.t.Placeholder = m.placeholder + strings.Repeat(" ", utils.MaxInt(0, m.t.Width-len(m.placeholder)+1))

	v := zone.Mark("input", m.t.View())

	return style.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Center, method, " ", v, " ", send, " ", save,
		),
	)
}
