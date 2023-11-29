package footer

import (
	"restman/app"
	"restman/components/config"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	container = lipgloss.NewStyle().
			BorderForeground(config.COLOR_SUBTLE).
			Background(config.COLOR_SUBTLE).
			Border(lipgloss.NormalBorder()).
			BorderBottom(false).
			BorderTop(false)

	stopwatchStyle = lipgloss.NewStyle().
			Background(config.COLOR_HIGHLIGHT).
			Padding(0, 1)
)

// model represents the properties of the UI.
type model struct {
	stopwatch  stopwatch.Model
	height     int
	width      int
	url        string
	loading    bool
	spinner    spinner.Model
	statusCode int
}

// New creates a new instance of the UI.
func New() model {
	return model{
		stopwatch: stopwatch.NewWithInterval(time.Millisecond),
	}
}

// Init intializes the UI.
func (m model) Init() tea.Cmd {
	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.Line

	return tea.Batch(m.stopwatch.Init(), m.spinner.Tick)
}

// Update handles all UI interactions.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width - 2

	case app.OnLoadingMsg:
		m.url = msg.Url
		m.loading = true
		return m, tea.Batch(m.stopwatch.Reset(), m.stopwatch.Start(), m.spinner.Tick)

	case app.OnResponseMsg:
		m.loading = false
		m.statusCode = msg.Response.StatusCode
		return m, m.stopwatch.Stop()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	}

	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}

// View returns a string representation of the UI.
func (m model) View() string {
	var status string
	var color string
	if m.loading {
		status = "󰞉 LOADING"
		color = "#F59E0B"
	} else if m.statusCode > 0 {
		status = "󰞉 STATUS: " + strconv.Itoa(m.statusCode)
		if m.statusCode >= 200 && m.statusCode < 300 {
			color = "#34D399"
		} else if m.statusCode >= 300 && m.statusCode < 400 {
			color = "#F59E0B"
		} else if m.statusCode >= 400 && m.statusCode < 500 {
			color = "#F97316"
		} else if m.statusCode >= 500 && m.statusCode < 600 {
			color = "#EF4444"
		}
	} else {
		status = " READY TO GO"
		color = "#666666"
	}

	statusToRender := lipgloss.NewStyle().
		Background(lipgloss.Color(color)).
		Foreground(config.COLOR_SUBTLE).
		Padding(0, 1).
		Render(status)

	statusWidth := lipgloss.Width(statusToRender)

	return container.Width(m.width).Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			statusToRender,
			lipgloss.PlaceHorizontal(
				m.width-statusWidth,
				lipgloss.Right,
				stopwatchStyle.Render(" "+m.stopwatch.View()),
				lipgloss.WithWhitespaceBackground(config.COLOR_SUBTLE),
			),
		),
	)

}
