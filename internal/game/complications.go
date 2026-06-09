package game

import (
	"errors"
	"fmt"
	"strings"
)

const (
	ComplicationBribeCost             = 75
	ComplicationClinicCost            = 60
	ComplicationMinorStressGain       = 1
	ComplicationMajorStressGain       = 2
	ComplicationStressRelief          = 2
	ComplicationMinorHeatGain         = 1
	ComplicationMajorHeatGain         = 2
	ComplicationMinorIntegrityLoss    = 1
	ComplicationModerateIntegrityLoss = 3
	ComplicationMajorIntegrityLoss    = 5
	ComplicationMinorCargoLoss        = 10
	ComplicationModerateCargoLoss     = 20
	ComplicationMajorCargoLoss        = 35
	ComplicationMinorDelayTurns       = 1
	ComplicationMajorDelayTurns       = 2
)

var (
	ErrComplicationNotFound       = errors.New("complication not found")
	ErrComplicationAlreadyHandled = errors.New("complication already handled")
	ErrComplicationChoiceNotFound = errors.New("complication choice not found")
)

func FirstPlayableComplicationDefinitions() []ComplicationDefinition {
	return []ComplicationDefinition{
		{
			Type:        ComplicationCheckpoint,
			Title:       "Checkpoint",
			Prompt:      "A checkpoint opens across the route and starts pulling runners out of flow.",
			RiskTag:     "authority pressure",
			Description: "Talk, bribe, reroute, or abandon choices can resolve checkpoint pressure.",
			Choices: []ComplicationChoice{
				choice(ChoiceTalkThrough, "Talk through", "Lean on runner talk and cover story."),
				choice(ChoiceReroute, "Reroute", "Leave the checkpoint and take a messier path."),
				choice(ChoiceBribe, "Bribe", "Spend credits to keep the line moving."),
				choice(ChoiceAbandon, "Abandon", "Cut the job before the checkpoint closes."),
			},
		},
		{
			Type:        ComplicationScannerSweep,
			Title:       "Scanner Sweep",
			Prompt:      "A scanner sweep rolls through the corridor and starts matching cargo tags.",
			RiskTag:     "scan exposure",
			Description: "Hide cargo, rush through, or spoof tag choices can resolve scanner exposure.",
			Choices: []ComplicationChoice{
				choice(ChoiceHideCargo, "Hide cargo", "Slow down and reduce scan exposure."),
				choice(ChoiceRushThrough, "Rush through", "Push speed before the sweep tightens."),
				choice(ChoiceSpoofTag, "Spoof tag", "Try to poison the scanner match."),
			},
		},
		{
			Type:        ComplicationRunnerPanic,
			Title:       "Runner Panic",
			Prompt:      "The runner's voice breaks over comms and their route discipline starts to slip.",
			RiskTag:     "stress spike",
			Description: "Reassure, order forward, or abort choices can resolve runner panic.",
			Choices: []ComplicationChoice{
				choice(ChoiceReassure, "Reassure", "Spend time calming the runner down."),
				choice(ChoiceOrderForward, "Order forward", "Force the runner to keep moving."),
				choice(ChoiceAbort, "Abort", "Pull the runner out before they break."),
			},
		},
		{
			Type:        ComplicationCargoLeak,
			Title:       "Cargo Leak",
			Prompt:      "The package reports a breach and the cargo condition starts dropping fast.",
			RiskTag:     "cargo integrity",
			Description: "Continue, seek help, or dump cargo choices can resolve a cargo leak.",
			Choices: []ComplicationChoice{
				choice(ChoiceContinue, "Continue", "Keep moving and accept cargo risk."),
				choice(ChoiceSeekClinic, "Seek clinic", "Divert toward help and lose time."),
				choice(ChoiceDumpCargo, "Dump cargo", "Protect the runner by losing the package."),
			},
		},
		{
			Type:        ComplicationSignalLoss,
			Title:       "Signal Loss",
			Prompt:      "The route feed breaks into static and the runner stops receiving clean directions.",
			RiskTag:     "lost contact",
			Description: "Trust route, wait, or reroute choices can resolve signal loss.",
			Choices: []ComplicationChoice{
				choice(ChoiceTrustRoute, "Trust route", "Let the runner follow the last known path."),
				choice(ChoiceWait, "Wait", "Hold position until the signal clears."),
				choice(ChoiceReroute, "Reroute", "Switch to a route with cleaner contact."),
			},
		},
		{
			Type:        ComplicationGangToll,
			Title:       "Gang Toll",
			Prompt:      "A local crew blocks the handoff lane and demands payment before the runner passes.",
			RiskTag:     "street pressure",
			Description: "Pay toll, threaten, reroute, or call favor choices can resolve gang pressure.",
			Choices: []ComplicationChoice{
				choice(ChoicePayToll, "Pay toll", "Spend credits to keep the route calm."),
				choice(ChoiceThreaten, "Threaten", "Push back and risk escalation."),
				choice(ChoiceReroute, "Reroute", "Avoid the crew through a slower side path."),
				choice(ChoiceCallFavor, "Call favor", "Lean on faction reputation to clear the lane."),
			},
		},
		{
			Type:        ComplicationClientTerms,
			Title:       "Client Changes Terms",
			Prompt:      "The client updates the drop terms mid-run and tries to move the goalposts.",
			RiskTag:     "contract pressure",
			Description: "Accept terms, renegotiate, refuse terms, or abandon choices can resolve client pressure.",
			Choices: []ComplicationChoice{
				choice(ChoiceAcceptTerms, "Accept terms", "Take the change and preserve the relationship."),
				choice(ChoiceRenegotiate, "Renegotiate", "Push for a cleaner deal before delivery."),
				choice(ChoiceRefuseTerms, "Refuse terms", "Hold the original contract line."),
				choice(ChoiceAbandon, "Abandon", "Cut the job before the client extracts more."),
			},
		},
		{
			Type:        ComplicationDroneTail,
			Title:       "Drone Tail",
			Prompt:      "A surveillance drone starts matching the runner's turns and altitude shifts.",
			RiskTag:     "aerial surveillance",
			Description: "Lose tail, jam signal, hide cargo, or shelter choices can resolve drone pressure.",
			Choices: []ComplicationChoice{
				choice(ChoiceLoseTail, "Lose tail", "Push movement and break visual contact."),
				choice(ChoiceJamSignal, "Jam signal", "Interfere with the drone link."),
				choice(ChoiceHideCargo, "Hide cargo", "Protect the package while the drone passes."),
				choice(ChoiceShelter, "Shelter", "Duck under cover and wait it out."),
			},
		},
		{
			Type:        ComplicationWitnessRefuses,
			Title:       "Witness Refuses",
			Prompt:      "The witness cargo stops cooperating and refuses to enter the next stretch.",
			RiskTag:     "human cargo",
			Description: "Coax witness, reassure, sedate witness, or abort choices can resolve witness refusal.",
			Choices: []ComplicationChoice{
				choice(ChoiceCoaxWitness, "Coax witness", "Talk them into moving without force."),
				choice(ChoiceReassure, "Reassure", "Slow down and lower panic."),
				choice(ChoiceSedateWitness, "Sedate witness", "Force the issue and accept medical risk."),
				choice(ChoiceAbort, "Abort", "Extract the runner and lose the job."),
			},
		},
		{
			Type:        ComplicationDataTrace,
			Title:       "Data Trace",
			Prompt:      "The data shard starts phoning home and leaving a live trace through the route.",
			RiskTag:     "trace exposure",
			Description: "Scrub trace, burn node, decoy packet, or rush through choices can resolve data trace exposure.",
			Choices: []ComplicationChoice{
				choice(ChoiceScrubTrace, "Scrub trace", "Spend time cleaning the signature."),
				choice(ChoiceBurnNode, "Burn node", "Sacrifice infrastructure to cut the trace."),
				choice(ChoiceDecoyPacket, "Decoy packet", "Feed the trace a false path."),
				choice(ChoiceRushThrough, "Rush through", "Finish before the trace completes."),
			},
		},
		{
			Type:        ComplicationCurfewDrop,
			Title:       "Curfew Drop",
			Prompt:      "A sudden curfew locks down the destination district before the handoff is complete.",
			RiskTag:     "district lockdown",
			Description: "Use safehouse, break curfew, wait, or reroute choices can resolve curfew pressure.",
			Choices: []ComplicationChoice{
				choice(ChoiceUseSafehouse, "Use safehouse", "Hold cargo in a protected stopover."),
				choice(ChoiceBreakCurfew, "Break curfew", "Push through the lockdown."),
				choice(ChoiceWait, "Wait", "Hold position until patrol rhythm opens."),
				choice(ChoiceReroute, "Reroute", "Move through a less watched district edge."),
			},
		},
		{
			Type:        ComplicationRivalCourier,
			Title:       "Rival Courier",
			Prompt:      "A rival courier appears on the same contract vector and tries to beat the runner to the drop.",
			RiskTag:     "competitive pressure",
			Description: "Race courier, block courier, share route, or bribe choices can resolve rival pressure.",
			Choices: []ComplicationChoice{
				choice(ChoiceRaceCourier, "Race courier", "Prioritize speed and beat them to the drop."),
				choice(ChoiceBlockCourier, "Block courier", "Interfere with their route."),
				choice(ChoiceShareRoute, "Share route", "Negotiate a split and reduce escalation."),
				choice(ChoiceBribe, "Bribe", "Pay them to leave the job alone."),
			},
		},
	}
}

