package main

import (
	"encoding/json"
	"fmt"
	"os"
	"restman/app"
	"restman/components/collections"
	"restman/components/config"
	"restman/components/footer"
	"restman/components/url"

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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		{

			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "tab":
				m, cmd := m.Next()
				return m, cmd
			}
			if m.focused != "" {
				m.tui.ModelMap[m.focused], cmd = m.tui.ModelMap[m.focused].Update(msg)
			}
		}
	case tea.WindowSizeMsg:
		m.tui.UpdateSize(msg)
	}

  var cmdF tea.Cmd
	m.tui.ModelMap["footer"], cmdF = m.tui.ModelMap["footer"].Update(msg)
  cmds = append(cmds, cmdF)

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.tui.View()
}

type box struct {
	title    string
	focused  bool
	body     string
	width    int
	height   int
	viewport viewport.Model
}

// satisfy the tea.Model interface
func (b box) Init() tea.Cmd {
	b.viewport = viewport.New(10, 10)
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

	case config.WindowFocusedMsg:
		b.focused = msg.State

	case config.AppStateChanged:
		b.body = msg.State.Body

		// Create an intersting JSON object to marshal in a pretty format
		f := colorjson.NewFormatter()
		f.Indent = 2

		var obj interface{}
		json.Unmarshal([]byte(b.body), &obj)
		s, _ := f.Marshal(obj)

		b.viewport.SetContent(string(s))

	}
	var cmds []tea.Cmd
	var cmd tea.Cmd
	b.viewport, cmd = b.viewport.Update(msg)

	cmds = append(cmds, cmd)
	return b, tea.Batch(cmds...)
}

func (b box) View() string {
	buttons := lipgloss.JoinHorizontal(lipgloss.Left,
		listHeader("Response"), " ",
		listHeader("Results"), " ",
		fmt.Sprintf("%3.f%%", b.viewport.ScrollPercent()*100),
	)

	b.viewport.Width = b.width - 2
	b.viewport.Height = b.height - 4

	content := lipgloss.JoinVertical(lipgloss.Left,
		buttons,
		b.viewport.View(),
	)

	if b.focused {
		return testStyleFocused.Render(content)
	}
	return testStyle.Render(content)
}
