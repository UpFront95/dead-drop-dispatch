package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"dead-drop-dispatch/internal/game"
)

const (
	focusCity = iota
	focusRunners
	focusJobs
	focusMessages
	focusDetail
)

type DashboardView struct {
	State    game.GameState
	Width    int
	Height   int
	Focused  int
	ShowHelp bool
	Styles   Styles
}

func RenderDashboard(view DashboardView) tea.View {
	width := max(view.Width, 80)
	height := max(view.Height, 24)
	styles := view.Styles
	if styles.Base.GetForeground() == nil {
		styles = NewStyles()
	}

	header := renderHeader(view.State, width, styles)
	footer := renderFooter(view.ShowHelp, width, styles)
	bodyHeight := height - lipgloss.Height(header) - lipgloss.Height(footer)
	if bodyHeight < 20 {
		bodyHeight = 20
	}

	gap := 1
	leftW := 46
	midW := 42
	rightW := width - leftW - midW - gap*2
	if rightW < 36 {
		rightW = 36
		midW = width - leftW - rightW - gap*2
	}

	topH := 16
	bottomH := 20
	spacerH := bodyHeight - topH - bottomH
	if spacerH < 1 {
		spacerH = 1
		bottomH = bodyHeight - topH - spacerH
	}
	if bottomH < 8 {
		bottomH = 8
		topH = bodyHeight - bottomH
		spacerH = 0
	}

	city := panel("CITY SECTOR", renderCity(view.State, styles), leftW, topH, view.Focused == focusCity, styles)
	runners := panel("RUNNERS", renderRunners(view.State, styles), midW, topH, view.Focused == focusRunners, styles)
	jobs := panel("JOB BOARD", renderJobs(view.State, styles), rightW, topH, view.Focused == focusJobs, styles)

	messages := panel("MESSAGE FEED", renderMessages(view.State, styles), leftW+midW+gap, bottomH, view.Focused == focusMessages, styles)
	detail := panel("DETAIL", renderDetail(view.State, styles), rightW, bottomH, view.Focused == focusDetail, styles)

	top := lipgloss.JoinHorizontal(lipgloss.Top, city, strings.Repeat(" ", gap), runners, strings.Repeat(" ", gap), jobs)
	bottom := lipgloss.JoinHorizontal(lipgloss.Top, messages, strings.Repeat(" ", gap), detail)
	body := lipgloss.JoinVertical(lipgloss.Left, top, blankLines(width, spacerH), bottom)

	rendered := styles.Base.Width(width).Height(height).Render(lipgloss.JoinVertical(lipgloss.Left, header, body, footer))
	result := tea.NewView(rendered)
	result.AltScreen = true
	result.WindowTitle = "Dead Drop Dispatch"
	return result
}

func renderHeader(state game.GameState, width int, styles Styles) string {
	title := styles.Title.Render("DEAD DROP DISPATCH")
	status := fmt.Sprintf("NIGHT %d/%d  TURN %d/%d  CRED %04d  HEAT %02d  INTEGRITY %03d",
		state.Night,
		state.RunNights,
		state.Turn,
		state.TurnsPerNight,
		state.Credits,
		state.Heat,
		state.DispatchIntegrity,
	)
	line := lipgloss.JoinHorizontal(lipgloss.Center, title, "  ", styles.Status.Render(status))
	return lipgloss.NewStyle().Width(width).Render(line)
}

func renderFooter(showHelp bool, width int, styles Styles) string {
	text := "tab focus   ? help   q quit"
	if showHelp {
		text = "tab/shift+tab cycle panels   arrows/hjkl pending   enter pending   space pending   ? hide help   q quit"
	}
	return styles.Help.Width(width).Render(text)
}

