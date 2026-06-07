package game_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"dead-drop-dispatch/internal/content"
	"dead-drop-dispatch/internal/game"
)

func TestFixedMessageResponseActionsCoverMSG01Actions(t *testing.T) {
	actions := game.FixedMessageResponseActions()
	if got, want := len(actions), len(msg01ResponseActions); got != want {
		t.Fatalf("action count = %d, want %d", got, want)
	}

	gotIDs := make([]game.MessageResponseActionID, 0, len(actions))
	seen := map[game.MessageResponseActionID]bool{}
	for _, action := range actions {
		gotIDs = append(gotIDs, action.ID)
		if action.ID == "" {
			t.Fatal("action id should not be empty")
		}
		if seen[action.ID] {
			t.Fatalf("duplicate action id %s", action.ID)
		}
		seen[action.ID] = true
		if action.Label == "" {
			t.Fatalf("%s label should not be empty", action.ID)
		}
		if action.Description == "" {
			t.Fatalf("%s description should not be empty", action.ID)
		}
		assertValidMessageAudiences(t, action)
	}

	if !reflect.DeepEqual(gotIDs, msg01ResponseActions) {
		t.Fatalf("action ids = %#v, want %#v", gotIDs, msg01ResponseActions)
	}
}

func TestMessageResponseActionForFindsKnownAction(t *testing.T) {
	action, ok := game.MessageResponseActionFor(game.ResponseAskMoreInfo)
	if !ok {
		t.Fatal("MessageResponseActionFor returned ok=false")
	}
	if action.Label != "Ask more info" {
		t.Fatalf("label = %q, want Ask more info", action.Label)
	}
}

func TestMessageResponseActionForRejectsUnknownAction(t *testing.T) {
	_, ok := game.MessageResponseActionFor(game.MessageResponseActionID("unknown"))
	if ok {
		t.Fatal("MessageResponseActionFor returned ok=true for unknown action")
	}
}

func TestMessageResponseActionsForAudiences(t *testing.T) {
	tests := []struct {
		audience game.MessageAudience
		want     []game.MessageResponseActionID
	}{
		{
			audience: game.MessageAudienceClient,
			want: []game.MessageResponseActionID{
				game.ResponseRefuse,
				game.ResponseAskMorePay,
				game.ResponseAskMoreInfo,
				game.ResponseThreaten,
				game.ResponseStall,
				game.ResponseDeceive,
				game.ResponseCancel,
				game.ResponseAccept,
			},
		},
		{
			audience: game.MessageAudienceRunner,
			want: []game.MessageResponseActionID{
				game.ResponseAskMoreInfo,
				game.ResponseThreaten,
				game.ResponseReassure,
				game.ResponseStall,
				game.ResponseDeceive,
				game.ResponseCancel,
				game.ResponseAccept,
			},
		},
		{
			audience: game.MessageAudienceFaction,
			want: []game.MessageResponseActionID{
				game.ResponseRefuse,
				game.ResponseAskMorePay,
				game.ResponseAskMoreInfo,
				game.ResponseThreaten,
				game.ResponseStall,
				game.ResponseDeceive,
				game.ResponseCancel,
				game.ResponseAccept,
			},
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.audience), func(t *testing.T) {
			actions := game.MessageResponseActionsFor(tt.audience)
			got := messageResponseActionIDs(actions)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("actions = %#v, want %#v", got, tt.want)
			}
			for _, actionID := range tt.want {
				if !game.MessageResponseActionAllowed(tt.audience, actionID) {
					t.Fatalf("%s should allow %s", tt.audience, actionID)
				}
			}
		})
	}
}

func TestMessageResponseActionAllowedRejectsUnsupportedAction(t *testing.T) {
	if game.MessageResponseActionAllowed(game.MessageAudienceRunner, game.ResponseAskMorePay) {
		t.Fatal("runner audience should not allow ask_more_pay")
	}
	if game.MessageResponseActionAllowed(game.MessageAudienceClient, game.ResponseReassure) {
		t.Fatal("client audience should not allow reassure")
	}
	if game.MessageResponseActionAllowed(game.MessageAudienceClient, game.MessageResponseActionID("unknown")) {
		t.Fatal("unknown action should not be allowed")
	}
}

