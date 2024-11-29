package collections

import (
	"restman/app"
	"restman/components/config"
	"restman/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

const (
	NONE         = "none"
	BASIC_AUTH   = "basic_auth"
	BEARER_TOKEN = "bearer_token"
	API_KEY      = "api_key"
)

type Authentication struct {
	inputs         []textinput.Model
	focused        int
	method         string
	footer         Footer
	numberOfInputs int
	collection     *app.Collection
	errors         []string
	mode           string
}

func NewAuthentication(collection *app.Collection) Authentication {
	mode := "Create"
	okText := "Create"
	if collection.ID != "" {
		mode = "Edit"
		okText = "Save"
	}

	method := NONE
	if collection.Auth != nil {
		method = collection.Auth.Type
	}

	auth := Authentication{
		mode:           mode,
		inputs:         make([]textinput.Model, 0),
		focused:        1,
		method:         method,
		numberOfInputs: 0,
		footer:         Footer{CancelText: "Back", OkText: okText, Width: 70},
		collection:     collection,
	}
	auth.setBasedOnMethod()

	return auth
}

// Init initializes the popup.
func (c Authentication) Init() tea.Cmd {
	return textinput.Blink
}

func (c *Authentication) setBasedOnMethod() {
	username := ""
	password := ""
	headerName := ""
	headerValue := ""
	bearerToken := ""
	if c.collection.Auth != nil {
		username = c.collection.Auth.Username
		password = c.collection.Auth.Password
		headerName = c.collection.Auth.HeaderName
		headerValue = c.collection.Auth.HeaderValue
		bearerToken = c.collection.Auth.Token

	}

	switch c.method {
	case NONE:
		c.inputs = make([]textinput.Model, 0)
		c.focused = 1
	case BASIC_AUTH:
		c.inputs = make([]textinput.Model, 2)
		c.inputs[0] = textinput.New()
		c.inputs[0].Placeholder = "username"
		c.inputs[0].Prompt = "  "
		c.inputs[0].SetValue(username)
		c.inputs[0].PromptStyle = config.InputStyle

		c.inputs[1] = textinput.New()
		c.inputs[1].Placeholder = "password"
		c.inputs[1].Prompt = "󰌆  "
		c.inputs[1].SetValue(password)
		c.focused = 0
	case BEARER_TOKEN:
		c.inputs = make([]textinput.Model, 1)
		c.inputs[0] = textinput.New()
		c.inputs[0].Placeholder = "token"
		c.inputs[0].Prompt = "󰌆  "
		c.inputs[0].SetValue(bearerToken)
		c.focused = 0
	case API_KEY:
		c.inputs = make([]textinput.Model, 2)
		c.inputs[0] = textinput.New()
		c.inputs[0].Placeholder = "header name"
		c.inputs[0].Focus()
		c.inputs[0].Prompt = "  "
		c.inputs[0].SetValue(headerName)

		c.inputs[1] = textinput.New()
		c.inputs[1].Placeholder = "value"
		c.inputs[1].Prompt = "󰌆  "
		c.inputs[1].SetValue(headerValue)
		c.focused = 0
	}
	c.numberOfInputs = len(c.inputs)
}

func (c *Authentication) nextMethod() {
	switch c.method {
	case NONE:
		c.method = BASIC_AUTH
	case BASIC_AUTH:
		c.method = BEARER_TOKEN
	case BEARER_TOKEN:
		c.method = API_KEY
	case API_KEY:
		c.method = NONE
	}
	c.setBasedOnMethod()
}

// Update handles messages.
func (c Authentication) Update(msg tea.Msg) (Authentication, tea.Cmd) {
	numOfInputs := c.numberOfInputs + 2

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {

		case tea.KeyShiftTab, tea.KeyCtrlP:
			c.focused = (c.focused - 1 + numOfInputs) % numOfInputs

		case tea.KeyTab, tea.KeyCtrlN:
			c.focused = (c.focused + 1) % numOfInputs

		case tea.KeyEnter:
			if c.focused == numOfInputs-2 {
				return c, func() tea.Msg { return SetStepMsg{0} }
				// handle save
			} else if c.focused == numOfInputs-1 {
				c.errors = c.collection.ValidatePartial("name", "baseUrl", "auth")
				if len(c.errors) == 0 {
					// TODO: refacor to use this logic in app.GetInstance().SaveCollection()
					//       and get rid of edit/create.go
					if c.mode == "Create" {
						return c, tea.Batch(
							app.GetInstance().CreateCollection(*c.collection),
							func() tea.Msg { return CreateResultMsg{false} })
					} else {
						return c, tea.Batch(
							app.GetInstance().UpdateCollection(*c.collection),
							func() tea.Msg { return CreateResultMsg{false} })
					}
				}
			}

		case tea.KeyCtrlT:
			c.nextMethod()
			numOfInputs = c.numberOfInputs + 2
		}
	}

	// cancel and ok button logic
	c.footer.CancelFocused = false
	c.footer.OkFocused = false

	if c.focused == numOfInputs-2 {
		c.footer.CancelFocused = true
	}

	if c.focused == numOfInputs-1 {
		c.footer.OkFocused = true
	}

	var cmds []tea.Cmd = make([]tea.Cmd, len(c.inputs))
	for i := range c.inputs {
		c.inputs[i].Blur()
	}

	if c.focused < len(c.inputs) {
		c.inputs[c.focused].Focus()

		for i := range c.inputs {
			c.inputs[i], cmds[i] = c.inputs[i].Update(msg)
		}
	}

	// set collection auth value and validate
	c.collection.Auth = &app.Auth{
		Type: c.method,
	}
	if c.method == BASIC_AUTH {
		c.collection.Auth.Username = c.inputs[0].Value()
		c.collection.Auth.Password = c.inputs[1].Value()
	} else if c.method == BEARER_TOKEN {
		c.collection.Auth.Token = c.inputs[0].Value()
	} else if c.method == API_KEY {
		c.collection.Auth.HeaderName = c.inputs[0].Value()
		c.collection.Auth.HeaderValue = c.inputs[1].Value()
	}
	c.errors = c.collection.ValidatePartial("name", "baseUrl", "auth")

	return c, tea.Batch(cmds...)
}

func (c Authentication) GetMethodName() string {
	switch c.method {
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
func (c Authentication) View() string {
	switchView := lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Authentication method: ",
		zone.Mark("auth_method", methodStyle.Padding(0, 1).Render(c.GetMethodName()+" ")),
	)

	var inputs string
	switch c.method {
	case BASIC_AUTH:
		inputs = lipgloss.JoinVertical(
			lipgloss.Left,
			" ",
			config.InputStyle.Render(c.inputs[0].View()),
			" ",
			config.InputStyle.Render(c.inputs[1].View()),
			" ",
		)
	case NONE:
		inputs = config.EmptyMessageStyle.Render("No authentication")
	case BEARER_TOKEN:
		inputs = lipgloss.JoinVertical(
			lipgloss.Left,
			" ",
			config.InputStyle.Render(c.inputs[0].View()),
			" ",
			" ",
			" ",
		)
	case API_KEY:
		inputs = lipgloss.JoinVertical(
			lipgloss.Left,
			" ",
			config.InputStyle.Render(c.inputs[0].View()),
			" ",
			config.InputStyle.Render(c.inputs[1].View()),
			" ",
		)

	}

	header := Header{Steps{Current: 1}, c.mode}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header.View(),
		"",
		switchView,
		inputs,
		utils.RenderErrors(c.errors),
		c.footer.View(),
	)
}
