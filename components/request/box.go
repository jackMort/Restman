package request

import (
	"net/url"
	"restman/app"
	"restman/components/auth"
	"restman/components/config"
	"restman/components/headers"
	"restman/components/params"
	"strings"

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
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Border(lipgloss.NormalBorder()).UnsetBorderTop()
	tabGap            = inactiveTabStyle.
				BorderTop(false).
				BorderLeft(false).
				BorderRight(false)

	emptyMessage = lipgloss.NewStyle().Padding(2, 2).Foreground(config.COLOR_GRAY)

	listHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(config.COLOR_SUBTLE).
			Render
)

type Request struct {
	title     string
	focused   bool
	body      string
	width     int
	height    int
	Tabs      []string
	activeTab int
	content   tea.Model
	call      *app.Call
}

func New() Request {
	return Request{
		title: "Params",
		Tabs:  []string{"Params", "Headers", "Auth", "Body"},
	}
}

// satisfy the tea.Model interface
func (b Request) Init() tea.Cmd {
	b.activeTab = 0
	return nil
}

func (b Request) GetContent() tea.Model {
	if b.activeTab == 2 {
		return auth.New(b.width, b.call)
	} else if b.activeTab == 3 {
		body := ""
		if b.call != nil {
			body = b.call.Data
		}
		return NewBody(body, b.width-2, b.height-4)
	}
	return nil
}

func (b Request) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case app.CallSelectedMsg:
		b.call = msg.Call
		// TODO:
		// b.body = msg.Tab.Results
		b.content = b.GetContent()

	case tea.WindowSizeMsg:
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
	cmds = append(cmds, cmd)

	if b.content != nil {
		b.content, cmd = b.content.Update(msg)
		cmds = append(cmds, cmd)
	}

	return b, tea.Batch(cmds...)
}

func (b *Request) SetActiveTab(tab int) {
	b.activeTab = tab
	b.content = b.GetContent()
}

func (b Request) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	if b.focused {
		inactiveTabStyle = inactiveTabStyle.BorderForeground(config.COLOR_HIGHLIGHT)
		activeTabStyle = activeTabStyle.BorderForeground(config.COLOR_HIGHLIGHT)
		windowStyle = windowStyle.BorderForeground(config.COLOR_HIGHLIGHT)
		tabGap = tabGap.BorderForeground(config.COLOR_HIGHLIGHT)
	} else {
		inactiveTabStyle = inactiveTabStyle.BorderForeground(config.COLOR_SUBTLE)
		activeTabStyle = activeTabStyle.BorderForeground(config.COLOR_SUBTLE)
		windowStyle = windowStyle.BorderForeground(config.COLOR_SUBTLE)
		tabGap = tabGap.BorderForeground(config.COLOR_SUBTLE)
	}

	for i, t := range b.Tabs {
		var style lipgloss.Style
		isFirst, isActive := i == 0, i == b.activeTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
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
	renderedTabs = append(renderedTabs, tabGap.Render(strings.Repeat(" ", b.width-43)))

	windowStyle = windowStyle.Height(b.height - 4)

	style := inactiveTabStyle
	border, _, _, _, _ := style.GetBorder()
	border.Right = " "
	border.BottomRight = "┐"
	style = style.Border(border).BorderTop(false).BorderLeft(false)
	renderedTabs = append(renderedTabs, style.Render(" "))
	row := lipgloss.JoinHorizontal(lipgloss.Bottom, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	var content string
	if b.activeTab == 0 {
		content = emptyMessage.Render("No url params")
		if b.call != nil {

			u, err := url.Parse(b.call.Url)
			if err == nil && b.call.Url != "" {
				m, _ := url.ParseQuery(u.RawQuery)
				if len(m) > 0 {

					table := params.New(m, b.width, b.height)
					content = lipgloss.NewStyle().
						UnsetBold().
						Render(
							table.View(),
						)
				}
			}
		}
	} else if b.activeTab == 1 {
		h := []string{}
		if b.call != nil {
			h = b.call.Headers
		}
		table := headers.New(h, b.width-2, b.width)
		content = lipgloss.NewStyle().
			UnsetBold().
			Render(
				table.View(),
			)
	} else if b.activeTab == 2 {
		content = b.content.View()
	} else if b.activeTab == 3 {
		content = b.content.View()
	} else {
		content = emptyMessage.Render("Not implemented yet")
	}

	doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(content))
	return docStyle.Render(doc.String())
}