func TestFixedMessageResponseActionsReturnsCopies(t *testing.T) {
	actions := game.FixedMessageResponseActions()
	actions[0].Label = "mutated action"
	actions[0].Audiences[0] = game.MessageAudienceRunner

	action, ok := game.MessageResponseActionFor(game.ResponseRefuse)
	if !ok {
		t.Fatal("missing refuse action")
	}
	if action.Label != "Refuse" {
		t.Fatalf("label = %q, want Refuse", action.Label)
	}
	if got, want := action.Audiences[0], game.MessageAudienceClient; got != want {
		t.Fatalf("audience = %q, want %q", got, want)
	}
}

func TestResolveMessageResponseAppliesClientAskMorePay(t *testing.T) {
	state := messageResponseState(game.Message{
		ID:        "msg-1",
		Turn:      game.FirstTurn,
		From:      "client",
		Subject:   "bad terms",
		Body:      "Need a discount.",
		Audience:  game.MessageAudienceClient,
		Status:    game.MessageOpen,
		Responses: game.MessageResponseActionsFor(game.MessageAudienceClient),
	})
	startCredits := state.Credits
	startLogs := len(state.EventLog)

	resolution, err := game.ResolveMessageResponse(&state, "msg-1", game.ResponseAskMorePay)
	if err != nil {
		t.Fatalf("ResolveMessageResponse returned error: %v", err)
	}

	if got, want := state.Messages[0].Status, game.MessageResolved; got != want {
		t.Fatalf("status = %q, want %q", got, want)
	}
	if got, want := state.Messages[0].ResolvedBy, game.ResponseAskMorePay; got != want {
		t.Fatalf("resolved by = %q, want %q", got, want)
	}
	if got, want := state.Credits, startCredits+game.MessageResponseNegotiationCreditGain; got != want {
		t.Fatalf("credits = %d, want %d", got, want)
	}
	if got, want := resolution.CreditsDelta, game.MessageResponseNegotiationCreditGain; got != want {
		t.Fatalf("credits delta = %d, want %d", got, want)
	}
	if got, want := resolution.Effects, []string{"credits +25"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("effects = %#v, want %#v", got, want)
	}
	if got, want := state.Messages[0].ResolutionEffects, []string{"credits +25"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("stored effects = %#v, want %#v", got, want)
	}
	if got, want := len(state.EventLog), startLogs+2; got != want {
		t.Fatalf("event log count = %d, want %d", got, want)
	}
	if got := state.EventLog[len(state.EventLog)-1].Text; !strings.Contains(got, "Resolved message bad terms with Ask more pay") || !strings.Contains(got, "Effects: credits +25.") {
		t.Fatalf("response log = %q, want summary and effects", got)
	}
	if got, want := state.Messages[len(state.Messages)-1].Subject, "response resolved"; got != want {
		t.Fatalf("report subject = %q, want %q", got, want)
	}
}

func TestResolveMessageResponseAppliesRunnerReassure(t *testing.T) {
	state := messageResponseState(game.Message{
		ID:             "msg-1",
		Turn:           game.FirstTurn,
		From:           "runner",
		Subject:        "nerves",
		Body:           "Runner is wavering.",
		Audience:       game.MessageAudienceRunner,
		Status:         game.MessageOpen,
		Responses:      game.MessageResponseActionsFor(game.MessageAudienceRunner),
		TargetRunnerID: "mira_vale",
	})
	state.Runners[0].Stress = 3
	state.Runners[0].Loyalty = 2

	resolution, err := game.ResolveMessageResponse(&state, "msg-1", game.ResponseReassure)
	if err != nil {
		t.Fatalf("ResolveMessageResponse returned error: %v", err)
	}

	if got, want := state.Runners[0].Stress, 2; got != want {
		t.Fatalf("stress = %d, want %d", got, want)
	}
	if got, want := state.Runners[0].Loyalty, 3; got != want {
		t.Fatalf("loyalty = %d, want %d", got, want)
	}
	if got, want := resolution.RunnerStressDelta, -game.MessageResponseRunnerStressRelief; got != want {
		t.Fatalf("stress delta = %d, want %d", got, want)
	}
	if got, want := resolution.RunnerLoyaltyDelta, game.MessageResponseRunnerLoyaltyGain; got != want {
		t.Fatalf("loyalty delta = %d, want %d", got, want)
	}
}

