package collections

import (
	"restman/app"
	"restman/components/config"
	"restman/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	TITLE_IDX = iota
	BASE_URL_IDX
	CANCEL_IDX
	OK_IDX
)

var (
	inputStyle = lipgloss.NewStyle().
			Foreground(config.COLOR_FOREGROUND)

	general = lipgloss.NewStyle().
		UnsetAlign().
		Padding(0, 1, 0, 1).
		Foreground(config.COLOR_FOREGROUND).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(config.COLOR_HIGHLIGHT)
)

const NUM_OF_INPUTS = 4

// Create is a popup that presents a yes/no choice to the user.
type BasicInfo struct {
	mode       string
	focused    int
	err        error
	footer     Footer
	header     Header
	errors     []string
	collection *app.Collection
	inputs     []textinput.Model
}

func NewBasicInfo(collection *app.Collection) BasicInfo {
	var inputs []textinput.Model = make([]textinput.Model, 2)
	inputs[TITLE_IDX] = textinput.New()
	inputs[TITLE_IDX].Placeholder = "My Collection"
	inputs[TITLE_IDX].Focus()
	inputs[TITLE_IDX].Prompt = ""
	inputs[TITLE_IDX].SetValue(collection.Name)

	inputs[BASE_URL_IDX] = textinput.New()
	inputs[BASE_URL_IDX].Placeholder = "https://sampleapi.com/api/v1"
	inputs[BASE_URL_IDX].Prompt = ""
	inputs[BASE_URL_IDX].SetValue(collection.BaseUrl)

	mode := "Create"
	if collection.ID != "" {
		mode = "Edit"
	}

	return BasicInfo{
		mode:       mode,
		focused:    0,
		inputs:     inputs,
		collection: collection,
		footer:     Footer{CancelText: "Cancel", OkText: "Next", Width: 70},
	}
}

// Init initializes the popup.
func (c BasicInfo) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages.
func (c BasicInfo) Update(msg tea.Msg) (BasicInfo, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(c.inputs))

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {

		case tea.KeyShiftTab, tea.KeyCtrlP:
			c.focused = (c.focused - 1 + NUM_OF_INPUTS) % NUM_OF_INPUTS

		case tea.KeyTab, tea.KeyCtrlN:
			c.focused = (c.focused + 1) % NUM_OF_INPUTS

		case tea.KeyEnter:
			if c.focused == OK_IDX {
				return c, func() tea.Msg { return SetStepMsg{1} }
			} else if c.focused == CANCEL_IDX {
				return c, func() tea.Msg { return CreateResultMsg{false} }
			}

		}
	}

	// cancel and ok button logic
	c.footer.CancelFocused = false
	c.footer.OkFocused = false

	if c.focused == CANCEL_IDX {
		c.footer.CancelFocused = true
	}

	if c.focused == OK_IDX {
		c.footer.OkFocused = true
	}

	for i := range c.inputs {
		c.inputs[i].Blur()
	}

	if c.focused < len(c.inputs) {
		c.inputs[c.focused].Focus()
	}

	for i := range c.inputs {
		c.inputs[i], cmds[i] = c.inputs[i].Update(msg)
	}

	// fill the values and get possible errors
	c.collection.Name = c.inputs[TITLE_IDX].Value()
	c.collection.BaseUrl = c.inputs[BASE_URL_IDX].Value()
	c.errors = c.collection.ValidatePartial("name", "baseUrl")

	return c, tea.Batch(cmds...)
}

func (c BasicInfo) View() string {
	inputs := lipgloss.JoinVertical(
		lipgloss.Left,
		inputStyle.Width(30).Render("Title:"),
		c.inputs[TITLE_IDX].View(),
		" ",

		inputStyle.Width(30).Render("Base URL:"),
		c.inputs[BASE_URL_IDX].View(),
		" ",
	)

	header := Header{Steps{Current: 0}, c.mode}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header.View(),
		"",
		inputs,
		utils.RenderErrors(c.errors),
		c.footer.View(),
	)

}
