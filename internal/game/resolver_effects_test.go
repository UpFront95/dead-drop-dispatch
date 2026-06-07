package game

import (
	"math/rand"
	"strings"
	"testing"
)

func TestApplyJobResultAppliesOutcomeEffects(t *testing.T) {
	state := GameState{
		Turn:              3,
		Credits:           100,
		Heat:              MaximumHeat - 1,
		DispatchIntegrity: 4,
		Runners: []Runner{{
			ID:      "runner-1",
			Name:    "Runner One",
			State:   RunnerOnJob,
			Stress:  MaxRunnerStress - 1,
			Loyalty: MinRunnerLoyalty,
		}},
		Factions: []Faction{{
			ID:         "faction-1",
			Reputation: 2,
			Suspicion:  4,
		}},
		ActiveJobs: []ActiveJob{{
			JobID: "job-1",
			Job: Job{
				ID:      "job-1",
				Faction: "faction-1",
			},
		}},
	}
	result := JobResult{
		JobID:                 "job-1",
		JobTitle:              "Hard Drop",
		RunnerID:              "runner-1",
		Outcome:               OutcomeIntercepted,
		Payout:                50,
		HeatGain:              DetectionHeatGain + InterceptionHeatGain,
		StressGain:            InjuryStressGain + BundleStressGain,
		CargoIntegrity:        0,
		CargoIntegrityLoss:    CargoFailureIntegrityLoss,
		DispatchIntegrityLoss: FailedJobDispatchIntegrityLoss + CargoDamageDispatchIntegrityLoss,
		Detection:             true,
		Injury:                true,
		CargoDamage:           true,
		Interception:          true,
		Summary:               "Runner One returned from Hard Drop: intercepted.",
	}

	applyJobResult(&state, result)

	if got, want := state.Credits, 150; got != want {
		t.Fatalf("credits = %d, want %d", got, want)
	}
	if got, want := state.Heat, MaximumHeat; got != want {
		t.Fatalf("heat = %d, want %d", got, want)
	}
	if got, want := state.DispatchIntegrity, DispatchIntegrityFailure; got != want {
		t.Fatalf("dispatch integrity = %d, want %d", got, want)
	}
	if got, want := state.Runners[0].Stress, MaxRunnerStress; got != want {
		t.Fatalf("runner stress = %d, want %d", got, want)
	}
	if got, want := state.Runners[0].State, RunnerInjured; got != want {
		t.Fatalf("runner state = %q, want %q", got, want)
	}
	if got, want := state.Runners[0].Recovery, 2; got != want {
		t.Fatalf("runner recovery = %d, want %d", got, want)
	}
	if got, want := state.Runners[0].Loyalty, MinRunnerLoyalty; got != want {
		t.Fatalf("runner loyalty = %d, want %d", got, want)
	}
	if got, want := state.Factions[0].Reputation, 2; got != want {
		t.Fatalf("faction reputation = %d, want %d", got, want)
	}
	if got, want := state.Factions[0].Suspicion, 6; got != want {
		t.Fatalf("faction suspicion = %d, want %d", got, want)
	}
	if got, want := len(state.EventLog), 2; got != want {
		t.Fatalf("event log count = %d, want %d", got, want)
	}
}

func TestApplyJobResultRewardsSuccessfulFactionAndReadiesRunner(t *testing.T) {
	state := GameState{
		Turn:              1,
		Credits:           100,
		Heat:              0,
		DispatchIntegrity: StartingDispatchIntegrity,
		Runners: []Runner{{
			ID:      "runner-1",
			State:   RunnerOnJob,
			Stress:  0,
			Loyalty: 3,
		}},
		Factions: []Faction{{
			ID:         "faction-1",
			Reputation: 1,
			Suspicion:  1,
		}},
		ActiveJobs: []ActiveJob{{
			JobID: "job-1",
			Job: Job{
				ID:      "job-1",
				Faction: "faction-1",
			},
		}},
	}
	result := JobResult{
		JobID:          "job-1",
		JobTitle:       "Clean Drop",
		RunnerID:       "runner-1",
		Outcome:        OutcomeSuccess,
		Payout:         25,
		StressGain:     BaseJobStressGain,
		CargoIntegrity: DefaultCargoIntegrity,
		Summary:        "Clean.",
	}

	applyJobResult(&state, result)

	if got, want := state.Runners[0].State, RunnerReady; got != want {
		t.Fatalf("runner state = %q, want %q", got, want)
	}
	if got, want := state.Factions[0].Reputation, 2; got != want {
		t.Fatalf("faction reputation = %d, want %d", got, want)
	}
	if got, want := state.Factions[0].Suspicion, 1; got != want {
		t.Fatalf("faction suspicion = %d, want %d", got, want)
	}
	if got, want := state.DispatchIntegrity, StartingDispatchIntegrity; got != want {
		t.Fatalf("dispatch integrity = %d, want %d", got, want)
	}
}

