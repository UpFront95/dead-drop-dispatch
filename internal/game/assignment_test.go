package game_test

import (
	"errors"
	"testing"

	"dead-drop-dispatch/internal/content"
	"dead-drop-dispatch/internal/game"
)

func TestAcceptJobMovesAvailableJobToAccepted(t *testing.T) {
	state := content.InitialGameState(42)
	job := state.AvailableJobs[0]

	if err := game.AcceptJob(&state, job.ID); err != nil {
		t.Fatalf("AcceptJob returned error: %v", err)
	}

	if containsJob(state.AvailableJobs, job.ID) {
		t.Fatalf("available jobs still contains accepted job %s", job.ID)
	}
	if !containsJob(state.AcceptedJobs, job.ID) {
		t.Fatalf("accepted jobs missing %s", job.ID)
	}
	if len(state.EventLog) == 0 || state.EventLog[len(state.EventLog)-1].Text == "" {
		t.Fatal("accepting a job should append a log entry")
	}
}

func TestAssignAcceptedJobMovesJobToActiveAndMarksRunnerBusy(t *testing.T) {
	state := content.InitialGameState(42)
	job := state.AvailableJobs[0]
	runner := state.Runners[0]
	route := job.Routes[0]

	if err := game.AcceptJob(&state, job.ID); err != nil {
		t.Fatalf("AcceptJob returned error: %v", err)
	}
	if err := game.AssignAcceptedJob(&state, job.ID, runner.ID, route.ID); err != nil {
		t.Fatalf("AssignAcceptedJob returned error: %v", err)
	}

	if containsJob(state.AcceptedJobs, job.ID) {
		t.Fatalf("accepted jobs still contains assigned job %s", job.ID)
	}
	if got, want := len(state.ActiveJobs), 1; got != want {
		t.Fatalf("active jobs = %d, want %d", got, want)
	}
	active := state.ActiveJobs[0]
	if active.JobID != job.ID || active.RunnerID != runner.ID || active.RouteID != route.ID {
		t.Fatalf("active job = %+v, want job %s runner %s route %s", active, job.ID, runner.ID, route.ID)
	}
	if state.Runners[0].State != game.RunnerOnJob {
		t.Fatalf("runner state = %q, want %q", state.Runners[0].State, game.RunnerOnJob)
	}
	if got, want := len(state.Bundles), 1; got != want {
		t.Fatalf("bundles = %d, want %d", got, want)
	}
	if state.ActiveJobs[0].Job.ID != job.ID {
		t.Fatalf("active job snapshot id = %q, want %q", state.ActiveJobs[0].Job.ID, job.ID)
	}
	if state.ActiveJobs[0].Route.ID != route.ID {
		t.Fatalf("active route snapshot id = %q, want %q", state.ActiveJobs[0].Route.ID, route.ID)
	}
}

func TestAssignAcceptedJobBundlesSecondCompatibleJob(t *testing.T) {
	state := assignmentState(
		testJob("job-1", game.CargoDataShard, "northline", "floodglass", game.RouteServiceTunnels),
		testJob("job-2", game.CargoMedicalCooler, "floodglass", "port_kestrel", game.RouteServiceTunnels),
	)
	runner := state.Runners[0]

	for _, job := range append([]game.Job(nil), state.AvailableJobs...) {
		if err := game.AcceptJob(&state, job.ID); err != nil {
			t.Fatalf("AcceptJob(%s) returned error: %v", job.ID, err)
		}
	}
	if err := game.AssignAcceptedJob(&state, "job-1", runner.ID, "job-1-r1"); err != nil {
		t.Fatalf("first AssignAcceptedJob returned error: %v", err)
	}
	if err := game.AssignAcceptedJob(&state, "job-2", runner.ID, "job-2-r1"); err != nil {
		t.Fatalf("second AssignAcceptedJob returned error: %v", err)
	}

	if got, want := len(state.ActiveJobs), 2; got != want {
		t.Fatalf("active jobs = %d, want %d", got, want)
	}
	if got, want := len(state.Bundles), 1; got != want {
		t.Fatalf("bundles = %d, want %d", got, want)
	}
	if got, want := len(state.Bundles[0].Jobs), 2; got != want {
		t.Fatalf("bundle jobs = %d, want %d", got, want)
	}
	if !containsString(state.Bundles[0].Penalties, "destination complexity") {
		t.Fatalf("bundle penalties = %#v, want destination complexity", state.Bundles[0].Penalties)
	}
}

