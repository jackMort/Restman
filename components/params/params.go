package params

import (
	"net/url"
	"restman/app"
	"restman/components/config"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyKey   = "key"
	columnKeyValue = "value"
)

var styleBase = lipgloss.NewStyle().
	Foreground(config.COLOR_FOREGROUND).
	Bold(false).
	BorderForeground(config.COLOR_SUBTLE)

type Model struct {
	width       int
	height      int
	call        *app.Call
	simpleTable table.Model
	items       url.Values
}

func New(call *app.Call, width int, height int) Model {
	items := make(map[string][]string)
	if call != nil {
		u, err := url.Parse(call.Url)
		if err == nil && call.Url != "" {
			items, _ = url.ParseQuery(u.RawQuery)
		}
	}

	// sort items
	keys := make([]string, 0, len(items))
	for k := range items {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	rows := make([]table.Row, 0, len(items))
	for _, key := range keys {
		row := table.NewRow(table.RowData{
			columnKeyKey:   " " + key,
			columnKeyValue: " " + strings.Join(items[key], ", "),
		})
		rows = append(rows, row)
	}

	return Model{
		call:   call,
		items:  items,
		width:  width,
		height: height,
		simpleTable: table.New([]table.Column{
			table.NewColumn(columnKeyKey, " Key", 20),
			table.NewColumn(columnKeyValue, " Value", width-25),
		}).WithRows(rows).BorderRounded().
			WithBaseStyle(styleBase).
			Focused(true),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.simpleTable, cmd = m.simpleTable.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if len(m.items) == 0 {
		return config.EmptyMessageStyle.Padding(2, 2).Render("No query params defined. You can " + config.LinkStyle.Underline(true).Render("add param"))
	}
	return m.simpleTable.View() + "\n" + config.LinkStyle.Padding(0, 1).Render("+ Add Param")
}
