package url

import (
	"fmt"
	"io"
	"os"
	"restman/app"
	"restman/components/config"
	"restman/utils"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	normal = lipgloss.NewStyle().
		Bold(true).
		Border(lipgloss.NormalBorder()).
		BorderForeground(config.COLOR_SUBTLE).
		PaddingLeft(1)

	focused = lipgloss.NewStyle().
		Bold(true).
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

type url struct {
	focused bool
	width   int
	method  string
	t       textinput.Model
}

func New() url {
	t := textinput.New()
	t.Prompt = ""
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
	case tea.WindowSizeMsg:
		normal.Width(msg.Width - 2)
		focused.Width(msg.Width - 2)
		m.width = msg.Width

		method := methodStyle.Render(" " + m.method + " ")
		send := buttonStyle.Render(" SEND ")

		m.t.Width = m.width - lipgloss.Width(method) - lipgloss.Width(send) - 7
		m.t.Placeholder = "http://google.pl" + strings.Repeat(" ", m.t.Width-15)

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
			app.SetUrl(m.t.Value())

			params := utils.HTTPRequestParams{
				Method:   "GET",
				URL:      "https://api.publicapis.org/entries",
				Username: "u",
				Password: "p",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}
			response, err := utils.MakeRequest(params)
			if err != nil {
				fmt.Println("Error making request:", err)
				os.Exit(1)
			}
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				os.Exit(1)
			}
			app.SetResponse(string(body), response.StatusCode)

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
	if app.Application.SelectedCall != nil {

		m.t.Prompt = app.Application.SelectedCall.Endpoint
	}

	m.t.Width = m.width - lipgloss.Width(method) - lipgloss.Width(send) - 7
	m.t.Placeholder = "http://google.pl" + strings.Repeat(" ", m.t.Width-15-len(m.t.Prompt))

	v := m.t.View()

	return style.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Center, method, " ", v, " ", send,
		),
	)
}
