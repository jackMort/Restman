package collections

import (
	"restman/app"
	"restman/components/config"
	"restman/components/overlay"
	"restman/utils"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

const (
	NAME_IDX = iota
	COLLECTION_IDX
	CLOSE_IDX
	SAVE_IDX
)

// AddToCollectionResultMsg is the message sent when a choice is made.
type AddToCollectionResultMsg struct {
	Result bool
}

// AddToCollection is a popup that presents a yes/no choice to the user.
type AddToCollection struct {
	inputs   []textinput.Model
	focused  int
	save     bool
	errors   []string
	bgRaw    string
	width    int
	startRow int
	startCol int
}

func NewAddToCollection(bgRaw string, width int, vWidth int) AddToCollection {
	var inputs []textinput.Model = make([]textinput.Model, 2)

	inputs[COLLECTION_IDX] = textinput.New()
	inputs[COLLECTION_IDX].Placeholder = "Collection"
	inputs[COLLECTION_IDX].Prompt = "󱞩 "
	inputs[COLLECTION_IDX].ShowSuggestions = true
	inputs[COLLECTION_IDX].KeyMap.AcceptSuggestion = key.NewBinding(
		key.WithKeys("enter"),
	)

	inputs[NAME_IDX] = textinput.New()
	inputs[NAME_IDX].Placeholder = "/hello"
	inputs[NAME_IDX].Prompt = "󱞩 "
	inputs[NAME_IDX].Width = 35

	if app.GetInstance().SelectedCollection != nil {
		inputs[COLLECTION_IDX].SetValue(app.GetInstance().SelectedCollection.Name)
	}

	collections := app.GetInstance().Collections
	suggestions := make([]string, len(collections))
	for i, c := range collections {
		suggestions[i] = c.Name
	}
	inputs[COLLECTION_IDX].SetSuggestions(suggestions)

	return AddToCollection{
		bgRaw:    bgRaw,
		startRow: 3,
		width:    width,
		startCol: vWidth - width - 4,
		inputs:   inputs,
		focused:  NAME_IDX,
	}
}
func (c AddToCollection) Name() string {
	return c.inputs[NAME_IDX].Value()
}

func (c AddToCollection) CollectionName() string {
	return c.inputs[COLLECTION_IDX].Value()
}

func (c AddToCollection) SetUrl(url string) {
	c.inputs[NAME_IDX].SetValue(url)
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
			if c.focused == CLOSE_IDX {
				return c, c.makeChoice()
			} else if c.focused == SAVE_IDX {
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

	c.Validate()

	return c, tea.Batch(cmds...)
}

func (c *AddToCollection) Validate() {
	c.errors = make([]string, 0)
	if c.inputs[COLLECTION_IDX].Value() == "" {
		c.errors = append(c.errors, "Collection is required")
	}
	if c.inputs[NAME_IDX].Value() == "" {
		c.errors = append(c.errors, "Request name is required")
	}
}

// View renders the popup.
func (c AddToCollection) View() string {
	okButtonStyle := config.ButtonStyle
	cancelButtonStyle := config.ButtonStyle
	if c.focused == SAVE_IDX {
		okButtonStyle = config.ActiveButtonStyle
	} else if c.focused == CLOSE_IDX {
		cancelButtonStyle = config.ActiveButtonStyle
	}

	okButton := zone.Mark("add_to_collection_save", okButtonStyle.Render("Save"))
	cancelButton := zone.Mark("add_to_collection_cancel", cancelButtonStyle.Render("Cancel"))

	buttons := lipgloss.PlaceHorizontal(
		c.width,
		lipgloss.Right,
		lipgloss.JoinHorizontal(lipgloss.Right, cancelButton, " ", okButton),
	)

	header := config.BoxHeader.Render("Add to collection")

	inputs := lipgloss.JoinVertical(
		lipgloss.Left,
		config.LabelStyle.Width(30).Render("Request Name:"),
		config.InputStyle.Render(c.inputs[NAME_IDX].View()),
		" ",
		config.LabelStyle.Width(30).Render("Collection:"),
		config.InputStyle.Render(c.inputs[COLLECTION_IDX].View()),
		" ",
		utils.RenderErrors(c.errors),
		buttons,
	)

	ui := lipgloss.JoinVertical(lipgloss.Left, header, " ", inputs)

	content := general.Render(ui)
	return overlay.PlaceOverlay(c.startCol, c.startRow, content, c.bgRaw)
}

func (c AddToCollection) makeChoice() tea.Cmd {
	return func() tea.Msg { return AddToCollectionResultMsg{c.save} }
}
