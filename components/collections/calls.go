package collections

import (
	"fmt"
	"io"
	"restman/app"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().Faint(true)
	selectedItemStyle = itemStyle.Copy().Faint(false)
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(3)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	str := list.Item(listItem).(app.Call).Title()
	method := list.Item(listItem).(app.Call).MethodShortView()

	style := itemStyle
	if index == m.Index() {
		style = selectedItemStyle
	}

	// TODO make some elipsis for long text
	style = style.MaxWidth(25)

	item := style.Render("󰑂 " + method + " " + style.Render(str))

	fmt.Fprint(w, item)
}

type callModel struct {
	list       list.Model
	collection *app.Collection
}

func NewCallModel() callModel {
	// Make initial list of items
	groceryList := list.New([]list.Item{}, itemDelegate{}, 0, 0)
	groceryList.Styles.Title = titleStyle
	groceryList.Styles.TitleBar = titleBarStyle
	groceryList.Help.ShowAll = true
	groceryList.SetShowHelp(false)
	groceryList.SetShowTitle(false)
	groceryList.SetShowStatusBar(false)
	groceryList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("ctrl+h"),
				key.WithHelp("ctrl+h", "back to collections"),
			),
		}
	}

	return callModel{
		list: groceryList,
	}
}

func (m callModel) Init() tea.Cmd {
	return nil
}

func (m callModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case app.FetchCollectionsSuccessMsg:
		if m.collection != nil {
			items := []list.Item{}
			for i, collection := range msg.Collections {
				if collection.ID == m.collection.ID {
					for _, call := range msg.Collections[i].Calls {
						items = append(items, call)
					}
				}
			}
			return m, m.list.SetItems(items)
		}

	case app.CollectionSelectedMsg:
		m.collection = msg.Collection

		items := []list.Item{}
		for _, call := range msg.Collection.Calls {
			items = append(items, call)
		}
		return m, m.list.SetItems(items)

	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h-2, msg.Height-v-6)

	case tea.KeyMsg:

		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "ctrl+h":
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
