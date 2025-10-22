package styles

import "github.com/charmbracelet/lipgloss"

var (
	ColorSuccess = lipgloss.Color("#00C853") // green
	ColorError   = lipgloss.Color("#FF5252") // red
	ColorInfo    = lipgloss.Color("#40C4FF") // blue
	ColorWarn    = lipgloss.Color("#FFD740") // yellow
	ColorAccent  = lipgloss.Color("#7D56F4") // purple

	Base = lipgloss.NewStyle().
		Padding(0, 1).
		Margin(0, 0, 0, 0)

	Success = Base.Copy().Foreground(ColorSuccess).Bold(true)
	Error   = Base.Copy().Foreground(ColorError).Bold(true)
	Info    = Base.Copy().Foreground(ColorInfo).Bold(true)
	Warn    = Base.Copy().Foreground(ColorWarn).Bold(true)

	// Icons
	IconSuccess = "✔"
	IconError   = "✖"
	IconInfo    = "ℹ"
	IconWarn    = "⚠"

	Header = Base.Copy().
		Foreground(ColorAccent).
		Bold(true).
		Underline(true)
)