func TestResolveMessageResponseAppliesFactionThreaten(t *testing.T) {
	state := messageResponseState(game.Message{
		ID:              "msg-1",
		Turn:            game.FirstTurn,
		From:            "faction",
		Subject:         "pressure",
		Body:            "Faction wants a concession.",
		Audience:        game.MessageAudienceFaction,
		Status:          game.MessageOpen,
		Responses:       game.MessageResponseActionsFor(game.MessageAudienceFaction),
		TargetFactionID: "helix_municipal_security",
	})
	startSuspicion := state.Factions[0].Suspicion

	resolution, err := game.ResolveMessageResponse(&state, "msg-1", game.ResponseThreaten)
	if err != nil {
		t.Fatalf("ResolveMessageResponse returned error: %v", err)
	}

	if got, want := state.Heat, game.MessageResponseMinorHeatGain; got != want {
		t.Fatalf("heat = %d, want %d", got, want)
	}
	if got, want := state.Factions[0].Suspicion, startSuspicion+game.MessageResponseFactionShift; got != want {
		t.Fatalf("suspicion = %d, want %d", got, want)
	}
	if got, want := resolution.HeatDelta, game.MessageResponseMinorHeatGain; got != want {
		t.Fatalf("heat delta = %d, want %d", got, want)
	}
	if got, want := resolution.FactionSuspicionDelta, game.MessageResponseFactionShift; got != want {
		t.Fatalf("suspicion delta = %d, want %d", got, want)
	}
}

func TestResolveMessageResponseUsesAudienceDefaultsWhenMessageHasNoSnapshot(t *testing.T) {
	state := messageResponseState(game.Message{
		ID:       "msg-1",
		Turn:     game.FirstTurn,
		From:     "client",
		Subject:  "wait",
		Body:     "Hold a minute.",
		Audience: game.MessageAudienceClient,
	})

	_, err := game.ResolveMessageResponse(&state, "msg-1", game.ResponseStall)
	if err != nil {
		t.Fatalf("ResolveMessageResponse returned error: %v", err)
	}

	if got, want := state.Messages[0].Status, game.MessageResolved; got != want {
		t.Fatalf("status = %q, want %q", got, want)
	}
}

func TestResolveMessageResponseRejectsUnsupportedAction(t *testing.T) {
	state := messageResponseState(game.Message{
		ID:        "msg-1",
		Turn:      game.FirstTurn,
		From:      "runner",
		Subject:   "pay",
		Body:      "Runner asks about terms.",
		Audience:  game.MessageAudienceRunner,
		Status:    game.MessageOpen,
		Responses: game.MessageResponseActionsFor(game.MessageAudienceRunner),
	})

	_, err := game.ResolveMessageResponse(&state, "msg-1", game.ResponseAskMorePay)

	if !errors.Is(err, game.ErrMessageResponseNotFound) {
		t.Fatalf("error = %v, want %v", err, game.ErrMessageResponseNotFound)
	}
	if got, want := state.Messages[0].Status, game.MessageOpen; got != want {
		t.Fatalf("status = %q, want %q", got, want)
	}
}

func TestResolveMessageResponseRejectsResolvedMessage(t *testing.T) {
	state := messageResponseState(game.Message{
		ID:        "msg-1",
		Turn:      game.FirstTurn,
		From:      "client",
		Subject:   "done",
		Body:      "Already handled.",
		Audience:  game.MessageAudienceClient,
		Status:    game.MessageResolved,
		Responses: game.MessageResponseActionsFor(game.MessageAudienceClient),
	})

	_, err := game.ResolveMessageResponse(&state, "msg-1", game.ResponseAccept)

	if !errors.Is(err, game.ErrMessageAlreadyHandled) {
		t.Fatalf("error = %v, want %v", err, game.ErrMessageAlreadyHandled)
	}
}

