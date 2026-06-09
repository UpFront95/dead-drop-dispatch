package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

	"dead-drop-dispatch/internal/content"
)

func TestRenderRouteTopologyHighlightsSelectedJobAndRoute(t *testing.T) {
	state := content.InitialGameState(42)
	job := state.AvailableJobs[0]

	rendered := RenderRouteTopology(RouteTopologyView{
		State:              state,
		Job:                &job,
		SelectedRouteIndex: 1,
		Width:              58,
		Height:             15,
		Styles:             NewStyles(),
	})

	output := ansi.Strip(rendered)
	for _, want := range []string{
		"District topology",
		"[O] Northline surveil",
		"[D] Crown Verge surveil",
		job.Title,
		"Northline -> Crown Verge",
		"Route: Service tunnels",
		"Path: Northline -> Crown Verge",
		"Traits: concealed, slow, watched destination",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("rendered topology missing %q in\n%s", want, output)
		}
	}

	if got := lineCount(output); got > 15 {
		t.Fatalf("rendered topology height = %d, want <= 15", got)
	}
	if got := maxLineWidth(output); got > 58 {
		t.Fatalf("rendered topology width = %d, want <= 58", got)
	}
}

func TestRenderRouteTopologyWithoutSelectedJobStillShowsDistricts(t *testing.T) {
	rendered := RenderRouteTopology(RouteTopologyView{
		State:  content.InitialGameState(42),
		Width:  52,
		Height: 8,
		Styles: NewStyles(),
	})

	output := ansi.Strip(rendered)
	for _, want := range []string{
		"District topology",
		"[ ] Northline surveil",
		"[ ] Crown Verge surveil",
		"[ ] Floodglass signal",
		"[ ] Port Kestrel danger",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("rendered topology missing %q in\n%s", want, output)
		}
	}
	if strings.Contains(output, "Route:") {
		t.Fatal("topology without selected job should not render route summary")
	}

	if got := maxLineWidth(output); got > 52 {
		t.Fatalf("rendered topology width = %d, want <= 52", got)
	}
}

func TestRenderRouteTopologyFitsSuppliedHeight(t *testing.T) {
	state := content.InitialGameState(42)
	job := state.AvailableJobs[0]

	rendered := RenderRouteTopology(RouteTopologyView{
		State:  state,
		Job:    &job,
		Width:  44,
		Height: 6,
		Styles: NewStyles(),
	})

	output := ansi.Strip(rendered)
	if got := lineCount(output); got != 6 {
		t.Fatalf("rendered topology height = %d, want 6", got)
	}
	if got := maxLineWidth(output); got > 44 {
		t.Fatalf("rendered topology width = %d, want <= 44", got)
	}
	if strings.Contains(output, "Path:") {
		t.Fatal("height-limited topology should crop lower summary lines")
	}
}
