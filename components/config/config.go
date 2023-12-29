package config

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// App Colors
	COLOR_SUBTLE     = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	COLOR_HIGHLIGHT  = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	COLOR_SPECIAL    = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	COLOR_FOREGROUND = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#ffffff"}
	COLOR_GRAY       = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#666666"}
	COLOR_WHITE      = lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"}
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

var methodColors = map[string]string{
	GET:    "#43BF6D",
	POST:   "#FFB454",
	PUT:    "#F2C94C",
	DELETE: "#F25C54",
}

var BoxHeader = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderBottom(true).
	BorderForeground(COLOR_SUBTLE)

var BoxDescription = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#545454")).
	Italic(true)

var ButtonStyle = lipgloss.NewStyle().
	Foreground(COLOR_FOREGROUND).
	Background(COLOR_GRAY).
	Padding(0, 2)

var ActiveButtonStyle = ButtonStyle.Copy().
	Foreground(COLOR_FOREGROUND).
	Background(COLOR_HIGHLIGHT).
	Underline(true)

var EmptyMessageStyle = lipgloss.NewStyle().
	Padding(2, 0).
	Foreground(COLOR_GRAY)

var MethodStyleShort = lipgloss.NewStyle().
	Bold(false)

var MethodStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(COLOR_FOREGROUND).
	Background(COLOR_HIGHLIGHT).
	Padding(0, 1)

var FullscreenStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(COLOR_HIGHLIGHT).
	Align(lipgloss.Center).
	AlignVertical(lipgloss.Center).
	Padding(1)

var Methods = map[string]string{
	"GET":    MethodStyle.Copy().Background(lipgloss.Color(methodColors["GET"])).Render("GET"),
	"POST":   MethodStyle.Copy().Background(lipgloss.Color(methodColors["POST"])).Render("POST"),
	"PUT":    MethodStyle.Copy().Background(lipgloss.Color(methodColors["PUT"])).Render("PUT"),
	"DELETE": MethodStyle.Copy().Background(lipgloss.Color(methodColors["DELETE"])).Render("DELETE"),
}

var MethodsShort = map[string]string{
	"GET":    MethodStyleShort.Copy().Foreground(lipgloss.Color(methodColors["GET"])).Render("GET"),
	"POST":   MethodStyleShort.Copy().Foreground(lipgloss.Color(methodColors["POST"])).Render("POS"),
	"PUT":    MethodStyleShort.Copy().Foreground(lipgloss.Color(methodColors["PUT"])).Render("PUT"),
	"DELETE": MethodStyleShort.Copy().Foreground(lipgloss.Color(methodColors["DELETE"])).Render("DEL"),
}

type WindowFocusedMsg struct {
	State bool
}