func TestApplyJobResultQueuesComplication(t *testing.T) {
	state := GameState{
		Turn:              2,
		Night:             1,
		Credits:           100,
		Heat:              0,
		DispatchIntegrity: StartingDispatchIntegrity,
		Runners: []Runner{{
			ID:    "runner-1",
			State: RunnerOnJob,
		}},
		ActiveJobs: []ActiveJob{{
			JobID: "job-1",
			Job:   Job{ID: "job-1"},
		}},
	}
	result := JobResult{
		JobID:            "job-1",
		JobTitle:         "Signal Job",
		RunnerID:         "runner-1",
		RunnerName:       "Runner One",
		Outcome:          OutcomePartial,
		StressGain:       BaseJobStressGain,
		CargoIntegrity:   DefaultCargoIntegrity,
		Complication:     true,
		ComplicationType: ComplicationSignalLoss,
		Summary:          "Signal trouble.",
	}

	applyJobResult(&state, result)

	if got, want := len(state.Complications), 1; got != want {
		t.Fatalf("complication count = %d, want %d", got, want)
	}
	if state.Complications[0].Status != ComplicationPending {
		t.Fatalf("status = %q, want %q", state.Complications[0].Status, ComplicationPending)
	}
	if got, want := len(state.EventLog), 2; got != want {
		t.Fatalf("event log count = %d, want %d", got, want)
	}
	if got, want := len(state.Messages), 1; got != want {
		t.Fatalf("message count = %d, want %d", got, want)
	}
	if got, want := state.Messages[0].Subject, "complication opened"; got != want {
		t.Fatalf("message subject = %q, want %q", got, want)
	}
	if state.Messages[0].Body == "" {
		t.Fatal("complication message body should not be empty")
	}
}

