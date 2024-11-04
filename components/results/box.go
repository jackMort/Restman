package results

import (
	"encoding/json"
	"restman/app"
	"restman/components/config"
	"restman/utils"
	"strconv"
	"strings"

	"github.com/TylerBrock/colorjson"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	inactiveStyle = lipgloss.NewStyle().BorderForeground(config.COLOR_SUBTLE).Border(lipgloss.NormalBorder())
	activeStyle   = lipgloss.NewStyle().BorderForeground(config.COLOR_HIGHLIGHT).Border(lipgloss.NormalBorder())
	emptyMessage  = lipgloss.NewStyle().Padding(2, 2).Foreground(config.COLOR_GRAY)
	statusStyle   = lipgloss.NewStyle().Padding(0, 1).Background(config.COLOR_GRAY)
)

type Results struct {
	title     string
	focused   bool
	body      string
	width     int
	height    int
	viewport  viewport.Model
	Tabs      []string
	activeTab int
	content   tea.Model
	call      *app.Call
	status    int
}

func New() Results {
	return Results{
		title: "Results",
		Tabs:  []string{"Response", "Headers", "Cookies", "Statistics"},
	}
}

// satisfy the tea.Model interface
func (b Results) Init() tea.Cmd {
	b.viewport = viewport.New(10, 10)
	b.activeTab = 0
	return nil
}

func (b Results) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case app.CallSelectedMsg:
		b.call = msg.Call

	case app.OnResponseMsg:
		if msg.Body != "" {
			f := colorjson.NewFormatter()
			f.Indent = 2

			var obj interface{}
			json.Unmarshal([]byte(msg.Body), &obj)
			if obj == nil {
				b.body = msg.Body
			} else {
				s, _ := f.Marshal(obj)
				b.body = string(s)
			}
			// prepend line numbers to each line
			lines := utils.SplitLines(b.body)
			numberOfLines := len(lines)
			maxDigits := len(strconv.Itoa(numberOfLines))
			for i, line := range lines {
				// pad line number with spaces
				linenr := strconv.Itoa(i + 1)
				line = strings.Repeat(" ", maxDigits-len(linenr)) + linenr + "  " + line
				lines[i] = lipgloss.NewStyle().Foreground(config.COLOR_GRAY).Render(line) + "\n"
			}
			b.body = strings.Join(lines, "")
			b.viewport.SetContent(string(b.body))
			b.status = msg.Response.StatusCode
		}

	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+l":
			b.activeTab = min(b.activeTab+1, len(b.Tabs)-1)

		case "ctrl+h":
			b.activeTab = max(b.activeTab-1, 0)
		}

	case config.WindowFocusedMsg:
		b.focused = msg.State

	}
	var cmds []tea.Cmd
	var cmd tea.Cmd
	b.viewport, cmd = b.viewport.Update(msg)
	cmds = append(cmds, cmd)

	if b.content != nil {
		b.content, cmd = b.content.Update(msg)
		cmds = append(cmds, cmd)
	}

	return b, tea.Batch(cmds...)
}

func (b *Results) SetActiveTab(tab int) {
	b.activeTab = tab
}

func (b Results) View() string {
	var style lipgloss.Style
	if b.focused {
		style = activeStyle
	} else {
		style = inactiveStyle
	}

	b.viewport.Width = b.width - 2
	b.viewport.Height = b.height - 4

	var content string
	if b.body != "" {
		content = b.viewport.View()
	} else {
		icon := `
   ____
  /\___\
 /\ \___\
 \ \/ / /
  \/_/_/
`

		message := lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().Foreground(config.COLOR_HIGHLIGHT).Render(icon),
			"Not sent yet")

		center := lipgloss.PlaceHorizontal(b.viewport.Width, lipgloss.Center, message)
		content = lipgloss.NewStyle().
			Foreground(config.COLOR_GRAY).
			Bold(true).
			Render(lipgloss.PlaceVertical(b.viewport.Height-1, lipgloss.Center, center))
	}

	header := "Response"
	if b.status != 0 {
		header += " " + statusStyle.Render(strconv.Itoa(b.status))
	}
	return style.Render(" " + header + "\n\n" + content)
}
