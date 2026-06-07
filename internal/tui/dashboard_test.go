package tui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"dead-drop-dispatch/internal/content"
	"dead-drop-dispatch/internal/game"
)

func TestRenderDashboardTargetSizeContent(t *testing.T) {
	view := RenderDashboard(DashboardView{
		State:  content.InitialGameState(42),
		Width:  TargetWidth,
		Height: TargetHeight,
		Styles: NewStyles(),
	})

	content := ansi.Strip(view.Content)
	for _, want := range []string{
		"DEAD DROP DISPATCH",
		"DASHBOARD",
		"EQUIPMENT",
		"HELP",
		"CITY SECTOR",
		"JOB BOARD",
		"RUNNERS",
		"MESSAGE FEED",
		"DETAIL",
		"NEXT accept job",
		"ACTION accept highlighted job",
		"RISK",
		"ROUTE",
		">Northline",
		"Northline",
		"Ashgate Yard",
		"Mira Vale",
		" r",
		"f:",
		"Desk is live. City is listening.",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("rendered dashboard missing %q", want)
		}
	}

	if got := lineCount(content); got != TargetHeight {
		t.Fatalf("rendered dashboard height = %d, want %d", got, TargetHeight)
	}

	if got := maxLineWidth(content); got > TargetWidth {
		t.Fatalf("rendered dashboard width = %d, want <= %d", got, TargetWidth)
	}
}

func TestRenderDashboardDistrictBriefingStaysInCityPanel(t *testing.T) {
	view := RenderDashboard(DashboardView{
		State:             content.InitialGameState(42),
		Width:             TargetWidth,
		Height:            TargetHeight,
		Focused:           focusCity,
		SelectedDistrict:  1,
		ShowDistrictBrief: true,
		Styles:            NewStyles(),
	})

	content := ansi.Strip(view.Content)
	for _, want := range []string{
		"CITY SECTOR",
		"Floodglass",
		"Low streets, tunnels",
		"Control: CLINIC",
		"SURV 2  TRAF 4  DNGR 3  SGNL 2",
		"Pressure: crowded",
		"Signal: weak comms",
		"Jobs touch district:",
		"esc returns to sector list.",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("rendered district briefing missing %q", want)
		}
	}

	if got := lineCount(content); got != TargetHeight {
		t.Fatalf("rendered district briefing height = %d, want %d", got, TargetHeight)
	}

	if got := maxLineWidth(content); got > TargetWidth {
		t.Fatalf("rendered district briefing width = %d, want <= %d", got, TargetWidth)
	}
}

func TestRenderDashboardFocusedJobsFitsTargetSize(t *testing.T) {
	view := RenderDashboard(DashboardView{
		State:   content.InitialGameState(42),
		Width:   TargetWidth,
		Height:  TargetHeight,
		Focused: focusJobs,
		Styles:  NewStyles(),
	})

	content := ansi.Strip(view.Content)
	for _, want := range []string{
		"Route options",
		"Factors:",
		"Exact risk stays hidden.",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("rendered focused jobs dashboard missing %q", want)
		}
	}

	if got := lineCount(content); got != TargetHeight {
		t.Fatalf("rendered focused jobs dashboard height = %d, want %d", got, TargetHeight)
	}

	if got := maxLineWidth(content); got > TargetWidth {
		t.Fatalf("rendered focused jobs dashboard width = %d, want <= %d", got, TargetWidth)
	}
}

func TestRenderDashboardRunnersUseFixedStatusAndStatRows(t *testing.T) {
	view := RenderDashboard(DashboardView{
		State:  content.InitialGameState(42),
		Width:  TargetWidth,
		Height: TargetHeight,
		Styles: NewStyles(),
	})

	content := ansi.Strip(view.Content)
	for _, want := range []string{
		">Mira Vale        ready",
		" Kaito Senn       ready",
		" Vex Calder       ready",
		"SPD 5  STL 3  NRV 3  TLK 2",
		"LOY 4  STR 0  CAP 0/2",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("rendered runners missing %q in\n%s", want, runnerLines(content))
		}
	}
	if strings.Contains(content, "\n  ready") {
		t.Fatalf("runner status should stay on the name row")
	}
}

