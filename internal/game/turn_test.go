package game_test

import (
	"testing"

	"dead-drop-dispatch/internal/content"
	"dead-drop-dispatch/internal/game"
)

func TestAdvanceTurnPhaseResolvesActiveJobsAndMovesToReports(t *testing.T) {
	state := assignedState(t, 42)

	advance := game.AdvanceTurnPhase(&state)

	if got, want := advance.From, game.PhaseDispatch; got != want {
		t.Fatalf("from phase = %q, want %q", got, want)
	}
	if got, want := advance.To, game.PhaseReports; got != want {
		t.Fatalf("to phase = %q, want %q", got, want)
	}
	if got, want := len(advance.Results), 1; got != want {
		t.Fatalf("results = %d, want %d", got, want)
	}
	if advance.WaitingForInput {
		t.Fatal("advance should not wait when no complication is pending")
	}
	if got, want := len(state.ActiveJobs), 0; got != want {
		t.Fatalf("active jobs = %d, want %d", got, want)
	}
	if got, want := len(state.LastResults), 1; got != want {
		t.Fatalf("last results = %d, want %d", got, want)
	}
}

func TestAdvanceTurnPhaseResolvesWhenAllRunnersAreOnJobs(t *testing.T) {
	state := assignedState(t, 42)
	for i := range state.Runners {
		state.Runners[i].State = game.RunnerOnJob
	}

	advance := game.AdvanceTurnPhase(&state)

	if got, want := advance.To, game.PhaseReports; got != want {
		t.Fatalf("to phase = %q, want %q", got, want)
	}
	if got, want := advance.Status.State, game.RunInProgress; got != want {
		t.Fatalf("run state = %q, want %q", got, want)
	}
	if got, want := len(advance.Results), 1; got != want {
		t.Fatalf("results = %d, want %d", got, want)
	}
	if got, want := len(state.ActiveJobs), 0; got != want {
		t.Fatalf("active jobs = %d, want %d", got, want)
	}
	if got, want := state.Phase, game.PhaseReports; got != want {
		t.Fatalf("state phase = %q, want %q", got, want)
	}
}

func TestAdvanceTurnPhasePausesForPendingComplications(t *testing.T) {
	state := content.InitialGameState(42)
	state.Complications = []game.Complication{{
		ID:     "cmp-1",
		Status: game.ComplicationPending,
	}}

	advance := game.AdvanceTurnPhase(&state)

	if got, want := advance.To, game.PhaseComplications; got != want {
		t.Fatalf("to phase = %q, want %q", got, want)
	}
	if !advance.WaitingForInput {
		t.Fatal("advance should wait for complication response")
	}
	if got, want := state.Turn, game.FirstTurn; got != want {
		t.Fatalf("turn = %d, want %d", got, want)
	}
}

func TestAdvanceTurnPhaseContinuesAfterComplicationsResolve(t *testing.T) {
	state := content.InitialGameState(42)
	state.Phase = game.PhaseComplications
	state.Complications = []game.Complication{{
		ID:     "cmp-1",
		Status: game.ComplicationResolved,
	}}

	advance := game.AdvanceTurnPhase(&state)

	if got, want := advance.To, game.PhaseReports; got != want {
		t.Fatalf("to phase = %q, want %q", got, want)
	}
	if advance.WaitingForInput {
		t.Fatal("advance should not wait after complications are resolved")
	}
}

func TestAdvanceTurnPhaseRefreshesJobsWhenMessagesAdvance(t *testing.T) {
	state := content.InitialGameState(42)
	state.Phase = game.PhaseMessages
	state.Turn = 2
	state.AvailableJobs = []game.Job{{ID: "stale"}}
	state.AcceptedJobs = []game.Job{{ID: "accepted"}}
	state.ActiveJobs = []game.ActiveJob{{JobID: "active"}}

	advance := game.AdvanceTurnPhase(&state)

	if got, want := advance.To, game.PhaseJobs; got != want {
		t.Fatalf("to phase = %q, want %q", got, want)
	}
	if got, want := advance.JobsGenerated, game.DefaultJobsPerTurn; got != want {
		t.Fatalf("jobs generated = %d, want %d", got, want)
	}
	if got, want := len(state.AvailableJobs), game.DefaultJobsPerTurn; got != want {
		t.Fatalf("available jobs = %d, want %d", got, want)
	}
	if state.AvailableJobs[0].ID == "stale" {
		t.Fatal("available jobs should be refreshed")
	}
	if got, want := len(state.AcceptedJobs), 1; got != want {
		t.Fatalf("accepted jobs = %d, want %d", got, want)
	}
	if got, want := len(state.ActiveJobs), 1; got != want {
		t.Fatalf("active jobs = %d, want %d", got, want)
	}
}

