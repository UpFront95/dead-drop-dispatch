package tui

import (
	"strings"
	"testing"

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
		"Desk is live. City is listening.",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("rendered dashboard missing %q", want)
		}
	}

	if got := lineCount(content); got > TargetHeight {
		t.Fatalf("rendered dashboard height = %d, want <= %d", got, TargetHeight)
	}
}

func lineCount(value string) int {
	if value == "" {
		return 0
	}
	return strings.Count(value, "\n") + 1
}
