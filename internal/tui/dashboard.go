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
	State              game.GameState
	Width              int
	Height             int
	Focused            int
	ActiveTab          int
	SelectedDistrict   int
	SelectedJobIndex   int
	SelectedRunner     int
	SelectedRouteIndex int
	SelectedMessage    int
	SelectedResponse   int
	MessageScroll      int
	Notice             string
	ShowDistrictBrief  bool
	ShowHelp           bool
	Styles             Styles
}

var dashboardTabs = []string{
	"DASHBOARD",
	"JOBS",
	"ROUTING",
	"RUNNERS",
	"EQUIPMENT",
	"HELP",
}

func RenderDashboard(view DashboardView) tea.View {
	width := max(view.Width, 80)
	height := max(view.Height, 24)
	styles := view.Styles
	if styles.Base.GetForeground() == nil {
		styles = NewStyles()
	}

	header := renderHeader(view, width, styles)
	tabs := renderTabs(dashboardTabs, view.ActiveTab, width, styles)
	footer := renderFooter(view.ShowHelp, width, styles)
	bodyHeight := height - lipgloss.Height(header) - lipgloss.Height(tabs) - lipgloss.Height(footer)
	if bodyHeight < 20 {
		bodyHeight = 20
	}

	body := blankLines(width, bodyHeight)
	if view.ActiveTab == 0 {
		body = renderDashboardBody(view, width, bodyHeight, styles)
	}

	rendered := styles.Base.Width(width).Height(height).Render(lipgloss.JoinVertical(lipgloss.Left, header, tabs, body, footer))
	result := tea.NewView(rendered)
	result.AltScreen = true
	result.WindowTitle = "Dead Drop Dispatch"
	return result
}

func renderDashboardBody(view DashboardView, width int, bodyHeight int, styles Styles) string {
	gap := 1
	leftW := 40
	midW := 48
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

	city := panel("CITY SECTOR", renderCity(view.State, view.SelectedDistrict, view.ShowDistrictBrief, styles), leftW, topH, view.Focused == focusCity, styles)
	runners := panel("RUNNERS", renderRunners(view.State, view.SelectedRunner, styles), midW, topH, view.Focused == focusRunners, styles)
	jobs := panel("JOB BOARD", renderJobs(view.State, view.SelectedJobIndex, styles), rightW, topH, view.Focused == focusJobs, styles)

	messages := panel("MESSAGE FEED", renderMessages(view.State, view.SelectedMessage, view.MessageScroll, panelBodyHeight(styles.Panel, bottomH), styles), leftW+midW+gap, bottomH, view.Focused == focusMessages, styles)
	detail := panel("DETAIL", renderDetail(view.State, view.Focused, view.SelectedJobIndex, view.SelectedRunner, view.SelectedRouteIndex, view.SelectedMessage, view.SelectedResponse, view.Notice, styles), rightW, bottomH, view.Focused == focusDetail, styles)

	top := lipgloss.JoinHorizontal(lipgloss.Top, city, strings.Repeat(" ", gap), runners, strings.Repeat(" ", gap), jobs)
	bottom := lipgloss.JoinHorizontal(lipgloss.Top, messages, strings.Repeat(" ", gap), detail)
	return lipgloss.JoinVertical(lipgloss.Left, top, renderActionStrip(view, width, spacerH, styles), bottom)
}

func renderTabs(labels []string, active int, width int, styles Styles) string {
	if len(labels) == 0 {
		return ""
	}
	if active < 0 || active >= len(labels) {
		active = 0
	}

	renderedTabs := make([]string, 0, len(labels))
	for i, label := range labels {
		style := styles.Tab
		if i == active {
			style = styles.TabActive
		}
		border, _, _, _, _ := style.GetBorder()
		isFirst, isLast, isActive := i == 0, i == len(labels)-1, i == active
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(label))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	if lipgloss.Width(row) > width {
		row = lipgloss.NewStyle().MaxWidth(width).Render(row)
	}
	return padLinesRight(row, width)
}

func padLinesRight(value string, width int) string {
	lines := strings.Split(value, "\n")
	for i, line := range lines {
		lineWidth := lipgloss.Width(line)
		if lineWidth < width {
			lines[i] = line + strings.Repeat(" ", width-lineWidth)
		}
	}
	return strings.Join(lines, "\n")
}

