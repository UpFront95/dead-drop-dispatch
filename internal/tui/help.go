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
			{keys: "enter", text: "accept job, assign runner, respond, or open city brief"},
			{keys: "r", text: "cycle route options, or cycle message responses in Message Feed"},
			{keys: "space", text: "advance the turn phase or resolve active runs"},
			{keys: "esc", text: "close brief/help, then cancel pending accepted job"},
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
	leftW := min(52, max(42, width/3))
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
		styles.Accent.Render("Stat glossary"),
		statHelpRow("N", "Night", "run day; survive seven", width, styles),
		statHelpRow("T", "Turn", "dispatch cycle this night", width, styles),
		statHelpRow("C", "Credits", "cash for costs/upgrades", width, styles),
		statHelpRow("H", "Heat", "attention; max burns desk", width, styles),
		statHelpRow("I", "Integrity", "desk health; zero loses", width, styles),
		statHelpRow("RNR", "Runners", "ready/busy/injured/out", width, styles),
		statHelpRow("JOB", "Jobs", "board/accepted/active", width, styles),
		statHelpRow("DUE", "Deadline", "nearest job deadline", width, styles),
		statHelpRow("CG", "Cargo", "lowest known integrity", width, styles),
		styles.PanelText.Render(""),
		styles.Accent.Render("District stats"),
		statHelpRow("SURV", "Surveillance", "detection risk", width, styles),
		statHelpRow("TRAF", "Traffic", "crowds, delay, cover", width, styles),
		statHelpRow("DNGR", "Danger", "injury pressure", width, styles),
		statHelpRow("SGNL", "Signal", "comms and intel quality", width, styles),
		styles.PanelText.Render(""),
		styles.Accent.Render("Runner stats"),
		statHelpRow("SPD", "Speed", "beats deadlines", width, styles),
		statHelpRow("STL", "Stealth", "avoids detection", width, styles),
		statHelpRow("NRV", "Nerve", "holds under pressure", width, styles),
		statHelpRow("TLK", "Talk", "handles people/gates", width, styles),
		statHelpRow("LOY", "Loyalty", "resists betrayal pressure", width, styles),
		statHelpRow("STR", "Stress", "strain; high is risky", width, styles),
		statHelpRow("CAP", "Capacity", "jobs carried out of 2", width, styles),
	}
	return strings.Join(lines, "\n")
}

func statHelpRow(code string, name string, effect string, width int, styles Styles) string {
	label := styles.InlineCode.Width(5).Render(clipText(code, 5))
	text := clipText(name+": "+effect, max(1, width-7))
	return label + styles.PanelText.Render(" ") + styles.PanelText.Render(text)
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
			lines = append(lines, key+styles.PanelText.Render("  ")+description)
		}
	}
	return strings.Join(lines, "\n")
}
