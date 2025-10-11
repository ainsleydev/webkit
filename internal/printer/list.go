package printer

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"

	"github.com/ainsleydev/webkit/internal/printer/styles"
)

// List prints a simple bullet list with brand colors.
func (c *Console) List(items ...any) {
	enumeratorStyle := lipgloss.NewStyle().Foreground(styles.ColorAccent).MarginRight(1)
	itemStyle := lipgloss.NewStyle().Foreground(styles.ColorInfo)

	l := list.New(items...).
		Enumerator(list.Alphabet).
		EnumeratorStyle(enumeratorStyle).
		ItemStyle(itemStyle)

	c.write(fmt.Sprint(l))
}
