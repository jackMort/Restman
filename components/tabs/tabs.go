package tabs

import (
	"encoding/json"
	"restman/app"
	"restman/components/config"
	"restman/utils"
	"strings"

	"github.com/TylerBrock/colorjson"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var (
	normal = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderBottom(false).
		Padding(0, 1).
		BorderForeground(config.COLOR_SUBTLE).
		Foreground(config.COLOR_SUBTLE)

	focused = normal.Copy().
		BorderForeground(config.COLOR_HIGHLIGHT).
		Foreground(config.COLOR_HIGHLIGHT).
		Bold(true)

	plus = normal.Copy()

	more = normal.Copy().
		Border(lipgloss.HiddenBorder()).
		BorderBottom(false).
		Foreground(config.COLOR_GRAY)
)

type Tabs struct {
	height  int
	width   int
	tabs    []Tab
	focused int
}

func New() Tabs {
	return Tabs{
		tabs:    []Tab{},
		focused: 0,
	}
}

func (m Tabs) Init() tea.Cmd {
	return nil
}

func (m *Tabs) AddTab() (tea.Model, tea.Cmd) {
	m.tabs = append(m.tabs, NewTab())
	return m.setFocused(len(m.tabs) - 1)
}

func (m *Tabs) setFocused(index int) (tea.Model, tea.Cmd) {
	m.focused = index
	tab := m.tabs[m.focused]
	return m, func() tea.Msg {
		return TabFocusedMsg{Tab: &tab}
	}
}

func (m *Tabs) removeTab(index int) (tea.Model, tea.Cmd) {
	if len(m.tabs) > 1 {
		m.tabs = append(m.tabs[:index], m.tabs[index+1:]...)
	}
	if m.focused >= index {
		newIndex := m.focused - 1
		if newIndex < 0 {
			newIndex = 0
		}
		return m.setFocused(newIndex)
	}
	return m, nil
}

func (m *Tabs) setName(index int, name string) {
	m.tabs[index].Name = name
}

func (m Tabs) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case app.OnResponseMsg:
		if msg.Body != "" {
			_, index := m.GetTab(msg.Call)
			f := colorjson.NewFormatter()
			f.Indent = 2

			var obj interface{}
			json.Unmarshal([]byte(msg.Body), &obj)
			if obj == nil {
				m.tabs[index].Results = msg.Body
			} else {
				s, _ := f.Marshal(obj)
				m.tabs[index].Results = string(s)
			}
			return m.setFocused(index)
		}
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			if zone.Get("add-tab").InBounds(msg) {
				return m.AddTab()
			}

			for i := range m.tabs {
				if zone.Get(utils.Join("tab-", i)).InBounds(msg) {
					return m.setFocused(i)
				}
			}

			for i := range m.tabs {
				if zone.Get(utils.Join("remove-tab-", i)).InBounds(msg) {
					return m.removeTab(i)
				}
			}
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	case app.CallSelectedMsg:
		return m.GetOrCreateTab(msg.Call)

	case app.CallUpdatedMsg:
		tab, index := m.GetTab(msg.Call)
		if tab != nil {
			m.tabs[index].Name = msg.Call.Title()
			m.tabs[index].Call = msg.Call
			return m.setFocused(index)
		}
	}

	var cmd tea.Cmd
	return m, cmd
}

func (m *Tabs) GetOrCreateTab(call *app.Call) (tea.Model, tea.Cmd) {
	for i, tab := range m.tabs {
		if tab.Call != nil && tab.Call.ID == call.ID {
			return m.setFocused(i)
		}
	}
	m.tabs = append(m.tabs, NewTabWithCall(call))
	return m.setFocused(len(m.tabs) - 1)
}

func (m *Tabs) GetTab(call *app.Call) (*Tab, int) {
	for i, tab := range m.tabs {
		if tab.Call != nil && tab.Call.ID == call.ID {
			return &tab, i
		}
	}
	return nil, 0
}

func (m Tabs) View() string {
	tabs := []string{}

	add := zone.Mark("add-tab", plus.Render(""))
	for i, tab := range m.tabs {
		close := zone.Mark(utils.Join("remove-tab-", i), "󰅙")
		icon := " "
		if tab.Call.Collection() == nil {
			icon = ""
		}
		title := zone.Mark(utils.Join("tab-", i), icon+" "+tab.Name)

		style := normal
		if m.focused == i {
			style = focused
		}
		newTab := style.Render(lipgloss.JoinHorizontal(lipgloss.Left, title, " ", close))

		// if string length of newTab is greater than width of m.width append "..."
		tmpTabs := tabs
		tmpTabs = append(tmpTabs, newTab)
		tmpTabs = append(tmpTabs, add)

		rendered := lipgloss.JoinHorizontal(
			lipgloss.Left,
			tmpTabs...,
		)
		finalWidth := lipgloss.Width(rendered)
		if finalWidth > m.width {
			tmpTabs := tabs
			tmpTabs = append(tmpTabs, add)
			rendered = lipgloss.JoinHorizontal(
				lipgloss.Left,
				tmpTabs...,
			)
			finalWidth = lipgloss.Width(rendered)
			if m.width-finalWidth < 7 {
				finalWidth = finalWidth - lipgloss.Width(tabs[len(tabs)-1])
				tabs = tabs[:len(tabs)-1]
			}
			// add spacer
			tabs = append(tabs, more.Render("..."))
			count := m.width - finalWidth - 7
			if count > 0 {
				tabs = append(tabs, strings.Repeat(" ", count))
			}
			break
		} else {
			tabs = append(tabs, newTab)
		}
	}

	tabs = append(tabs, add)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		tabs...,
	)
}
