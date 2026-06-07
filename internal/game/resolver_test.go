package game_test

import (
	"reflect"
	"testing"

	"dead-drop-dispatch/internal/content"
	"dead-drop-dispatch/internal/game"
)

func TestResolveActiveJobsClearsWorkAndRecordsResults(t *testing.T) {
	state := assignedState(t, 42)
	startCredits := state.Credits

	results := game.ResolveActiveJobs(&state)

	if got, want := len(results), 1; got != want {
		t.Fatalf("results = %d, want %d", got, want)
	}
	if got, want := len(state.ActiveJobs), 0; got != want {
		t.Fatalf("active jobs = %d, want %d", got, want)
	}
	if got, want := len(state.Bundles), 0; got != want {
		t.Fatalf("bundles = %d, want %d", got, want)
	}
	if state.Runners[0].State == game.RunnerOnJob {
		t.Fatalf("runner state still on job after resolution")
	}
	if len(state.LastResults) != len(results) {
		t.Fatalf("last results = %d, want %d", len(state.LastResults), len(results))
	}
	if len(state.Messages) == 0 || state.Messages[len(state.Messages)-1].From != "after-action" {
		t.Fatalf("resolution should append after-action message, got %#v", state.Messages)
	}
	if len(state.EventLog) == 0 || state.EventLog[len(state.EventLog)-1].Text == "" {
		t.Fatal("resolution should append event log entry")
	}
	if state.Credits < startCredits {
		t.Fatalf("credits = %d, want at least %d", state.Credits, startCredits)
	}
}

func TestResolveActiveJobsIsDeterministicForFixedSeed(t *testing.T) {
	first := assignedState(t, 1701)
	second := assignedState(t, 1701)

	firstResults := game.ResolveActiveJobs(&first)
	secondResults := game.ResolveActiveJobs(&second)

	if !reflect.DeepEqual(firstResults, secondResults) {
		t.Fatalf("results differ for fixed seed:\nfirst=%#v\nsecond=%#v", firstResults, secondResults)
	}
	if first.Credits != second.Credits || first.Heat != second.Heat || first.DispatchIntegrity != second.DispatchIntegrity {
		t.Fatalf("state totals differ: first=%+v second=%+v", first, second)
	}
}

func TestResolveActiveJobsAppliesBundlePressure(t *testing.T) {
	state := content.InitialGameState(42)
	state.AvailableJobs = []game.Job{
		testJob("job-1", game.CargoDataShard, "northline", "floodglass", game.RouteServiceTunnels),
		testJob("job-2", game.CargoMedicalCooler, "floodglass", "port_kestrel", game.RouteServiceTunnels),
	}
	runnerID := state.Runners[0].ID
	for _, job := range append([]game.Job(nil), state.AvailableJobs...) {
		if err := game.AcceptJob(&state, job.ID); err != nil {
			t.Fatalf("AcceptJob(%s) returned error: %v", job.ID, err)
		}
	}
	if err := game.AssignAcceptedJob(&state, "job-1", runnerID, "job-1-r1"); err != nil {
		t.Fatalf("first AssignAcceptedJob returned error: %v", err)
	}
	if err := game.AssignAcceptedJob(&state, "job-2", runnerID, "job-2-r1"); err != nil {
		t.Fatalf("second AssignAcceptedJob returned error: %v", err)
	}

	results := game.ResolveActiveJobs(&state)

	if got, want := len(results), 2; got != want {
		t.Fatalf("results = %d, want %d", got, want)
	}
	for _, result := range results {
		if !hasFactor(result.Factors, "bundle pressure") {
			t.Fatalf("result factors = %#v, want bundle pressure", result.Factors)
		}
	}
}

func assignedState(t *testing.T, seed int64) game.GameState {
	t.Helper()
	state := content.InitialGameState(seed)
	job := state.AvailableJobs[0]
	runner := state.Runners[0]
	route := job.Routes[0]
	if err := game.AcceptJob(&state, job.ID); err != nil {
		t.Fatalf("AcceptJob returned error: %v", err)
	}
	if err := game.AssignAcceptedJob(&state, job.ID, runner.ID, route.ID); err != nil {
		t.Fatalf("AssignAcceptedJob returned error: %v", err)
	}
	return state
}

func hasFactor(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
