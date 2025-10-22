package printer

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/ainsleydev/webkit/internal/printer/styles"
)

// Table prints a styled table with headers and rows.
func (c *Console) Table(headers []string, rows [][]string) {
	headerStyle := lipgloss.NewStyle().Inherit(styles.Header)
	cellStyle := lipgloss.NewStyle().Padding(0, 1)

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(styles.ColorAccent)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return headerStyle
			case row%2 == 0:
				return cellStyle.Foreground(lipgloss.Color("245"))
			default:
				return cellStyle.Foreground(lipgloss.Color("241"))
			}
		}).
		Headers(headers...).
		Rows(rows...)

	c.write(fmt.Sprint(t))
}
