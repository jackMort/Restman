package main

import (
	"restman/app"
	"restman/components/collections"
	"restman/components/config"
	"restman/components/popup"
	"restman/components/request"
	"restman/components/url"
	"restman/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	boxer "github.com/treilik/bubbleboxer"
)

var windows = []string{"collections", "url", "request", "results"}

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
	rootCmd.Flags().StringP("url", "u", "", "Url")
	rootCmd.Flags().StringP("data", "d", "", "Data")
	rootCmd.Flags().StringP("request", "X", "GET", "HTTP method")
	rootCmd.Flags().StringArrayP("header", "H", []string{}, "HTTP header")

	rootCmd.Execute()
}

type Model struct {
	tui         boxer.Boxer
	focused     string
	popup       tea.Model
	collections collections.Collections
	initialCall *app.Call
	width       int
	height      int
}

func (m Model) Init() tea.Cmd {
	var (
		focusCmd      tea.Cmd
		initalCallCmd tea.Cmd
		runCmd        tea.Cmd
	)

	m.focused = "url"
	m.tui.ModelMap[m.focused], focusCmd = m.tui.ModelMap[m.focused].Update(config.WindowFocusedMsg{State: true})

	// set initial call if provided
	if m.initialCall != nil && m.initialCall.IsValid() {
		initalCallCmd = func() tea.Msg {
			return app.CallSelectedMsg{Call: m.initialCall}
		}
		runCmd = app.GetInstance().GetResponse(m.initialCall)
	}

	return tea.Sequence(
		app.GetInstance().ReadCollectionsFromJSON(),
		focusCmd,
		initalCallCmd,
		runCmd,
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
	case "tabs":
		m.focused = "url"
	case "url":
		m.focused = "request"
	case "request":
		m.focused = "results"
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

func (m Model) getRequestPane() *request.Request {
	middle := m.tui.ModelMap["request"].(request.Request)
	return &middle
}

func (m Model) AddToCollection() tea.Cmd {
	url := m.getUrlPane()
	coll := m.popup.(collections.AddToCollection)
	call := url.Call()
	call.Url = url.Value()
	call.Method = url.Method()
	return app.GetInstance().AddToCollection(coll.CollectionName(), call)
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

	case app.CallSelectedMsg:
		m.SetFocused("url")

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
				url := m.getUrlPane()
				_, cmd = url.Submit()
				m.tui.ModelMap["url"] = url
				return m, cmd

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
				m.SetFocused("results")
			} else if zone.Get("tab_Params").InBounds(msg) {
				m.SetFocused("request")
				request := m.getRequestPane()
				request.SetActiveTab(0)
				m.tui.ModelMap["request"] = request

			} else if zone.Get("tab_Headers").InBounds(msg) {
				m.SetFocused("request")
				request := m.getRequestPane()
				request.SetActiveTab(1)
				m.tui.ModelMap["request"] = request

			} else if zone.Get("tab_Auth").InBounds(msg) {
				m.SetFocused("request")
				request := m.getRequestPane()
				request.SetActiveTab(2)
				m.tui.ModelMap["request"] = request

			} else if zone.Get("tab_Body").InBounds(msg) {
				m.SetFocused("request")
				request := m.getRequestPane()
				request.SetActiveTab(3)
				m.tui.ModelMap["request"] = request

			} else if zone.Get("collections_minify").InBounds(msg) {
				m.tui.ModelMap["collections"], cmd = m.tui.ModelMap["collections"].(collections.Collections).SetMinified(true)
				m.tui.UpdateSize(tea.WindowSizeMsg{Width: m.tui.LayoutTree.GetWidth(), Height: m.tui.LayoutTree.GetHeight()})
				return m.SetFocused("url")
			} else if zone.Get("collections").InBounds(msg) {
				return m.SetFocused("collections")
			}
		}

		// TODO: refactor to use in the same way as AddToCollection
	case collections.CreateResultMsg:
		m.popup = nil
		return m, cmd

	case collections.AddToCollectionResultMsg:
		if msg.Result {
			cmd = m.AddToCollection()
		}
		m.popup = nil
		return m, cmd
	}

	// If we are showing a popup, we need to update the popup
	if m.popup != nil {
		m.popup, cmd = m.popup.Update(msg)
		return m, cmd
	}
	// ----------------------------

	switch msg := msg.(type) {
	case app.SetFocusMsg:
		m.SetFocused(msg.Item)

	case tea.KeyMsg:
		{
			switch msg.String() {
			case "q":
				if m.SizeIsTooSmall() {
					return m, tea.Quit
				}

			case "ctrl+c":
				if m.SizeIsTooSmall() {
					return m, tea.Quit
				}

				width := 100
				m.popup = popup.NewChoice(m.View(), width, "Are you sure, you want to quit?", false)
				return m, m.popup.Init()

			case "ctrl+n":
				m.popup = collections.NewCreate(m.View(), utils.MinInt(70, 100))
				return m, m.popup.Init()

			case "ctrl+a":
				m.popup = NewHelp(m.View(), 70)
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
		m.width = msg.Width
		m.height = msg.Height
		m.tui.UpdateSize(msg)

	default:
		var cmd tea.Cmd
		for key, element := range m.tui.ModelMap {
			m.tui.ModelMap[key], cmd = element.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) SizeIsTooSmall() bool {
	return m.width < 140 || m.height < 30
}

func (m Model) View() string {
	if m.SizeIsTooSmall() {
		return config.FullscreenStyle.
			Width(m.width - 2).
			Height(m.height - 2).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					config.BoxHeader.Render("Restman "+version),
					"Please resize the window to at least 140x30"),
			)
	}

	if m.popup != nil {
		return m.popup.View()
	}
	return zone.Scan(m.tui.View())
}
