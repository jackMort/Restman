package params

import (
	"net/url"
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
	simpleTable table.Model
}

func New(items url.Values, width int, height int) Model {
	colWidth := (width - 4) / 2

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
		width:  width,
		height: height,
		simpleTable: table.New([]table.Column{
			table.NewColumn(columnKeyKey, " Key", colWidth),
			table.NewColumn(columnKeyValue, " Value", colWidth),
		}).WithRows(rows).BorderRounded().
			WithBaseStyle(styleBase),
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
	return m.simpleTable.View()
}
