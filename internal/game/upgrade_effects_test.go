package game

import "testing"

func TestSignalRelayRevealsMoreRouteIntel(t *testing.T) {
	riskFactors := []string{"weak signal", "traffic pressure", "violent district"}

	withoutRelay := JobIntelReport(GameState{}, riskFactors)
	withRelay := JobIntelReport(GameState{PurchasedUpgrades: []UpgradeID{UpgradeSignalRelay}}, riskFactors)

	if got, want := len(withoutRelay.Claims), DefaultIntelClaimsPerJob; got != want {
		t.Fatalf("base claims = %d, want %d", got, want)
	}
	if got, want := len(withRelay.Claims), DefaultIntelClaimsPerJob+SignalRelayIntelClaimBonus; got != want {
		t.Fatalf("relay claims = %d, want %d", got, want)
	}
	if got, want := len(withRelay.OmittedTags), 0; got != want {
		t.Fatalf("relay omitted tags = %d, want %d", got, want)
	}
}

func TestSafehouseImprovesStressRecovery(t *testing.T) {
	base := recoveryState()
	upgraded := recoveryState()
	upgraded.PurchasedUpgrades = []UpgradeID{UpgradeSafehouse}

	RecoverRunners(&base)
	RecoverRunners(&upgraded)

	if got, want := base.Runners[0].Stress, 4; got != want {
		t.Fatalf("base stress = %d, want %d", got, want)
	}
	if got, want := upgraded.Runners[0].Stress, 3; got != want {
		t.Fatalf("safehouse stress = %d, want %d", got, want)
	}
}

func TestClinicFavorShortensInjuryRecovery(t *testing.T) {
	base := GameState{Runners: []Runner{{ID: "runner-1", State: RunnerOnJob}}}
	upgraded := GameState{
		PurchasedUpgrades: []UpgradeID{UpgradeClinicFavor},
		Runners:           []Runner{{ID: "runner-1", State: RunnerOnJob}},
	}
	result := JobResult{
		JobID:          "job-1",
		RunnerID:       "runner-1",
		Outcome:        OutcomePartial,
		StressGain:     0,
		CargoIntegrity: DefaultCargoIntegrity,
		Injury:         true,
		Summary:        "hurt.",
	}

	applyJobResult(&base, result)
	applyJobResult(&upgraded, result)

	if got, want := base.Runners[0].Recovery, DefaultInjuryRecoveryTurns; got != want {
		t.Fatalf("base recovery = %d, want %d", got, want)
	}
	if got, want := upgraded.Runners[0].Recovery, ClinicFavorInjuryRecoveryTurns; got != want {
		t.Fatalf("clinic recovery = %d, want %d", got, want)
	}
}

func TestFakeCredentialPrinterImprovesCheckpointOutcome(t *testing.T) {
	state := GameState{
		Credits:           100,
		DispatchIntegrity: StartingDispatchIntegrity,
		PurchasedUpgrades: []UpgradeID{UpgradeFakeCredentialPrinter},
		Runners:           []Runner{{ID: "runner-1", Stress: 2}},
		Complications: []Complication{{
			ID:       "cmp-1",
			Type:     ComplicationCheckpoint,
			Status:   ComplicationPending,
			Title:    "Checkpoint",
			JobTitle: "Hard Drop",
			RunnerID: "runner-1",
			Choices:  []ComplicationChoice{choice(ChoiceBribe, "Bribe", "")},
		}},
	}

	resolution, err := ResolveComplicationChoice(&state, "cmp-1", ChoiceBribe)
	if err != nil {
		t.Fatalf("ResolveComplicationChoice returned error: %v", err)
	}

	if got, want := resolution.CreditsDelta, -50; got != want {
		t.Fatalf("credits delta = %d, want %d", got, want)
	}
	if got, want := state.Credits, 50; got != want {
		t.Fatalf("credits = %d, want %d", got, want)
	}
	if !containsString(resolution.Effects, "fake credentials reduced bribe") {
		t.Fatalf("effects = %#v, want fake credential effect", resolution.Effects)
	}
}

func TestRiskUpgradesReduceContrabandAndDataTraceRisk(t *testing.T) {
	state := GameState{Districts: []District{{ID: "dest", Surveillance: 4, Danger: 4}}}
	runner := Runner{Stealth: 2, Talk: 1, Nerve: 2, Loyalty: 1}
	route := Route{Type: RouteDroneCorridor}
	contraband := Job{Cargo: CargoContrabandPackage, Destination: "dest", Faction: "faction-1"}
	data := Job{Cargo: CargoDataShard, Destination: "dest", Faction: "faction-1"}

	baseContrabandDetection := detectionRisk(state, contraband, route, runner, 1)
	state.PurchasedUpgrades = []UpgradeID{UpgradeDeadDropLocker}
	upgradedContrabandDetection := detectionRisk(state, contraband, route, runner, 1)
	if upgradedContrabandDetection >= baseContrabandDetection {
		t.Fatalf("contraband detection risk = %d, want below base %d", upgradedContrabandDetection, baseContrabandDetection)
	}

	baseDataComplication := complicationRisk(GameState{Districts: state.Districts}, data, route, runner, 1, true, false)
	state.PurchasedUpgrades = []UpgradeID{UpgradeScrambler}
	upgradedDataComplication := complicationRisk(state, data, route, runner, 1, true, false)
	if upgradedDataComplication >= baseDataComplication {
		t.Fatalf("data complication risk = %d, want below base %d", upgradedDataComplication, baseDataComplication)
	}
}

func recoveryState() GameState {
	return GameState{
		Runners: []Runner{
			{ID: "ready", State: RunnerReady, Stress: 5},
			{ID: "injured", State: RunnerInjured, Stress: 5, Recovery: 1},
		},
	}
}