func renderHeader(view DashboardView, width int, styles Styles) string {
	state := view.State
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
	left := lipgloss.JoinHorizontal(lipgloss.Center, title, "  ", styles.Status.Render(status))
	prompt := styles.Status.Render(currentActionPrompt(view))
	space := width - lipgloss.Width(left) - lipgloss.Width(prompt)
	if space < 2 {
		line := lipgloss.JoinHorizontal(lipgloss.Center, title, "  ", styles.Status.Render(status))
		return lipgloss.NewStyle().Width(width).Render(line)
	}
	line := lipgloss.JoinHorizontal(lipgloss.Center, left, strings.Repeat(" ", space), prompt)
	return lipgloss.NewStyle().Width(width).Render(line)
}

func renderActionStrip(view DashboardView, width int, height int, styles Styles) string {
	if height <= 0 {
		return ""
	}
	text := compactActionLine(view)
	line := styles.Help.Width(width).Render(clipText(text, width-2))
	if height == 1 {
		return line
	}
	return lipgloss.JoinVertical(lipgloss.Left, line, blankLines(width, height-1))
}

func currentActionPrompt(view DashboardView) string {
	state := view.State
	if state.Phase == game.PhaseGameOver {
		status := game.EvaluateRunStatus(state)
		if status.State == game.RunWon {
			return "NEXT run complete"
		}
		return "NEXT run lost"
	}
	if hasPendingComplication(state) {
		return "NEXT resolve complication"
	}
	if len(state.AcceptedJobs) > 0 {
		return "NEXT assign runner"
	}
	if len(state.ActiveJobs) > 0 && state.Phase == game.PhaseDispatch {
		return "NEXT resolve runs"
	}
	switch state.Phase {
	case game.PhaseMessages:
		return "NEXT review messages"
	case game.PhaseJobs:
		return "NEXT review jobs"
	case game.PhaseDispatch:
		if len(state.AvailableJobs) > 0 {
			return "NEXT accept job"
		}
		return "NEXT await postings"
	case game.PhaseComplications:
		return "NEXT resolve complication"
	case game.PhaseReports:
		return "NEXT file reports"
	case game.PhaseCityUpdate:
		return "NEXT city update"
	default:
		return "NEXT dispatch"
	}
}

func compactActionLine(view DashboardView) string {
	state := view.State
	parts := []string{"ACTION " + currentActionText(view)}
	job, ok := currentDecisionJob(view)
	if ok {
		if len(job.RiskFactors) > 0 {
			parts = append(parts, "RISK "+formatFactorsShort(job.RiskFactors))
		}
		if route, routeOK := currentDecisionRoute(view, job); routeOK {
			parts = append(parts, "ROUTE "+formatRouteDetail(route))
		}
	} else if len(state.LastResults) > 0 {
		last := state.LastResults[len(state.LastResults)-1]
		parts = append(parts, fmt.Sprintf("LAST %s: %s", last.JobTitle, last.Outcome))
	}
	if view.Notice != "" {
		parts = append(parts, "NOTE "+view.Notice)
	}
	if message, ok := currentMessage(view.State, view.SelectedMessage); ok && message.Status != game.MessageResolved && message.Audience != "" {
		responses := messageResponseOptions(message)
		if len(responses) > 0 {
			response := responses[clampIndex(view.SelectedResponse, len(responses))]
			parts = append(parts, "REPLY "+response.Label)
		}
	}
	return strings.Join(parts, "  |  ")
}

func currentActionText(view DashboardView) string {
	state := view.State
	if hasPendingComplication(state) {
		return "resolve pending complication"
	}
	if len(state.AcceptedJobs) > 0 {
		return "select runner and route"
	}
	if len(state.ActiveJobs) > 0 && state.Phase == game.PhaseDispatch {
		return "resolve active runs"
	}
	switch state.Phase {
	case game.PhaseMessages:
		return "refresh job board"
	case game.PhaseJobs:
		return "enter dispatch"
	case game.PhaseDispatch:
		if len(state.AvailableJobs) > 0 {
			return "accept highlighted job"
		}
		return "advance dispatch"
	case game.PhaseReports:
		return "file reports"
	case game.PhaseCityUpdate:
		return "run city update"
	case game.PhaseGameOver:
		return "run ended"
	default:
		return "review dashboard"
	}
}

func currentDecisionJob(view DashboardView) (game.Job, bool) {
	if len(view.State.AcceptedJobs) > 0 {
		return view.State.AcceptedJobs[0], true
	}
	if len(view.State.AvailableJobs) > 0 {
		return view.State.AvailableJobs[clampIndex(view.SelectedJobIndex, len(view.State.AvailableJobs))], true
	}
	return game.Job{}, false
}

