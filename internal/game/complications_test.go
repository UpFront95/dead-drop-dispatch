package game_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"dead-drop-dispatch/internal/content"
	"dead-drop-dispatch/internal/game"
)

func TestFirstPlayableComplicationDefinitionsCoverMVPTypes(t *testing.T) {
	definitions := game.FirstPlayableComplicationDefinitions()
	if got, want := len(definitions), len(mvpComplicationTypes); got != want {
		t.Fatalf("definition count = %d, want %d", got, want)
	}

	gotTypes := make([]game.ComplicationType, 0, len(definitions))
	seenTypes := map[game.ComplicationType]bool{}
	for _, definition := range definitions {
		gotTypes = append(gotTypes, definition.Type)
		if definition.Type == "" {
			t.Fatal("definition type should not be empty")
		}
		if seenTypes[definition.Type] {
			t.Fatalf("duplicate definition type %s", definition.Type)
		}
		seenTypes[definition.Type] = true
		if definition.Title == "" {
			t.Fatalf("%s title should not be empty", definition.Type)
		}
		if definition.Prompt == "" {
			t.Fatalf("%s prompt should not be empty", definition.Type)
		}
		if definition.RiskTag == "" {
			t.Fatalf("%s risk tag should not be empty", definition.Type)
		}
		if definition.Description == "" {
			t.Fatalf("%s description should not be empty", definition.Type)
		}
		if len(definition.Choices) == 0 {
			t.Fatalf("%s choices should not be empty", definition.Type)
		}
		assertValidChoices(t, definition)
	}

	if !reflect.DeepEqual(gotTypes, mvpComplicationTypes) {
		t.Fatalf("types = %#v, want %#v", gotTypes, mvpComplicationTypes)
	}
}

func TestComplicationDefinitionForCoversEveryMVPType(t *testing.T) {
	for _, complicationType := range mvpComplicationTypes {
		t.Run(string(complicationType), func(t *testing.T) {
			definition, ok := game.ComplicationDefinitionFor(complicationType)
			if !ok {
				t.Fatalf("missing definition for %s", complicationType)
			}
			if definition.Type != complicationType {
				t.Fatalf("type = %q, want %q", definition.Type, complicationType)
			}
			assertValidChoices(t, definition)
		})
	}
}

func TestFirstPlayableComplicationDefinitionsIncludeFixedChoices(t *testing.T) {
	tests := []struct {
		complicationType game.ComplicationType
		wantChoices      []game.ComplicationChoiceID
	}{
		{
			complicationType: game.ComplicationCheckpoint,
			wantChoices: []game.ComplicationChoiceID{
				game.ChoiceTalkThrough,
				game.ChoiceReroute,
				game.ChoiceBribe,
				game.ChoiceAbandon,
			},
		},
		{
			complicationType: game.ComplicationScannerSweep,
			wantChoices: []game.ComplicationChoiceID{
				game.ChoiceHideCargo,
				game.ChoiceRushThrough,
				game.ChoiceSpoofTag,
			},
		},
		{
			complicationType: game.ComplicationRunnerPanic,
			wantChoices: []game.ComplicationChoiceID{
				game.ChoiceReassure,
				game.ChoiceOrderForward,
				game.ChoiceAbort,
			},
		},
		{
			complicationType: game.ComplicationCargoLeak,
			wantChoices: []game.ComplicationChoiceID{
				game.ChoiceContinue,
				game.ChoiceSeekClinic,
				game.ChoiceDumpCargo,
			},
		},
		{
			complicationType: game.ComplicationSignalLoss,
			wantChoices: []game.ComplicationChoiceID{
				game.ChoiceTrustRoute,
				game.ChoiceWait,
				game.ChoiceReroute,
			},
		},
		{
			complicationType: game.ComplicationGangToll,
			wantChoices: []game.ComplicationChoiceID{
				game.ChoicePayToll,
				game.ChoiceThreaten,
				game.ChoiceReroute,
				game.ChoiceCallFavor,
			},
		},
		{
			complicationType: game.ComplicationClientTerms,
			wantChoices: []game.ComplicationChoiceID{
				game.ChoiceAcceptTerms,
				game.ChoiceRenegotiate,
				game.ChoiceRefuseTerms,
				game.ChoiceAbandon,
			},
		},
		{
			complicationType: game.ComplicationDroneTail,
			wantChoices: []game.ComplicationChoiceID{
				game.ChoiceLoseTail,
				game.ChoiceJamSignal,
				game.ChoiceHideCargo,
				game.ChoiceShelter,
			},
		},
		{
			complicationType: game.ComplicationWitnessRefuses,
			wantChoices: []game.ComplicationChoiceID{
				game.ChoiceCoaxWitness,
				game.ChoiceReassure,
				game.ChoiceSedateWitness,
				game.ChoiceAbort,
			},
		},
		{
			complicationType: game.ComplicationDataTrace,
			wantChoices: []game.ComplicationChoiceID{
				game.ChoiceScrubTrace,
				game.ChoiceBurnNode,
				game.ChoiceDecoyPacket,
				game.ChoiceRushThrough,
			},
		},
		{
			complicationType: game.ComplicationCurfewDrop,
			wantChoices: []game.ComplicationChoiceID{
				game.ChoiceUseSafehouse,
				game.ChoiceBreakCurfew,
				game.ChoiceWait,
				game.ChoiceReroute,
			},
		},
		{
			complicationType: game.ComplicationRivalCourier,
			wantChoices: []game.ComplicationChoiceID{
				game.ChoiceRaceCourier,
				game.ChoiceBlockCourier,
				game.ChoiceShareRoute,
				game.ChoiceBribe,
			},
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.complicationType), func(t *testing.T) {
			definition, ok := game.ComplicationDefinitionFor(tt.complicationType)
			if !ok {
				t.Fatalf("missing definition for %s", tt.complicationType)
			}
			gotChoices := choiceIDs(definition.Choices)
			if !reflect.DeepEqual(gotChoices, tt.wantChoices) {
				t.Fatalf("choices = %#v, want %#v", gotChoices, tt.wantChoices)
			}
			for _, choice := range definition.Choices {
				if choice.Label == "" {
					t.Fatalf("choice %s label should not be empty", choice.ID)
				}
				if choice.Description == "" {
					t.Fatalf("choice %s description should not be empty", choice.ID)
				}
			}
		})
	}
}

