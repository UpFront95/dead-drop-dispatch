package tui

import "charm.land/lipgloss/v2"

const (
	TargetWidth  = 140
	TargetHeight = 40
)

type Styles struct {
	Base       lipgloss.Style
	Title      lipgloss.Style
	Status     lipgloss.Style
	Panel      lipgloss.Style
	PanelFocus lipgloss.Style
	PanelText  lipgloss.Style
	PanelTitle lipgloss.Style
	Muted      lipgloss.Style
	Accent     lipgloss.Style
	Warning    lipgloss.Style
	Critical   lipgloss.Style
	Help       lipgloss.Style
	InlineCode lipgloss.Style
	Divider    lipgloss.Style
}

func NewStyles() Styles {
	return NewStylesForBackground(true)
}

func NewStylesForBackground(dark bool) Styles {
	if !dark {
		return lightStyles()
	}
	return darkStyles()
}

func darkStyles() Styles {
	return Styles{
		Base:       lipgloss.NewStyle().Foreground(lipgloss.Color("#D8D8E6")).Background(lipgloss.Color("#2B2B31")),
		Title:      lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4FD8")).Background(lipgloss.Color("#2B2B31")).Bold(true),
		Status:     lipgloss.NewStyle().Foreground(lipgloss.Color("#E8E8F2")).Background(lipgloss.Color("#3A3A44")).Padding(0, 1),
		Panel:      lipgloss.NewStyle().Foreground(lipgloss.Color("#D8D8E6")).Background(lipgloss.Color("#32323A")).Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#4A4A57")).Padding(0, 1),
		PanelFocus: lipgloss.NewStyle().Foreground(lipgloss.Color("#E8E8F2")).Background(lipgloss.Color("#32323A")).Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#00D9FF")).Padding(0, 1),
		PanelText:  lipgloss.NewStyle().Foreground(lipgloss.Color("#D8D8E6")),
		PanelTitle: lipgloss.NewStyle().Foreground(lipgloss.Color("#D8D8E6")).Bold(true),
		Muted:      lipgloss.NewStyle().Foreground(lipgloss.Color("#B8B8C8")),
		Accent:     lipgloss.NewStyle().Foreground(lipgloss.Color("#7EE787")),
		Warning:    lipgloss.NewStyle().Foreground(lipgloss.Color("#E6C97A")),
		Critical:   lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B81")),
		Help:       lipgloss.NewStyle().Foreground(lipgloss.Color("#B8B8C8")).Background(lipgloss.Color("#32323A")).Padding(0, 1),
		InlineCode: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF7AE6")),
		Divider:    lipgloss.NewStyle().Foreground(lipgloss.Color("#4A4A57")),
	}
}

func lightStyles() Styles {
	return Styles{
		Base:       lipgloss.NewStyle().Foreground(lipgloss.Color("#2B2B31")).Background(lipgloss.Color("#E8E8F2")),
		Title:      lipgloss.NewStyle().Foreground(lipgloss.Color("#B60091")).Background(lipgloss.Color("#E8E8F2")).Bold(true),
		Status:     lipgloss.NewStyle().Foreground(lipgloss.Color("#2B2B31")).Background(lipgloss.Color("#B8B8C8")).Padding(0, 1),
		Panel:      lipgloss.NewStyle().Foreground(lipgloss.Color("#2B2B31")).Background(lipgloss.Color("#D8D8E6")).Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#8C8C9E")).Padding(0, 1),
		PanelFocus: lipgloss.NewStyle().Foreground(lipgloss.Color("#2B2B31")).Background(lipgloss.Color("#D8D8E6")).Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#0084A3")).Padding(0, 1),
		PanelText:  lipgloss.NewStyle().Foreground(lipgloss.Color("#2B2B31")),
		PanelTitle: lipgloss.NewStyle().Foreground(lipgloss.Color("#2B2B31")).Bold(true),
		Muted:      lipgloss.NewStyle().Foreground(lipgloss.Color("#5B5B6A")),
		Accent:     lipgloss.NewStyle().Foreground(lipgloss.Color("#007F4E")),
		Warning:    lipgloss.NewStyle().Foreground(lipgloss.Color("#8A6500")),
		Critical:   lipgloss.NewStyle().Foreground(lipgloss.Color("#B00028")),
		Help:       lipgloss.NewStyle().Foreground(lipgloss.Color("#2B2B31")).Background(lipgloss.Color("#B8B8C8")).Padding(0, 1),
		InlineCode: lipgloss.NewStyle().Foreground(lipgloss.Color("#B60091")),
		Divider:    lipgloss.NewStyle().Foreground(lipgloss.Color("#8C8C9E")),
	}
}
