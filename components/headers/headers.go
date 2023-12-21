package headers

import (
	"restman/components/config"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyKey   = "key"
	columnKeyValue = "value"
)

var (
	styleBase = lipgloss.NewStyle().
		Foreground(config.COLOR_FOREGROUND).
		Bold(false).
		BorderForeground(config.COLOR_SUBTLE)
)

type Model struct {
	width       int
	height      int
	simpleTable table.Model
}

func New(headers []string, width int, height int) Model {
	rows := make([]table.Row, 0, len(headers))
	for _, header := range headers {
		// split : to get key and value
		parts := strings.Split(header, ":")
		row := table.NewRow(table.RowData{
			columnKeyKey:   " " + parts[0],
			columnKeyValue: " " + parts[1],
		})
		rows = append(rows, row)
	}

	return Model{
		width:  width,
		height: height,
		simpleTable: table.New([]table.Column{
			table.NewColumn(columnKeyKey, " Key", 20),
			table.NewColumn(columnKeyValue, " Value", width - 23),
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
	return m.simpleTable.View()
}