func currentDecisionRoute(view DashboardView, job game.Job) (game.Route, bool) {
	if len(job.Routes) == 0 {
		return game.Route{}, false
	}
	return job.Routes[clampIndex(view.SelectedRouteIndex, len(job.Routes))], true
}

func hasPendingComplication(state game.GameState) bool {
	for _, complication := range state.Complications {
		if complication.Status == game.ComplicationPending {
			return true
		}
	}
	return false
}

func renderFooter(showHelp bool, width int, styles Styles) string {
	text := "tab focus   [ ] tabs   j/k select   enter accept/assign   r route   space resolve   ? more   q quit"
	if showHelp {
		text = "shift+tab prev panel   1-6 jump tabs   arrows move   [ and ] tabs   enter accept/assign   r route   space resolve   ? less"
	}
	return styles.Help.Width(width).Render(text)
}

func renderCity(state game.GameState, selected int, briefing bool, styles Styles) string {
	if len(state.Districts) == 0 {
		return styles.Muted.Render("No district telemetry.")
	}
	selected = clampIndex(selected, len(state.Districts))
	if briefing {
		return renderDistrictBriefing(state, state.Districts[selected], styles)
	}

	var b strings.Builder
	for i, district := range state.Districts {
		marker := " "
		if i == selected {
			marker = ">"
		}
		name := styles.Accent.Width(24).Render(clipText(district.Name, 24))
		faction := styles.InlineCode.Width(9).Align(lipgloss.Right).Render(formatFactionControl(district.FactionControl))
		fmt.Fprintf(&b, "%s%s%s%s\n", styles.PanelText.Render(marker), name, styles.PanelText.Render(" "), faction)
		fmt.Fprintln(&b, styles.PanelText.Render(fmt.Sprintf("  SURV %d   TRAF %d   DNGR %d   SGNL %d",
			district.Surveillance,
			district.Traffic,
			district.Danger,
			district.SignalQuality,
		)))
	}
	return strings.TrimRight(b.String(), "\n")
}

func renderDistrictBriefing(state game.GameState, district game.District, styles Styles) string {
	lines := []string{
		styles.Accent.Render(district.Name),
	}
	for _, line := range wrapText(district.Description, 36, 2) {
		lines = append(lines, styles.Muted.Render(line))
	}
	lines = append(lines,
		styles.PanelText.Render(" "),
		styles.PanelText.Render(fmt.Sprintf("Control: %s", formatFactionControl(district.FactionControl))),
		styles.PanelText.Render(fmt.Sprintf("SURV %d  TRAF %d  DNGR %d  SGNL %d",
			district.Surveillance,
			district.Traffic,
			district.Danger,
			district.SignalQuality,
		)),
		styles.PanelText.Render(" "),
		styles.Accent.Render("Briefing"),
		styles.PanelText.Render("Pressure: "+districtPressureSummary(district)),
		styles.PanelText.Render(fmt.Sprintf("Jobs touch district: %d", districtJobCount(state, district.ID))),
		styles.PanelText.Render(" "),
		styles.Muted.Render("esc returns to sector list."),
	)
	return strings.Join(lines, "\n")
}

