package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const helpTabIndex = 3

type HelpView struct {
	Width  int
	Height int
	Styles Styles
}

type helpSection struct {
	title string
	rows  []helpRow
}

type helpRow struct {
	keys string
	text string
}

var helpSections = []helpSection{
	{
		title: "Navigation",
		rows: []helpRow{
			{keys: "tab", text: "focus next dashboard panel"},
			{keys: "shift+tab", text: "focus previous dashboard panel"},
			{keys: "[ and ]", text: "move to previous or next tab"},
			{keys: "1-4", text: "jump to Dashboard, Routing, Equipment, Help"},
			{keys: "up/down or j/k", text: "move the focused panel selection"},
		},
	},
	{
		title: "Actions",
		rows: []helpRow{
			{keys: "enter", text: "accept a job, assign a runner, respond to a message, or open a city brief"},
			{keys: "r", text: "cycle route options, or cycle message responses in Message Feed"},
			{keys: "space", text: "advance the turn phase or resolve active runs"},
			{keys: "esc", text: "go back from a city brief or help, then cancel the pending accepted job"},
		},
	},
	{
		title: "System",
		rows: []helpRow{
			{keys: "?", text: "toggle compact footer help"},
			{keys: "q", text: "quit"},
		},
	},
}

func RenderHelpView(view HelpView) tea.View {
	width := max(view.Width, 80)
	height := max(view.Height, 24)
	styles := view.Styles
	if styles.Base.GetForeground() == nil {
		styles = NewStyles()
	}

	rendered := styles.Base.Width(width).Height(height).Render(renderHelpSurface(width, height, styles))
	result := tea.NewView(rendered)
	result.AltScreen = true
	result.WindowTitle = "Dead Drop Dispatch Help"
	return result
}

func renderHelpSurface(width int, height int, styles Styles) string {
	innerW := max(1, width-4)
	bodyH := max(1, height-2)
	body := renderHelpBody(innerW, bodyH, styles)
	return lipgloss.NewStyle().Width(width).Height(height).Padding(1, 2).Render(body)
}

func renderHelpBody(width int, height int, styles Styles) string {
	gap := 1
	leftW := min(35, max(30, width/4))
	rightW := width - leftW - gap
	if rightW < 32 {
		rightW = 32
		leftW = max(32, width-rightW-gap)
	}

	overview := renderHelpOverview(max(1, leftW-4), styles)
	controls := renderHelpControls(max(1, rightW-4), styles)
	left := panel("HELP", overview, leftW, height, false, styles)
	right := panel("CONTROLS", controls, rightW, height, false, styles)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, strings.Repeat(" ", gap), right)
}

func renderHelpOverview(width int, styles Styles) string {
	lines := []string{
		styles.Accent.Render("Dashboard-first dispatch"),
		styles.PanelText.Render("Most work happens in the dashboard panels."),
		styles.PanelText.Render("Routing has a map view for job geography."),
		styles.PanelText.Render("Equipment stays available for upgrade work."),
		styles.PanelText.Render(""),
		styles.Accent.Render("Panel jobs"),
		styles.PanelText.Render("CITY SECTOR opens district briefs."),
		styles.PanelText.Render("JOB BOARD accepts posted contracts."),
		styles.PanelText.Render("RUNNERS assigns accepted jobs."),
		styles.PanelText.Render("MESSAGE FEED answers open contacts."),
		styles.PanelText.Render("DETAIL shows routes, runner load, and replies."),
		styles.PanelText.Render(""),
		styles.Muted.Render(clipText("Use ? for footer-sized help without leaving the dashboard.", width)),
	}
	return strings.Join(lines, "\n")
}

func renderHelpControls(width int, styles Styles) string {
	lines := []string{}
	for sectionIndex, section := range helpSections {
		if sectionIndex > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, styles.Accent.Render(section.title))
		for _, row := range section.rows {
			key := styles.InlineCode.Width(16).Render(clipText(row.keys, 16))
			description := styles.PanelText.Render(clipText(row.text, max(1, width-18)))
			lines = append(lines, key+"  "+description)
		}
	}
	return strings.Join(lines, "\n")
}
