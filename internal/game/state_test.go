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

	if got, want := len(state.Districts), 5; got != want {
		t.Fatalf("district count = %d, want %d", got, want)
	}
	if got, want := len(state.Runners), 3; got != want {
		t.Fatalf("runner count = %d, want %d", got, want)
	}
	if got, want := len(state.Factions), 4; got != want {
		t.Fatalf("faction count = %d, want %d", got, want)
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

func assertStatRange(t *testing.T, name string, value int) {
	t.Helper()
	if value < 1 || value > 5 {
		t.Fatalf("%s = %d, want 1..5", name, value)
	}
}
