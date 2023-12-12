package tabs

import (
	"restman/components/config"
	"restman/utils"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var (
	normal = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderBottom(false).
		Padding(0, 1).
		BorderForeground(config.COLOR_GRAY).
		Foreground(config.COLOR_GRAY)

	focused = normal.Copy().
		BorderForeground(config.COLOR_HIGHLIGHT).
		Foreground(config.COLOR_HIGHLIGHT)

	plus = normal.Copy()

	more = normal.Copy().
		Border(lipgloss.HiddenBorder()).
		BorderBottom(false).
		Foreground(config.COLOR_GRAY)
)

type model struct {
	height  int
	width   int
	tabs    []Tab
	focused int
}

func New() model {
	tab := NewTab()
	return model{
		tabs:    []Tab{tab},
		focused: 0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) AddTab() {
	m.tabs = append(m.tabs, NewTab())
	m.setFocused(len(m.tabs) - 1)
}

func (m *model) setFocused(index int) {
	m.focused = index
}

func (m *model) removeTab(index int) {
	if len(m.tabs) > 1 {
		m.tabs = append(m.tabs[:index], m.tabs[index+1:]...)
	}
	if m.focused >= index {
		m.setFocused(m.focused - 1)
	}
}

func (m *model) setName(index int, name string) {
	m.tabs[index].Name = name
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			if zone.Get("add-tab").InBounds(msg) {
				m.AddTab()
			}

			for i := range m.tabs {
				if zone.Get(utils.Join("tab-", i)).InBounds(msg) {
					m.setFocused(i)
				}
			}

			for i := range m.tabs {
				if zone.Get(utils.Join("remove-tab-", i)).InBounds(msg) {
					m.removeTab(i)
				}
			}
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	var cmd tea.Cmd
	return m, cmd
}

func (m model) View() string {
	tabs := []string{}

	add := zone.Mark("add-tab", plus.Render(""))
	for i, tab := range m.tabs {
		close := zone.Mark(utils.Join("remove-tab-", i), "󰅙")
		title := zone.Mark(utils.Join("tab-", i), " "+tab.Name)

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
