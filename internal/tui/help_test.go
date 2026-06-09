package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

	"dead-drop-dispatch/internal/content"
)

func TestRenderHelpViewShowsCurrentControls(t *testing.T) {
	view := RenderHelpView(HelpView{
		Width:  TargetWidth,
		Height: TargetHeight,
		Styles: NewStyles(),
	})

	content := ansi.Strip(view.Content)
	for _, want := range []string{
		"HELP",
		"CONTROLS",
		"Dashboard-first dispatch",
		"tab",
		"focus next dashboard panel",
		"shift+tab",
		"[ and ]",
		"1-4",
		"up/down or j/k",
		"enter",
		"accept a job, assign a runner, respond to a message, or open a city brief",
		"r",
		"cycle route options",
		"space",
		"advance the turn phase or resolve active runs",
		"esc",
		"go back from a city brief or help, then cancel the pending accepted job",
		"?",
		"toggle compact footer help",
		"q",
		"quit",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("rendered help view missing %q", want)
		}
	}

	if got := lineCount(content); got != TargetHeight {
		t.Fatalf("rendered help view height = %d, want %d", got, TargetHeight)
	}

	if got := maxLineWidth(content); got > TargetWidth {
		t.Fatalf("rendered help view width = %d, want <= %d", got, TargetWidth)
	}
}

func TestRenderDashboardHelpTabUsesHelpViewContent(t *testing.T) {
	view := RenderDashboard(DashboardView{
		State:     content.InitialGameState(42),
		Width:     TargetWidth,
		Height:    TargetHeight,
		ActiveTab: helpTabIndex,
		Styles:    NewStyles(),
	})

	content := ansi.Strip(view.Content)
	for _, want := range []string{
		"DEAD DROP DISPATCH",
		"HELP",
		"CONTROLS",
		"Dashboard-first dispatch",
		"[ and ]",
		"space",
		"advance the turn phase or resolve active runs",
		"esc",
		"cancel the pending accepted job",
		"tab focus   [ and ] tabs   j/k select   enter accept/assign   r route   space resolve   ? more   q quit",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("rendered dashboard help tab missing %q", want)
		}
	}

	for _, unwanted := range []string{
		"Desk is live. City is listening.",
		"Northline",
	} {
		if strings.Contains(content, unwanted) {
			t.Fatalf("rendered dashboard help tab unexpectedly contained %q", unwanted)
		}
	}

	if got := lineCount(content); got != TargetHeight {
		t.Fatalf("rendered dashboard help tab height = %d, want %d", got, TargetHeight)
	}

	if got := maxLineWidth(content); got > TargetWidth {
		t.Fatalf("rendered dashboard help tab width = %d, want <= %d", got, TargetWidth)
	}
}
