package styles

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Base          lipgloss.Style
	Bold          lipgloss.Style
	Muted         lipgloss.Style
	StatusRunning lipgloss.Style
	StatusError   lipgloss.Style
	StatusStopped lipgloss.Style
	BorderActive  lipgloss.Style
	TitleBar      lipgloss.Style
}

func DefaultTheme() Theme {
	base := lipgloss.NewStyle()
	return Theme{
		Base:          base,
		Bold:          base.Bold(true),
		Muted:         base.Foreground(lipgloss.Color("240")),
		StatusRunning: base.Foreground(lipgloss.Color("82")),
		StatusError:   base.Foreground(lipgloss.Color("196")),
		StatusStopped: base.Foreground(lipgloss.Color("240")),
		BorderActive: base.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")),
		TitleBar: lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1).
			Bold(true),
	}
}

func CatppuccinTheme() Theme {
	base := lipgloss.NewStyle()
	return Theme{
		Base:          base,
		Bold:          base.Bold(true),
		Muted:         base.Foreground(lipgloss.Color("#585b70")),
		StatusRunning: base.Foreground(lipgloss.Color("#a6e3a1")), // Green
		StatusError:   base.Foreground(lipgloss.Color("#f38ba8")), // Red
		StatusStopped: base.Foreground(lipgloss.Color("#585b70")), // Overlay0
		BorderActive: base.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#cba6f7")), // Mauve
		TitleBar: lipgloss.NewStyle().
			Background(lipgloss.Color("#89b4fa")). // Blue
			Foreground(lipgloss.Color("#1e1e2e")). // Base
			Padding(0, 1).
			Bold(true),
	}
}

var ActiveTheme = DefaultTheme()
