package game

import "fmt"

type TurnAdvance struct {
	From            GamePhase `json:"from"`
	To              GamePhase `json:"to"`
	Results         []JobResult
	Status          RunStatus `json:"status"`
	Recovery        RunnerRecoveryReport
	JobsGenerated   int    `json:"jobs_generated"`
	NightChanged    bool   `json:"night_changed"`
	WaitingForInput bool   `json:"waiting_for_input"`
	Summary         string `json:"summary"`
}

func AdvanceTurnPhase(state *GameState) TurnAdvance {
	if state.Phase == "" {
		state.Phase = PhaseDispatch
	}

	from := state.Phase
	status := EvaluateRunStatus(*state)
	if status.State != RunInProgress || state.Phase == PhaseGameOver {
		if status.State != RunInProgress {
			transitionToGameOver(state, status)
		} else {
			state.Phase = PhaseGameOver
		}
		return turnAdvance(from, state.Phase, nil, status, RunnerRecoveryReport{}, 0, false, false, status.Summary)
	}

	if hasPendingComplications(*state) && state.Phase != PhaseComplications {
		state.Phase = PhaseComplications
		return turnAdvance(from, state.Phase, nil, status, RunnerRecoveryReport{}, 0, false, true, "Resolve pending complications before the turn can continue.")
	}

	var results []JobResult
	recovery := RunnerRecoveryReport{}
	jobsGenerated := 0
	nightChanged := false
	waiting := false
	summary := ""

	switch state.Phase {
	case PhaseMessages:
		appendTurnLog(state, "Switchboard traffic reviewed.")
		jobsGenerated = len(RefreshAvailableJobs(state, DefaultJobsPerTurn))
		state.Phase = PhaseJobs
		summary = fmt.Sprintf("Messages reviewed. Job board refreshed with %d postings.", jobsGenerated)
	case PhaseJobs:
		appendTurnLog(state, "Job board reviewed.")
		state.Phase = PhaseDispatch
		summary = "Jobs shown. Assign runners from the dashboard."
	case PhaseDispatch:
		if len(state.ActiveJobs) == 0 {
			waiting = true
			summary = "Assign at least one runner before resolving the turn."
			break
		}
		state.Phase = PhaseResolve
		results = ResolveActiveJobs(state)
		state.Phase = nextPhaseAfterResolution(*state)
		waiting = state.Phase == PhaseComplications
		summary = resolutionAdvanceSummary(results, waiting)
	case PhaseResolve:
		if len(state.ActiveJobs) > 0 {
			results = ResolveActiveJobs(state)
		}
		state.Phase = nextPhaseAfterResolution(*state)
		waiting = state.Phase == PhaseComplications
		summary = resolutionAdvanceSummary(results, waiting)
	case PhaseComplications:
		if hasPendingComplications(*state) {
			waiting = true
			summary = "Resolve pending complications before reports can be filed."
			break
		}
		state.Phase = PhaseReports
		summary = "Complications resolved. Reports are ready."
	case PhaseReports:
		appendTurnLog(state, "After-action reports filed.")
		state.Phase = PhaseCityUpdate
		summary = "Reports filed. City state update is next."
	case PhaseCityUpdate:
		nightChanged, recovery = applyCityUpdate(state)
		state.Phase = PhaseMessages
		summary = clockAdvanceSummary(*state, nightChanged)
	default:
		state.Phase = PhaseDispatch
		summary = "Dispatch phase restored."
	}

	status = EvaluateRunStatus(*state)
	if status.State != RunInProgress {
		transitionToGameOver(state, status)
		summary = status.Summary
		waiting = false
	}

	return turnAdvance(from, state.Phase, results, status, recovery, jobsGenerated, nightChanged, waiting, summary)
}

func transitionToGameOver(state *GameState, status RunStatus) {
	state.Phase = PhaseGameOver
	appendRunEndReport(state, status)
}

func appendRunEndReport(state *GameState, status RunStatus) {
	if status.State == RunInProgress || hasRunEndMessage(*state, status.Reason) {
		return
	}

	state.EventLog = append(state.EventLog, LogEntry{
		Turn: state.Turn,
		Text: fmt.Sprintf("Run ended (%s): %s", status.Reason, status.Summary),
	})

	subject := "run lost"
	if status.State == RunWon {
		subject = "run complete"
	}
	state.Messages = append(state.Messages, Message{
		ID:      runEndMessageID(status.Reason),
		Turn:    state.Turn,
		From:    "switchboard",
		Subject: subject,
		Body:    status.Summary,
		Status:  MessageResolved,
	})
}

func hasRunEndMessage(state GameState, reason RunEndReason) bool {
	id := runEndMessageID(reason)
	for _, message := range state.Messages {
		if message.ID == id {
			return true
		}
	}
	return false
}

