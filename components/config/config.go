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
)

var BoxHeader = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderBottom(true).
	BorderForeground(COLOR_SUBTLE).
	Render

type WindowFocusedMsg struct {
	State bool
}

type AppStateChanged struct {
	State app.App
}