func choice(id ComplicationChoiceID, label string, description string) ComplicationChoice {
	return ComplicationChoice{
		ID:          id,
		Label:       label,
		Description: description,
	}
}

func ComplicationDefinitionFor(complicationType ComplicationType) (ComplicationDefinition, bool) {
	for _, definition := range FirstPlayableComplicationDefinitions() {
		if definition.Type == complicationType {
			return definition, true
		}
	}
	return ComplicationDefinition{}, false
}

func QueueComplicationFromResult(state *GameState, result JobResult) (Complication, bool) {
	if !result.Complication || result.ComplicationType == ComplicationNone {
		return Complication{}, false
	}
	definition, ok := ComplicationDefinitionFor(result.ComplicationType)
	if !ok {
		return Complication{}, false
	}

	complication := Complication{
		ID:                 nextComplicationID(*state),
		Type:               result.ComplicationType,
		Status:             ComplicationPending,
		Title:              definition.Title,
		Prompt:             definition.Prompt,
		Choices:            append([]ComplicationChoice(nil), definition.Choices...),
		Turn:               state.Turn,
		Night:              state.Night,
		JobID:              result.JobID,
		JobTitle:           result.JobTitle,
		RunnerID:           result.RunnerID,
		RunnerName:         result.RunnerName,
		FactionID:          result.FactionID,
		Outcome:            result.Outcome,
		CargoIntegrity:     result.CargoIntegrity,
		CargoIntegrityLoss: result.CargoIntegrityLoss,
		DelayTurns:         delayTurnsFromResult(result),
		Factors:            append([]string(nil), result.Factors...),
		Summary:            complicationSummary(definition, result),
	}
	state.Complications = append(state.Complications, complication)
	return complication, true
}