func TestAdvanceTurnPhaseFilesReportsAndStartsNextTurnMessages(t *testing.T) {
	state := content.InitialGameState(42)
	state.Phase = game.PhaseReports
	startMessages := len(state.Messages)

	advance := game.AdvanceTurnPhase(&state)
	if got, want := advance.To, game.PhaseCityUpdate; got != want {
		t.Fatalf("to phase after reports = %q, want %q", got, want)
	}

	advance = game.AdvanceTurnPhase(&state)
	if got, want := advance.To, game.PhaseMessages; got != want {
		t.Fatalf("to phase after city update = %q, want %q", got, want)
	}
	if got, want := state.Turn, game.FirstTurn+1; got != want {
		t.Fatalf("turn = %d, want %d", got, want)
	}
	if got, want := state.Night, game.FirstNight; got != want {
		t.Fatalf("night = %d, want %d", got, want)
	}
	if advance.NightChanged {
		t.Fatal("night should not change before the configured turn limit")
	}
	if got, want := len(state.Messages), startMessages+1; got != want {
		t.Fatalf("message count = %d, want %d", got, want)
	}
	if got, want := state.Messages[len(state.Messages)-1].Subject, "turn brief"; got != want {
		t.Fatalf("last message subject = %q, want %q", got, want)
	}
}

func TestAdvanceTurnPhaseRecoversRunnersDuringCityUpdate(t *testing.T) {
	state := content.InitialGameState(42)
	state.Phase = game.PhaseCityUpdate
	state.Runners[0].State = game.RunnerReady
	state.Runners[0].Stress = 4
	state.Runners[1].State = game.RunnerInjured
	state.Runners[1].Stress = 5
	state.Runners[1].Recovery = 2
	state.Runners[2].State = game.RunnerInjured
	state.Runners[2].Stress = 3
	state.Runners[2].Recovery = 1
	startLog := len(state.EventLog)

	advance := game.AdvanceTurnPhase(&state)

	if got, want := state.Runners[0].Stress, 3; got != want {
		t.Fatalf("ready runner stress = %d, want %d", got, want)
	}
	if got, want := state.Runners[1].Recovery, 1; got != want {
		t.Fatalf("injured runner recovery = %d, want %d", got, want)
	}
	if got, want := state.Runners[1].State, game.RunnerInjured; got != want {
		t.Fatalf("injured runner state = %q, want %q", got, want)
	}
	if got, want := state.Runners[2].Recovery, 0; got != want {
		t.Fatalf("readying runner recovery = %d, want %d", got, want)
	}
	if got, want := state.Runners[2].State, game.RunnerReady; got != want {
		t.Fatalf("readying runner state = %q, want %q", got, want)
	}
	if got, want := advance.Recovery.StressRecovered, 1; got != want {
		t.Fatalf("stress recovered = %d, want %d", got, want)
	}
	if got, want := advance.Recovery.InjuryTicks, 2; got != want {
		t.Fatalf("injury ticks = %d, want %d", got, want)
	}
	if got, want := advance.Recovery.RunnersReadied, 1; got != want {
		t.Fatalf("runners readied = %d, want %d", got, want)
	}
	if got, want := len(state.EventLog), startLog+2; got != want {
		t.Fatalf("event log entries = %d, want %d", got, want)
	}
}

func TestAdvanceTurnPhaseRollsOverAfterFixedTurnsPerNight(t *testing.T) {
	state := content.InitialGameState(42)
	state.Phase = game.PhaseCityUpdate
	state.Turn = state.TurnsPerNight
	startMessages := len(state.Messages)

	advance := game.AdvanceTurnPhase(&state)

	if got, want := advance.To, game.PhaseMessages; got != want {
		t.Fatalf("to phase = %q, want %q", got, want)
	}
	if got, want := state.Night, game.FirstNight+1; got != want {
		t.Fatalf("night = %d, want %d", got, want)
	}
	if got, want := state.Turn, game.FirstTurn; got != want {
		t.Fatalf("turn = %d, want %d", got, want)
	}
	if !advance.NightChanged {
		t.Fatal("advance should report night rollover")
	}
	if got, want := len(state.Messages), startMessages+1; got != want {
		t.Fatalf("message count = %d, want %d", got, want)
	}
	if got, want := state.Messages[len(state.Messages)-1].Subject, "night brief"; got != want {
		t.Fatalf("last message subject = %q, want %q", got, want)
	}
}

