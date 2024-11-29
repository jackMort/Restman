package auth

import (
	"restman/app"
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
	call    *app.Call
	width   int
}

func New(width int, call *app.Call) Model {
	var inputs []textinput.Model = make([]textinput.Model, 5)
	inputs[USERNAME_IDX] = textinput.New()
	inputs[USERNAME_IDX].Placeholder = "username"
	inputs[USERNAME_IDX].Focus()
	inputs[USERNAME_IDX].Prompt = "  "
	inputs[USERNAME_IDX].Width = width - 12

	inputs[PASSWORD_IDX] = textinput.New()
	inputs[PASSWORD_IDX].Placeholder = "password"
	inputs[PASSWORD_IDX].Prompt = "󰌆  "
	inputs[PASSWORD_IDX].Width = width - 12

	inputs[TOKEN_IDX] = textinput.New()
	inputs[TOKEN_IDX].Placeholder = "token"
	inputs[TOKEN_IDX].Prompt = "󰌆  "
	inputs[TOKEN_IDX].Width = width - 12

	inputs[API_KEY_IDX] = textinput.New()
	inputs[API_KEY_IDX].Placeholder = "header name"
	inputs[API_KEY_IDX].Focus()
	inputs[API_KEY_IDX].Prompt = "  "
	inputs[API_KEY_IDX].Width = width - 12

	inputs[API_VALUE_IDX] = textinput.New()
	inputs[API_VALUE_IDX].Placeholder = "value"
	inputs[API_VALUE_IDX].Prompt = "󰌆  "
	inputs[API_VALUE_IDX].Width = width - 12

	method := INHERIT
	if call != nil && call.Auth != nil {
		method = call.Auth.Type

		switch call.Auth.Type {
		case "basic_auth":
			inputs[USERNAME_IDX].SetValue(call.Auth.Username)
			inputs[PASSWORD_IDX].SetValue(call.Auth.Password)
		case "bearer_token":
			inputs[TOKEN_IDX].SetValue(call.Auth.Token)
		case "api_key":
			inputs[API_KEY_IDX].SetValue(call.Auth.HeaderName)
			inputs[API_VALUE_IDX].SetValue(call.Auth.HeaderValue)
		}
	}

	return Model{
		inputs:  inputs,
		focused: 0,
		method:  method,
		call:    call,
		width:   width,
	}
}

// Init initializes the popup.
func (c Model) Init() tea.Cmd {
	return textinput.Blink
}

// nextInput focuses the next input field
func (c *Model) nextInput() {
	switch c.method {
	case INHERIT:
		c.focused = 0
	case BASIC_AUTH:
		c.focused = (c.focused + 1) % 2
	case BEARER_TOKEN:
		c.focused = TOKEN_IDX
	case API_KEY:
		if c.focused == API_KEY_IDX {
			c.focused = API_VALUE_IDX
		} else {
			c.focused = API_KEY_IDX
		}
	}
}

// prevInput focuses the previous input field
func (c *Model) prevInput() {
	c.focused--
	// Wrap around
	if c.focused < 0 {
		c.focused = len(c.inputs) - 1
	}
}

func (c *Model) nextMethod() tea.Cmd {
	switch c.method {
	case INHERIT:
		c.method = NONE
		c.focused = 0
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
		c.focused = 0
	}
	if c.call != nil {
		return app.GetInstance().SetCallAuthType(c.call, c.method)
	}
	return nil
}

func (c Model) getKey() string {
	switch c.focused {
	case USERNAME_IDX:
		return "username"
	case PASSWORD_IDX:
		return "password"
	case TOKEN_IDX:
		return "token"
	case API_KEY_IDX:
		return "header_name"
	case API_VALUE_IDX:
		return "header_value"
	}
	return ""
}

// Update handles messages.
func (c Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(c.inputs))
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			if zone.Get("auth_method").InBounds(msg) {
				c.nextMethod()
			}
		}
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {

		case tea.KeyCtrlT:
			// cycle over methods
			return c, c.nextMethod()

		case tea.KeyEnter:
			c.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return c, c.makeChoice()
		case tea.KeyShiftTab, tea.KeyCtrlP:
			c.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			c.nextInput()
		default:
			var cmd tea.Cmd
			c.inputs[c.focused], cmd = c.inputs[c.focused].Update(msg)

			key := c.getKey()
			value := c.inputs[c.focused].Value()

			if c.call != nil {
				return c, tea.Sequence(
					cmd,
					app.GetInstance().SetCallAuthValue(c.call, key, value),
				)
			}
			return c, nil

		}
	}

	for i := range c.inputs {
		c.inputs[i].Blur()
	}
	c.inputs[c.focused].Focus()

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
		if c.call == nil || c.call.Collection() == nil {
			inputs = config.EmptyMessageStyle.Width(c.width - 8).Render("Please save this request in any collection to inherit the authorization")
		} else {

			inputs = config.EmptyMessageStyle.Render("Inherited from collection")
		}
	case BASIC_AUTH:
		inputs = lipgloss.JoinVertical(
			lipgloss.Left,
			" ",
			config.InputStyle.Render(c.inputs[USERNAME_IDX].View()),
			" ",
			config.InputStyle.Render(c.inputs[PASSWORD_IDX].View()),
			" ",
		)
	case NONE:
		inputs = config.EmptyMessageStyle.Render("No authentication")
	case BEARER_TOKEN:
		inputs = lipgloss.JoinVertical(
			lipgloss.Left,
			" ",
			config.InputStyle.Render(c.inputs[TOKEN_IDX].View()),
			" ",
		)
	case API_KEY:
		inputs = lipgloss.JoinVertical(
			lipgloss.Left,
			" ",
			config.InputStyle.Render(c.inputs[API_KEY_IDX].View()),
			" ",
			config.InputStyle.Render(c.inputs[API_VALUE_IDX].View()),
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