func TestComplicationDefinitionForFindsKnownDefinition(t *testing.T) {
	definition, ok := game.ComplicationDefinitionFor(game.ComplicationScannerSweep)
	if !ok {
		t.Fatal("ComplicationDefinitionFor returned ok=false")
	}
	if definition.Title != "Scanner Sweep" {
		t.Fatalf("title = %q, want Scanner Sweep", definition.Title)
	}
}

func TestComplicationDefinitionForRejectsUnknownDefinition(t *testing.T) {
	_, ok := game.ComplicationDefinitionFor(game.ComplicationType("unknown"))
	if ok {
		t.Fatal("ComplicationDefinitionFor returned ok=true for unknown type")
	}
}

func TestQueueComplicationFromResultAppendsPendingComplication(t *testing.T) {
	state := content.InitialGameState(42)
	state.Complications = nil
	result := game.JobResult{
		JobID:              "job-1",
		JobTitle:           "Signal Job",
		RunnerID:           "runner-1",
		RunnerName:         "Runner One",
		FactionID:          "faction-1",
		Outcome:            game.OutcomePartial,
		CargoIntegrity:     65,
		CargoIntegrityLoss: 35,
		Delay:              true,
		Complication:       true,
		ComplicationType:   game.ComplicationSignalLoss,
		Factors:            []string{"signal noise", "delayed"},
	}

	complication, ok := game.QueueComplicationFromResult(&state, result)

	if !ok {
		t.Fatal("QueueComplicationFromResult returned ok=false")
	}
	if got, want := len(state.Complications), 1; got != want {
		t.Fatalf("complication count = %d, want %d", got, want)
	}
	if complication.ID == "" {
		t.Fatal("complication id should not be empty")
	}
	if complication.Status != game.ComplicationPending {
		t.Fatalf("status = %q, want %q", complication.Status, game.ComplicationPending)
	}
	if complication.Type != game.ComplicationSignalLoss {
		t.Fatalf("type = %q, want %q", complication.Type, game.ComplicationSignalLoss)
	}
	if complication.Title == "" {
		t.Fatal("complication title should not be empty")
	}
	if complication.Prompt == "" {
		t.Fatal("complication prompt should not be empty")
	}
	if got, want := choiceIDs(complication.Choices), []game.ComplicationChoiceID{
		game.ChoiceTrustRoute,
		game.ChoiceWait,
		game.ChoiceReroute,
	}; !reflect.DeepEqual(got, want) {
		t.Fatalf("choices = %#v, want %#v", got, want)
	}
	if complication.JobID != result.JobID || complication.RunnerID != result.RunnerID {
		t.Fatalf("complication context = %+v, want job %s runner %s", complication, result.JobID, result.RunnerID)
	}
	if got, want := complication.FactionID, result.FactionID; got != want {
		t.Fatalf("faction id = %q, want %q", got, want)
	}
	if got, want := complication.CargoIntegrity, 65; got != want {
		t.Fatalf("cargo integrity = %d, want %d", got, want)
	}
	if got, want := complication.CargoIntegrityLoss, 35; got != want {
		t.Fatalf("cargo integrity loss = %d, want %d", got, want)
	}
	if got, want := complication.DelayTurns, 1; got != want {
		t.Fatalf("delay turns = %d, want %d", got, want)
	}
	if len(complication.Factors) != len(result.Factors) {
		t.Fatalf("factors = %#v, want %#v", complication.Factors, result.Factors)
	}
}

