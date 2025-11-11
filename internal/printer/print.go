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
	writer       io.Writer
	infoWriter   io.Writer // Separate writer for informational messages (can be io.Discard)
}

// New returns a new Console instance that writes to
// the given io.Writer.
func New(w io.Writer) *Console {
	return &Console{
		writer:     w,
		infoWriter: w, // By default, both writers are the same
	}
}

// NewSilent returns a new Console instance that suppresses
// informational output but preserves content output.
func NewSilent(w io.Writer) *Console {
	return &Console{
		writer:     w,
		infoWriter: io.Discard,
	}
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

// Println writes plain, unstyled text to the console, with a linebreak.
// Uses infoWriter (can be suppressed in silent mode).
func (c *Console) Println(msg string) {
	c.writeInfo(msg)
	c.lineBreakInfo()
}

// Printf writes plain, unstyled text to the console,  with formatting.
func (c *Console) Printf(msg string, args ...any) {
	c.write(fmt.Sprintf(msg, args...))
}

// Success prints a success message with a checkmark icon and success color.
// Uses infoWriter (can be suppressed in silent mode).
func (c *Console) Success(msg string) {
	c.writeInfo(styles.Success.Render(fmt.Sprintf("%s %s", styles.IconSuccess, msg)))
	c.lineBreakInfo()
}

// Error prints an error message with a cross icon and error color.
// Always outputs (uses main writer, not suppressed in silent mode).
func (c *Console) Error(msg string) {
	c.write(styles.Error.Render(fmt.Sprintf("%s %s", styles.IconError, msg)))
	c.LineBreak()
}

// Info prints an informational message with an info icon and color.
// Uses infoWriter (can be suppressed in silent mode).
func (c *Console) Info(msg string) {
	c.writeInfo(styles.Info.Render(fmt.Sprintf("%s %s", styles.IconInfo, msg)))
	c.lineBreakInfo()
}

// Warn prints a warning message with a warning icon and color.
// Uses infoWriter (can be suppressed in silent mode).
func (c *Console) Warn(msg string) {
	c.writeInfo(styles.Warn.Render(fmt.Sprintf("%s %s", styles.IconWarn, msg)))
	c.lineBreakInfo()
}

// LineBreak prints \n to the writer.
func (c *Console) LineBreak() {
	c.write("\n")
}

func (c *Console) write(s string) {
	if c.writer == nil { // Guard check
		return
	}
	_, _ = io.WriteString(c.writer, s)
}

func (c *Console) writeInfo(s string) {
	if c.infoWriter == nil { // Guard check
		return
	}
	_, _ = io.WriteString(c.infoWriter, s)
}

func (c *Console) lineBreakInfo() {
	c.writeInfo("\n")
}