func renderJobs(state game.GameState, selected int, styles Styles) string {
	if len(state.AvailableJobs) == 0 {
		lines := renderNoJobsState(state, styles)
		if len(state.AcceptedJobs) > 0 {
			lines = append(lines, styles.PanelText.Render(" "), styles.Accent.Render("Accepted"))
			for _, job := range state.AcceptedJobs {
				lines = append(lines, styles.PanelText.Render("  "+job.Title))
			}
		}
		return strings.Join(lines, "\n")
	}

	var b strings.Builder
	districts := districtNames(state)
	for i, job := range state.AvailableJobs {
		factor := "none"
		if len(job.RiskFactors) > 0 {
			factor = job.RiskFactors[0]
		}
		marker := " "
		if i == clampIndex(selected, len(state.AvailableJobs)) {
			marker = ">"
		}
		fmt.Fprintf(&b, "%s%s%s\n",
			styles.PanelText.Render(marker),
			styles.Accent.Render(job.Title),
			styles.PanelText.Render(fmt.Sprintf("  %s", formatCargo(job.Cargo))),
		)
		fmt.Fprintln(&b, styles.PanelText.Render(fmt.Sprintf("  %s -> %s",
			districts[job.Origin],
			districts[job.Destination],
		)))
		fmt.Fprintln(&b, styles.PanelText.Render(fmt.Sprintf("  pay %d  due %dT  r%d",
			job.Payout,
			job.DeadlineTurns,
			len(job.Routes),
		))+styles.PanelText.Render("  f:")+styles.Warning.Render(shortFactor(factor)))
		if job.ID != state.AvailableJobs[len(state.AvailableJobs)-1].ID {
			fmt.Fprintln(&b, styles.PanelText.Render(" "))
		}
	}
	if len(state.AcceptedJobs) > 0 {
		fmt.Fprintln(&b, styles.PanelText.Render(" "))
		fmt.Fprintln(&b, styles.Accent.Render("Accepted"))
		for _, job := range state.AcceptedJobs {
			fmt.Fprintln(&b, styles.PanelText.Render("  "+job.Title))
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func renderNoJobsState(state game.GameState, styles Styles) []string {
	lines := []string{
		styles.Muted.Render("No contracts posted."),
		styles.PanelText.Render(" "),
	}
	switch state.Phase {
	case game.PhaseMessages:
		lines = append(lines, styles.PanelText.Render("Review messages to refresh the board."))
	case game.PhaseJobs:
		lines = append(lines, styles.PanelText.Render("Fresh postings are pending on the wire."))
	case game.PhaseDispatch:
		if len(state.AcceptedJobs) > 0 || len(state.ActiveJobs) > 0 {
			lines = append(lines, styles.PanelText.Render("Work is already on the desk."))
		} else {
			lines = append(lines, styles.PanelText.Render("Advance the phase for new postings."))
		}
	case game.PhaseGameOver:
		lines = append(lines, styles.PanelText.Render("Run is closed. No new work."))
	default:
		lines = append(lines, styles.PanelText.Render("Dispatch wire is quiet."))
	}
	return lines
}

func renderDetail(state game.GameState, focused int, selectedJob int, selectedRunner int, selectedRoute int, selectedMessage int, selectedResponse int, notice string, styles Styles) string {
	if len(state.AcceptedJobs) > 0 {
		return renderAcceptedJobDetail(state, state.AcceptedJobs[0], selectedRunner, selectedRoute, notice, styles)
	}
	if len(state.AvailableJobs) > 0 && (focused == focusJobs || focused == focusDetail) {
		return renderJobDetail(state, state.AvailableJobs[clampIndex(selectedJob, len(state.AvailableJobs))], styles)
	}
	if len(state.Runners) > 0 && focused == focusRunners {
		return renderRunnerDetail(state, state.Runners[clampIndex(selectedRunner, len(state.Runners))], notice, styles)
	}
	if focused == focusMessages {
		return renderMessageDetail(state, selectedMessage, selectedResponse, notice, styles)
	}

	lines := []string{
		styles.Accent.Render("Desk state"),
		styles.PanelText.Render(fmt.Sprintf("Phase: %s", state.Phase)),
		styles.PanelText.Render(fmt.Sprintf("Seed: %d", state.RandomSeed)),
		styles.PanelText.Render(fmt.Sprintf("Active runs: %d", len(state.ActiveJobs))),
		styles.PanelText.Render(" "),
		styles.Accent.Render("Current focus"),
		styles.PanelText.Render("Use tab to inspect panels."),
		styles.PanelText.Render(" "),
		styles.Muted.Render("Exact risk math stays off-screen."),
		styles.Muted.Render("Route factors will appear here."),
	}
	if notice != "" {
		lines = append([]string{styles.Warning.Render(notice), styles.PanelText.Render(" ")}, lines...)
	}
	if len(state.LastResults) > 0 {
		lines = append(lines, styles.PanelText.Render(" "), styles.Accent.Render("Last results"))
		for _, result := range state.LastResults {
			lines = append(lines, styles.PanelText.Render(fmt.Sprintf("  %s: %s", result.JobTitle, result.Outcome)))
		}
	}
	return strings.Join(lines, "\n")
}

func renderAcceptedJobDetail(state game.GameState, job game.Job, selectedRunner int, selectedRoute int, notice string, styles Styles) string {
	lines := []string{}
	if notice != "" {
		lines = append(lines, styles.Warning.Render(notice), styles.PanelText.Render(" "))
	}
	lines = append(lines,
		styles.Accent.Render("Pending assignment"),
		styles.PanelText.Render(job.Title),
		styles.PanelText.Render(formatCargo(job.Cargo)),
		styles.PanelText.Render(" "),
		styles.Accent.Render("Route"),
	)
	routeIndex := clampIndex(selectedRoute, len(job.Routes))
	for i, route := range job.Routes {
		marker := " "
		if i == routeIndex {
			marker = ">"
		}
		lines = append(lines, styles.PanelText.Render(marker+" "+formatRouteDetail(route)))
	}
	if runner, ok := selectedPendingRunner(state, selectedRunner); ok {
		load := runnerLoad(state, runner.ID)
		if load > 0 {
			lines = append(lines, styles.PanelText.Render(" "))
			if load >= game.MaxJobsPerRunner {
				lines = append(lines, styles.Critical.Render(fmt.Sprintf("%s bundle full (%d/%d).", runner.Name, load, game.MaxJobsPerRunner)))
			} else {
				lines = append(lines, styles.Warning.Render(fmt.Sprintf("Will bundle with %s (%d/%d).", runner.Name, load, game.MaxJobsPerRunner)))
			}
		}
	}
	lines = append(lines, styles.PanelText.Render(" "), styles.Muted.Render("Select a runner, press enter."))
	return strings.Join(lines, "\n")
}

func renderJobDetail(state game.GameState, job game.Job, styles Styles) string {
	districts := districtNames(state)
	lines := []string{
		styles.Accent.Render(job.Title),
		styles.Muted.Render(clipText(job.ClientMessage, 44)),
		styles.PanelText.Render(clipText(fmt.Sprintf("%s -> %s", districts[job.Origin], districts[job.Destination]), 44)),
		styles.PanelText.Render(fmt.Sprintf("%s  pay %d  due %dT", formatCargo(job.Cargo), job.Payout, job.DeadlineTurns)),
	}
	if len(job.RiskFactors) > 0 {
		lines = append(lines, styles.PanelText.Render("Factors: ")+styles.Warning.Render(formatFactorsShort(job.RiskFactors)))
	}
	lines = append(lines, styles.PanelText.Render(" "), styles.Accent.Render("Route options"))
	for _, route := range job.Routes {
		lines = append(lines, styles.PanelText.Render(formatRouteDetail(route)))
	}
	lines = append(lines, styles.PanelText.Render(" "), styles.Muted.Render("Exact risk stays hidden. Factors stay visible."))
	return strings.Join(lines, "\n")
}

func renderRunnerDetail(state game.GameState, runner game.Runner, notice string, styles Styles) string {
	lines := []string{}
	if notice != "" {
		lines = append(lines, styles.Warning.Render(notice), styles.PanelText.Render(" "))
	}
	lines = append(lines,
		styles.Accent.Render(runner.Name),
		styles.PanelText.Render(runner.Style),
		styles.PanelText.Render(fmt.Sprintf("SPD %d  STL %d  NRV %d  TLK %d", runner.Speed, runner.Stealth, runner.Nerve, runner.Talk)),
		styles.PanelText.Render(fmt.Sprintf("LOY %d  STR %d  CAP %d/%d", runner.Loyalty, runner.Stress, runnerLoad(state, runner.ID), game.MaxJobsPerRunner)),
		styles.PanelText.Render(" "),
	)

	bundle := runnerBundle(state, runner.ID)
	if len(bundle.Jobs) == 0 {
		lines = append(lines, renderNoActiveAssignmentState(state, runner, styles)...)
		return strings.Join(lines, "\n")
	}

	lines = append(lines, styles.Accent.Render(fmt.Sprintf("Bundle %d/%d", len(bundle.Jobs), game.MaxJobsPerRunner)))
	for i, active := range bundle.Jobs {
		lines = append(lines, styles.PanelText.Render(fmt.Sprintf("  B%d %s", i+1, clipText(active.Job.Title, 19)))+styles.Muted.Render(" / "+formatRouteDetail(active.Route)))
	}
	if len(bundle.Penalties) > 0 {
		lines = append(lines, styles.PanelText.Render(" "), styles.Accent.Render("Bundle pressure"))
		for _, penalty := range bundle.Penalties {
			lines = append(lines, styles.Warning.Render("  "+penalty))
		}
	}
	return strings.Join(lines, "\n")
}

func renderRunners(state game.GameState, selected int, styles Styles) string {
	var b strings.Builder
	const runnerNameWidth = 17
	for i, runner := range state.Runners {
		stateText := styles.Accent.Width(9).Render(string(runner.State))
		if runner.State != game.RunnerReady {
			stateText = styles.Warning.Width(9).Render(string(runner.State))
		}
		marker := " "
		if i == clampIndex(selected, len(state.Runners)) {
			marker = ">"
		}
		fmt.Fprintf(&b, "%s%s%s\n",
			styles.PanelText.Render(marker),
			styles.Accent.Width(runnerNameWidth).Render(clipText(runner.Name, runnerNameWidth)),
			stateText,
		)
		fmt.Fprintln(&b, styles.PanelText.Render(fmt.Sprintf("  SPD %-2d STL %-2d NRV %-2d TLK %-2d", runner.Speed, runner.Stealth, runner.Nerve, runner.Talk)))
		fmt.Fprintln(&b, styles.PanelText.Render(fmt.Sprintf("  LOY %-2d STR %-2d CAP %d/%d", runner.Loyalty, runner.Stress, runnerLoad(state, runner.ID), game.MaxJobsPerRunner)))
		bundle := runnerBundle(state, runner.ID)
		for bundleIndex, active := range bundle.Jobs {
			fmt.Fprintln(&b, styles.Muted.Render(fmt.Sprintf("  B%d %s", bundleIndex+1, clipText(active.Job.Title, 31))))
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func renderMessages(state game.GameState, selected int, scroll int, visibleLines int, styles Styles) string {
	if len(state.Messages) == 0 {
		return strings.Join(renderNoMessagesState(state, styles), "\n")
	}

	lines := []string{}
	selected = clampIndex(selected, len(state.Messages))
	for i, message := range state.Messages {
		marker := " "
		if i == selected {
			marker = ">"
		}
		status := ""
		if message.Status == game.MessageResolved {
			status = " done"
		} else if message.Audience != "" {
			status = " open"
		}
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Left,
			styles.PanelText.Render(marker),
			styles.PanelText.Render(fmt.Sprintf("[%02d] ", message.Turn)),
			styles.Accent.Render(message.From),
			styles.PanelText.Render(" / "),
			styles.PanelText.Render(message.Subject),
			styles.Muted.Render(status),
		))
		lines = append(lines, styles.PanelText.Render(fmt.Sprintf("     %s", message.Body)))
	}
	return strings.Join(scrolledLines(lines, scroll, visibleLines), "\n")
}

func renderMessageDetail(state game.GameState, selected int, selectedResponse int, notice string, styles Styles) string {
	message, ok := currentMessage(state, selected)
	if !ok {
		return strings.Join(renderNoMessagesState(state, styles), "\n")
	}

	lines := []string{}
	if notice != "" {
		lines = append(lines, styles.Warning.Render(notice), styles.PanelText.Render(" "))
	}
	lines = append(lines,
		styles.Accent.Render(message.Subject),
		styles.PanelText.Render(fmt.Sprintf("From: %s", message.From)),
		styles.PanelText.Render(fmt.Sprintf("Audience: %s", formatMessageAudience(message.Audience))),
		styles.PanelText.Render(" "),
	)
	for _, line := range wrapText(message.Body, 40, 2) {
		lines = append(lines, styles.PanelText.Render(line))
	}
	if message.Status == game.MessageResolved {
		lines = append(lines, styles.PanelText.Render(" "), styles.Accent.Render("Resolved"))
		if message.ResolvedBy != "" {
			lines = append(lines, styles.PanelText.Render("By: "+formatResponseID(message.ResolvedBy)))
		}
		if message.Summary != "" {
			lines = append(lines, styles.Muted.Render(clipText(message.Summary, 40)))
		}
		for _, effect := range message.ResolutionEffects {
			lines = append(lines, styles.Warning.Render("  "+effect))
		}
		return strings.Join(lines, "\n")
	}

	responses := messageResponseOptions(message)
	if len(responses) == 0 {
		lines = append(lines, styles.PanelText.Render(" "), styles.Muted.Render("No response required."))
		return strings.Join(lines, "\n")
	}

	lines = append(lines, styles.PanelText.Render(" "), styles.Accent.Render("Responses"))
	selectedResponse = clampIndex(selectedResponse, len(responses))
	for i, response := range responses {
		marker := " "
		if i == selectedResponse {
			marker = ">"
		}
		lines = append(lines, styles.PanelText.Render(marker+" "+response.Label))
		lines = append(lines, styles.Muted.Render("  "+clipText(response.Description, 38)))
	}
	lines = append(lines, styles.PanelText.Render(" "), styles.Muted.Render("r cycles response, enter sends."))
	return strings.Join(lines, "\n")
}

func renderNoActiveAssignmentState(state game.GameState, runner game.Runner, styles Styles) []string {
	if len(state.AcceptedJobs) > 0 && runner.State == game.RunnerReady {
		return []string{
			styles.Muted.Render("No active assignment."),
			styles.PanelText.Render("Ready for the pending job."),
		}
	}
	if runner.State != game.RunnerReady {
		return []string{
			styles.Muted.Render("No active assignment."),
			styles.PanelText.Render("Runner is not ready for dispatch."),
		}
	}
	return []string{
		styles.Muted.Render("No active assignment."),
		styles.PanelText.Render("Accept a job, then assign here."),
	}
}

func renderNoMessagesState(state game.GameState, styles Styles) []string {
	lines := []string{
		styles.Muted.Render("No messages."),
		styles.PanelText.Render(" "),
	}
	switch state.Phase {
	case game.PhaseReports:
		lines = append(lines, styles.PanelText.Render("Reports will file on the next advance."))
	case game.PhaseCityUpdate:
		lines = append(lines, styles.PanelText.Render("City update will queue the next brief."))
	case game.PhaseMessages:
		lines = append(lines, styles.PanelText.Render("Switchboard is quiet for this turn."))
	case game.PhaseGameOver:
		lines = append(lines, styles.PanelText.Render("Run is closed."))
	default:
		lines = append(lines, styles.PanelText.Render("No calls need a response."))
	}
	return lines
}

func selectedPendingRunner(state game.GameState, selected int) (game.Runner, bool) {
	if len(state.Runners) == 0 {
		return game.Runner{}, false
	}
	runner := state.Runners[clampIndex(selected, len(state.Runners))]
	if runnerLoad(state, runner.ID) > 0 && runner.State == game.RunnerOnJob {
		return runner, true
	}
	return game.Runner{}, false
}

func currentMessage(state game.GameState, selected int) (game.Message, bool) {
	if len(state.Messages) == 0 {
		return game.Message{}, false
	}
	return state.Messages[clampIndex(selected, len(state.Messages))], true
}

func messageResponseOptions(message game.Message) []game.MessageResponseAction {
	if message.Status == game.MessageResolved || message.Audience == "" {
		return nil
	}
	if len(message.Responses) > 0 {
		return message.Responses
	}
	return game.MessageResponseActionsFor(message.Audience)
}

func panel(title string, body string, width int, height int, focused bool, styles Styles) string {
	style := styles.Panel
	if focused {
		style = styles.PanelFocus
	}

	frameW, frameH := style.GetFrameSize()
	contentW := max(1, width-frameW)
	bodyH := max(0, height-frameH-2)
	body = strings.Join(fitLines(strings.Split(body, "\n"), bodyH), "\n")
	content := styles.PanelTitle.Render(title) + "\n" + styles.Divider.Render(strings.Repeat("─", contentW)) + "\n" + body
	return style.Width(width).Height(height).Render(content)
}

func panelBodyHeight(style lipgloss.Style, height int) int {
	_, frameH := style.GetFrameSize()
	return max(0, height-frameH-2)
}

func scrolledLines(lines []string, scroll int, visible int) []string {
	if visible <= 0 {
		return nil
	}
	if len(lines) <= visible {
		return lines
	}
	maxScroll := len(lines) - visible
	scroll = clampInt(scroll, 0, maxScroll)
	return lines[scroll : scroll+visible]
}

func fitLines(lines []string, height int) []string {
	if height <= 0 {
		return nil
	}
	if len(lines) > height {
		return lines[:height]
	}
	return lines
}

func districtNames(state game.GameState) map[game.DistrictID]string {
	names := make(map[game.DistrictID]string, len(state.Districts))
	for _, district := range state.Districts {
		names[district.ID] = district.Name
	}
	return names
}

func formatCargo(cargo game.CargoType) string {
	switch cargo {
	case game.CargoDataShard:
		return "data shard"
	case game.CargoMedicalCooler:
		return "medical cooler"
	case game.CargoWitness:
		return "witness"
	case game.CargoContrabandPackage:
		return "contraband"
	case game.CargoCorporatePrototype:
		return "prototype"
	default:
		return strings.ReplaceAll(string(cargo), "_", " ")
	}
}

func formatMessageAudience(audience game.MessageAudience) string {
	if audience == "" {
		return "none"
	}
	return strings.ReplaceAll(string(audience), "_", " ")
}

func formatResponseID(responseID game.MessageResponseActionID) string {
	if response, ok := game.MessageResponseActionFor(responseID); ok {
		return response.Label
	}
	return strings.ReplaceAll(string(responseID), "_", " ")
}

func shortFactor(factor string) string {
	replacer := strings.NewReplacer(
		"destination surveillance", "dest surv",
		"traffic pressure", "traffic",
		"violent district", "danger",
		"weak signal", "signal",
		"cargo integrity", "cargo",
		"client urgency", "urgent",
		"corporate trace", "corp trace",
		"witness nerves", "witness",
		"checkpoint exposure", "checkpoint",
		"corporate trackers", "trackers",
		"intercept interest", "intercept",
		"union politics", "union",
		"security audit", "audit",
		"bad witnesses", "witnesses",
		"attention magnet", "heat sink",
		"destination complexity", "complex",
		"curfew patrols", "curfew",
		"medical spoilage", "spoilage",
		"betrayal risk", "betrayal",
		"unclear package", "unclear",
	)
	return replacer.Replace(factor)
}

func formatFactorsShort(factors []string) string {
	limit := min(len(factors), 3)
	short := make([]string, 0, limit)
	for _, factor := range factors[:limit] {
		short = append(short, shortFactor(factor))
	}
	return strings.Join(short, ", ")
}

func districtPressureSummary(district game.District) string {
	pressure := []string{}
	if district.Surveillance >= 4 {
		pressure = append(pressure, "watched")
	}
	if district.Traffic >= 4 {
		pressure = append(pressure, "crowded")
	}
	if district.Danger >= 4 {
		pressure = append(pressure, "volatile")
	}
	if len(pressure) == 0 {
		return "manageable"
	}
	return strings.Join(pressure, ", ")
}

func districtSignalSummary(district game.District) string {
	switch {
	case district.SignalQuality >= 4:
		return "clean comms"
	case district.SignalQuality <= 2:
		return "weak comms"
	default:
		return "patchy but usable"
	}
}

func districtJobCount(state game.GameState, districtID game.DistrictID) int {
	count := 0
	for _, job := range state.AvailableJobs {
		if job.Origin == districtID || job.Destination == districtID {
			count++
		}
	}
	for _, job := range state.AcceptedJobs {
		if job.Origin == districtID || job.Destination == districtID {
			count++
		}
	}
	for _, active := range state.ActiveJobs {
		if active.Job.Origin == districtID || active.Job.Destination == districtID {
			count++
		}
	}
	return count
}

func formatRouteDetail(route game.Route) string {
	traits := route.Traits
	if len(traits) > 2 {
		traits = traits[:2]
	}
	return clipText(fmt.Sprintf("%s  %dT  %s", route.Name, route.TimeCost, strings.Join(traits, ", ")), 44)
}

func clipText(value string, width int) string {
	if lipgloss.Width(value) <= width {
		return value
	}
	if width <= 1 {
		return value[:0]
	}
	return value[:width-1] + "…"
}

func wrapText(value string, width int, maxLines int) []string {
	if maxLines <= 0 {
		return nil
	}
	words := strings.Fields(value)
	if len(words) == 0 {
		return []string{""}
	}

	lines := []string{}
	current := ""
	for _, word := range words {
		next := word
		if current != "" {
			next = current + " " + word
		}
		if lipgloss.Width(next) <= width {
			current = next
			continue
		}
		lines = append(lines, clipText(current, width))
		current = word
		if len(lines) == maxLines {
			return lines
		}
	}
	if current != "" && len(lines) < maxLines {
		lines = append(lines, clipText(current, width))
	}
	return lines
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clampIndex(index int, length int) int {
	if length <= 0 {
		return 0
	}
	if index < 0 {
		return 0
	}
	if index >= length {
		return length - 1
	}
	return index
}

func clampInt(value int, minValue int, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func runnerLoad(state game.GameState, runnerID game.RunnerID) int {
	count := 0
	for _, active := range state.ActiveJobs {
		if active.RunnerID == runnerID {
			count++
		}
	}
	return count
}

func runnerBundle(state game.GameState, runnerID game.RunnerID) game.Bundle {
	for _, bundle := range state.Bundles {
		if bundle.RunnerID == runnerID {
			return bundle
		}
	}
	return game.Bundle{}
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