func TestResolveMessageResponseRejectsMissingAudience(t *testing.T) {
	state := messageResponseState(game.Message{
		ID:      "msg-1",
		Turn:    game.FirstTurn,
		From:    "switchboard",
		Subject: "line check",
		Body:    "No action context.",
	})

	_, err := game.ResolveMessageResponse(&state, "msg-1", game.ResponseAccept)

	if !errors.Is(err, game.ErrMessageAudienceNotSet) {
		t.Fatalf("error = %v, want %v", err, game.ErrMessageAudienceNotSet)
	}
}

func TestResolveMessageResponseRejectsUnknownMessage(t *testing.T) {
	state := messageResponseState(game.Message{
		ID:       "msg-1",
		Turn:     game.FirstTurn,
		From:     "client",
		Subject:  "terms",
		Body:     "Handle this.",
		Audience: game.MessageAudienceClient,
	})

	_, err := game.ResolveMessageResponse(&state, "missing", game.ResponseAccept)

	if !errors.Is(err, game.ErrMessageNotFound) {
		t.Fatalf("error = %v, want %v", err, game.ErrMessageNotFound)
	}
}

func TestResolveMessageResponseRejectsActionMissingFromSnapshot(t *testing.T) {
	state := messageResponseState(game.Message{
		ID:        "msg-1",
		Turn:      game.FirstTurn,
		From:      "client",
		Subject:   "narrow terms",
		Body:      "Only one reply is on the table.",
		Audience:  game.MessageAudienceClient,
		Status:    game.MessageOpen,
		Responses: []game.MessageResponseAction{mustMessageResponseAction(t, game.ResponseAccept)},
	})
	startHeat := state.Heat
	startMessages := len(state.Messages)
	startLogs := len(state.EventLog)

	_, err := game.ResolveMessageResponse(&state, "msg-1", game.ResponseStall)

	if !errors.Is(err, game.ErrMessageResponseNotFound) {
		t.Fatalf("error = %v, want %v", err, game.ErrMessageResponseNotFound)
	}
	if got, want := state.Messages[0].Status, game.MessageOpen; got != want {
		t.Fatalf("status = %q, want %q", got, want)
	}
	if state.Heat != startHeat {
		t.Fatalf("heat = %d, want unchanged %d", state.Heat, startHeat)
	}
	if len(state.Messages) != startMessages {
		t.Fatalf("message count = %d, want %d", len(state.Messages), startMessages)
	}
	if len(state.EventLog) != startLogs {
		t.Fatalf("event log count = %d, want %d", len(state.EventLog), startLogs)
	}
}