func TestRenderDashboardEmptyStatesAreActionable(t *testing.T) {
	state := content.InitialGameState(42)
	state.AvailableJobs = nil
	state.Messages = nil
	state.Phase = game.PhaseMessages

	view := RenderDashboard(DashboardView{
		State:   state,
		Width:   TargetWidth,
		Height:  TargetHeight,
		Focused: focusRunners,
		Styles:  NewStyles(),
	})

	content := ansi.Strip(view.Content)
	for _, want := range []string{
		"No contracts posted.",
		"Review messages to refresh the board.",
		"No messages.",
		"Switchboard is quiet for this turn.",
		"No active assignment.",
		"Accept a job, then assign here.",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("rendered empty state missing %q", want)
		}
	}

	if got := lineCount(content); got != TargetHeight {
		t.Fatalf("rendered empty dashboard height = %d, want %d", got, TargetHeight)
	}

	if got := maxLineWidth(content); got > TargetWidth {
		t.Fatalf("rendered empty dashboard width = %d, want <= %d", got, TargetWidth)
	}
}

func TestRenderDashboardShowsBundleMarkers(t *testing.T) {
	state := bundledDashboardState()

	view := RenderDashboard(DashboardView{
		State:          state,
		Width:          TargetWidth,
		Height:         TargetHeight,
		Focused:        focusRunners,
		SelectedRunner: 0,
		Styles:         NewStyles(),
	})

	content := ansi.Strip(view.Content)
	for _, want := range []string{
		"B1 Data handoff",
		"B2 Med cooler",
		"Bundle 2/2",
		"destination complexity",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("rendered bundled dashboard missing %q", want)
		}
	}

	if got := lineCount(content); got != TargetHeight {
		t.Fatalf("rendered bundled dashboard height = %d, want %d", got, TargetHeight)
	}

	if got := maxLineWidth(content); got > TargetWidth {
		t.Fatalf("rendered bundled dashboard width = %d, want <= %d", got, TargetWidth)
	}
}

func TestRenderDashboardShowsPendingBundleCue(t *testing.T) {
	state := bundledDashboardState()
	state.ActiveJobs = state.ActiveJobs[:1]
	state.Bundles[0].Jobs = state.Bundles[0].Jobs[:1]
	state.Bundles[0].Penalties = nil
	state.AcceptedJobs = []game.Job{testDashboardJob("job-3", "Rush followup", "floodglass", "port_kestrel")}

	view := RenderDashboard(DashboardView{
		State:          state,
		Width:          TargetWidth,
		Height:         TargetHeight,
		Focused:        focusDetail,
		SelectedRunner: 0,
		Styles:         NewStyles(),
	})

	content := ansi.Strip(view.Content)
	for _, want := range []string{
		"Pending assignment",
		"Will bundle with Mira Vale (1/2).",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("rendered pending bundle cue missing %q", want)
		}
	}
}

