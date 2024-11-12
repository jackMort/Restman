package headers

import (
	"restman/app"
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

var styleBase = lipgloss.NewStyle().
	Foreground(config.COLOR_FOREGROUND).
	Bold(false).
	BorderForeground(config.COLOR_SUBTLE)

type Model struct {
	width       int
	height      int
	simpleTable table.Model
	call        *app.Call
}

func GetRows(headers []string) []table.Row {
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
	return rows
}

func New(call *app.Call, width int, height int) Model {

	headers := []string{}
	if call != nil {
		headers = call.Headers
	}

	return Model{
		call:   call,
		width:  width,
		height: height,
		simpleTable: table.New([]table.Column{
			table.NewColumn(columnKeyKey, " Key", 20),
			table.NewColumn(columnKeyValue, " Value", width-23),
		}).WithRows(GetRows(headers)).BorderRounded().
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

	switch msg := msg.(type) {
	case tea.KeyMsg:
		{
			switch msg.String() {
			case "x":
				key := strings.TrimSpace(m.simpleTable.HighlightedRow().Data[columnKeyKey].(string))
				headers := []string{}
				for _, header := range m.call.Headers {
					if key != strings.Split(header, ":")[0] {
						headers = append(headers, header)
					}
				}
				m.call.Headers = headers

				cmd := func() tea.Msg {
					return app.CallUpdatedMsg{Call: m.call}
				}
				m.simpleTable = m.simpleTable.WithRows(GetRows(m.call.Headers))
				cmds = append(cmds, cmd)
			}
		}
	case app.CallUpdatedMsg:
		m.call = msg.Call

	}

	m.simpleTable, cmd = m.simpleTable.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	content := config.EmptyMessageStyle.Padding(2, 2).Render("No headers defined.")
	if m.call != nil && len(m.call.Headers) > 0 {
		content = m.simpleTable.View()
	}
	return content + "\n" + config.LinkStyle.Padding(0, 1).Foreground(config.COLOR_LINK).Render("+ Add Header")
}
