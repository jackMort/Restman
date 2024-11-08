package collections

import (
	"restman/app"
	"restman/components/config"
	"restman/components/overlay"
	"restman/components/popup"
	"restman/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	TITLE_IDX = iota
	BASE_URL_IDX
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

// CreateResultMsg is the message sent when a choice is made.
type CreateResultMsg struct {
	Result bool
}

// Create is a popup that presents a yes/no choice to the user.
type Create struct {
	overlay popup.Overlay
	inputs  []textinput.Model
	focused int
	err     error
	bgRaw   string
}

func NewCreate(bgRaw string, width int) Create {
	var inputs []textinput.Model = make([]textinput.Model, 2)
	inputs[TITLE_IDX] = textinput.New()
	inputs[TITLE_IDX].Placeholder = "My Collection"
	inputs[TITLE_IDX].Focus()
	inputs[TITLE_IDX].Prompt = ""

	inputs[BASE_URL_IDX] = textinput.New()
	inputs[BASE_URL_IDX].Placeholder = "https://sampleapi.com/api/v1"
	inputs[BASE_URL_IDX].Prompt = ""

	return Create{
		bgRaw:   bgRaw,
		overlay: popup.NewOverlay(bgRaw, width, 13),
		inputs:  inputs,
		focused: 0,
	}
}

// Init initializes the popup.
func (c Create) Init() tea.Cmd {
	return textinput.Blink
}

// nextInput focuses the next input field
func (c *Create) nextInput() {
	c.focused = (c.focused + 1) % (len(c.inputs) + 2)
}

// prevInput focuses the previous input field
func (c *Create) prevInput() {
	c.focused--
	// Wrap around
	if c.focused < 0 {
		c.focused = len(c.inputs) + 1
	}
}

// Update handles messages.
func (c Create) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(c.inputs))

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {
		case tea.KeyEnter:
			if c.focused == 2 {
				return c, tea.Batch(
					app.GetInstance().CreateCollection(
						c.inputs[TITLE_IDX].Value(),
						c.inputs[BASE_URL_IDX].Value(),
					),
					c.makeChoice(),
				)
			} else if c.focused == 3 {
				return c, c.makeChoice()
			} else {
				c.nextInput()
			}
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
		if c.focused < len(c.inputs) {
			c.inputs[c.focused].Focus()
		}
	}

	for i := range c.inputs {
		c.inputs[i], cmds[i] = c.inputs[i].Update(msg)
	}

	return c, tea.Batch(cmds...)
}

// View renders the popup.
func (c Create) View() string {
	okButtonStyle := config.ButtonStyle
	cancelButtonStyle := config.ButtonStyle
	if c.focused == 2 {
		okButtonStyle = config.ActiveButtonStyle
	} else if c.focused == 3 {
		cancelButtonStyle = config.ActiveButtonStyle
	}

	okButton := okButtonStyle.Render("Save")
	cancelButton := cancelButtonStyle.Render("Cancel")
	buttons := lipgloss.PlaceHorizontal(
		c.overlay.Width(),
		lipgloss.Right,
		lipgloss.JoinHorizontal(lipgloss.Right, okButton, " ", cancelButton),
	)

	header := config.BoxHeader.Render("Create collection")

	inputs := lipgloss.JoinVertical(
		lipgloss.Left,

		inputStyle.Width(30).Render("Title:"),
		c.inputs[TITLE_IDX].View(),
		" ",

		inputStyle.Width(30).Render("Base URL:"),
		c.inputs[BASE_URL_IDX].View(),
		" ",
		" ",
		buttons,
	)

	ui := lipgloss.JoinVertical(lipgloss.Left, header, " ", inputs)
	dialog := lipgloss.Place(c.overlay.Width()-2, c.overlay.Height()-2, lipgloss.Left, lipgloss.Top, ui)

	content := general.Render(dialog)

	startCol, startRow := utils.GetStartColRow(content, c.bgRaw)
	return overlay.PlaceOverlay(startCol, startRow, content, c.bgRaw)
}

func (c Create) makeChoice() tea.Cmd {
	return func() tea.Msg { return CreateResultMsg{false} }
}