func TestResultCargoIntegrity(t *testing.T) {
	tests := []struct {
		name        string
		outcome     JobOutcome
		cargoDamage bool
		want        int
	}{
		{name: "success", outcome: OutcomeSuccess, want: DefaultCargoIntegrity},
		{name: "partial damage", outcome: OutcomePartial, cargoDamage: true, want: DefaultCargoIntegrity - CargoDamageIntegrityLoss},
		{name: "failed", outcome: OutcomeFailed, want: 0},
		{name: "intercepted", outcome: OutcomeIntercepted, cargoDamage: true, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resultCargoIntegrity(tt.outcome, tt.cargoDamage)
			if got != tt.want {
				t.Fatalf("cargo integrity = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestResultSummaryExplainsDeliveryOutcome(t *testing.T) {
	tests := []struct {
		name   string
		result JobResult
		want   []string
	}{
		{
			name: "success",
			result: JobResult{
				JobTitle:   "Clean Drop",
				RunnerName: "Runner One",
				Outcome:    OutcomeSuccess,
			},
			want: []string{"completed Clean Drop cleanly", "cargo intact"},
		},
		{
			name: "partial",
			result: JobResult{
				JobTitle:    "Late Cooler",
				RunnerName:  "Runner One",
				Outcome:     OutcomePartial,
				Payout:      50,
				Delay:       true,
				CargoDamage: true,
				Factors:     []string{"weak signal", "medical spoilage", "traffic pressure"},
			},
			want: []string{"delivery came in rough", "route delay forced a late handoff", "cargo was damaged in transit", "client cut the payout", "Factors: weak signal, medical spoilage, traffic pressure."},
		},
		{
			name: "failed injury",
			result: JobResult{
				JobTitle:   "Bad Witness",
				RunnerName: "Runner One",
				Outcome:    OutcomeFailed,
				Injury:     true,
			},
			want: []string{"could not complete Bad Witness", "runner was hurt during the run", "No payout cleared"},
		},
		{
			name: "intercepted",
			result: JobResult{
				JobTitle:     "Data Trace",
				RunnerName:   "Runner One",
				Outcome:      OutcomeIntercepted,
				Detection:    true,
				Interception: true,
			},
			want: []string{"lost Data Trace to an intercept", "route drew heat", "Cargo is gone and heat climbs"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resultSummary(tt.result)
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Fatalf("summary = %q, want substring %q", got, want)
				}
			}
		})
	}
}

func TestResolveActiveJobCanProduceComplication(t *testing.T) {
	state := highComplicationState()
	active := ActiveJob{
		JobID:    "job-1",
		RunnerID: "runner-1",
		RouteID:  "route-1",
		Job: Job{
			ID:            "job-1",
			Title:         "Witness Trouble",
			Cargo:         CargoWitness,
			Destination:   "danger-zone",
			DeadlineTurns: 1,
			Payout:        100,
			Faction:       "faction-1",
			RiskFactors:   []string{"volatile passenger"},
		},
		Route: Route{
			ID:        "route-1",
			Type:      RouteDroneCorridor,
			Name:      "Bad air",
			Districts: []DistrictID{"safe-zone", "danger-zone"},
			Traits:    []string{"exposed signal"},
			TimeCost:  3,
		},
	}

	var result JobResult
	for seed := int64(1); seed <= 1000; seed++ {
		result = resolveActiveJob(&state, active, 2, rand.New(rand.NewSource(seed)))
		if result.Complication {
			break
		}
	}

	if !result.Complication {
		t.Fatal("expected high-risk job to produce a complication within fixed seed range")
	}
	if result.ComplicationType == ComplicationNone {
		t.Fatal("complication type should be set")
	}
	if !containsString(result.Factors, "complication: "+string(result.ComplicationType)) {
		t.Fatalf("factors = %#v, want complication factor %q", result.Factors, result.ComplicationType)
	}
	if result.Outcome == "" {
		t.Fatal("outcome should be explicit")
	}
}

func TestChooseComplicationType(t *testing.T) {
	runner := Runner{Stress: 0}
	job := Job{Cargo: CargoCorporatePrototype}

	tests := []struct {
		name        string
		job         Job
		route       Route
		runner      Runner
		detection   bool
		delay       bool
		cargoDamage bool
		want        ComplicationType
	}{
		{
			name:        "cargo leak",
			job:         Job{Cargo: CargoMedicalCooler},
			route:       Route{Type: RouteFloodline},
			cargoDamage: true,
			want:        ComplicationCargoLeak,
		},
		{
			name:      "data trace",
			job:       Job{Cargo: CargoDataShard},
			route:     Route{Type: RouteDroneCorridor},
			detection: true,
			want:      ComplicationDataTrace,
		},
		{
			name:      "scanner sweep",
			job:       job,
			route:     Route{Type: RouteDroneCorridor},
			detection: true,
			want:      ComplicationScannerSweep,
		},
		{
			name:   "runner panic",
			job:    Job{Cargo: CargoWitness},
			route:  Route{Type: RouteMarketWeave},
			runner: runner,
			want:   ComplicationRunnerPanic,
		},
		{
			name:   "signal loss",
			job:    job,
			route:  Route{Type: RouteServiceTunnels},
			runner: runner,
			want:   ComplicationSignalLoss,
		},
		{
			name:   "checkpoint",
			job:    job,
			route:  Route{Type: RouteMarketWeave},
			runner: runner,
			want:   ComplicationCheckpoint,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := chooseComplicationType(tt.job, tt.route, tt.runner, tt.detection, tt.delay, tt.cargoDamage)
			if got != tt.want {
				t.Fatalf("complication type = %q, want %q", got, tt.want)
			}
		})
	}
}

func highComplicationState() GameState {
	return GameState{
		Districts: []District{{
			ID:           "danger-zone",
			Surveillance: 5,
			Danger:       5,
		}},
		Runners: []Runner{{
			ID:      "runner-1",
			Name:    "Runner One",
			Nerve:   1,
			Stealth: 1,
			Talk:    1,
			Loyalty: 3,
			Stress:  6,
		}},
		Factions: []Faction{{
			ID:        "faction-1",
			Suspicion: 4,
		}},
	}
}
