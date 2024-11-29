package config

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

var version = "dev"

var (
	// App Colors
	COLOR_SUBTLE     = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	COLOR_HIGHLIGHT  = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	COLOR_SPECIAL    = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	COLOR_FOREGROUND = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#ffffff"}
	COLOR_GRAY       = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#666666"}
	COLOR_WHITE      = lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"}
	COLOR_LINK       = lipgloss.AdaptiveColor{Light: "#6C9EF8", Dark: "#6C9EF8"}
	COLOR_ERROR      = lipgloss.AdaptiveColor{Light: "#F25C54", Dark: "#F25C54"}
	COLOR_LIGHTER    = lipgloss.AdaptiveColor{Light: "#9f9f9f", Dark: "#9f9f9f"}
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

var ActiveButtonStyle = ButtonStyle.
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

var LinkStyle = lipgloss.NewStyle().
	Bold(true).
	Underline(true).
	Foreground(COLOR_LINK)

var ErrorStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(COLOR_ERROR)

var InputStyle = lipgloss.NewStyle().
	BorderForeground(COLOR_SUBTLE).
	Foreground(COLOR_FOREGROUND)

var LabelStyle = lipgloss.NewStyle().
	Foreground(COLOR_LIGHTER)

var Methods = map[string]string{
	"GET":    MethodStyle.Background(lipgloss.Color(methodColors["GET"])).Render("GET"),
	"POST":   MethodStyle.Background(lipgloss.Color(methodColors["POST"])).Render("POST"),
	"PUT":    MethodStyle.Background(lipgloss.Color(methodColors["PUT"])).Render("PUT"),
	"DELETE": MethodStyle.Background(lipgloss.Color(methodColors["DELETE"])).Render("DELETE"),
}

var MethodsShort = map[string]string{
	"GET":    MethodStyleShort.Foreground(lipgloss.Color(methodColors["GET"])).Render("GET"),
	"POST":   MethodStyleShort.Foreground(lipgloss.Color(methodColors["POST"])).Render("POS"),
	"PUT":    MethodStyleShort.Foreground(lipgloss.Color(methodColors["PUT"])).Render("PUT"),
	"DELETE": MethodStyleShort.Foreground(lipgloss.Color(methodColors["DELETE"])).Render("DEL"),
}

type WindowFocusedMsg struct {
	State bool
}

type KeyMap struct {
	Up                key.Binding
	Down              key.Binding
	Left              key.Binding
	Right             key.Binding
	Help              key.Binding
	Quit              key.Binding
	NewCollection     key.Binding
	ChangeActivePanel key.Binding
	Save              key.Binding
	ChangeToggle      key.Binding
}

func SetVersion(v string) {
	version = v
}

func GetVersion() string {
	return version
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right, k.ChangeActivePanel, k.Help, k.Quit},
		{k.NewCollection, k.Save, k.ChangeToggle},
	}
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("ctrl+h", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctr+c", "quit"),
	),
	NewCollection: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "new collection"),
	),
	ChangeActivePanel: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "change panel"),
	),
	Save: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "save"),
	),
	ChangeToggle: key.NewBinding(
		key.WithKeys("ctrl+t"),
		key.WithHelp("ctrl+t", "change toggle"),
	),
}
