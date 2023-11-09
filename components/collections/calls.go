package collections

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().Italic(true).Bold(false)
	selectedItemStyle = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(3)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type callItem struct {
	path string
}
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	str := list.Item(listItem).(callItem).Title()

  fn := func(s ...string) string {
		return itemStyle.Render(" 󱂛 " + strings.Join(s, " "))
	}
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(" 󱂛 " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (i callItem) Title() string       { return i.path }
func (i callItem) Description() string { return "GET" }
func (i callItem) FilterValue() string { return i.path }

type callModel struct {
	list list.Model
}

func NewCallModel() callModel {

	// Make initial list of items
	items := []list.Item{
		callItem{path: "/us/90210"},
		callItem{path: "/us/ma/belmont"},
		callItem{path: "/us/ma/boston"},
	}

	groceryList := list.New(items, itemDelegate{}, 0, 0)
	groceryList.Title = "  segments-api"
	groceryList.Styles.Title = titleStyle
	groceryList.Styles.TitleBar = titleBarStyle
	groceryList.Help.ShowAll = true
	groceryList.SetShowHelp(false)
  groceryList.SetShowTitle(false)
	groceryList.SetShowStatusBar(false)

	return callModel{
		list: groceryList,
	}
}

func (m callModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m callModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h-2, msg.Height-v-5)
		println(msg.Height)
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m callModel) View() string {
	return appStyle.Render(m.list.View())
}
