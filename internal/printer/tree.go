package printer

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"

	"github.com/ainsleydev/webkit/internal/printer/styles"
)

// Tree prints a directory-style tree view.
func (c *Console) Tree(root string, children ...any) {
	t := tree.
		Root(root).
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(lipgloss.NewStyle().Foreground(styles.ColorAccent)).
		RootStyle(lipgloss.NewStyle().Foreground(styles.ColorAccent).Bold(true)).
		ItemStyle(lipgloss.NewStyle().Foreground(styles.ColorInfo))

	for _, child := range children {
		switch v := child.(type) {
		case string:
			t = t.Child(v)
		case *tree.Tree:
			t = t.Child(v)
		}
	}

	c.write(fmt.Sprint(t))
}