func TestAdvanceTurnPhaseProgressesPastConfiguredRunNights(t *testing.T) {
	state := content.InitialGameState(42)
	state.Phase = game.PhaseCityUpdate
	state.Night = state.RunNights
	state.Turn = state.TurnsPerNight
	state.Credits = game.VictoryCreditTarget + game.NightlyOperatingCost
	startMessages := len(state.Messages)
	startLog := len(state.EventLog)

	advance := game.AdvanceTurnPhase(&state)

	if got, want := state.Night, game.DefaultRunNights+1; got != want {
		t.Fatalf("night = %d, want %d", got, want)
	}
	if got, want := state.Turn, game.FirstTurn; got != want {
		t.Fatalf("turn = %d, want %d", got, want)
	}
	if !advance.NightChanged {
		t.Fatal("advance should report final night rollover")
	}
	if got, want := advance.To, game.PhaseGameOver; got != want {
		t.Fatalf("to phase = %q, want %q", got, want)
	}
	if got, want := advance.Status.State, game.RunWon; got != want {
		t.Fatalf("run state = %q, want %q", got, want)
	}
	if got, want := advance.Status.Reason, game.RunEndVictory; got != want {
		t.Fatalf("run reason = %q, want %q", got, want)
	}
	if got, want := len(state.Messages), startMessages+1; got != want {
		t.Fatalf("message count = %d, want %d after run completion", got, want)
	}
	if got, want := state.Messages[len(state.Messages)-1].Subject, "run complete"; got != want {
		t.Fatalf("last message subject = %q, want %q", got, want)
	}
	if got, want := len(state.EventLog), startLog+4; got != want {
		t.Fatalf("event log entries = %d, want %d", got, want)
	}
}

func TestAdvanceTurnPhaseEndsFinalNightWithCreditShortfall(t *testing.T) {
	state := content.InitialGameState(42)
	state.Phase = game.PhaseCityUpdate
	state.Night = state.RunNights
	state.Turn = state.TurnsPerNight
	state.Credits = game.VictoryCreditTarget + game.NightlyOperatingCost - 1

	advance := game.AdvanceTurnPhase(&state)

	if got, want := advance.To, game.PhaseGameOver; got != want {
		t.Fatalf("to phase = %q, want %q", got, want)
	}
	if got, want := advance.Status.State, game.RunLost; got != want {
		t.Fatalf("run state = %q, want %q", got, want)
	}
	if got, want := advance.Status.Reason, game.RunEndShortfall; got != want {
		t.Fatalf("run reason = %q, want %q", got, want)
	}
	if got, want := state.Messages[len(state.Messages)-1].Subject, "run lost"; got != want {
		t.Fatalf("last message subject = %q, want %q", got, want)
	}
}

func TestAdvanceTurnPhaseStopsImmediatelyOnFailureCondition(t *testing.T) {
	state := content.InitialGameState(42)
	state.Phase = game.PhaseMessages
	state.Heat = game.MaximumHeat

	advance := game.AdvanceTurnPhase(&state)

	if got, want := advance.From, game.PhaseMessages; got != want {
		t.Fatalf("from phase = %q, want %q", got, want)
	}
	if got, want := advance.To, game.PhaseGameOver; got != want {
		t.Fatalf("to phase = %q, want %q", got, want)
	}
	if got, want := advance.Status.Reason, game.RunEndBurned; got != want {
		t.Fatalf("run reason = %q, want %q", got, want)
	}
	if advance.WaitingForInput {
		t.Fatal("game over should not wait for input")
	}
}

func TestAdvanceTurnPhaseDoesNotDuplicateRunEndReport(t *testing.T) {
	state := content.InitialGameState(42)
	state.Phase = game.PhaseMessages
	state.DispatchIntegrity = game.DispatchIntegrityFailure

	game.AdvanceTurnPhase(&state)
	messages := len(state.Messages)
	logs := len(state.EventLog)

	advance := game.AdvanceTurnPhase(&state)

	if got, want := advance.To, game.PhaseGameOver; got != want {
		t.Fatalf("to phase = %q, want %q", got, want)
	}
	if got, want := len(state.Messages), messages; got != want {
		t.Fatalf("message count = %d, want %d", got, want)
	}
	if got, want := len(state.EventLog), logs; got != want {
		t.Fatalf("event log entries = %d, want %d", got, want)
	}
}

func TestAdvanceTurnPhaseAppliesEndOfNightOperatingCostAndHeatDecay(t *testing.T) {
	state := content.InitialGameState(42)
	state.Phase = game.PhaseCityUpdate
	state.Turn = state.TurnsPerNight
	state.Credits = 800
	state.Heat = 5

	advance := game.AdvanceTurnPhase(&state)

	if !advance.NightChanged {
		t.Fatal("expected night to change")
	}

	// StartingCredits - NightlyOperatingCost = 800 - 150 = 650
	if got, want := state.Credits, 650; got != want {
		t.Fatalf("credits = %d, want %d", got, want)
	}

	// Heat - NightlyHeatDecay = 5 - 1 = 4
	if got, want := state.Heat, 4; got != want {
		t.Fatalf("heat = %d, want %d", got, want)
	}

	// Let's verify logs
	foundCostLog := false
	foundDecayLog := false
	for _, entry := range state.EventLog {
		if entry.Text == "Credits -150: nightly operating costs." {
			foundCostLog = true
		}
		if entry.Text == "Heat decayed by 1 (current: 4)." {
			foundDecayLog = true
		}
	}

	if !foundCostLog {
		t.Fatal("expected operating cost log entry")
	}
	if !foundDecayLog {
		t.Fatal("expected heat decay log entry")
	}
}