func runEndMessageID(reason RunEndReason) MessageID {
	return MessageID(fmt.Sprintf("run-end-%s", reason))
}

func turnAdvance(from GamePhase, to GamePhase, results []JobResult, status RunStatus, recovery RunnerRecoveryReport, jobsGenerated int, nightChanged bool, waiting bool, summary string) TurnAdvance {
	return TurnAdvance{
		From:            from,
		To:              to,
		Results:         results,
		Status:          status,
		Recovery:        recovery,
		JobsGenerated:   jobsGenerated,
		NightChanged:    nightChanged,
		WaitingForInput: waiting,
		Summary:         summary,
	}
}

func nextPhaseAfterResolution(state GameState) GamePhase {
	if hasPendingComplications(state) {
		return PhaseComplications
	}
	return PhaseReports
}

func resolutionAdvanceSummary(results []JobResult, waiting bool) string {
	if len(results) == 0 {
		if waiting {
			return "Pending complications need a response."
		}
		return "No active jobs resolved."
	}
	if waiting {
		return "Runs resolved. Pending complications need a response."
	}
	if len(results) == 1 {
		return fmt.Sprintf("Resolved %s: %s.", results[0].JobTitle, results[0].Outcome)
	}
	return "Resolved active runs."
}

func hasPendingComplications(state GameState) bool {
	for _, complication := range state.Complications {
		if complication.Status == ComplicationPending {
			return true
		}
	}
	return false
}

func applyCityUpdate(state *GameState) (bool, RunnerRecoveryReport) {
	appendTurnLog(state, "City state updated.")
	recovery := RecoverRunners(state)
	appendRunnerRecoveryLog(state, recovery)
	nightChanged := advanceRunClock(state)
	if nightChanged {
		applyEndofNightUpdate(state)
	}
	if !runComplete(*state) {
		appendTurnBriefMessage(state, nightChanged)
	}
	return nightChanged, recovery
}

func applyEndofNightUpdate(state *GameState) {
	_, _ = ApplyOperatingCost(state, NightlyOperatingCost, "nightly operating costs")

	if state.Heat > StartingHeat {
		decay := NightlyHeatDecay
		if state.Heat-StartingHeat < decay {
			decay = state.Heat - StartingHeat
		}
		state.Heat -= decay
		appendTurnLog(state, fmt.Sprintf("Heat decayed by %d (current: %d).", decay, state.Heat))
	}
}

func advanceRunClock(state *GameState) bool {
	state.Turn++
	if state.Turn <= state.TurnsPerNight {
		return false
	}
	state.Night++
	state.Turn = FirstTurn
	appendTurnLog(state, fmt.Sprintf("Night %d begins.", state.Night))
	return true
}

func clockAdvanceSummary(state GameState, nightChanged bool) string {
	if runComplete(state) {
		return "The seven-night run clock is complete."
	}
	if nightChanged {
		return fmt.Sprintf("Night %d begins. Messages are arriving.", state.Night)
	}
	return fmt.Sprintf("Turn %d begins. Messages are arriving.", state.Turn)
}

func appendTurnBriefMessage(state *GameState, nightChanged bool) {
	subject := "turn brief"
	body := fmt.Sprintf("Night %d, turn %d is live. Review messages, jobs, and active pressure.", state.Night, state.Turn)
	if nightChanged {
		subject = "night brief"
		body = fmt.Sprintf("Night %d is live. Review messages, jobs, and active pressure.", state.Night)
	}
	state.Messages = append(state.Messages, Message{
		ID:       MessageID(fmt.Sprintf("brief-n%02d-t%02d", state.Night, state.Turn)),
		Turn:     state.Turn,
		From:     "switchboard",
		Subject:  subject,
		Body:     body,
		Audience: MessageAudienceRunner,
		Status:   MessageOpen,
		Responses: []MessageResponseAction{
			mustMessageResponseAction(ResponseAskMoreInfo),
			mustMessageResponseAction(ResponseStall),
			mustMessageResponseAction(ResponseAccept),
		},
	})
}

func mustMessageResponseAction(actionID MessageResponseActionID) MessageResponseAction {
	action, _ := MessageResponseActionFor(actionID)
	return action
}

func appendTurnLog(state *GameState, text string) {
	state.EventLog = append(state.EventLog, LogEntry{
		Turn: state.Turn,
		Text: text,
	})
}

func appendRunnerRecoveryLog(state *GameState, recovery RunnerRecoveryReport) {
	if recovery.StressRecovered == 0 && recovery.InjuryTicks == 0 && recovery.RunnersReadied == 0 {
		return
	}
	state.EventLog = append(state.EventLog, LogEntry{
		Turn: state.Turn,
		Text: fmt.Sprintf("Runner recovery: stress -%d, injuries ticked %d, runners ready %d.", recovery.StressRecovered, recovery.InjuryTicks, recovery.RunnersReadied),
	})
}
