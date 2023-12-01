package help

import (
	"restman/components/config"
	"restman/components/popup"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	general = lipgloss.NewStyle().
		UnsetAlign().
		Padding(0, 1, 0, 1).
		Foreground(config.COLOR_FOREGROUND).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(config.COLOR_HIGHLIGHT)
)

type Help struct {
	overlay popup.Overlay
}

func NewHelp(bgRaw string, width int) Help {
	return Help{
		overlay: popup.NewOverlay(bgRaw, width, 13),
	}
}

func (c Help) Init() tea.Cmd {
	return nil
}

func (c Help) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
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
	iconStyle := lipgloss.NewStyle().Foreground(config.COLOR_SUBTLE)
	ui := iconStyle.Render(lipgloss.JoinVertical(lipgloss.Left, icon, ""))
	dialog := lipgloss.Place(c.overlay.Width()-2, c.overlay.Height()-2, lipgloss.Left, lipgloss.Top, ui)

	return c.overlay.WrapView(general.Render(dialog))
}