func nextComplicationID(state GameState) ComplicationID {
	return ComplicationID(fmt.Sprintf("cmp-%02d-%02d-%02d", state.Night, state.Turn, len(state.Complications)+1))
}

func complicationSummary(definition ComplicationDefinition, result JobResult) string {
	return fmt.Sprintf("%s complication on %s with %s.", definition.Title, result.JobTitle, result.RunnerName)
}

func ResolveComplicationChoice(state *GameState, complicationID ComplicationID, choiceID ComplicationChoiceID) (ComplicationResolution, error) {
	index := findComplicationIndex(state.Complications, complicationID)
	if index < 0 {
		return ComplicationResolution{}, ErrComplicationNotFound
	}
	if state.Complications[index].Status != ComplicationPending {
		return ComplicationResolution{}, ErrComplicationAlreadyHandled
	}
	if !complicationHasChoice(state.Complications[index], choiceID) {
		return ComplicationResolution{}, ErrComplicationChoiceNotFound
	}

	complication := &state.Complications[index]
	resolution, err := applyComplicationChoiceEffects(state, *complication, choiceID)
	if err != nil {
		return ComplicationResolution{}, err
	}

	complication.Status = ComplicationResolved
	complication.ResolvedBy = choiceID
	complication.CargoIntegrity = clampInt(complication.CargoIntegrity+resolution.CargoIntegrityDelta, 0, DefaultCargoIntegrity)
	complication.CargoIntegrityLoss = DefaultCargoIntegrity - complication.CargoIntegrity
	complication.DelayTurns += resolution.DelayTurns
	complication.ResolutionEffects = append([]string(nil), resolution.Effects...)
	complication.Summary = resolvedComplicationSummary(*complication, choiceID)
	resolution.ComplicationID = complication.ID
	resolution.ChoiceID = choiceID
	resolution.Summary = complication.Summary
	report := complicationResolutionReport(resolution)
	state.EventLog = append(state.EventLog, LogEntry{
		Turn: state.Turn,
		Text: report,
	})
	appendComplicationResolutionMessage(state, report)
	return resolution, nil
}

