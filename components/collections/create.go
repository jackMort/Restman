package collections

import (
	"restman/app"
	"restman/components/overlay"
	"restman/components/popup"
	"restman/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const ()

// CreateResultMsg is the message sent when a choice is made.
type CreateResultMsg struct {
	Result bool
}

type SetStepMsg struct {
	Step int
}

// Create is a popup that presents a yes/no choice to the user.
type Create struct {
	overlay      popup.Overlay
	current_step int
	bgRaw        string

	basicInfoStep      BasicInfo
	authenticationStep Authentication
}

func NewCreate(bgRaw string, width int) Create {
	collection := app.NewCollection()
	return Create{
		bgRaw:              bgRaw,
		overlay:            popup.NewOverlay(bgRaw, width, 13),
		current_step:       0,
		basicInfoStep:      NewBasicInfo(&collection),
		authenticationStep: NewAuthentication(&collection),
	}
}

func (c Create) Init() tea.Cmd {
	return nil
}

// Update handles messages.
// func (c Create) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	if msg, ok := msg.(tea.KeyMsg); ok {
// 		switch msg.Type {
// 		case tea.KeyCtrlC, tea.KeyEsc:
// 			return c, c.makeChoice()
// 		case tea.KeyEnter:
// 			if c.focused == c.GetNumberOfInputs()-1 {
// 				if c.current_step < 1 {
// 					c.current_step++
// 					c.focused = 0
// 				}
// 			} else if c.focused == c.GetNumberOfInputs()-2 {
// 				if c.current_step > 0 {
// 					c.current_step--
// 					c.focused = 0
// 				}
// 			}
// 			c.steps.Current = c.current_step
// 			c.content = c.getStepComponent()
// 		case tea.KeyShiftTab, tea.KeyCtrlP:
// 			return c, c.prevInput()
// 		case tea.KeyTab, tea.KeyCtrlN:
// 			return c, c.nextInput()
// 		}
// 	}
// 	// update the current step
// 	content, _ := c.content.Update(msg)
// 	c.content = content.(Step)
//
// 	// update the steps
// 	steps, _ := c.steps.Update(msg)
// 	c.steps = steps.(Steps)
//
// 	return c, nil
// }

func (c Create) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case SetStepMsg:
		c.current_step = msg.Step

	case tea.KeyMsg:

		switch msg.Type {

		// close the popup on escape
		case tea.KeyEsc:
			return c, c.makeChoice()

		}

	}

	// update the current step
	if c.current_step == 0 {
		var cmd tea.Cmd
		c.basicInfoStep, cmd = c.basicInfoStep.Update(msg)
		cmds = append(cmds, cmd)
	}

	// update the current step
	if c.current_step == 1 {
		var cmd tea.Cmd
		c.authenticationStep, cmd = c.authenticationStep.Update(msg)
		cmds = append(cmds, cmd)
	}

	return c, tea.Batch(cmds...)
}

func (c Create) View() string {
	var formView string
	if c.current_step == 0 {
		formView = c.basicInfoStep.View()
	} else if c.current_step == 1 {
		formView = c.authenticationStep.View()
	}

	content := general.Render(
		lipgloss.Place(c.overlay.Width()-2, c.overlay.Height()-2, lipgloss.Left, lipgloss.Top, formView),
	)

	startCol, startRow := utils.GetStartColRow(content, c.bgRaw)
	return overlay.PlaceOverlay(startCol, startRow, content, c.bgRaw)
}

func (c Create) makeChoice() tea.Cmd {
	return func() tea.Msg { return CreateResultMsg{false} }
}