func TestResolveMessageResponseDeterministicOutcomes(t *testing.T) {
	tests := []struct {
		name                      string
		actionID                  game.MessageResponseActionID
		audience                  game.MessageAudience
		targetRunnerID            game.RunnerID
		targetFactionID           game.FactionID
		wantCreditsDelta          int
		wantHeatDelta             int
		wantDispatchIntegrityDiff int
		wantRunnerStressDelta     int
		wantRunnerLoyaltyDelta    int
		wantFactionReputationDiff int
		wantFactionSuspicionDiff  int
		wantEffects               []string
		wantLogDelta              int
	}{
		{
			name:                      "refuse faction",
			actionID:                  game.ResponseRefuse,
			audience:                  game.MessageAudienceFaction,
			targetFactionID:           "helix_municipal_security",
			wantDispatchIntegrityDiff: -game.MessageResponseMinorIntegrityLoss,
			wantFactionSuspicionDiff:  game.MessageResponseFactionShift,
			wantEffects:               []string{"dispatch integrity -1", "faction suspicion +1"},
			wantLogDelta:              1,
		},
		{
			name:                      "ask more pay",
			actionID:                  game.ResponseAskMorePay,
			audience:                  game.MessageAudienceClient,
			targetFactionID:           "helix_municipal_security",
			wantCreditsDelta:          game.MessageResponseNegotiationCreditGain,
			wantFactionReputationDiff: -game.MessageResponseFactionShift,
			wantEffects:               []string{"credits +25", "faction reputation -1"},
			wantLogDelta:              2,
		},
		{
			name:         "ask more info",
			actionID:     game.ResponseAskMoreInfo,
			audience:     game.MessageAudienceClient,
			wantEffects:  []string{"intel requested"},
			wantLogDelta: 1,
		},
		{
			name:                     "threaten runner and faction",
			actionID:                 game.ResponseThreaten,
			audience:                 game.MessageAudienceRunner,
			targetRunnerID:           "mira_vale",
			targetFactionID:          "helix_municipal_security",
			wantHeatDelta:            game.MessageResponseMinorHeatGain,
			wantRunnerStressDelta:    game.MessageResponseRunnerStressGain,
			wantRunnerLoyaltyDelta:   -game.MessageResponseRunnerLoyaltyLoss,
			wantFactionSuspicionDiff: game.MessageResponseFactionShift,
			wantEffects: []string{
				"heat +1",
				"runner stress +1",
				"runner loyalty -1",
				"faction suspicion +1",
			},
			wantLogDelta: 1,
		},
		{
			name:                   "reassure runner",
			actionID:               game.ResponseReassure,
			audience:               game.MessageAudienceRunner,
			targetRunnerID:         "mira_vale",
			wantRunnerStressDelta:  -game.MessageResponseRunnerStressRelief,
			wantRunnerLoyaltyDelta: game.MessageResponseRunnerLoyaltyGain,
			wantEffects:            []string{"runner stress -1", "runner loyalty +1"},
			wantLogDelta:           1,
		},
		{
			name:          "stall",
			actionID:      game.ResponseStall,
			audience:      game.MessageAudienceClient,
			wantHeatDelta: game.MessageResponseMinorHeatGain,
			wantEffects:   []string{"heat +1"},
			wantLogDelta:  1,
		},
		{
			name:                      "deceive faction",
			actionID:                  game.ResponseDeceive,
			audience:                  game.MessageAudienceFaction,
			targetFactionID:           "helix_municipal_security",
			wantHeatDelta:             game.MessageResponseMinorHeatGain,
			wantDispatchIntegrityDiff: -game.MessageResponseMinorIntegrityLoss,
			wantFactionSuspicionDiff:  game.MessageResponseFactionShift,
			wantEffects:               []string{"heat +1", "dispatch integrity -1", "faction suspicion +1"},
			wantLogDelta:              1,
		},
		{
			name:                      "cancel faction",
			actionID:                  game.ResponseCancel,
			audience:                  game.MessageAudienceFaction,
			targetFactionID:           "helix_municipal_security",
			wantDispatchIntegrityDiff: -game.MessageResponseMajorIntegrityLoss,
			wantFactionSuspicionDiff:  game.MessageResponseFactionShift,
			wantEffects:               []string{"dispatch integrity -2", "faction suspicion +1"},
			wantLogDelta:              1,
		},
		{
			name:                      "accept faction",
			actionID:                  game.ResponseAccept,
			audience:                  game.MessageAudienceFaction,
			targetFactionID:           "helix_municipal_security",
			wantFactionReputationDiff: game.MessageResponseFactionShift,
			wantEffects:               []string{"faction reputation +1"},
			wantLogDelta:              1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := messageResponseState(game.Message{
				ID:              "msg-1",
				Turn:            game.FirstTurn,
				From:            string(tt.audience),
				Subject:         tt.name,
				Body:            "Resolve this.",
				Audience:        tt.audience,
				Status:          game.MessageOpen,
				Responses:       game.MessageResponseActionsFor(tt.audience),
				TargetRunnerID:  tt.targetRunnerID,
				TargetFactionID: tt.targetFactionID,
			})
			state.Runners[0].Stress = 4
			state.Runners[0].Loyalty = 3
			state.Factions[0].Reputation = 2
			state.Factions[0].Suspicion = 1
			startCredits := state.Credits
			startHeat := state.Heat
			startDispatchIntegrity := state.DispatchIntegrity
			startRunnerStress := state.Runners[0].Stress
			startRunnerLoyalty := state.Runners[0].Loyalty
			startFactionReputation := state.Factions[0].Reputation
			startFactionSuspicion := state.Factions[0].Suspicion
			startLogs := len(state.EventLog)

			resolution, err := game.ResolveMessageResponse(&state, "msg-1", tt.actionID)
			if err != nil {
				t.Fatalf("ResolveMessageResponse returned error: %v", err)
			}

			if got, want := state.Messages[0].Status, game.MessageResolved; got != want {
				t.Fatalf("status = %q, want %q", got, want)
			}
			if got, want := state.Messages[0].ResolvedBy, tt.actionID; got != want {
				t.Fatalf("resolved by = %q, want %q", got, want)
			}
			if got, want := resolution.ActionID, tt.actionID; got != want {
				t.Fatalf("resolution action = %q, want %q", got, want)
			}
			if got, want := resolution.Effects, tt.wantEffects; !reflect.DeepEqual(got, want) {
				t.Fatalf("effects = %#v, want %#v", got, want)
			}
			if got, want := state.Messages[0].ResolutionEffects, tt.wantEffects; !reflect.DeepEqual(got, want) {
				t.Fatalf("stored effects = %#v, want %#v", got, want)
			}
			if got, want := state.Credits, startCredits+tt.wantCreditsDelta; got != want {
				t.Fatalf("credits = %d, want %d", got, want)
			}
			if got, want := state.Heat, startHeat+tt.wantHeatDelta; got != want {
				t.Fatalf("heat = %d, want %d", got, want)
			}
			if got, want := state.DispatchIntegrity, startDispatchIntegrity+tt.wantDispatchIntegrityDiff; got != want {
				t.Fatalf("dispatch integrity = %d, want %d", got, want)
			}
			if got, want := state.Runners[0].Stress, startRunnerStress+tt.wantRunnerStressDelta; got != want {
				t.Fatalf("runner stress = %d, want %d", got, want)
			}
			if got, want := state.Runners[0].Loyalty, startRunnerLoyalty+tt.wantRunnerLoyaltyDelta; got != want {
				t.Fatalf("runner loyalty = %d, want %d", got, want)
			}
			if got, want := state.Factions[0].Reputation, startFactionReputation+tt.wantFactionReputationDiff; got != want {
				t.Fatalf("faction reputation = %d, want %d", got, want)
			}
			if got, want := state.Factions[0].Suspicion, startFactionSuspicion+tt.wantFactionSuspicionDiff; got != want {
				t.Fatalf("faction suspicion = %d, want %d", got, want)
			}
			if got, want := len(state.EventLog), startLogs+tt.wantLogDelta; got != want {
				t.Fatalf("event log count = %d, want %d", got, want)
			}
			if got, want := state.Messages[len(state.Messages)-1].Subject, "response resolved"; got != want {
				t.Fatalf("report subject = %q, want %q", got, want)
			}
		})
	}
}

