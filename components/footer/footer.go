package footer

import (
	"restman/components/config"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistakenelf/teacup/statusbar"
)

// model represents the properties of the UI.
type model struct {
	statusbar statusbar.Model
	stopwatch stopwatch.Model
	height    int
}

// New creates a new instance of the UI.
func New() model {
	sb := statusbar.New(
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#F25D94", Dark: "#F25D94"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#3c3836"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#3c3836"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#6124DF", Dark: "#6124DF"},
		},
	)

	return model{
		statusbar: sb,
		stopwatch: stopwatch.New(),
	}
}

// Init intializes the UI.
func (m model) Init() tea.Cmd {
	return m.stopwatch.Init()
}

// Update handles all UI interactions.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.statusbar.SetSize(msg.Width)
	case config.AppStateChanged:
		cmd := m.stopwatch.Start()
		return m, cmd
	}

	var cmd tea.Cmd
	m.statusbar.SetContent("󰞉 STATUS: 200", "https://zippopotamus.us/us/90210", m.stopwatch.View(), " 805 ms")
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}

// View returns a string representation of the UI.
func (m model) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.statusbar.View(),
	)
}
