package collections

import (
	"restman/components/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Footer struct {
	CancelFocused bool
	CancelText    string
	OkFocused     bool
	OkText        string

	Width int
}

func (f Footer) Init() tea.Cmd {
	return nil
}

func (f Footer) Update(msg tea.Msg) (Footer, tea.Cmd) {
	return f, nil
}

func (f Footer) View() string {
	okButtonStyle := config.ButtonStyle
	cancelButtonStyle := config.ButtonStyle
	if f.CancelFocused {
		cancelButtonStyle = config.ActiveButtonStyle
	} else if f.OkFocused {
		okButtonStyle = config.ActiveButtonStyle
	}

	okButton := okButtonStyle.Render(f.OkText)
	cancelButton := cancelButtonStyle.Render(f.CancelText)

	return lipgloss.PlaceHorizontal(
		f.Width,
		lipgloss.Right,
		lipgloss.JoinHorizontal(lipgloss.Right, cancelButton, " ", okButton),
	)
}