func TestQueueComplicationFromResultSnapshotsChoices(t *testing.T) {
	state := content.InitialGameState(42)
	result := game.JobResult{
		JobID:            "job-1",
		JobTitle:         "Signal Job",
		RunnerID:         "runner-1",
		RunnerName:       "Runner One",
		Outcome:          game.OutcomePartial,
		Complication:     true,
		ComplicationType: game.ComplicationSignalLoss,
	}

	complication, ok := game.QueueComplicationFromResult(&state, result)
	if !ok {
		t.Fatal("QueueComplicationFromResult returned ok=false")
	}
	complication.Choices[0].Label = "mutated snapshot"

	definition, ok := game.ComplicationDefinitionFor(game.ComplicationSignalLoss)
	if !ok {
		t.Fatal("missing signal loss definition")
	}
	if got, want := definition.Choices[0].Label, "Trust route"; got != want {
		t.Fatalf("definition choice label = %q, want %q", got, want)
	}
}

func TestResolveComplicationChoiceAppliesCheckpointBribe(t *testing.T) {
	state := complicationResolutionState(game.ComplicationCheckpoint)
	state.Credits = game.ComplicationBribeCost
	state.Heat = 3

	resolution, err := game.ResolveComplicationChoice(&state, "cmp-1", game.ChoiceBribe)
	if err != nil {
		t.Fatalf("ResolveComplicationChoice returned error: %v", err)
	}

	if got, want := state.Complications[0].Status, game.ComplicationResolved; got != want {
		t.Fatalf("status = %q, want %q", got, want)
	}
	if got, want := state.Complications[0].ResolvedBy, game.ChoiceBribe; got != want {
		t.Fatalf("resolved by = %q, want %q", got, want)
	}
	if got, want := state.Credits, 0; got != want {
		t.Fatalf("credits = %d, want %d", got, want)
	}
	if got, want := state.Heat, 2; got != want {
		t.Fatalf("heat = %d, want %d", got, want)
	}
	if got, want := len(state.EventLog), 2; got != want {
		t.Fatalf("event log count = %d, want %d", got, want)
	}
	if got := state.EventLog[len(state.EventLog)-1].Text; !strings.Contains(got, "Effects: credits -75, heat -1.") {
		t.Fatalf("resolution log = %q, want effects", got)
	}
	if got, want := len(state.Messages), 1; got != want {
		t.Fatalf("message count = %d, want %d", got, want)
	}
	if got, want := state.Messages[0].From, "after-action"; got != want {
		t.Fatalf("message from = %q, want %q", got, want)
	}
	if got, want := state.Messages[0].Subject, "complication resolved"; got != want {
		t.Fatalf("message subject = %q, want %q", got, want)
	}
	if got := state.Messages[0].Body; !strings.Contains(got, "Resolved Checkpoint") || !strings.Contains(got, "Effects: credits -75, heat -1.") {
		t.Fatalf("message body = %q, want summary and effects", got)
	}
	if got, want := resolution.Effects, []string{"credits -75", "heat -1"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("effects = %#v, want %#v", got, want)
	}
	if got, want := resolution.CreditsDelta, -game.ComplicationBribeCost; got != want {
		t.Fatalf("credits delta = %d, want %d", got, want)
	}
	if got, want := resolution.HeatDelta, -game.ComplicationMinorHeatGain; got != want {
		t.Fatalf("heat delta = %d, want %d", got, want)
	}
	if got, want := state.Complications[0].ResolutionEffects, []string{"credits -75", "heat -1"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("stored effects = %#v, want %#v", got, want)
	}
}

func TestResolveComplicationChoiceAppliesRunnerPanicOrder(t *testing.T) {
	state := complicationResolutionState(game.ComplicationRunnerPanic)
	state.Runners[0].Stress = game.MaxRunnerStress - 1
	state.Runners[0].Loyalty = 1

	_, err := game.ResolveComplicationChoice(&state, "cmp-1", game.ChoiceOrderForward)
	if err != nil {
		t.Fatalf("ResolveComplicationChoice returned error: %v", err)
	}

	if got, want := state.Runners[0].Stress, game.MaxRunnerStress; got != want {
		t.Fatalf("stress = %d, want %d", got, want)
	}
	if got, want := state.Runners[0].Loyalty, game.MinRunnerLoyalty; got != want {
		t.Fatalf("loyalty = %d, want %d", got, want)
	}
}

