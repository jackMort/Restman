package collections

import (
	"restman/app"
	"restman/components/config"
	"restman/components/popup"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

const (
	COLLECTION_IDX = iota
	URL_IDX
)

// AddToCollectionResultMsg is the message sent when a choice is made.
type AddToCollectionResultMsg struct {
	Result bool
}

// AddToCollection is a popup that presents a yes/no choice to the user.
type AddToCollection struct {
	overlay popup.Overlay
	inputs  []textinput.Model
	focused int
	save    bool
	err     error
}

func NewAddToCollection(bgRaw string, width int, vWidth int) AddToCollection {
	var inputs []textinput.Model = make([]textinput.Model, 2)

	inputs[COLLECTION_IDX] = textinput.New()
	inputs[COLLECTION_IDX].Placeholder = "Collection"
	inputs[COLLECTION_IDX].Prompt = ""
	inputs[COLLECTION_IDX].ShowSuggestions = true
	inputs[COLLECTION_IDX].KeyMap.AcceptSuggestion = key.NewBinding(
		key.WithKeys("enter"),
	)

	inputs[URL_IDX] = textinput.New()
	inputs[URL_IDX].Placeholder = "https://sampleapi.com/api/v1"
	inputs[URL_IDX].Prompt = ""

	if app.GetInstance().SelectedCollection != nil {
		inputs[COLLECTION_IDX].SetValue(app.GetInstance().SelectedCollection.Name)
		inputs[URL_IDX].Focus()
	} else {
		inputs[COLLECTION_IDX].Focus()
	}
	collections := app.GetInstance().Collections
	suggestions := make([]string, len(collections))
	for i, c := range collections {
		suggestions[i] = c.Name
	}
	inputs[COLLECTION_IDX].SetSuggestions(suggestions)

	return AddToCollection{
		overlay: popup.NewOverlayOnPosition(bgRaw, width, 13, 5, vWidth-width-4),
		inputs:  inputs,
		focused: 0,
	}
}

func (c AddToCollection) CollectionName() string {
	return c.inputs[COLLECTION_IDX].Value()
}

func (c AddToCollection) SetUrl(url string) {
	c.inputs[URL_IDX].SetValue(url)
}

// Init initializes the popup.
func (c AddToCollection) Init() tea.Cmd {
	return textinput.Blink
}

// nextInput focuses the next input field
func (c *AddToCollection) nextInput() {
	c.focused = (c.focused + 1) % (len(c.inputs) + 2)
}

// prevInput focuses the previous input field
func (c *AddToCollection) prevInput() {
	c.focused--
	// Wrap around
	if c.focused < 0 {
		c.focused = len(c.inputs) + 1
	}
}

// Update handles messages.
func (c AddToCollection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(c.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if c.focused == 3 {
				return c, c.makeChoice()
			} else if c.focused == 2 {
				c.save = true
				return c, c.makeChoice()
			} 
		case tea.KeyCtrlC, tea.KeyEsc:
			return c, c.makeChoice()
		case tea.KeyShiftTab, tea.KeyCtrlK:
			c.prevInput()
		case tea.KeyTab, tea.KeyCtrlJ:
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
func (c AddToCollection) View() string {
	okButtonStyle := config.ButtonStyle
	cancelButtonStyle := config.ButtonStyle
	if c.focused == 2 {
		okButtonStyle = config.ActiveButtonStyle
	} else if c.focused == 3 {
		cancelButtonStyle = config.ActiveButtonStyle
	}

	okButton := zone.Mark("add_to_collection_save", okButtonStyle.Render("Save"))
	cancelButton := zone.Mark("add_to_collection_cancel", cancelButtonStyle.Render("Cancel"))

	buttons := lipgloss.PlaceHorizontal(
		c.overlay.Width(),
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Right, okButton, " ", cancelButton),
	)

	header := config.BoxHeader.Render("Add to collection")

	inputs := lipgloss.JoinVertical(
		lipgloss.Left,
		inputStyle.Width(30).Render("Collection:"),
		c.inputs[COLLECTION_IDX].View(),
		" ",

		inputStyle.Width(30).Render("URL:"),
		c.inputs[URL_IDX].View(),
		" ",
		" ",
		buttons,
	)

	ui := lipgloss.JoinVertical(lipgloss.Left, header, " ", inputs)
	dialog := lipgloss.Place(c.overlay.Width()-2, c.overlay.Height()-2, lipgloss.Left, lipgloss.Top, ui)

	return c.overlay.WrapView(general.Render(dialog))
}

func (c AddToCollection) makeChoice() tea.Cmd {
	return func() tea.Msg { return AddToCollectionResultMsg{c.save} }
}
