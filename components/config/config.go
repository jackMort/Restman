package config

import (
	"restman/app"

	"github.com/charmbracelet/lipgloss"
)

var (
	VERSION = "v0.0.1"

	// App Colors
	COLOR_SUBTLE     = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	COLOR_HIGHLIGHT  = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	COLOR_SPECIAL    = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	COLOR_FOREGROUND = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#ffffff"}
	COLOR_GRAY       = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#666666"}
)

var BoxHeader = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderBottom(true).
	BorderForeground(COLOR_SUBTLE)

var BoxDescription = lipgloss.NewStyle().
	Foreground(COLOR_SUBTLE).
	Italic(true)

var ButtonStyle = lipgloss.NewStyle().
	Foreground(COLOR_FOREGROUND).
	Background(COLOR_GRAY).
	Padding(0, 2)

var ActiveButtonStyle = ButtonStyle.Copy().
	Foreground(COLOR_FOREGROUND).
	Background(COLOR_HIGHLIGHT).
	Underline(true)

type WindowFocusedMsg struct {
	State bool
}

type AppStateChanged struct {
	State app.App
}
