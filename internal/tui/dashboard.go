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
	focusJobs
	focusRunners
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
	width := clamp(view.Width, 80, TargetWidth)
	height := clamp(view.Height, 24, TargetHeight)
	styles := view.Styles
	if styles.Base.GetForeground() == nil {
		styles = NewStyles()
	}

	header := renderHeader(view.State, width, styles)
	footer := renderFooter(view.ShowHelp, width, styles)
	bodyHeight := height - lipgloss.Height(header) - lipgloss.Height(footer) - 2
	if bodyHeight < 18 {
		bodyHeight = 18
	}

	gap := 1
	leftW := 36
	midW := 48
	rightW := width - leftW - midW - gap*2
	if rightW < 28 {
		rightW = 28
		midW = width - leftW - rightW - gap*2
	}

	topH := 16
	bottomH := bodyHeight - topH
	if bottomH < 8 {
		bottomH = 8
		topH = bodyHeight - bottomH
	}

	city := panel("CITY SECTOR", renderCity(view.State, styles), leftW, topH, view.Focused == focusCity, styles)
	jobs := panel("JOB BOARD", renderJobs(view.State, styles), midW, topH, view.Focused == focusJobs, styles)
	runners := panel("RUNNERS", renderRunners(view.State, styles), rightW, topH, view.Focused == focusRunners, styles)

	messages := panel("MESSAGE FEED", renderMessages(view.State, styles), leftW+midW+gap, bottomH, view.Focused == focusMessages, styles)
	detail := panel("DETAIL", renderDetail(view.State, styles), rightW, bottomH, view.Focused == focusDetail, styles)

	top := lipgloss.JoinHorizontal(lipgloss.Top, city, strings.Repeat(" ", gap), jobs, strings.Repeat(" ", gap), runners)
	bottom := lipgloss.JoinHorizontal(lipgloss.Top, messages, strings.Repeat(" ", gap), detail)
	body := lipgloss.JoinVertical(lipgloss.Left, top, bottom)

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
		fmt.Fprintf(&b, "%s\n", styles.Accent.Render(district.Name))
		fmt.Fprintf(&b, "  surv %d  traffic %d  danger %d  signal %d\n",
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
		fmt.Fprintf(&b, "%s  %s\n", styles.Accent.Render(runner.Name), stateText)
		fmt.Fprintf(&b, "  spd %d stl %d nrv %d talk %d\n", runner.Speed, runner.Stealth, runner.Nerve, runner.Talk)
		fmt.Fprintf(&b, "  loyalty %d stress %d\n", runner.Loyalty, runner.Stress)
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

	frameW, frameH := style.GetFrameSize()
	innerW := max(1, width-frameW)
	innerH := max(1, height-frameH)
	content := styles.PanelTitle.Render(title) + "\n" + styles.Divider.Render(strings.Repeat("─", innerW)) + "\n" + body
	return style.Width(innerW).Height(innerH).Render(content)
}

func clamp(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
