package collections

import (
	"fmt"
	"io"
	"restman/app"
	"restman/components/config"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var (
	itemStyle         = lipgloss.NewStyle().Faint(true).Bold(false).Foreground(config.COLOR_GRAY)
	selectedItemStyle = lipgloss.NewStyle().Faint(false).Bold(true).Foreground(config.COLOR_FOREGROUND)
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	maxWidth := m.Width()

	str := list.Item(listItem).(app.Call).Title()
	method := list.Item(listItem).(app.Call).MethodShortView()
	methodName := listItem.(app.Call).Method
	prefix := " " + method + " "
	prefixWidth := len(methodName) + 3

	style := itemStyle
	if index == m.Index() {
		style = selectedItemStyle
		prefix = lipgloss.NewStyle().Foreground(config.COLOR_SPECIAL).Render("") + prefix
	} else {
		prefix = " " + prefix
	}

	// truncate str if it's too long
	if len(str) > maxWidth-prefixWidth {
		str = str[:maxWidth-prefixWidth-1] + "…"
	}
	item := style.Render(prefix + style.Render(str))

	fmt.Fprint(w, item)
}

type callModel struct {
	list       list.Model
	collection *app.Collection
}

func NewCallModel() callModel {
	// Make initial list of items
	callsList := list.New([]list.Item{}, itemDelegate{}, 0, 0)
	callsList.Styles.Title = titleStyle
	callsList.Styles.TitleBar = titleBarStyle
	callsList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("esc", "go back"),
			),
		}
	}
	callsList.DisableQuitKeybindings()
	callsList.SetShowHelp(false)

	return callModel{
		list: callsList,
	}
}

func (m callModel) Init() tea.Cmd {
	return nil
}

func (m callModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case app.CollectionSelectedMsg:
		m.collection = msg.Collection

		items := []list.Item{}
		for _, call := range msg.Collection.Calls {
			items = append(items, call)
		}

		m.list.Title = zone.Mark("collections_minify", "󰅁 "+m.collection.Name)
		return m, m.list.SetItems(items)

	case tea.WindowSizeMsg:
		x, y := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-x-4, msg.Height-y-2)

	case tea.KeyMsg:

		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "esc":
			return m, app.GetInstance().SetSelectedCollection(nil)

		case "enter":
			i, _ := m.list.SelectedItem().(app.Call)
			return m, app.GetInstance().SetSelectedCall(&i)
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m callModel) View() string {
	return appStyle.Render(m.list.View())
}