func TestResolveComplicationChoiceAppliesCargoDump(t *testing.T) {
	state := complicationResolutionState(game.ComplicationCargoLeak)
	state.DispatchIntegrity = 4

	resolution, err := game.ResolveComplicationChoice(&state, "cmp-1", game.ChoiceDumpCargo)
	if err != nil {
		t.Fatalf("ResolveComplicationChoice returned error: %v", err)
	}

	if got, want := state.DispatchIntegrity, game.DispatchIntegrityFailure; got != want {
		t.Fatalf("dispatch integrity = %d, want %d", got, want)
	}
	if got, want := state.Factions[0].Suspicion, 2; got != want {
		t.Fatalf("faction suspicion = %d, want %d", got, want)
	}
	if got, want := state.Complications[0].CargoIntegrity, 0; got != want {
		t.Fatalf("cargo integrity = %d, want %d", got, want)
	}
	if got, want := state.Complications[0].CargoIntegrityLoss, game.DefaultCargoIntegrity; got != want {
		t.Fatalf("cargo integrity loss = %d, want %d", got, want)
	}
	if got, want := resolution.CargoIntegrityDelta, -game.DefaultCargoIntegrity; got != want {
		t.Fatalf("cargo integrity delta = %d, want %d", got, want)
	}
	if got, want := resolution.DispatchIntegrityDelta, -game.ComplicationMajorIntegrityLoss; got != want {
		t.Fatalf("dispatch integrity delta = %d, want %d", got, want)
	}
	if got, want := resolution.FactionSuspicionDelta, 1; got != want {
		t.Fatalf("faction suspicion delta = %d, want %d", got, want)
	}
}

func TestResolveComplicationChoiceAppliesDelayCost(t *testing.T) {
	state := complicationResolutionState(game.ComplicationSignalLoss)
	state.Complications[0].DelayTurns = 1

	resolution, err := game.ResolveComplicationChoice(&state, "cmp-1", game.ChoiceWait)
	if err != nil {
		t.Fatalf("ResolveComplicationChoice returned error: %v", err)
	}

	if got, want := resolution.DelayTurns, 1; got != want {
		t.Fatalf("resolution delay turns = %d, want %d", got, want)
	}
	if got, want := state.Complications[0].DelayTurns, 2; got != want {
		t.Fatalf("stored delay turns = %d, want %d", got, want)
	}
}

func TestResolveComplicationChoiceCanShiftFactionReputation(t *testing.T) {
	state := complicationResolutionState(game.ComplicationScannerSweep)

	resolution, err := game.ResolveComplicationChoice(&state, "cmp-1", game.ChoiceSpoofTag)
	if err != nil {
		t.Fatalf("ResolveComplicationChoice returned error: %v", err)
	}

	if got, want := state.Factions[0].Reputation, 1; got != want {
		t.Fatalf("faction reputation = %d, want %d", got, want)
	}
	if got, want := resolution.FactionReputationDelta, 1; got != want {
		t.Fatalf("faction reputation delta = %d, want %d", got, want)
	}
}

func TestResolveComplicationChoiceRejectsInvalidChoice(t *testing.T) {
	state := complicationResolutionState(game.ComplicationScannerSweep)

	_, err := game.ResolveComplicationChoice(&state, "cmp-1", game.ChoiceBribe)

	if !errors.Is(err, game.ErrComplicationChoiceNotFound) {
		t.Fatalf("error = %v, want %v", err, game.ErrComplicationChoiceNotFound)
	}
	if got, want := state.Complications[0].Status, game.ComplicationPending; got != want {
		t.Fatalf("status = %q, want %q", got, want)
	}
}

func TestResolveComplicationChoiceRejectsResolvedComplication(t *testing.T) {
	state := complicationResolutionState(game.ComplicationSignalLoss)
	_, err := game.ResolveComplicationChoice(&state, "cmp-1", game.ChoiceWait)
	if err != nil {
		t.Fatalf("ResolveComplicationChoice returned error: %v", err)
	}

	_, err = game.ResolveComplicationChoice(&state, "cmp-1", game.ChoiceWait)

	if !errors.Is(err, game.ErrComplicationAlreadyHandled) {
		t.Fatalf("error = %v, want %v", err, game.ErrComplicationAlreadyHandled)
	}
}

