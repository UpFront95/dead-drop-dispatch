package game_test

import (
	"reflect"
	"testing"

	"dead-drop-dispatch/internal/content"
	"dead-drop-dispatch/internal/game"
)

func TestGenerateJobsIsDeterministic(t *testing.T) {
	state := content.InitialGameState(42)

	first := game.GenerateJobs(state, content.JobTemplates(), game.DefaultJobsPerTurn)
	second := game.GenerateJobs(state, content.JobTemplates(), game.DefaultJobsPerTurn)

	if !reflect.DeepEqual(first, second) {
		t.Fatal("generated jobs are not deterministic for the same state and seed")
	}
}

func TestGeneratedJobsAreValid(t *testing.T) {
	state := content.InitialGameState(42)
	jobs := state.AvailableJobs

	if got, want := len(jobs), game.DefaultJobsPerTurn; got != want {
		t.Fatalf("job count = %d, want %d", got, want)
	}

	districts := map[game.DistrictID]bool{}
	for _, district := range state.Districts {
		districts[district.ID] = true
	}

	factions := map[game.FactionID]bool{}
	for _, faction := range state.Factions {
		factions[faction.ID] = true
	}

	for _, job := range jobs {
		if job.ID == "" {
			t.Fatal("generated job missing id")
		}
		if job.TemplateID == "" {
			t.Fatalf("job %s missing template id", job.ID)
		}
		if job.Title == "" {
			t.Fatalf("job %s missing title", job.ID)
		}
		if job.ClientMessage == "" {
			t.Fatalf("job %s missing client message", job.ID)
		}
		if !districts[job.Origin] {
			t.Fatalf("job %s has unknown origin %s", job.ID, job.Origin)
		}
		if !districts[job.Destination] {
			t.Fatalf("job %s has unknown destination %s", job.ID, job.Destination)
		}
		if job.Origin == job.Destination {
			t.Fatalf("job %s origin and destination both %s", job.ID, job.Origin)
		}
		if !factions[job.Faction] {
			t.Fatalf("job %s has unknown faction %s", job.ID, job.Faction)
		}
		if job.DeadlineTurns < 1 || job.DeadlineTurns > state.TurnsPerNight {
			t.Fatalf("job %s deadline = %d, want 1..%d", job.ID, job.DeadlineTurns, state.TurnsPerNight)
		}
		if job.Payout <= 0 {
			t.Fatalf("job %s payout = %d, want positive", job.ID, job.Payout)
		}
		if len(job.RiskFactors) == 0 {
			t.Fatalf("job %s has no visible risk factors", job.ID)
		}
		if len(job.Routes) < 2 || len(job.Routes) > 4 {
			t.Fatalf("job %s route count = %d, want 2..4", job.ID, len(job.Routes))
		}
		for _, route := range job.Routes {
			if route.ID == "" {
				t.Fatalf("job %s has route missing id", job.ID)
			}
			if route.Name == "" {
				t.Fatalf("job %s route %s missing name", job.ID, route.ID)
			}
			if route.TimeCost < 1 {
				t.Fatalf("job %s route %s time cost = %d, want positive", job.ID, route.ID, route.TimeCost)
			}
			if len(route.Traits) == 0 {
				t.Fatalf("job %s route %s has no visible traits", job.ID, route.ID)
			}
		}
	}
}

func TestRefreshAvailableJobsUsesStateTemplatesAndCurrentClock(t *testing.T) {
	state := content.InitialGameState(42)
	state.Turn = 3
	state.AvailableJobs = []game.Job{{ID: "stale"}}
	state.AcceptedJobs = []game.Job{{ID: "accepted"}}
	state.ActiveJobs = []game.ActiveJob{{JobID: "active"}}
	startLog := len(state.EventLog)

	jobs := game.RefreshAvailableJobs(&state, game.DefaultJobsPerTurn)

	if got, want := len(jobs), game.DefaultJobsPerTurn; got != want {
		t.Fatalf("jobs = %d, want %d", got, want)
	}
	if got, want := len(state.AvailableJobs), game.DefaultJobsPerTurn; got != want {
		t.Fatalf("available jobs = %d, want %d", got, want)
	}
	if state.AvailableJobs[0].ID == "stale" {
		t.Fatal("refresh should replace stale available jobs")
	}
	for _, job := range state.AvailableJobs {
		if job.ID[:6] != "n01t03" {
			t.Fatalf("job id = %q, want current clock prefix n01t03", job.ID)
		}
	}
	if got, want := len(state.AcceptedJobs), 1; got != want {
		t.Fatalf("accepted jobs = %d, want %d", got, want)
	}
	if got, want := len(state.ActiveJobs), 1; got != want {
		t.Fatalf("active jobs = %d, want %d", got, want)
	}
	if got, want := len(state.EventLog), startLog+1; got != want {
		t.Fatalf("event log = %d, want %d", got, want)
	}
}
