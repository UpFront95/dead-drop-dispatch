package tui

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"

	"dead-drop-dispatch/internal/content"
)

func TestRenderDashboardTargetSizeContent(t *testing.T) {
	view := RenderDashboard(DashboardView{
		State:  content.InitialGameState(42),
		Width:  TargetWidth,
		Height: TargetHeight,
		Styles: NewStyles(),
	})

	content := view.Content
	for _, want := range []string{
		"DEAD DROP DISPATCH",
		"CITY SECTOR",
		"JOB BOARD",
		"RUNNERS",
		"MESSAGE FEED",
		"DETAIL",
		"Northline",
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

	content := view.Content
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