func applyComplicationChoiceEffects(state *GameState, complication Complication, choiceID ComplicationChoiceID) (ComplicationResolution, error) {
	resolution := ComplicationResolution{}
	addEffect := func(effect string) {
		resolution.Effects = append(resolution.Effects, effect)
	}

	switch complication.Type {
	case ComplicationCheckpoint:
		switch choiceID {
		case ChoiceTalkThrough:
			stressGain := checkpointStressGain(*state, ComplicationMinorStressGain)
			addRunnerStress(state, complication.RunnerID, stressGain)
			resolution.RunnerStressDelta += stressGain
			if stressGain > 0 {
				addEffect(fmt.Sprintf("runner stress +%d", stressGain))
			}
			if stressGain < ComplicationMinorStressGain {
				addEffect("fake credentials covered checkpoint")
			}
		case ChoiceReroute:
			addRunnerStress(state, complication.RunnerID, ComplicationMinorStressGain)
			loseDispatchIntegrity(state, ComplicationMinorIntegrityLoss)
			resolution.RunnerStressDelta += ComplicationMinorStressGain
			resolution.DispatchIntegrityDelta -= ComplicationMinorIntegrityLoss
			resolution.DelayTurns += ComplicationMinorDelayTurns
			addEffect("runner stress +1")
			addEffect("dispatch integrity -1")
			addEffect("delay +1")
		case ChoiceBribe:
			bribeCost := checkpointBribeCost(*state)
			if _, err := SpendBribe(state, bribeCost, fmt.Sprintf("%s checkpoint bribe", complication.JobTitle)); err != nil {
				return ComplicationResolution{}, err
			}
			state.Heat = maxInt(state.Heat-ComplicationMinorHeatGain, StartingHeat)
			resolution.CreditsDelta -= bribeCost
			resolution.HeatDelta -= ComplicationMinorHeatGain
			addEffect(fmt.Sprintf("credits -%d", bribeCost))
			addEffect("heat -1")
			if bribeCost < ComplicationBribeCost {
				addEffect("fake credentials reduced bribe")
			}
		case ChoiceAbandon:
			loseDispatchIntegrity(state, ComplicationModerateIntegrityLoss)
			addFactionSuspicion(state, complication.FactionID, 1)
			resolution.DispatchIntegrityDelta -= ComplicationModerateIntegrityLoss
			resolution.FactionSuspicionDelta++
			addEffect("dispatch integrity -3")
			addEffect("faction suspicion +1")
		}
	case ComplicationScannerSweep:
		switch choiceID {
		case ChoiceHideCargo:
			addRunnerStress(state, complication.RunnerID, ComplicationMinorStressGain)
			loseDispatchIntegrity(state, ComplicationMinorIntegrityLoss)
			resolution.RunnerStressDelta += ComplicationMinorStressGain
			resolution.DispatchIntegrityDelta -= ComplicationMinorIntegrityLoss
			resolution.DelayTurns += ComplicationMinorDelayTurns
			addEffect("runner stress +1")
			addEffect("dispatch integrity -1")
			addEffect("delay +1")
		case ChoiceRushThrough:
			addRunnerStress(state, complication.RunnerID, ComplicationMinorStressGain)
			gainHeat(state, ComplicationMajorHeatGain)
			resolution.RunnerStressDelta += ComplicationMinorStressGain
			resolution.HeatDelta += ComplicationMajorHeatGain
			resolution.CargoIntegrityDelta -= ComplicationMinorCargoLoss
			addEffect("runner stress +1")
			addEffect("heat +2")
			addEffect("cargo integrity -10")
		case ChoiceSpoofTag:
			loseDispatchIntegrity(state, ComplicationModerateIntegrityLoss)
			gainHeat(state, ComplicationMinorHeatGain)
			resolution.DispatchIntegrityDelta -= ComplicationModerateIntegrityLoss
			resolution.HeatDelta += ComplicationMinorHeatGain
			resolution.FactionReputationDelta++
			addFactionReputation(state, complication.FactionID, 1)
			addEffect("dispatch integrity -3")
			addEffect("heat +1")
			addEffect("faction reputation +1")
		}
	case ComplicationRunnerPanic:
		switch choiceID {
		case ChoiceReassure:
			addRunnerStress(state, complication.RunnerID, -ComplicationStressRelief)
			resolution.RunnerStressDelta -= ComplicationStressRelief
			resolution.DelayTurns += ComplicationMinorDelayTurns
			addEffect("runner stress -2")
			addEffect("delay +1")
		case ChoiceOrderForward:
			addRunnerStress(state, complication.RunnerID, ComplicationMajorStressGain)
			addRunnerLoyalty(state, complication.RunnerID, -1)
			resolution.RunnerStressDelta += ComplicationMajorStressGain
			resolution.RunnerLoyaltyDelta--
			resolution.CargoIntegrityDelta -= ComplicationMinorCargoLoss
			addEffect("runner stress +2")
			addEffect("runner loyalty -1")
			addEffect("cargo integrity -10")
		case ChoiceAbort:
			loseDispatchIntegrity(state, ComplicationModerateIntegrityLoss)
			addFactionSuspicion(state, complication.FactionID, 1)
			resolution.DispatchIntegrityDelta -= ComplicationModerateIntegrityLoss
			resolution.FactionSuspicionDelta++
			addEffect("dispatch integrity -3")
			addEffect("faction suspicion +1")
		}
	case ComplicationCargoLeak:
		switch choiceID {
		case ChoiceContinue:
			loseDispatchIntegrity(state, ComplicationModerateIntegrityLoss)
			resolution.DispatchIntegrityDelta -= ComplicationModerateIntegrityLoss
			resolution.CargoIntegrityDelta -= ComplicationMajorCargoLoss
			addEffect("dispatch integrity -3")
			addEffect("cargo integrity -35")
		case ChoiceSeekClinic:
			if _, err := SpendTreatment(state, ComplicationClinicCost, fmt.Sprintf("%s cargo leak response", complication.JobTitle)); err != nil {
				return ComplicationResolution{}, err
			}
			addRunnerStress(state, complication.RunnerID, ComplicationMinorStressGain)
			resolution.CreditsDelta -= ComplicationClinicCost
			resolution.RunnerStressDelta += ComplicationMinorStressGain
			resolution.CargoIntegrityDelta -= ComplicationMinorCargoLoss
			resolution.DelayTurns += ComplicationMajorDelayTurns
			addEffect("credits -60")
			addEffect("runner stress +1")
			addEffect("cargo integrity -10")
			addEffect("delay +2")
		case ChoiceDumpCargo:
			loseDispatchIntegrity(state, ComplicationMajorIntegrityLoss)
			addFactionSuspicion(state, complication.FactionID, 1)
			resolution.DispatchIntegrityDelta -= ComplicationMajorIntegrityLoss
			resolution.CargoIntegrityDelta -= DefaultCargoIntegrity
			resolution.FactionSuspicionDelta++
			addEffect("dispatch integrity -5")
			addEffect("cargo integrity -100")
			addEffect("faction suspicion +1")
		}
	case ComplicationSignalLoss:
		switch choiceID {
		case ChoiceTrustRoute:
			gainHeat(state, ComplicationMinorHeatGain)
			resolution.HeatDelta += ComplicationMinorHeatGain
			resolution.CargoIntegrityDelta -= ComplicationMinorCargoLoss
			addEffect("heat +1")
			addEffect("cargo integrity -10")
		case ChoiceWait:
			addRunnerStress(state, complication.RunnerID, ComplicationMinorStressGain)
			resolution.RunnerStressDelta += ComplicationMinorStressGain
			resolution.DelayTurns += ComplicationMinorDelayTurns
			addEffect("runner stress +1")
			addEffect("delay +1")
		case ChoiceReroute:
			loseDispatchIntegrity(state, ComplicationMinorIntegrityLoss)
			resolution.DispatchIntegrityDelta -= ComplicationMinorIntegrityLoss
			resolution.DelayTurns += ComplicationMinorDelayTurns
			addEffect("dispatch integrity -1")
			addEffect("delay +1")
		}
	}

	return resolution, nil
}