func TestAssignAcceptedJobRejectsIncompatibleBundle(t *testing.T) {
	state := assignmentState(
		testJob("job-1", game.CargoDataShard, "northline", "floodglass", game.RouteServiceTunnels),
		testJob("job-2", game.CargoWitness, "port_kestrel", "crown_verge", game.RouteMainArtery),
	)
	runner := state.Runners[0]

	for _, job := range append([]game.Job(nil), state.AvailableJobs...) {
		if err := game.AcceptJob(&state, job.ID); err != nil {
			t.Fatalf("AcceptJob(%s) returned error: %v", job.ID, err)
		}
	}
	if err := game.AssignAcceptedJob(&state, "job-1", runner.ID, "job-1-r1"); err != nil {
		t.Fatalf("first AssignAcceptedJob returned error: %v", err)
	}
	err := game.AssignAcceptedJob(&state, "job-2", runner.ID, "job-2-r1")
	if !errors.Is(err, game.ErrBundleIncompatible) {
		t.Fatalf("AssignAcceptedJob error = %v, want %v", err, game.ErrBundleIncompatible)
	}
}

func TestAssignAcceptedJobRejectsRunnerAtCapacity(t *testing.T) {
	state := assignmentState(
		testJob("job-1", game.CargoDataShard, "northline", "floodglass", game.RouteServiceTunnels),
		testJob("job-2", game.CargoDataShard, "floodglass", "port_kestrel", game.RouteServiceTunnels),
		testJob("job-3", game.CargoDataShard, "port_kestrel", "northline", game.RouteServiceTunnels),
	)
	runner := state.Runners[0]

	for _, job := range append([]game.Job(nil), state.AvailableJobs...) {
		if err := game.AcceptJob(&state, job.ID); err != nil {
			t.Fatalf("AcceptJob(%s) returned error: %v", job.ID, err)
		}
	}
	if err := game.AssignAcceptedJob(&state, "job-1", runner.ID, "job-1-r1"); err != nil {
		t.Fatalf("first AssignAcceptedJob returned error: %v", err)
	}
	if err := game.AssignAcceptedJob(&state, "job-2", runner.ID, "job-2-r1"); err != nil {
		t.Fatalf("second AssignAcceptedJob returned error: %v", err)
	}
	err := game.AssignAcceptedJob(&state, "job-3", runner.ID, "job-3-r1")
	if !errors.Is(err, game.ErrRunnerAtCapacity) {
		t.Fatalf("AssignAcceptedJob error = %v, want %v", err, game.ErrRunnerAtCapacity)
	}
}

func TestAssignAcceptedJobRejectsUnavailableRunner(t *testing.T) {
	state := content.InitialGameState(42)
	job := state.AvailableJobs[0]
	runner := state.Runners[0]
	route := job.Routes[0]
	state.Runners[0].State = game.RunnerInjured

	if err := game.AcceptJob(&state, job.ID); err != nil {
		t.Fatalf("AcceptJob returned error: %v", err)
	}
	err := game.AssignAcceptedJob(&state, job.ID, runner.ID, route.ID)
	if !errors.Is(err, game.ErrRunnerBusy) {
		t.Fatalf("AssignAcceptedJob error = %v, want %v", err, game.ErrRunnerBusy)
	}
}

func TestAssignAcceptedJobRejectsUnknownRoute(t *testing.T) {
	state := content.InitialGameState(42)
	job := state.AvailableJobs[0]
	runner := state.Runners[0]

	if err := game.AcceptJob(&state, job.ID); err != nil {
		t.Fatalf("AcceptJob returned error: %v", err)
	}
	err := game.AssignAcceptedJob(&state, job.ID, runner.ID, "not-a-route")
	if !errors.Is(err, game.ErrRouteNotFound) {
		t.Fatalf("AssignAcceptedJob error = %v, want %v", err, game.ErrRouteNotFound)
	}
}

func containsJob(jobs []game.Job, jobID string) bool {
	for _, job := range jobs {
		if job.ID == jobID {
			return true
		}
	}
	return false
}

func assignmentState(jobs ...game.Job) game.GameState {
	state := content.InitialGameState(42)
	state.AvailableJobs = jobs
	state.AcceptedJobs = nil
	state.ActiveJobs = nil
	state.Bundles = nil
	return state
}

func testJob(id string, cargo game.CargoType, origin game.DistrictID, destination game.DistrictID, routeType game.RouteType) game.Job {
	return game.Job{
		ID:            id,
		Title:         id,
		Cargo:         cargo,
		Origin:        origin,
		Destination:   destination,
		DeadlineTurns: 1,
		Payout:        100,
		Routes: []game.Route{{
			ID:        id + "-r1",
			Type:      routeType,
			Name:      "Test route",
			Districts: []game.DistrictID{origin, destination},
			TimeCost:  1,
		}},
	}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
