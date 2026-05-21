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
		Base:       lipgloss.NewStyle().Foreground(lipgloss.Color("#CBD5D1")),
		Title:      lipgloss.NewStyle().Foreground(lipgloss.Color("#E6FF7A")).Bold(true),
		Status:     lipgloss.NewStyle().Foreground(lipgloss.Color("#D6DDD8")).Background(lipgloss.Color("#17211F")).Padding(0, 1),
		Panel:      lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#33413D")).Padding(0, 1),
		PanelFocus: lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#E6FF7A")).Padding(0, 1),
		PanelTitle: lipgloss.NewStyle().Foreground(lipgloss.Color("#7FD7C4")).Bold(true),
		Muted:      lipgloss.NewStyle().Foreground(lipgloss.Color("#7B8984")),
		Accent:     lipgloss.NewStyle().Foreground(lipgloss.Color("#7FD7C4")),
		Warning:    lipgloss.NewStyle().Foreground(lipgloss.Color("#F2B36D")),
		Critical:   lipgloss.NewStyle().Foreground(lipgloss.Color("#F27979")),
		Help:       lipgloss.NewStyle().Foreground(lipgloss.Color("#B8C4BF")).Background(lipgloss.Color("#202B28")).Padding(0, 1),
		InlineCode: lipgloss.NewStyle().Foreground(lipgloss.Color("#E6FF7A")),
		Divider:    lipgloss.NewStyle().Foreground(lipgloss.Color("#33413D")),
	}
}

func lightStyles() Styles {
	return Styles{
		Base:       lipgloss.NewStyle().Foreground(lipgloss.Color("#20302C")),
		Title:      lipgloss.NewStyle().Foreground(lipgloss.Color("#435900")).Bold(true),
		Status:     lipgloss.NewStyle().Foreground(lipgloss.Color("#20302C")).Background(lipgloss.Color("#DDE8E3")).Padding(0, 1),
		Panel:      lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#AAB8B2")).Padding(0, 1),
		PanelFocus: lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#435900")).Padding(0, 1),
		PanelTitle: lipgloss.NewStyle().Foreground(lipgloss.Color("#176D5D")).Bold(true),
		Muted:      lipgloss.NewStyle().Foreground(lipgloss.Color("#687A73")),
		Accent:     lipgloss.NewStyle().Foreground(lipgloss.Color("#176D5D")),
		Warning:    lipgloss.NewStyle().Foreground(lipgloss.Color("#925F16")),
		Critical:   lipgloss.NewStyle().Foreground(lipgloss.Color("#9B2626")),
		Help:       lipgloss.NewStyle().Foreground(lipgloss.Color("#20302C")).Background(lipgloss.Color("#EDF3F0")).Padding(0, 1),
		InlineCode: lipgloss.NewStyle().Foreground(lipgloss.Color("#435900")),
		Divider:    lipgloss.NewStyle().Foreground(lipgloss.Color("#AAB8B2")),
	}
}
