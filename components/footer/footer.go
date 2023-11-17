package footer

import (
	"restman/app"
	"time"

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
	url       string
	loading   bool
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
		stopwatch: stopwatch.NewWithInterval(time.Millisecond),
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

	case app.OnLoadingMsg:
		m.url = msg.Url
		m.loading = true
		return m, tea.Batch(m.stopwatch.Reset(), m.stopwatch.Start())

	case app.OnResponseMsg:
		m.loading = false
		return m, m.stopwatch.Stop()

	}

	var cmd tea.Cmd

	var status string
	var color string
	if m.loading {
		status = "󰞉 LOADING"
		color = "#F59E0B"
	} else {
		status = "󰞉 STATUS: " + "SUCCESS"
		color = "#34D399"

	}

	m.statusbar.SetColors(
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: color, Dark: color},
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

	m.statusbar.SetContent(status, m.url, "", " "+m.stopwatch.View())
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