func findComplicationIndex(complications []Complication, complicationID ComplicationID) int {
	for i, complication := range complications {
		if complication.ID == complicationID {
			return i
		}
	}
	return -1
}

func complicationHasChoice(complication Complication, choiceID ComplicationChoiceID) bool {
	for _, choice := range complication.Choices {
		if choice.ID == choiceID {
			return true
		}
	}
	return false
}

func ComplicationChoiceRequiresConfirmation(choiceID ComplicationChoiceID) bool {
	switch choiceID {
	case ChoiceAbandon,
		ChoiceAbort,
		ChoiceBribe,
		ChoicePayToll,
		ChoiceSeekClinic,
		ChoiceContinue,
		ChoiceDumpCargo,
		ChoiceSpoofTag,
		ChoiceOrderForward,
		ChoiceThreaten,
		ChoiceSedateWitness,
		ChoiceBurnNode,
		ChoiceBreakCurfew,
		ChoiceBlockCourier:
		return true
	default:
		return false
	}
}

func addRunnerStress(state *GameState, runnerID RunnerID, delta int) {
	if runnerIndex := findRunnerIndex(state.Runners, runnerID); runnerIndex >= 0 {
		state.Runners[runnerIndex].Stress = clampInt(state.Runners[runnerIndex].Stress+delta, 0, MaxRunnerStress)
	}
}

