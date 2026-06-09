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
		"Stat glossary",
		"SURV",
		"Surveillance: detection risk",
		"TRAF",
		"Traffic: crowds, delay, cover",
		"DNGR",
		"Danger: injury pressure",
		"SGNL",
		"Signal: comms and intel quality",
		"SPD",
		"Speed: beats deadlines",
		"STL",
		"Stealth: avoids detection",
		"NRV",
		"Nerve: holds under pressure",
		"TLK",
		"Talk: handles people/gates",
		"LOY",
		"Loyalty: resists betrayal pressure",
		"STR",
		"Stress: strain; high is risky",
		"CAP",
		"Capacity: jobs carried out of 2",
		"RNR",
		"Runners: ready/busy/injured/out",
		"JOB",
		"Jobs: board/accepted/active",
		"DUE",
		"Deadline: nearest job deadline",
		"CG",
		"Cargo: lowest known integrity",
		"tab",
		"focus next dashboard panel",
		"shift+tab",
		"[ and ]",
		"1-4",
		"up/down or j/k",
		"enter",
		"accept job, assign runner, respond, or open city brief",
		"r",
		"cycle route options",
		"space",
		"advance the turn phase or resolve active runs",
		"esc",
		"close brief/help, then cancel pending accepted job",
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
		"Stat glossary",
		"SURV",
		"SPD",
		"RNR",
		"[ and ]",
		"space",
		"advance the turn phase or resolve active runs",
		"esc",
		"cancel pending accepted job",
		"[ and ] tabs   1-4 jump   ? more   esc cancel pending   q quit",
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
