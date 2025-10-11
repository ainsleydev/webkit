package printer

import (
	"fmt"
	"io"

	"github.com/ainsleydev/webkit/internal/printer/styles"
)

// See: https://github.com/hay-kot/scaffold/blob/main/internal/printer/printer.go
// See: https://github.com/charmbracelet/lipgloss

// Console provides simple, styled console output
// using lipgloss for consistent branding.
type Console struct {
	writer io.Writer
}

// New returns a new Console instance that writes to
// the given io.Writer.
func New(w io.Writer) *Console {
	return &Console{writer: w}
}

// SetWriter changes the output writer for the console.
// Useful for redirecting output in tests or file logs.
func (c *Console) SetWriter(w io.Writer) {
	c.writer = w
}

// Print writes plain, unstyled text to the console.
func (c *Console) Print(msg string) {
	c.write(msg)
}

// Success prints a success message with a checkmark icon and success color.
func (c *Console) Success(msg string) {
	c.write(styles.Success.Render(fmt.Sprintf("%s %s", styles.IconSuccess, msg)))
}

// Error prints an error message with a cross icon and error color.
func (c *Console) Error(msg string) {
	c.write(styles.Error.Render(fmt.Sprintf("%s %s", styles.IconError, msg)))
}

// Info prints an informational message with an info icon and color.
func (c *Console) Info(msg string) {
	c.write(styles.Info.Render(fmt.Sprintf("%s %s", styles.IconInfo, msg)))
}

// Warn prints a warning message with a warning icon and color.
func (c *Console) Warn(msg string) {
	c.write(styles.Warn.Render(fmt.Sprintf("%s %s", styles.IconWarn, msg)))
}

// LineBreak prints \n to the writer.
func (c *Console) LineBreak() {
	c.write("\n")
}

func (c *Console) write(s string) {
	_, _ = io.WriteString(c.writer, s+"\n")
}