func addRunnerLoyalty(state *GameState, runnerID RunnerID, delta int) {
	if runnerIndex := findRunnerIndex(state.Runners, runnerID); runnerIndex >= 0 {
		state.Runners[runnerIndex].Loyalty = maxInt(state.Runners[runnerIndex].Loyalty+delta, MinRunnerLoyalty)
	}
}

func addFactionSuspicion(state *GameState, factionID FactionID, delta int) {
	if factionIndex := findFactionIndex(state.Factions, factionID); factionIndex >= 0 {
		state.Factions[factionIndex].Suspicion += delta
	}
}

func addFactionReputation(state *GameState, factionID FactionID, delta int) {
	if factionIndex := findFactionIndex(state.Factions, factionID); factionIndex >= 0 {
		state.Factions[factionIndex].Reputation += delta
	}
}

func gainHeat(state *GameState, amount int) {
	state.Heat = clampInt(state.Heat+amount, StartingHeat, MaximumHeat)
}

func loseDispatchIntegrity(state *GameState, amount int) {
	state.DispatchIntegrity = maxInt(state.DispatchIntegrity-amount, DispatchIntegrityFailure)
}

func resolvedComplicationSummary(complication Complication, choiceID ComplicationChoiceID) string {
	label := string(choiceID)
	for _, choice := range complication.Choices {
		if choice.ID == choiceID {
			label = choice.Label
			break
		}
	}
	return fmt.Sprintf("Resolved %s on %s with %s.", complication.Title, complication.JobTitle, label)
}

func appendComplicationOpenedMessage(state *GameState, complication Complication) {
	state.Messages = append(state.Messages, Message{
		Turn:    state.Turn,
		From:    "after-action",
		Subject: "complication opened",
		Body:    complicationOpenReport(complication),
	})
}

func appendComplicationResolutionMessage(state *GameState, report string) {
	state.Messages = append(state.Messages, Message{
		Turn:    state.Turn,
		From:    "after-action",
		Subject: "complication resolved",
		Body:    report,
	})
}

func complicationOpenReport(complication Complication) string {
	parts := []string{complication.Summary}
	if complication.Prompt != "" {
		parts = append(parts, complication.Prompt)
	}
	if len(complication.Choices) > 0 {
		labels := make([]string, 0, len(complication.Choices))
		for _, choice := range complication.Choices {
			labels = append(labels, choice.Label)
		}
		parts = append(parts, "Choices: "+strings.Join(labels, ", ")+".")
	}
	return strings.Join(parts, " ")
}

func complicationResolutionReport(resolution ComplicationResolution) string {
	if len(resolution.Effects) == 0 {
		return resolution.Summary
	}
	return resolution.Summary + " Effects: " + strings.Join(resolution.Effects, ", ") + "."
}

func delayTurnsFromResult(result JobResult) int {
	if result.Delay {
		return ComplicationMinorDelayTurns
	}
	return 0
}
