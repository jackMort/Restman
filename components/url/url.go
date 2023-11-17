package url

import (
	"restman/app"
	"restman/components/config"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

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

type url struct {
	focused     bool
	width       int
	method      string
	t           textinput.Model
	defaultText string
	call        *app.Call
	collection  *app.Collection
}

func New() url {
	t := textinput.New()
	t.PromptStyle = promptStyle
	t.Prompt = "{COLLECTION_BASE_URL}"
	return url{
		t:      t,
		method: GET,
	}
}

func (m url) Init() tea.Cmd {
	return nil
}

func (m url) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case app.CollectionSelectedMsg:
		m.collection = msg.Collection
		if m.collection != nil {
			m.t.Prompt = m.collection.BaseUrl
		} else {
			m.call = nil
		}

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

		m.t.Width = m.width - lipgloss.Width(method) - lipgloss.Width(send) - 7
		m.t.Placeholder = "/some/endpoint" + strings.Repeat(" ", MaxInt(0, m.t.Width-13))

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
			url := m.collection.BaseUrl + m.t.Value()
			return m, app.GetResponse(url)

		case "ctrl+n":
			// cycle over methods
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
	}

	var cmd tea.Cmd
	newModel, cmd := m.t.Update(msg)
	m.t = newModel
	return m, cmd
}

func (m url) View() string {
	style := normal
	if m.focused {
		style = focused
	}
	methodStyle.Background(lipgloss.Color(methodColors[m.method]))
	method := methodStyle.Render(" " + m.method + " ")
	send := buttonStyle.Render(" SEND ")

	m.t.Width = m.width - lipgloss.Width(method) - lipgloss.Width(send) - 7 - len(m.t.Prompt)
	m.t.Placeholder = "/some/endpoint" + strings.Repeat(" ", MaxInt(0, m.t.Width-13))

	v := m.t.View()

	return style.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Center, method, " ", v, " ", send,
		),
	)
}
