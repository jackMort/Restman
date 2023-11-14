package popup

import (
	"restman/components/config"

	"github.com/charmbracelet/lipgloss"
)

// style is the style of the choice popup
type style struct {
	button       lipgloss.Style
	activeButton lipgloss.Style
	question     lipgloss.Style
	general      lipgloss.Style
}

// newStyle creates a new style for the choice popup
func newStyle(width, height int) style {
	buttonStyle := lipgloss.NewStyle().
		Foreground(config.COLOR_FOREGROUND).
		Background(config.COLOR_GRAY).
		Padding(0, 2).
		Margin(0, 1)

	activeButtonStyle := buttonStyle.Copy().
		Foreground(config.COLOR_FOREGROUND).
		Background(config.COLOR_HIGHLIGHT).
		Underline(true)

	general := lipgloss.NewStyle().
		Foreground(config.COLOR_FOREGROUND).
		Width(width - 2).
		Height(height - 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(config.COLOR_HIGHLIGHT)

	question := lipgloss.NewStyle().
		Width(width).
		Margin(1, 0).
		Bold(true).
		Align(lipgloss.Center)

	return style{
		button:       buttonStyle,
		activeButton: activeButtonStyle,
		question:     question,
		general:      general,
	}
}