func renderCity(state game.GameState, styles Styles) string {
	var b strings.Builder
	for _, district := range state.Districts {
		name := styles.Accent.Width(25).Render(district.Name)
		faction := styles.InlineCode.Width(9).Align(lipgloss.Right).Render(formatFactionControl(district.FactionControl))
		fmt.Fprintf(&b, "%s %s\n", name, faction)
		fmt.Fprintf(&b, "  SURV %d   TRAF %d   DANG %d   SIG %d\n",
			district.Surveillance,
			district.Traffic,
			district.Danger,
			district.SignalQuality,
		)
	}
	return strings.TrimRight(b.String(), "\n")
}

func renderJobs(state game.GameState, styles Styles) string {
	if len(state.AvailableJobs) == 0 {
		return strings.Join([]string{
			styles.Muted.Render("No contracts posted."),
			"",
			"Dispatch wire is quiet.",
			"Job generation comes next.",
		}, "\n")
	}

	var b strings.Builder
	for _, job := range state.AvailableJobs {
		fmt.Fprintf(&b, "%s  %s -> %s\n", styles.Accent.Render(job.Title), job.Origin, job.Destination)
		fmt.Fprintf(&b, "  pay %d  deadline %d turns  cargo %s\n", job.Payout, job.DeadlineTurns, job.Cargo)
	}
	return strings.TrimRight(b.String(), "\n")
}

func renderRunners(state game.GameState, styles Styles) string {
	var b strings.Builder
	for _, runner := range state.Runners {
		stateText := styles.Accent.Render(string(runner.State))
		if runner.State != game.RunnerReady {
			stateText = styles.Warning.Render(string(runner.State))
		}
		fmt.Fprintf(&b, "%s\n", styles.Accent.Render(runner.Name))
		fmt.Fprintf(&b, "  %-9s SPD %d  STL %d  NRV %d  TLK %d\n", stateText, runner.Speed, runner.Stealth, runner.Nerve, runner.Talk)
		fmt.Fprintf(&b, "  LOY %d  STR %d  CAP 0/%d\n", runner.Loyalty, runner.Stress, game.MaxJobsPerRunner)
	}
	return strings.TrimRight(b.String(), "\n")
}

func renderMessages(state game.GameState, styles Styles) string {
	if len(state.Messages) == 0 {
		return styles.Muted.Render("No messages.")
	}

	var b strings.Builder
	for _, message := range state.Messages {
		fmt.Fprintf(&b, "[%02d] %s / %s\n", message.Turn, styles.Accent.Render(message.From), message.Subject)
		fmt.Fprintf(&b, "     %s\n", message.Body)
	}
	return strings.TrimRight(b.String(), "\n")
}

func renderDetail(state game.GameState, styles Styles) string {
	lines := []string{
		styles.Accent.Render("Desk state"),
		fmt.Sprintf("Phase: %s", state.Phase),
		fmt.Sprintf("Seed: %d", state.RandomSeed),
		"",
		styles.Accent.Render("Current focus"),
		"Use tab to inspect panels.",
		"",
		styles.Muted.Render("Exact risk math stays off-screen."),
		styles.Muted.Render("Route factors will appear here."),
	}
	return strings.Join(lines, "\n")
}

func panel(title string, body string, width int, height int, focused bool, styles Styles) string {
	style := styles.Panel
	if focused {
		style = styles.PanelFocus
	}

	frameW, _ := style.GetFrameSize()
	contentW := max(1, width-frameW)
	content := styles.PanelTitle.Render(title) + "\n" + styles.Divider.Render(strings.Repeat("─", contentW)) + "\n" + body
	return style.Width(width).Height(height).Render(content)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func blankLines(width int, height int) string {
	if height <= 0 {
		return ""
	}
	lines := make([]string, height)
	for i := range lines {
		lines[i] = strings.Repeat(" ", width)
	}
	return strings.Join(lines, "\n")
}

func formatFactionControl(faction game.FactionID) string {
	switch faction {
	case "helix_municipal_security":
		return "HELIX"
	case "kestrel_dock_union":
		return "UNION"
	case "saint_orison_clinic_network":
		return "CLINIC"
	case "asterion_systems":
		return "ASTERION"
	default:
		return strings.ToUpper(string(faction))
	}
}
