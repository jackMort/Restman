package header

import (
	"restman/components/config"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	bg = lipgloss.NewStyle().Background(config.COLOR_SUBTLE)

	container = lipgloss.NewStyle().
			Padding(1).
			Background(config.COLOR_SUBTLE).
			BorderForeground(config.COLOR_SUBTLE)

	title = lipgloss.NewStyle().
		Bold(true).
		PaddingLeft(1).
		Foreground(config.COLOR_HIGHLIGHT)

	version = lipgloss.NewStyle().
		Bold(true).
		Background(config.COLOR_SUBTLE).
		PaddingRight(1)
)

// Bubble represents the properties of the UI.
type Header struct {
	title string
	width int
}

// New creates a new instance of the UI.
func New() Header {
	return Header{
		title: "RESTMAN",
	}
}

// Init intializes the UI.
func (h Header) Init() tea.Cmd {
	return nil
}

// Update handles all UI interactions.
func (h Header) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.width = msg.Width - 2
	}
	return h, nil
}

// View returns a string representation of the UI.
func (h Header) View() string {
	rTitle := title.Render(h.title)
	rVersion := version.Render(config.VERSION)
	wT := lipgloss.Width(rTitle)
	wV := lipgloss.Width(rVersion)

	return container.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Center, rTitle, strings.Repeat(bg.Render(" "), h.width-wT-wV), rVersion,
		),
	)
}
