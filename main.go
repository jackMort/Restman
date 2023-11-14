package main

import (
	"encoding/json"
	"fmt"
	"os"
	"restman/app"
	"restman/components/collections"
	"restman/components/config"
	"restman/components/footer"
	"restman/components/popup"
	"restman/components/url"
	"strings"

	"github.com/TylerBrock/colorjson"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	boxer "github.com/treilik/bubbleboxer"
)

var windows = []string{"collections", "url", "middle"}

var (
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

func stripErr(n boxer.Node, _ error) boxer.Node {
	return n
}

func main() {
	middle := box{title: "Results"}
	middle.Tabs = []string{"Results", "Params", "Headers", "Auth"}
	url := url.New()
	colBox := collections.New()
	footerBox := footer.New()

	// layout-tree defintion
	m := model{tui: boxer.Boxer{}}

	centerNode := boxer.CreateNoBorderNode()
	centerNode.VerticalStacked = true
	centerNode.SizeFunc = func(node boxer.Node, widthOrHeight int) []int {
		return []int{
			3,
			widthOrHeight - 3,
		}
	}
	centerNode.Children = []boxer.Node{
		stripErr(m.tui.CreateLeaf("url", url)),
		stripErr(m.tui.CreateLeaf("middle", middle)),
	}

	// middle Node
	middleNode := boxer.CreateNoBorderNode()
	middleNode.SizeFunc = func(node boxer.Node, widthOrHeight int) []int {
		fmt.Errorf("widthOrHeight: %d", widthOrHeight)
		return []int{
			30,
			widthOrHeight - 30,
		}
	}
	middleNode.Children = []boxer.Node{
		stripErr(m.tui.CreateLeaf("collections", colBox)),
		centerNode,
	}

	rootNode := boxer.CreateNoBorderNode()
	rootNode.VerticalStacked = true
	rootNode.SizeFunc = func(node boxer.Node, widthOrHeight int) []int {
		return []int{
			widthOrHeight - 1,
			1,
		}
	}
	rootNode.Children = []boxer.Node{
		middleNode,
		stripErr(m.tui.CreateLeaf("footer", footerBox)),
	}

	m.tui.LayoutTree = rootNode

	if f, err := tea.LogToFile("debug.log", "debug"); err != nil {
		fmt.Println("Couldn't open a file for logging:", err)
		os.Exit(1)
	} else {
		defer f.Close()
	}

	app.AddListener(m)

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}

type model struct {
	tui     boxer.Boxer
	focused string
	popup   tea.Model
}

func (m model) OnChange(ap app.App) {
	for key, element := range m.tui.ModelMap {
		m.tui.ModelMap[key], _ = element.Update(config.AppStateChanged{State: ap})
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Next() (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.focused != "" {
		var previous tea.Model = m.tui.ModelMap[m.focused]
		m.tui.ModelMap[m.focused], cmd = previous.Update(config.WindowFocusedMsg{State: false})
		cmds = append(cmds, cmd)
	}

	switch m.focused {
	case "collections":
		m.focused = "url"
	case "url":
		m.focused = "middle"
	case "middle":
		m.focused = "collections"
	default:
		m.focused = "collections"
	}

	m.tui.ModelMap[m.focused], cmd = m.tui.ModelMap[m.focused].Update(config.WindowFocusedMsg{State: true})

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// ----------------------------
	//        POPUP LOGIC
	switch msg := msg.(type) {
	case popup.ChoiceResultMsg:
		m.popup = nil
		if msg.Result {
			return m, tea.Quit
		}
	}

	// If we are showing a popup, we need to update the popup
	if m.popup != nil {
		m.popup, cmd = m.popup.Update(msg)
		return m, cmd
	}
	// ----------------------------

	switch msg := msg.(type) {
	case tea.KeyMsg:
		{

			switch msg.String() {
			case "q", "ctrl+c":
				width := 60
				m.popup = popup.NewChoice(m.View(), width, "Are you sure, you want to quit?", false)
				return m, m.popup.Init()
			case "tab":
				m, cmd := m.Next()
				return m, cmd

			case "p":
			}
			if m.focused != "" {
				m.tui.ModelMap[m.focused], cmd = m.tui.ModelMap[m.focused].Update(msg)
			}
		}
	case tea.WindowSizeMsg:
		m.tui.UpdateSize(msg)

	case popup.ChoiceResultMsg:
		m.popup = nil
		if msg.Result {
			return m, tea.Quit
		}
	}

	var cmdF tea.Cmd
	m.tui.ModelMap["footer"], cmdF = m.tui.ModelMap["footer"].Update(msg)
	cmds = append(cmds, cmdF)

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.popup != nil {
		return m.popup.View()
	}
	return m.tui.View()
}

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
)

type box struct {
	title     string
	focused   bool
	body      string
	width     int
	height    int
	viewport  viewport.Model
	Tabs      []string
	activeTab int
}

// satisfy the tea.Model interface
func (b box) Init() tea.Cmd {
	b.viewport = viewport.New(10, 10)
	b.activeTab = 1
	return nil
}

func (b box) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		testStyle.Width(msg.Width - 2)
		testStyle.Height(msg.Height - 2)
		testStyleFocused.Width(msg.Width - 2)
		testStyleFocused.Height(msg.Height - 2)
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

	case config.AppStateChanged:
		if msg.State.Body != "" {
			b.body = msg.State.Body

			// Create an intersting JSON object to marshal in a pretty format
			f := colorjson.NewFormatter()
			f.Indent = 2

			var obj interface{}
			json.Unmarshal([]byte(b.body), &obj)
			s, _ := f.Marshal(obj)
			b.viewport.SetContent(string(s))
		} else {

			b.body = "No response"
			b.viewport.SetContent(emptyMessage.Render(b.body))
		}

	}
	var cmds []tea.Cmd
	var cmd tea.Cmd
	b.viewport, cmd = b.viewport.Update(msg)

	cmds = append(cmds, cmd)
	return b, tea.Batch(cmds...)
}

func (b box) View() string {
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
		renderedTabs = append(renderedTabs, style.Render(t))
	}
	renderedTabs = append(renderedTabs, tabGap.Render(strings.Repeat(" ", b.width-46)))

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
		content = b.viewport.View()
	} else {
		content = emptyMessage.Render("No implemented yet")
	}

	doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(content))
	return docStyle.Render(doc.String())

	// b.viewport.Width = b.width - 2
	// b.viewport.Height = b.height - 4

	// content := lipgloss.JoinVertical(lipgloss.Left,
	// 	buttons,
	// 	b.viewport.View(),
	// )

	// if b.focused {
	// 	return testStyleFocused.Render(content)
	// }
	// return testStyle.Render(content)
}
