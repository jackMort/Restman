package main

import (
	"fmt"
	"os"
	"restman/app"
	"restman/components/collections"
	"restman/components/config"
	"restman/components/footer"
	"restman/components/help_popup"
	"restman/components/popup"
	"restman/components/results"
	"restman/components/tabs"
	"restman/components/url"
	"restman/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
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
	zone.NewGlobal()

	// layout-tree defintion
	m := Model{tui: boxer.Boxer{}, focused: "url"}

	url := url.New()
	middle := results.New()
	footerBox := footer.New()
	colBox := collections.New()
	tabs := tabs.New()

	centerNode := boxer.CreateNoBorderNode()
	centerNode.VerticalStacked = true
	centerNode.SizeFunc = func(node boxer.Node, widthOrHeight int) []int {
		return []int{
			2,
			3,
			widthOrHeight - 5,
		}
	}
	centerNode.Children = []boxer.Node{
		stripErr(m.tui.CreateLeaf("tabs", tabs)),
		stripErr(m.tui.CreateLeaf("url", url)),
		stripErr(m.tui.CreateLeaf("middle", middle)),
	}

	// middle Node
	middleNode := boxer.CreateNoBorderNode()
	middleNode.SizeFunc = func(node boxer.Node, widthOrHeight int) []int {
		gap := 30
		if m.tui.ModelMap["collections"].(collections.Collections).IsMinified() {
			gap = 6
		}
		return []int{
			gap,
			widthOrHeight - gap,
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

type Model struct {
	tui         boxer.Boxer
	focused     string
	popup       tea.Model
	collections collections.Collections
}

func (m Model) Init() tea.Cmd {
	var (
		cmd  tea.Cmd
		cmd2 tea.Cmd
	)

	m.focused = "url"
	m.tui.ModelMap[m.focused], cmd = m.tui.ModelMap[m.focused].Update(config.WindowFocusedMsg{State: true})

  tabs := m.tui.ModelMap["tabs"].(tabs.Tabs)
	m.tui.ModelMap["tabs"], cmd2 = tabs.AddTab()

	return tea.Batch(
		app.GetInstance().ReadCollectionsFromJSON(),
		cmd,
		cmd2,
	)
}

func (m *Model) Next() (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.focused != "" {
		var previous tea.Model = m.tui.ModelMap[m.focused]
		m.tui.ModelMap[m.focused], cmd = previous.Update(config.WindowFocusedMsg{State: false})
		cmds = append(cmds, cmd)
	}

	coll := m.tui.ModelMap["collections"].(collections.Collections)

	switch m.focused {
	case "collections":
		m.focused = "url"
	case "url":
		m.focused = "middle"
	default:
		if coll.IsMinified() {
			m.focused = "url"
		} else {
			m.focused = "collections"
		}
	}

	m.tui.ModelMap[m.focused], cmd = m.tui.ModelMap[m.focused].Update(config.WindowFocusedMsg{State: true})

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) SetFocused(newFocused string) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.focused != "" {
		var previous tea.Model = m.tui.ModelMap[m.focused]
		m.tui.ModelMap[m.focused], cmd = previous.Update(config.WindowFocusedMsg{State: false})
		cmds = append(cmds, cmd)
	}

	m.focused = newFocused
	m.tui.ModelMap[m.focused], cmd = m.tui.ModelMap[m.focused].Update(config.WindowFocusedMsg{State: true})

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) getCollectionsPane() *collections.Collections {
	collections := m.tui.ModelMap["collections"].(collections.Collections)
	return &collections
}

func (m Model) getUrlPane() *url.Url {
	url := m.tui.ModelMap["url"].(url.Url)
	return &url
}

func (m Model) getMiddlePane() *results.Middle {
	middle := m.tui.ModelMap["middle"].(results.Middle)
	return &middle
}

func (m Model) AddToCollection() tea.Cmd {
	url := m.getUrlPane()
	coll := m.popup.(collections.AddToCollection)
	return app.GetInstance().AddToCollection(coll.CollectionName(), url.Value(), url.Method())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

  case popup.ClosePopupMsg:
		m.popup = nil

	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			if zone.Get("method").InBounds(msg) {
				m.SetFocused("url")
				url := m.getUrlPane()
				url.CycleOverMethods()
				m.tui.ModelMap["url"] = url

			} else if zone.Get("input").InBounds(msg) {
				m.SetFocused("url")

			} else if zone.Get("send").InBounds(msg) {
				m.SetFocused("url")

			} else if zone.Get("save").InBounds(msg) {
				m.SetFocused("url")
				url := m.getUrlPane()

				coll := collections.NewAddToCollection(m.View(), 40, m.tui.LayoutTree.GetWidth())
				coll.SetUrl(url.Value())

				m.popup = coll
				return m, m.popup.Init()

			} else if zone.Get("add_to_collection_cancel").InBounds(msg) {
				m.popup = nil
				return m, nil

			} else if zone.Get("add_to_collection_save").InBounds(msg) {
				cmd := m.AddToCollection()

				m.popup = nil
				return m, cmd

			} else if zone.Get("collections_minified").InBounds(msg) {
				m.tui.ModelMap["collections"], cmd = m.tui.ModelMap["collections"].(collections.Collections).SetMinified(false)
				m.tui.UpdateSize(tea.WindowSizeMsg{Width: m.tui.LayoutTree.GetWidth(), Height: m.tui.LayoutTree.GetHeight()})
				return m.SetFocused("collections")

			} else if zone.Get("tab_Results").InBounds(msg) {
				m.SetFocused("middle")
				middle := m.getMiddlePane()
				middle.SetActiveTab(0)
				m.tui.ModelMap["middle"] = middle

			} else if zone.Get("tab_Params").InBounds(msg) {
				m.SetFocused("middle")
				middle := m.getMiddlePane()
				middle.SetActiveTab(1)
				m.tui.ModelMap["middle"] = middle

			} else if zone.Get("tab_Headers").InBounds(msg) {
				m.SetFocused("middle")
				middle := m.getMiddlePane()
				middle.SetActiveTab(2)
				m.tui.ModelMap["middle"] = middle

			} else if zone.Get("tab_Auth").InBounds(msg) {
				m.SetFocused("middle")
				middle := m.getMiddlePane()
				middle.SetActiveTab(3)
				m.tui.ModelMap["middle"] = middle

			} else if zone.Get("collections_minify").InBounds(msg) {
				m.tui.ModelMap["collections"], cmd = m.tui.ModelMap["collections"].(collections.Collections).SetMinified(true)
				m.tui.UpdateSize(tea.WindowSizeMsg{Width: m.tui.LayoutTree.GetWidth(), Height: m.tui.LayoutTree.GetHeight()})
				return m.SetFocused("url")
			} else if zone.Get("collections").InBounds(msg) {
				return m.SetFocused("collections")
			}
		}

	case collections.CreateResultMsg, collections.AddToCollectionResultMsg:
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
	case app.SetFocusMsg:
		fmt.Println("SetFocusMsg: ", msg.Item)
		m.SetFocused(msg.Item)

	case tea.KeyMsg:
		{

			switch msg.String() {
			case "q", "ctrl+c":
				width := 100
				m.popup = popup.NewChoice(m.View(), width, "Are you sure, you want to quit?", false)
				return m, m.popup.Init()

			case "ctrl+n":
				m.popup = collections.NewCreate(m.View(), utils.MinInt(70, 100))
				return m, m.popup.Init()

			case "?":
				m.popup = help.NewHelp(m.View(), 50)
				return m, m.popup.Init()

			case "ctrl+s":
				url := m.getUrlPane()

				coll := collections.NewAddToCollection(m.View(), 40, m.tui.LayoutTree.GetWidth())
				coll.SetUrl(url.Value())

				m.popup = coll
				return m, m.popup.Init()

			case "tab":
				m, cmd := m.Next()
				return m, cmd

			case "ctrl+f":
				coll := m.tui.ModelMap["collections"].(collections.Collections)
				minified := coll.IsMinified()
				m.tui.ModelMap["collections"], cmd = m.tui.ModelMap["collections"].(collections.Collections).SetMinified(!minified)
				m.tui.UpdateSize(tea.WindowSizeMsg{Width: m.tui.LayoutTree.GetWidth(), Height: m.tui.LayoutTree.GetHeight()})

				if minified && m.focused != "collections" {
					return m.SetFocused("collections")
				}
				if !minified && m.focused == "collections" {
					return m.SetFocused("url")
				}
				return m, nil

			default:
				if m.focused != "" {
					m.tui.ModelMap[m.focused], cmd = m.tui.ModelMap[m.focused].Update(msg)
				}

			}
		}
	case tea.WindowSizeMsg:
		m.tui.UpdateSize(msg)

	default:
		var cmd tea.Cmd
		for key, element := range m.tui.ModelMap {
			m.tui.ModelMap[key], cmd = element.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	var cmdF tea.Cmd
	m.tui.ModelMap["footer"], cmdF = m.tui.ModelMap["footer"].Update(msg)
	cmds = append(cmds, cmdF)

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.popup != nil {
		return m.popup.View()
	}
	return zone.Scan(m.tui.View())
}
