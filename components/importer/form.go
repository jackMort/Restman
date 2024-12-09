package importer

import (
	"restman/components/config"
	"restman/components/overlay"
	"restman/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ImportResultMsg struct {
	Result bool
	Url    string
}

const NUM_OF_INPUTS = 3

const (
	URL_IDX = iota
	CANCEL_IDX
	OK_IDX
)

var (
	general = lipgloss.NewStyle().
		UnsetAlign().
		Padding(0, 1, 0, 1).
		Foreground(config.COLOR_FOREGROUND).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(config.COLOR_HIGHLIGHT)
)

// Create is a popup that presents a yes/no choice to the user.
type Form struct {
	width   int
	focused int
	errors  []string
	input   textinput.Model
	bgRaw   string
}

func NewForm(bgRaw string, width int) Form {
	input := textinput.New()
	input.Placeholder = "https://sampleapi.com/api/v1"
	input.Prompt = "󱞩 "

	return Form{
		input: input,
		bgRaw: bgRaw,
		width: width,
	}
}

// Init initializes the popup.
func (c Form) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages.
func (c Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {

		case tea.KeyShiftTab, tea.KeyCtrlP:
			c.focused = (c.focused - 1 + NUM_OF_INPUTS) % NUM_OF_INPUTS

		case tea.KeyTab, tea.KeyCtrlN:
			c.focused = (c.focused + 1) % NUM_OF_INPUTS

		case tea.KeyEnter:
			if c.focused == OK_IDX {
				return c, func() tea.Msg { return ImportResultMsg{true, c.input.Value()} }
			} else if c.focused == CANCEL_IDX {
				return c, func() tea.Msg { return ImportResultMsg{false, ""} }
			}

		case tea.KeyEsc:
			return c, func() tea.Msg { return ImportResultMsg{false, ""} }

		}
	}

	input, cmd := c.input.Update(msg)
	if c.focused == URL_IDX {
		input.Focus()
	} else {
		input.Blur()
	}
	c.input = input
	return c, cmd

}

func (c Form) View() string {
	inputs := lipgloss.JoinVertical(
		lipgloss.Left,
		config.LabelStyle.Render("URL:"),
		config.InputStyle.Render(c.input.View()),
	)

	okButtonStyle := config.ButtonStyle
	cancelButtonStyle := config.ButtonStyle
	if c.focused == CANCEL_IDX {
		cancelButtonStyle = config.ActiveButtonStyle
	} else if c.focused == OK_IDX {
		okButtonStyle = config.ActiveButtonStyle
	}

	okButton := okButtonStyle.Render("Import")
	cancelButton := cancelButtonStyle.Render("Cancel")

	footer := lipgloss.PlaceHorizontal(
		c.width,
		lipgloss.Right,
		lipgloss.JoinHorizontal(lipgloss.Right, cancelButton, " ", okButton),
	)

	header := lipgloss.JoinVertical(
		lipgloss.Left,
		config.BoxHeader.Render("Import collection"),
	)

	formView := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		inputs,
		utils.RenderErrors(c.errors),
		footer,
	)

	content := general.Render(formView)

	startCol, startRow := utils.GetStartColRow(content, c.bgRaw)
	return overlay.PlaceOverlay(startCol, startRow, content, c.bgRaw)

}
