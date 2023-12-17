package help

import (
	"restman/components/config"
	"restman/components/popup"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	general = lipgloss.NewStyle().
		UnsetAlign().
    Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		Foreground(config.COLOR_FOREGROUND).
		BorderForeground(config.COLOR_HIGHLIGHT)
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit},                // second column
	}
}

var keys = keyMap{
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
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("ctr+c", "quit"),
	),
}

type Help struct {
	overlay popup.Overlay
	help    help.Model
	keys    keyMap
}

func NewHelp(bgRaw string, width int) Help {
	help := help.New()
	help.ShowAll = true
	return Help{
		help:    help,
		keys:    keys,
		overlay: popup.NewOverlay(bgRaw, width, 20),
	}
}

func (c Help) Init() tea.Cmd {
	return nil
}

func (c Help) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.help.Width = msg.Width

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return c, func() tea.Msg { return popup.ClosePopupMsg{} }
		}
	}

	return c, tea.Batch(cmds...)
}

func (c Help) View() string {
	icon := `██████╗ ███████╗███████╗████████╗███╗   ███╗ █████╗ ███╗   ██╗
██╔══██╗██╔════╝██╔════╝╚══██╔══╝████╗ ████║██╔══██╗████╗  ██║
██████╔╝█████╗  ███████╗   ██║   ██╔████╔██║███████║██╔██╗ ██║
██╔══██╗██╔══╝  ╚════██║   ██║   ██║╚██╔╝██║██╔══██║██║╚██╗██║
██║  ██║███████╗███████║   ██║   ██║ ╚═╝ ██║██║  ██║██║ ╚████║
╚═╝  ╚═╝╚══════╝╚══════╝   ╚═╝   ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝
                                                        v0.0.2
`
	helpView := c.help.View(c.keys)
	iconStyle := lipgloss.NewStyle().Foreground(config.COLOR_HIGHLIGHT)
	ui := iconStyle.Render(lipgloss.JoinVertical(
		lipgloss.Center,
		icon,
		general.Width(c.overlay.Width()).Render(helpView)),
	)
	dialog := lipgloss.Place(c.overlay.Width()-2, c.overlay.Height(), lipgloss.Left, lipgloss.Top, ui)

	return c.overlay.WrapView(dialog)
}
