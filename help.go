package main

import (
	"restman/components/config"
	"restman/components/popup"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var general = lipgloss.NewStyle().
	UnsetAlign().
	Padding(1, 2).
	Border(lipgloss.RoundedBorder()).
	Foreground(config.COLOR_FOREGROUND).
	BorderForeground(config.COLOR_HIGHLIGHT)

type Help struct {
	overlay popup.Overlay
	help    help.Model
	keys    config.KeyMap
}

func NewHelp(bgRaw string, width int) Help {
	help := help.New()
	help.ShowAll = true
	return Help{
		help:    help,
		keys:    config.Keys,
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
`
	helpView := c.help.View(c.keys)
	iconStyle := lipgloss.NewStyle().Foreground(config.COLOR_HIGHLIGHT)
	ui := iconStyle.Render(lipgloss.JoinVertical(
		lipgloss.Center,
		icon,
		" https://github.com/jackMort/Restman, version: "+version,
		"",
		general.Width(c.overlay.Width()).Render(helpView)),
	)
	dialog := lipgloss.Place(c.overlay.Width()-2, c.overlay.Height(), lipgloss.Left, lipgloss.Top, ui)

	return c.overlay.WrapView(dialog)
}