func runnerLines(content string) string {
	content = ansi.Strip(content)
	lines := []string{}
	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(line, "Mira Vale") ||
			strings.Contains(line, "Kaito Senn") ||
			strings.Contains(line, "Vex Calder") ||
			strings.Contains(line, "SPD ") ||
			strings.Contains(line, "LOY ") {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}

func bundledDashboardState() game.GameState {
	state := content.InitialGameState(42)
	state.AvailableJobs = nil
	state.AcceptedJobs = nil
	state.ActiveJobs = nil
	state.Bundles = nil
	state.Runners[0].State = game.RunnerOnJob

	first := testDashboardJob("job-1", "Data handoff", "northline", "floodglass")
	second := testDashboardJob("job-2", "Med cooler", "floodglass", "port_kestrel")
	firstActive := game.ActiveJob{
		JobID:    first.ID,
		RunnerID: state.Runners[0].ID,
		RouteID:  first.Routes[0].ID,
		Job:      first,
		Route:    first.Routes[0],
	}
	secondActive := game.ActiveJob{
		JobID:    second.ID,
		RunnerID: state.Runners[0].ID,
		RouteID:  second.Routes[0].ID,
		Job:      second,
		Route:    second.Routes[0],
	}
	state.ActiveJobs = []game.ActiveJob{firstActive, secondActive}
	state.Bundles = []game.Bundle{{
		RunnerID:  state.Runners[0].ID,
		Jobs:      []game.ActiveJob{firstActive, secondActive},
		Penalties: []string{"destination complexity"},
	}}
	return state
}

func testDashboardJob(id string, title string, origin game.DistrictID, destination game.DistrictID) game.Job {
	return game.Job{
		ID:            id,
		Title:         title,
		Cargo:         game.CargoDataShard,
		Origin:        origin,
		Destination:   destination,
		DeadlineTurns: 2,
		Payout:        100,
		Routes: []game.Route{{
			ID:        id + "-r1",
			Type:      game.RouteServiceTunnels,
			Name:      "Service tunnels",
			Districts: []game.DistrictID{origin, destination},
			TimeCost:  1,
			Traits:    []string{"low cameras", "tight access"},
		}},
	}
}

func TestRenderDashboardPlaceholderTabLeavesBodyBlank(t *testing.T) {
	view := RenderDashboard(DashboardView{
		State:     content.InitialGameState(42),
		Width:     TargetWidth,
		Height:    TargetHeight,
		ActiveTab: 4,
		Styles:    NewStyles(),
	})

	content := view.Content
	for _, want := range []string{
		"DEAD DROP DISPATCH",
		"EQUIPMENT",
		"tab focus   [ ] tabs   j/k select   enter accept/assign   r route   space resolve   ? more   q quit",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("rendered placeholder tab missing %q", want)
		}
	}

	for _, unwanted := range []string{
		"CITY SECTOR",
		"JOB BOARD",
		"MESSAGE FEED",
		"DETAIL",
		"Northline",
		"Desk is live. City is listening.",
	} {
		if strings.Contains(content, unwanted) {
			t.Fatalf("rendered placeholder tab unexpectedly contained %q", unwanted)
		}
	}

	if got := lineCount(content); got != TargetHeight {
		t.Fatalf("rendered placeholder tab height = %d, want %d", got, TargetHeight)
	}

	if got := maxLineWidth(content); got > TargetWidth {
		t.Fatalf("rendered placeholder tab width = %d, want <= %d", got, TargetWidth)
	}
}

func TestRenderDashboardHelpFooterShowsExpandedShortcuts(t *testing.T) {
	view := RenderDashboard(DashboardView{
		State:    content.InitialGameState(42),
		Width:    TargetWidth,
		Height:   TargetHeight,
		ShowHelp: true,
		Styles:   NewStyles(),
	})

	content := view.Content
	for _, want := range []string{
		"shift+tab prev panel",
		"1-6 jump tabs",
		"arrows move",
		"[ and ] tabs",
		"? less",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("rendered help footer missing %q", want)
		}
	}

	for _, repeated := range []string{
		"tab focus",
		"? more",
	} {
		if strings.Contains(content, repeated) {
			t.Fatalf("rendered help footer repeated compact hint %q", repeated)
		}
	}
}

func lineCount(value string) int {
	if value == "" {
		return 0
	}
	return strings.Count(value, "\n") + 1
}

func maxLineWidth(value string) int {
	maxWidth := 0
	for _, line := range strings.Split(value, "\n") {
		maxWidth = max(maxWidth, lipgloss.Width(line))
	}
	return maxWidth
}
