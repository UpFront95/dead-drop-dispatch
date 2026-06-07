package game

import (
	"errors"
	"fmt"
	"strings"
)

const (
	MessageResponseNegotiationCreditGain = 25
	MessageResponseMinorHeatGain         = 1
	MessageResponseMinorIntegrityLoss    = 1
	MessageResponseMajorIntegrityLoss    = 2
	MessageResponseRunnerStressGain      = 1
	MessageResponseRunnerStressRelief    = 1
	MessageResponseRunnerLoyaltyGain     = 1
	MessageResponseRunnerLoyaltyLoss     = 1
	MessageResponseFactionShift          = 1
)

var (
	ErrMessageNotFound         = errors.New("message not found")
	ErrMessageAlreadyHandled   = errors.New("message already handled")
	ErrMessageResponseNotFound = errors.New("message response not found")
	ErrMessageAudienceNotSet   = errors.New("message audience not set")
)

func FixedMessageResponseActions() []MessageResponseAction {
	return copyMessageResponseActions([]MessageResponseAction{
		responseAction(
			ResponseRefuse,
			"Refuse",
			"Reject the request or demand.",
			MessageAudienceClient,
			MessageAudienceFaction,
		),
		responseAction(
			ResponseAskMorePay,
			"Ask more pay",
			"Push for better compensation before accepting risk.",
			MessageAudienceClient,
			MessageAudienceFaction,
		),
		responseAction(
			ResponseAskMoreInfo,
			"Ask more info",
			"Request more details before committing.",
			MessageAudienceClient,
			MessageAudienceRunner,
			MessageAudienceFaction,
		),
		responseAction(
			ResponseThreaten,
			"Threaten",
			"Escalate pressure to force compliance.",
			MessageAudienceClient,
			MessageAudienceRunner,
			MessageAudienceFaction,
		),
		responseAction(
			ResponseReassure,
			"Reassure",
			"Calm the contact and lower immediate pressure.",
			MessageAudienceRunner,
		),
		responseAction(
			ResponseStall,
			"Stall",
			"Buy time without resolving the demand.",
			MessageAudienceClient,
			MessageAudienceRunner,
			MessageAudienceFaction,
		),
		responseAction(
			ResponseDeceive,
			"Deceive",
			"Mislead the contact to preserve options.",
			MessageAudienceClient,
			MessageAudienceRunner,
			MessageAudienceFaction,
		),
		responseAction(
			ResponseCancel,
			"Cancel",
			"End the current arrangement or pending action.",
			MessageAudienceClient,
			MessageAudienceRunner,
			MessageAudienceFaction,
		),
		responseAction(
			ResponseAccept,
			"Accept",
			"Commit to the presented terms or instruction.",
			MessageAudienceClient,
			MessageAudienceRunner,
			MessageAudienceFaction,
		),
	})
}

func MessageResponseActionFor(actionID MessageResponseActionID) (MessageResponseAction, bool) {
	for _, action := range FixedMessageResponseActions() {
		if action.ID == actionID {
			return action, true
		}
	}
	return MessageResponseAction{}, false
}

func MessageResponseActionsFor(audience MessageAudience) []MessageResponseAction {
	actions := FixedMessageResponseActions()
	filtered := make([]MessageResponseAction, 0, len(actions))
	for _, action := range actions {
		if messageResponseActionSupportsAudience(action, audience) {
			filtered = append(filtered, action)
		}
	}
	return filtered
}

func MessageResponseActionAllowed(audience MessageAudience, actionID MessageResponseActionID) bool {
	action, ok := MessageResponseActionFor(actionID)
	if !ok {
		return false
	}
	return messageResponseActionSupportsAudience(action, audience)
}

