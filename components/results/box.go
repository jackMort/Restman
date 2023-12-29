package results

import (
	"net/url"
	"restman/app"
	"restman/components/auth"
	"restman/components/config"
	"restman/components/headers"
	"restman/components/params"
	"restman/components/tabs"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle()
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Border(lipgloss.NormalBorder()).UnsetBorderTop()
	tabGap            = inactiveTabStyle.Copy().
				BorderTop(false).
				BorderLeft(false).
				BorderRight(false)

	emptyMessage = lipgloss.NewStyle().Padding(2, 2).Foreground(config.COLOR_GRAY)

	testStyle = lipgloss.NewStyle().
			Bold(true).
			Border(lipgloss.NormalBorder()).
			BorderForeground(config.COLOR_SUBTLE).
			PaddingLeft(1)

	testStyleFocused = lipgloss.NewStyle().
				Bold(true).
				Border(lipgloss.NormalBorder()).
				BorderForeground(config.COLOR_HIGHLIGHT).
				PaddingLeft(1)

	listHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(config.COLOR_SUBTLE).
			Render
)

type Middle struct {
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
}

func New() Middle {
	return Middle{
		title: "Results",
		Tabs:  []string{"Results", "Params", "Headers", "Auth", "Body"},
	}
}

// satisfy the tea.Model interface
func (b Middle) Init() tea.Cmd {
	b.viewport = viewport.New(10, 10)
	b.activeTab = 1
	return nil
}

func (b Middle) GetContent() tea.Model {
	if b.activeTab == 3 {
		return auth.New(b.width, b.call)
	} else if b.activeTab == 4 {
		body := ""
		if b.call != nil {
			body = b.call.Data
		}
		return NewBody(body, b.width-2, b.height-4)
	}
	return nil
}

func (b Middle) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tabs.TabFocusedMsg:
		b.call = msg.Tab.Call
		b.body = msg.Tab.Results
		b.viewport.SetContent(string(b.body))
		b.content = b.GetContent()

	case tea.WindowSizeMsg:
		testStyle.Width(msg.Width - 2)
		testStyle.Height(msg.Height - 2)
		testStyleFocused.Width(msg.Width - 2)
		testStyleFocused.Height(msg.Height - 2)
		b.width = msg.Width
		b.height = msg.Height
		b.content = b.GetContent()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+l":
			b.activeTab = min(b.activeTab+1, len(b.Tabs)-1)
			b.content = b.GetContent()

		case "ctrl+h":
			b.activeTab = max(b.activeTab-1, 0)
			b.content = b.GetContent()
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

func (b *Middle) SetActiveTab(tab int) {
	b.activeTab = tab
	b.content = b.GetContent()
}

func (b Middle) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	if b.focused {
		inactiveTabStyle.BorderForeground(config.COLOR_HIGHLIGHT)
		activeTabStyle.BorderForeground(config.COLOR_HIGHLIGHT)
		windowStyle.BorderForeground(config.COLOR_HIGHLIGHT)
		tabGap.BorderForeground(config.COLOR_HIGHLIGHT)
	} else {
		inactiveTabStyle.BorderForeground(config.COLOR_SUBTLE)
		activeTabStyle.BorderForeground(config.COLOR_SUBTLE)
		windowStyle.BorderForeground(config.COLOR_SUBTLE)
		tabGap.BorderForeground(config.COLOR_SUBTLE)
	}

	for i, t := range b.Tabs {
		var style lipgloss.Style
		isFirst, isActive := i == 0, i == b.activeTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		}

		style = style.Border(border)
		renderedTabs = append(renderedTabs, zone.Mark("tab_"+t, style.Render(t)))
	}
	renderedTabs = append(renderedTabs, tabGap.Render(strings.Repeat(" ", b.width-54)))

	windowStyle.Height(b.height - 4)

	style := inactiveTabStyle.Copy()
	border, _, _, _, _ := style.GetBorder()
	border.Right = " "
	border.BottomRight = "┐"
	style = style.Border(border).BorderTop(false).BorderLeft(false)
	renderedTabs = append(renderedTabs, style.Render(" "))
	row := lipgloss.JoinHorizontal(lipgloss.Bottom, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	b.viewport.Width = b.width - 2
	b.viewport.Height = b.height - 5

	var content string
	if b.activeTab == 0 {
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
				Render(lipgloss.PlaceVertical(b.viewport.Height, lipgloss.Center, center))
		}
	} else if b.activeTab == 1 {
		content = emptyMessage.Render("No url params")
		if b.call != nil {

			u, err := url.Parse(b.call.Url)
			if err == nil && b.call.Url != "" {
				m, _ := url.ParseQuery(u.RawQuery)
				if len(m) > 0 {

					table := params.New(m, b.viewport.Width, b.viewport.Height)
					content = lipgloss.NewStyle().
						UnsetBold().
						Render(
							table.View(),
						)
				}
			}
		}
	} else if b.activeTab == 2 {
		h := []string{}
		if b.call != nil {
			h = b.call.Headers
		}
		table := headers.New(h, b.viewport.Width, b.viewport.Height)
		content = lipgloss.NewStyle().
			UnsetBold().
			Render(
				table.View(),
			)
	} else if b.activeTab == 3 {
		content = b.content.View()
	} else if b.activeTab == 4 {
		content = b.content.View()
	} else {
		content = emptyMessage.Render("Not implemented yet")
	}

	doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(content))
	return docStyle.Render(doc.String())
}
