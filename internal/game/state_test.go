package game_test

import (
	"testing"

	"dead-drop-dispatch/internal/content"
	"dead-drop-dispatch/internal/game"
)

func TestInitialGameState(t *testing.T) {
	state := content.InitialGameState(42)

	if state.Turn != game.FirstTurn {
		t.Fatalf("turn = %d, want %d", state.Turn, game.FirstTurn)
	}
	if state.Night != game.FirstNight {
		t.Fatalf("night = %d, want %d", state.Night, game.FirstNight)
	}
	if state.Credits != game.StartingCredits {
		t.Fatalf("credits = %d, want %d", state.Credits, game.StartingCredits)
	}
	if state.Heat != game.StartingHeat {
		t.Fatalf("heat = %d, want %d", state.Heat, game.StartingHeat)
	}
	if state.DispatchIntegrity != game.StartingDispatchIntegrity {
		t.Fatalf("dispatch integrity = %d, want %d", state.DispatchIntegrity, game.StartingDispatchIntegrity)
	}
	if state.Phase != game.PhaseDispatch {
		t.Fatalf("phase = %q, want %q", state.Phase, game.PhaseDispatch)
	}
	if state.RandomSeed != 42 {
		t.Fatalf("seed = %d, want 42", state.RandomSeed)
	}
}

func TestInitialContentCounts(t *testing.T) {
	state := content.InitialGameState(42)

	if got, want := len(state.Districts), 6; got != want {
		t.Fatalf("district count = %d, want %d", got, want)
	}
	if got, want := len(state.Runners), 3; got != want {
		t.Fatalf("runner count = %d, want %d", got, want)
	}
	if got, want := len(state.Factions), 4; got != want {
		t.Fatalf("faction count = %d, want %d", got, want)
	}
	if got, want := len(state.JobTemplates), 10; got != want {
		t.Fatalf("job template count = %d, want %d", got, want)
	}
	if got, want := len(state.Messages), 1; got != want {
		t.Fatalf("message count = %d, want %d", got, want)
	}
	if got, want := len(state.EventLog), 1; got != want {
		t.Fatalf("event log count = %d, want %d", got, want)
	}
	if got, want := len(state.AvailableJobs), game.DefaultJobsPerTurn; got != want {
		t.Fatalf("available job count = %d, want %d", got, want)
	}
}

func TestInitialReferencesAreValid(t *testing.T) {
	state := content.InitialGameState(42)
	factions := map[game.FactionID]bool{}
	for _, faction := range state.Factions {
		factions[faction.ID] = true
	}

	for _, district := range state.Districts {
		if !factions[district.FactionControl] {
			t.Fatalf("district %s controls unknown faction %s", district.ID, district.FactionControl)
		}
		assertStatRange(t, "surveillance", district.Surveillance)
		assertStatRange(t, "traffic", district.Traffic)
		assertStatRange(t, "danger", district.Danger)
		assertStatRange(t, "signal quality", district.SignalQuality)
	}

	for _, runner := range state.Runners {
		if runner.State != game.RunnerReady {
			t.Fatalf("runner %s state = %q, want %q", runner.ID, runner.State, game.RunnerReady)
		}
		assertStatRange(t, "speed", runner.Speed)
		assertStatRange(t, "stealth", runner.Stealth)
		assertStatRange(t, "nerve", runner.Nerve)
		assertStatRange(t, "talk", runner.Talk)
		assertStatRange(t, "loyalty", runner.Loyalty)
		if runner.Stress != 0 {
			t.Fatalf("runner %s stress = %d, want 0", runner.ID, runner.Stress)
		}
	}
}

func TestEvaluateRunStatusReturnsInProgressForInitialState(t *testing.T) {
	state := content.InitialGameState(42)

	status := game.EvaluateRunStatus(state)

	if status.State != game.RunInProgress {
		t.Fatalf("state = %q, want %q", status.State, game.RunInProgress)
	}
	if status.Reason != game.RunEndNone {
		t.Fatalf("reason = %q, want %q", status.Reason, game.RunEndNone)
	}
}

