package footer

import (
	"restman/app"
	"restman/components/config"
	"restman/utils"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	container = lipgloss.NewStyle()
	// BorderForeground(config.COLOR_SUBTLE).
	// Background(config.COLOR_SUBTLE).
	// Border(lipgloss.NormalBorder()).
	// BorderBottom(false).
	// BorderTop(false)

	versionStyle = lipgloss.NewStyle().
			Foreground(config.COLOR_HIGHLIGHT)

	nameStyle = lipgloss.NewStyle().
			Foreground(config.COLOR_HIGHLIGHT).Underline(true)
)

// model represents the properties of the UI.
type model struct {
	stopwatch  stopwatch.Model
	height     int
	width      int
	url        string
	bytes      int64
	loading    bool
	statusCode int
	error      error
}

// New creates a new instance of the UI.
func New() model {
	return model{
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
		m.width = msg.Width

	case app.OnLoadingMsg:
		m.error = nil
		m.url = msg.Call.Url
		m.loading = true
		return m, tea.Sequence(m.stopwatch.Reset(), m.stopwatch.Start())

	case app.OnResponseMsg:
		if msg.Err == nil {
			m.statusCode = msg.Response.StatusCode
		}
		m.bytes = msg.Bytes
		m.error = msg.Err
		m.loading = false
		return m, m.stopwatch.Stop()
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
	} else if m.error != nil {
		status = " ERROR: " + m.error.Error()
		color = "#EF4444"
	} else if m.statusCode > 0 {
		status = "󰞉 STATUS: " + strconv.Itoa(m.statusCode)
		if m.bytes > 0 {
			status = status + "   SIZE: " + utils.ByteCountIEC(m.bytes)
		}

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
		Foreground(lipgloss.Color(color)).
		// Foreground(config.COLOR_SUBTLE).
		Padding(0, 1).
		Render(status)

	statusWidth := lipgloss.Width(statusToRender)

	return container.Width(m.width).Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			statusToRender,
			lipgloss.PlaceHorizontal(
				m.width-statusWidth-1,
				lipgloss.Right,
				nameStyle.Render("Restman")+

					versionStyle.Render(" v."+config.GetVersion()),
			),
		),
	)
}