func ResolveMessageResponse(state *GameState, messageID MessageID, actionID MessageResponseActionID) (MessageResponseResolution, error) {
	index := findMessageIndex(state.Messages, messageID)
	if index < 0 {
		return MessageResponseResolution{}, ErrMessageNotFound
	}
	if state.Messages[index].Status == MessageResolved {
		return MessageResponseResolution{}, ErrMessageAlreadyHandled
	}
	if state.Messages[index].Audience == "" {
		return MessageResponseResolution{}, ErrMessageAudienceNotSet
	}
	if !messageHasResponse(state.Messages[index], actionID) {
		return MessageResponseResolution{}, ErrMessageResponseNotFound
	}

	message := &state.Messages[index]
	resolution, err := applyMessageResponseEffects(state, *message, actionID)
	if err != nil {
		return MessageResponseResolution{}, err
	}

	message.Status = MessageResolved
	message.ResolvedBy = actionID
	message.ResolutionEffects = append([]string(nil), resolution.Effects...)
	message.Summary = resolvedMessageSummary(*message, actionID)
	resolution.MessageID = message.ID
	resolution.ActionID = actionID
	resolution.Summary = message.Summary

	report := messageResponseResolutionReport(resolution)
	state.EventLog = append(state.EventLog, LogEntry{
		Turn: state.Turn,
		Text: report,
	})
	appendMessageResponseReport(state, report)
	return resolution, nil
}

func responseAction(
	id MessageResponseActionID,
	label string,
	description string,
	audiences ...MessageAudience,
) MessageResponseAction {
	return MessageResponseAction{
		ID:          id,
		Label:       label,
		Description: description,
		Audiences:   append([]MessageAudience(nil), audiences...),
	}
}

func findMessageIndex(messages []Message, messageID MessageID) int {
	for i, message := range messages {
		if message.ID == messageID {
			return i
		}
	}
	return -1
}

func messageHasResponse(message Message, actionID MessageResponseActionID) bool {
	if len(message.Responses) == 0 {
		return MessageResponseActionAllowed(message.Audience, actionID)
	}
	for _, response := range message.Responses {
		if response.ID == actionID && MessageResponseActionAllowed(message.Audience, actionID) {
			return true
		}
	}
	return false
}