func TestEvaluateRunStatusDetectsVictory(t *testing.T) {
	state := content.InitialGameState(42)
	state.Night = state.RunNights
	state.Turn = state.TurnsPerNight + 1
	state.Credits = game.VictoryCreditTarget

	status := game.EvaluateRunStatus(state)

	if status.State != game.RunWon {
		t.Fatalf("state = %q, want %q", status.State, game.RunWon)
	}
	if status.Reason != game.RunEndVictory {
		t.Fatalf("reason = %q, want %q", status.Reason, game.RunEndVictory)
	}
}

func TestEvaluateRunStatusDetectsFinalCreditShortfall(t *testing.T) {
	state := content.InitialGameState(42)
	state.Night = state.RunNights
	state.Turn = state.TurnsPerNight + 1
	state.Credits = game.VictoryCreditTarget - 1

	status := game.EvaluateRunStatus(state)

	assertRunLost(t, status, game.RunEndShortfall)
}

func TestEvaluateRunStatusDetectsBankruptcy(t *testing.T) {
	state := content.InitialGameState(42)
	state.Credits = game.BankruptcyCreditFloor - 1

	status := game.EvaluateRunStatus(state)

	assertRunLost(t, status, game.RunEndBankrupt)
}

func TestEvaluateRunStatusDetectsBurned(t *testing.T) {
	state := content.InitialGameState(42)
	state.Heat = game.MaximumHeat

	status := game.EvaluateRunStatus(state)

	assertRunLost(t, status, game.RunEndBurned)
}

func TestEvaluateRunStatusDetectsCollapse(t *testing.T) {
	state := content.InitialGameState(42)
	state.DispatchIntegrity = game.DispatchIntegrityFailure

	status := game.EvaluateRunStatus(state)

	assertRunLost(t, status, game.RunEndCollapse)
}

func TestEvaluateRunStatusDetectsRosterLoss(t *testing.T) {
	state := content.InitialGameState(42)
	for i := range state.Runners {
		state.Runners[i].State = game.RunnerMissing
	}

	status := game.EvaluateRunStatus(state)

	assertRunLost(t, status, game.RunEndRosterLoss)
}

func TestEvaluateRunStatusDoesNotTreatRunnersOnJobsAsRosterLoss(t *testing.T) {
	state := content.InitialGameState(42)
	for i := range state.Runners {
		state.Runners[i].State = game.RunnerOnJob
	}
	state.ActiveJobs = []game.ActiveJob{{JobID: "active"}}

	status := game.EvaluateRunStatus(state)

	if status.State != game.RunInProgress {
		t.Fatalf("state = %q, want %q", status.State, game.RunInProgress)
	}
	if status.Reason != game.RunEndNone {
		t.Fatalf("reason = %q, want %q", status.Reason, game.RunEndNone)
	}
}

func TestEvaluateRunStatusDetectsFactionLockout(t *testing.T) {
	state := content.InitialGameState(42)
	for i := 0; i < game.FactionLockoutCount; i++ {
		state.Factions[i].Suspicion = game.HostileFactionSuspicion
	}

	status := game.EvaluateRunStatus(state)

	assertRunLost(t, status, game.RunEndFactionLockout)
}

func TestEvaluateRunStatusLossTakesPriorityOverVictory(t *testing.T) {
	state := content.InitialGameState(42)
	state.Night = state.RunNights
	state.Turn = state.TurnsPerNight + 1
	state.Credits = game.VictoryCreditTarget
	state.Heat = game.MaximumHeat

	status := game.EvaluateRunStatus(state)

	assertRunLost(t, status, game.RunEndBurned)
}

func assertStatRange(t *testing.T, name string, value int) {
	t.Helper()
	if value < 1 || value > 5 {
		t.Fatalf("%s = %d, want 1..5", name, value)
	}
}

func assertRunLost(t *testing.T, status game.RunStatus, reason game.RunEndReason) {
	t.Helper()
	if status.State != game.RunLost {
		t.Fatalf("state = %q, want %q", status.State, game.RunLost)
	}
	if status.Reason != reason {
		t.Fatalf("reason = %q, want %q", status.Reason, reason)
	}
	if status.Summary == "" {
		t.Fatal("summary should not be empty")
	}
}
