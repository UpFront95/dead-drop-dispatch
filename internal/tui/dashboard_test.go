package tui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"dead-drop-dispatch/internal/content"
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