var msg01ResponseActions = []game.MessageResponseActionID{
	game.ResponseRefuse,
	game.ResponseAskMorePay,
	game.ResponseAskMoreInfo,
	game.ResponseThreaten,
	game.ResponseReassure,
	game.ResponseStall,
	game.ResponseDeceive,
	game.ResponseCancel,
	game.ResponseAccept,
}

func assertValidMessageAudiences(t *testing.T, action game.MessageResponseAction) {
	t.Helper()
	if len(action.Audiences) == 0 {
		t.Fatalf("%s audiences should not be empty", action.ID)
	}
	seen := map[game.MessageAudience]bool{}
	for _, audience := range action.Audiences {
		if audience == "" {
			t.Fatalf("%s has empty audience", action.ID)
		}
		switch audience {
		case game.MessageAudienceClient, game.MessageAudienceRunner, game.MessageAudienceFaction:
		default:
			t.Fatalf("%s has unknown audience %s", action.ID, audience)
		}
		if seen[audience] {
			t.Fatalf("%s has duplicate audience %s", action.ID, audience)
		}
		seen[audience] = true
	}
}

func messageResponseActionIDs(actions []game.MessageResponseAction) []game.MessageResponseActionID {
	ids := make([]game.MessageResponseActionID, 0, len(actions))
	for _, action := range actions {
		ids = append(ids, action.ID)
	}
	return ids
}

func messageResponseState(message game.Message) game.GameState {
	state := content.InitialGameState(42)
	state.Messages = []game.Message{message}
	state.EventLog = nil
	return state
}

func mustMessageResponseAction(t *testing.T, actionID game.MessageResponseActionID) game.MessageResponseAction {
	t.Helper()
	action, ok := game.MessageResponseActionFor(actionID)
	if !ok {
		t.Fatalf("missing response action %s", actionID)
	}
	return action
}
