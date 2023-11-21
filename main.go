package main

import (
	"fmt"
	"os"
	"restman/app"
	"restman/components/collections"
	"restman/components/config"
	"restman/components/footer"
	"restman/components/popup"
	"restman/components/results"
	"restman/components/url"
	"restman/utils"

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
	middle := results.New()
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

func (m model) Init() tea.Cmd {
	return app.GetInstance().ReadCollectionsFromJSON()
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

	case collections.CreateResultMsg:
		m.popup = nil
		return m, nil
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
				width := 100
				m.popup = popup.NewChoice(m.View(), width, "Are you sure, you want to quit?", false)
				return m, m.popup.Init()

			case "ctrl+n":
				m.popup = collections.NewCreate(m.View(), utils.MinInt(m.tui.LayoutTree.GetWidth() - 30, 100))
				return m, m.popup.Init()

			case "tab":
				m, cmd := m.Next()
				return m, cmd

			default:
				if m.focused != "" {
					m.tui.ModelMap[m.focused], cmd = m.tui.ModelMap[m.focused].Update(msg)
				}

			}
		}
	case tea.WindowSizeMsg:
		m.tui.UpdateSize(msg)

	default:
		// TODO: is this make sense?
		for key, element := range m.tui.ModelMap {
			m.tui.ModelMap[key], _ = element.Update(msg)
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