func applyMessageResponseEffects(state *GameState, message Message, actionID MessageResponseActionID) (MessageResponseResolution, error) {
	resolution := MessageResponseResolution{}
	addEffect := func(effect string) {
		resolution.Effects = append(resolution.Effects, effect)
	}

	switch actionID {
	case ResponseRefuse:
		loseDispatchIntegrity(state, MessageResponseMinorIntegrityLoss)
		resolution.DispatchIntegrityDelta -= MessageResponseMinorIntegrityLoss
		addEffect("dispatch integrity -1")
		if message.TargetFactionID != "" {
			addFactionSuspicion(state, message.TargetFactionID, MessageResponseFactionShift)
			resolution.FactionSuspicionDelta += MessageResponseFactionShift
			addEffect("faction suspicion +1")
		}
	case ResponseAskMorePay:
		if _, err := ApplyPayout(state, MessageResponseNegotiationCreditGain, "response negotiation"); err != nil {
			return MessageResponseResolution{}, err
		}
		resolution.CreditsDelta += MessageResponseNegotiationCreditGain
		addEffect("credits +25")
		if message.TargetFactionID != "" {
			addFactionReputation(state, message.TargetFactionID, -MessageResponseFactionShift)
			resolution.FactionReputationDelta -= MessageResponseFactionShift
			addEffect("faction reputation -1")
		}
	case ResponseAskMoreInfo:
		addEffect("intel requested")
	case ResponseThreaten:
		gainHeat(state, MessageResponseMinorHeatGain)
		resolution.HeatDelta += MessageResponseMinorHeatGain
		addEffect("heat +1")
		if message.TargetRunnerID != "" {
			addRunnerStress(state, message.TargetRunnerID, MessageResponseRunnerStressGain)
			addRunnerLoyalty(state, message.TargetRunnerID, -MessageResponseRunnerLoyaltyLoss)
			resolution.RunnerStressDelta += MessageResponseRunnerStressGain
			resolution.RunnerLoyaltyDelta -= MessageResponseRunnerLoyaltyLoss
			addEffect("runner stress +1")
			addEffect("runner loyalty -1")
		}
		if message.TargetFactionID != "" {
			addFactionSuspicion(state, message.TargetFactionID, MessageResponseFactionShift)
			resolution.FactionSuspicionDelta += MessageResponseFactionShift
			addEffect("faction suspicion +1")
		}
	case ResponseReassure:
		if message.TargetRunnerID != "" {
			addRunnerStress(state, message.TargetRunnerID, -MessageResponseRunnerStressRelief)
			addRunnerLoyalty(state, message.TargetRunnerID, MessageResponseRunnerLoyaltyGain)
			resolution.RunnerStressDelta -= MessageResponseRunnerStressRelief
			resolution.RunnerLoyaltyDelta += MessageResponseRunnerLoyaltyGain
			addEffect("runner stress -1")
			addEffect("runner loyalty +1")
		}
	case ResponseStall:
		gainHeat(state, MessageResponseMinorHeatGain)
		resolution.HeatDelta += MessageResponseMinorHeatGain
		addEffect("heat +1")
	case ResponseDeceive:
		gainHeat(state, MessageResponseMinorHeatGain)
		loseDispatchIntegrity(state, MessageResponseMinorIntegrityLoss)
		resolution.HeatDelta += MessageResponseMinorHeatGain
		resolution.DispatchIntegrityDelta -= MessageResponseMinorIntegrityLoss
		addEffect("heat +1")
		addEffect("dispatch integrity -1")
		if message.TargetFactionID != "" {
			addFactionSuspicion(state, message.TargetFactionID, MessageResponseFactionShift)
			resolution.FactionSuspicionDelta += MessageResponseFactionShift
			addEffect("faction suspicion +1")
		}
	case ResponseCancel:
		loseDispatchIntegrity(state, MessageResponseMajorIntegrityLoss)
		resolution.DispatchIntegrityDelta -= MessageResponseMajorIntegrityLoss
		addEffect("dispatch integrity -2")
		if message.TargetFactionID != "" {
			addFactionSuspicion(state, message.TargetFactionID, MessageResponseFactionShift)
			resolution.FactionSuspicionDelta += MessageResponseFactionShift
			addEffect("faction suspicion +1")
		}
	case ResponseAccept:
		if message.TargetFactionID != "" {
			addFactionReputation(state, message.TargetFactionID, MessageResponseFactionShift)
			resolution.FactionReputationDelta += MessageResponseFactionShift
			addEffect("faction reputation +1")
		}
	}

	return resolution, nil
}

func resolvedMessageSummary(message Message, actionID MessageResponseActionID) string {
	label := string(actionID)
	for _, response := range message.Responses {
		if response.ID == actionID {
			label = response.Label
			break
		}
	}
	if label == string(actionID) {
		if response, ok := MessageResponseActionFor(actionID); ok {
			label = response.Label
		}
	}
	return fmt.Sprintf("Resolved message %s with %s.", message.Subject, label)
}

func messageResponseResolutionReport(resolution MessageResponseResolution) string {
	if len(resolution.Effects) == 0 {
		return resolution.Summary
	}
	return resolution.Summary + " Effects: " + strings.Join(resolution.Effects, ", ") + "."
}

func appendMessageResponseReport(state *GameState, report string) {
	state.Messages = append(state.Messages, Message{
		Turn:    state.Turn,
		From:    "switchboard",
		Subject: "response resolved",
		Body:    report,
	})
}

func messageResponseActionSupportsAudience(action MessageResponseAction, audience MessageAudience) bool {
	for _, candidate := range action.Audiences {
		if candidate == audience {
			return true
		}
	}
	return false
}

func copyMessageResponseActions(actions []MessageResponseAction) []MessageResponseAction {
	copied := make([]MessageResponseAction, len(actions))
	for i, action := range actions {
		copied[i] = action
		copied[i].Audiences = append([]MessageAudience(nil), action.Audiences...)
	}
	return copied
}
