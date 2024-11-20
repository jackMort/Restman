package collections

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"restman/app"
	"restman/components/config"
)

var (
	appStyle = lipgloss.NewStyle()

	titleStyle = lipgloss.NewStyle().
			BorderForeground(config.COLOR_SUBTLE)

	titleBarStyle = lipgloss.NewStyle()
)

type model struct {
	list         list.Model
	delegateKeys *delegateKeyMap
}

func NewModel() model {
	var (
		delegateKeys = newDelegateKeyMap()
	)

	// Make initial list of items
	items := []list.Item{}

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	collectionsList := list.New(items, delegate, 0, 0)
	collectionsList.Title = zone.Mark("collections_minify", " My Collections")
	collectionsList.Styles.Title = titleStyle
	collectionsList.Styles.TitleBar = titleBarStyle

	collectionsList.SetStatusBarItemName("collection", "collections")
	collectionsList.DisableQuitKeybindings()
	collectionsList.SetShowHelp(false)

	return model{
		list:         collectionsList,
		delegateKeys: delegateKeys,
	}
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case app.FetchCollectionsSuccessMsg:
		items := []list.Item{}
		for _, call := range msg.Collections {
			items = append(items, call)
		}
		return m, m.list.SetItems(items)

	case tea.WindowSizeMsg:
		x, y := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-x-4, msg.Height-y-2)

	case tea.KeyMsg:

		if m.list.FilterState() == list.Filtering {
			break
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.list.View())
}
