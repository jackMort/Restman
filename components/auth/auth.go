package auth

import (
	"restman/components/config"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

const (
	USERNAME_IDX = iota
	PASSWORD_IDX
	TOKEN_IDX
	API_KEY_IDX
	API_VALUE_IDX
)

const (
	INHERIT      = "inherit"
	NONE         = "none"
	BASIC_AUTH   = "basic_auth"
	BEARER_TOKEN = "bearer_token"
	API_KEY      = "api_key"
)

var (
	inputStyle = lipgloss.NewStyle().
			Foreground(config.COLOR_FOREGROUND).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(config.COLOR_SUBTLE)

	general = lipgloss.NewStyle().
		UnsetAlign().
		Padding(1, 2).
		Foreground(config.COLOR_FOREGROUND)

	methodStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 2).
			Foreground(config.COLOR_FOREGROUND).
			Background(config.COLOR_HIGHLIGHT)
)

type Model struct {
	inputs  []textinput.Model
	focused int
	method  string
}

func New(width int) Model {
	var inputs []textinput.Model = make([]textinput.Model, 5)
	inputs[USERNAME_IDX] = textinput.New()
	inputs[USERNAME_IDX].Placeholder = "username"
	inputs[USERNAME_IDX].Focus()
	inputs[USERNAME_IDX].Prompt = "  "
	inputs[USERNAME_IDX].Width = 40

	inputs[PASSWORD_IDX] = textinput.New()
	inputs[PASSWORD_IDX].Placeholder = "password"
	inputs[PASSWORD_IDX].Prompt = "󰌆  "

	inputs[TOKEN_IDX] = textinput.New()
	inputs[TOKEN_IDX].Placeholder = "token"
	inputs[TOKEN_IDX].Prompt = "󰌆  "

	inputs[API_KEY_IDX] = textinput.New()
	inputs[API_KEY_IDX].Placeholder = "header name"
	inputs[API_KEY_IDX].Focus()
	inputs[API_KEY_IDX].Prompt = "  "

	inputs[API_VALUE_IDX] = textinput.New()
	inputs[API_VALUE_IDX].Placeholder = "value"
	inputs[API_VALUE_IDX].Prompt = "󰌆  "

	return Model{
		inputs:  inputs,
		focused: 0,
		method:  INHERIT,
	}
}

// Init initializes the popup.
func (c Model) Init() tea.Cmd {
	return textinput.Blink
}

// nextInput focuses the next input field
func (c *Model) nextInput() {
	c.focused = (c.focused + 1) % len(c.inputs)
}

// prevInput focuses the previous input field
func (c *Model) prevInput() {
	c.focused--
	// Wrap around
	if c.focused < 0 {
		c.focused = len(c.inputs) - 1
	}
}

func (c *Model) nextMethod() {
	switch c.method {
	case INHERIT:
		c.method = NONE
	case NONE:
		c.method = BASIC_AUTH
		c.focused = USERNAME_IDX
	case BASIC_AUTH:
		c.method = BEARER_TOKEN
		c.focused = TOKEN_IDX
	case BEARER_TOKEN:
		c.method = API_KEY
		c.focused = API_KEY_IDX
	case API_KEY:
		c.method = INHERIT
	}
}

// Update handles messages.
func (c Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(c.inputs))
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			if zone.Get("auth_method").InBounds(msg) {
				c.nextMethod()
				return c, nil
			}
		}
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {

		case tea.KeyCtrlE:
			// cycle over methods
			c.nextMethod()

		case tea.KeyEnter:
			c.nextInput()
			// return c, tea.Batch(
			// 	app.GetInstance().CreateCollection(
			// 		c.inputs[TITLE_IDX].Value(),
			// 		c.inputs[BASE_URL_IDX].Value(),
			// 	),
			// 	c.makeChoice(),
			// )
		case tea.KeyCtrlC, tea.KeyEsc:
			return c, c.makeChoice()
		case tea.KeyShiftTab, tea.KeyCtrlP:
			c.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			c.nextInput()
		}
		for i := range c.inputs {
			c.inputs[i].Blur()
		}
		c.inputs[c.focused].Focus()
	}

	for i := range c.inputs {
		c.inputs[i], cmds[i] = c.inputs[i].Update(msg)
	}

	return c, tea.Batch(cmds...)
}

func (c Model) GetMethodName() string {
	switch c.method {
	case INHERIT:
		return "Inherit"
	case BASIC_AUTH:
		return "Basic Auth"
	case BEARER_TOKEN:
		return "Bearer Token"
	case API_KEY:
		return "API Key"
	}
	return "None"
}

// View renders the popup.
func (c Model) View() string {
	header := lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Authentication method: ",
		zone.Mark("auth_method", methodStyle.Render(c.GetMethodName()+" ")),
	)

	var inputs string
	switch c.method {
	case INHERIT:
		inputs = config.EmptyMessageStyle.Render("Inherited from collection")
	case BASIC_AUTH:
		inputs = lipgloss.JoinVertical(
			lipgloss.Left,
			" ",
			inputStyle.Render(c.inputs[USERNAME_IDX].View()),
			" ",
			inputStyle.Render(c.inputs[PASSWORD_IDX].View()),
			" ",
		)
	case NONE:
		inputs = config.EmptyMessageStyle.Render("No authentication")
	case BEARER_TOKEN:
		inputs = lipgloss.JoinVertical(
			lipgloss.Left,
			" ",
			inputStyle.Render(c.inputs[TOKEN_IDX].View()),
			" ",
		)
	case API_KEY:
		inputs = lipgloss.JoinVertical(
			lipgloss.Left,
			" ",
			inputStyle.Render(c.inputs[API_KEY_IDX].View()),
			" ",
			inputStyle.Render(c.inputs[API_VALUE_IDX].View()),
			" ",
		)

	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		inputs,
	)

	return general.Render(content)
}

// CreateResultMsg is the message sent when a choice is made.
type CreateResultMsg struct {
	Result bool
}

func (c Model) makeChoice() tea.Cmd {
	return func() tea.Msg { return CreateResultMsg{false} }
}