func TestResolveComplicationChoiceRejectsInsufficientBribeCredits(t *testing.T) {
	state := complicationResolutionState(game.ComplicationCheckpoint)
	state.Credits = game.ComplicationBribeCost - 1

	_, err := game.ResolveComplicationChoice(&state, "cmp-1", game.ChoiceBribe)

	if !errors.Is(err, game.ErrInsufficientCredits) {
		t.Fatalf("error = %v, want %v", err, game.ErrInsufficientCredits)
	}
	if got, want := state.Complications[0].Status, game.ComplicationPending; got != want {
		t.Fatalf("status = %q, want %q", got, want)
	}
}

func TestQueueComplicationFromResultIgnoresUnknownComplicationType(t *testing.T) {
	state := content.InitialGameState(42)

	_, ok := game.QueueComplicationFromResult(&state, game.JobResult{
		JobID:            "job-1",
		Complication:     true,
		ComplicationType: game.ComplicationType("unknown"),
	})

	if ok {
		t.Fatal("QueueComplicationFromResult returned ok=true")
	}
	if len(state.Complications) != 0 {
		t.Fatalf("complication count = %d, want 0", len(state.Complications))
	}
}

func TestQueueComplicationFromResultIgnoresOrdinaryResult(t *testing.T) {
	state := content.InitialGameState(42)
	startCount := len(state.Complications)

	_, ok := game.QueueComplicationFromResult(&state, game.JobResult{
		JobID:   "job-1",
		Outcome: game.OutcomeSuccess,
	})

	if ok {
		t.Fatal("QueueComplicationFromResult returned ok=true")
	}
	if len(state.Complications) != startCount {
		t.Fatalf("complication count = %d, want %d", len(state.Complications), startCount)
	}
}

var mvpComplicationTypes = []game.ComplicationType{
	game.ComplicationCheckpoint,
	game.ComplicationScannerSweep,
	game.ComplicationRunnerPanic,
	game.ComplicationCargoLeak,
	game.ComplicationSignalLoss,
	game.ComplicationGangToll,
	game.ComplicationClientTerms,
	game.ComplicationDroneTail,
	game.ComplicationWitnessRefuses,
	game.ComplicationDataTrace,
	game.ComplicationCurfewDrop,
	game.ComplicationRivalCourier,
}

func assertValidChoices(t *testing.T, definition game.ComplicationDefinition) {
	t.Helper()
	seen := map[game.ComplicationChoiceID]bool{}
	for _, choice := range definition.Choices {
		if choice.ID == "" {
			t.Fatalf("%s has empty choice id", definition.Type)
		}
		if seen[choice.ID] {
			t.Fatalf("%s has duplicate choice id %s", definition.Type, choice.ID)
		}
		seen[choice.ID] = true
		if choice.Label == "" {
			t.Fatalf("%s choice %s label should not be empty", definition.Type, choice.ID)
		}
		if choice.Description == "" {
			t.Fatalf("%s choice %s description should not be empty", definition.Type, choice.ID)
		}
	}
}

func choiceIDs(choices []game.ComplicationChoice) []game.ComplicationChoiceID {
	ids := make([]game.ComplicationChoiceID, 0, len(choices))
	for _, choice := range choices {
		ids = append(ids, choice.ID)
	}
	return ids
}

func complicationResolutionState(complicationType game.ComplicationType) game.GameState {
	definition, ok := game.ComplicationDefinitionFor(complicationType)
	if !ok {
		panic("unknown complication type")
	}
	return game.GameState{
		Turn:              4,
		Night:             1,
		Credits:           200,
		Heat:              0,
		DispatchIntegrity: game.StartingDispatchIntegrity,
		Runners: []game.Runner{{
			ID:      "runner-1",
			Name:    "Courier Nine",
			State:   game.RunnerReady,
			Stress:  2,
			Loyalty: 3,
		}},
		Factions: []game.Faction{{
			ID:        "faction-1",
			Name:      "Faction One",
			Suspicion: 1,
		}},
		Complications: []game.Complication{{
			ID:                 "cmp-1",
			Type:               complicationType,
			Status:             game.ComplicationPending,
			Title:              definition.Title,
			Prompt:             definition.Prompt,
			Choices:            definition.Choices,
			JobID:              "job-1",
			JobTitle:           "Bad Medicine",
			RunnerID:           "runner-1",
			RunnerName:         "Courier Nine",
			FactionID:          "faction-1",
			Outcome:            game.OutcomePartial,
			CargoIntegrity:     game.DefaultCargoIntegrity,
			CargoIntegrityLoss: 0,
			Summary:            "pending",
		}},
	}
}
