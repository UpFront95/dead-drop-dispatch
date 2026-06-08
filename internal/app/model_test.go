package app

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"dead-drop-dispatch/internal/game"
	"dead-drop-dispatch/internal/tui"
)

func TestNewModelUsesInitialState(t *testing.T) {
	model := New(99)

	if model.width != tui.TargetWidth {
		t.Fatalf("width = %d, want %d", model.width, tui.TargetWidth)
	}
	if model.height != tui.TargetHeight {
		t.Fatalf("height = %d, want %d", model.height, tui.TargetHeight)
	}
	if model.state.RandomSeed != 99 {
		t.Fatalf("seed = %d, want 99", model.state.RandomSeed)
	}
	if len(model.state.Districts) != 6 {
		t.Fatalf("district count = %d, want 6", len(model.state.Districts))
	}
}

func TestModelUpdatesActiveTab(t *testing.T) {
	model := New(99)

	updated, _ := model.Update(keyPress("5"))
	model = updated.(Model)
	if model.tab != ScreenEquipment {
		t.Fatalf("tab = %d, want %d", model.tab, ScreenEquipment)
	}

	updated, _ = model.Update(keyPress("]"))
	model = updated.(Model)
	if model.tab != ScreenHelp {
		t.Fatalf("tab after ] = %d, want %d", model.tab, ScreenHelp)
	}

	updated, _ = model.Update(keyPress("]"))
	model = updated.(Model)
	if model.tab != ScreenDashboard {
		t.Fatalf("tab after wrap = %d, want %d", model.tab, ScreenDashboard)
	}

	updated, _ = model.Update(keyPress("["))
	model = updated.(Model)
	if model.tab != ScreenHelp {
		t.Fatalf("tab after reverse wrap = %d, want %d", model.tab, ScreenHelp)
	}
}

func TestModelTabKeyStillCyclesPanelFocus(t *testing.T) {
	model := New(99)

	updated, _ := model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyTab}))
	model = updated.(Model)

	if model.focused != PanelRunners {
		t.Fatalf("focused = %d, want %d", model.focused, PanelRunners)
	}
	if model.tab != ScreenDashboard {
		t.Fatalf("tab = %d, want %d", model.tab, ScreenDashboard)
	}
}

func TestModelOpensAndClosesDistrictBriefing(t *testing.T) {
	model := New(99)
	model.focused = PanelCity

	updated, _ := model.Update(keyPress("down"))
	model = updated.(Model)
	if model.selectedDistrict != 1 {
		t.Fatalf("selected district = %d, want 1", model.selectedDistrict)
	}

	updated, _ = model.Update(keyPress("enter"))
	model = updated.(Model)
	if !model.showDistrictBrief {
		t.Fatal("district briefing should open")
	}

	updated, _ = model.Update(keyPress("esc"))
	model = updated.(Model)
	if model.showDistrictBrief {
		t.Fatal("district briefing should close on esc")
	}
	if model.selectedDistrict != 1 {
		t.Fatalf("selected district after close = %d, want 1", model.selectedDistrict)
	}
}

func TestModelScrollsMessageFeedWhenFocused(t *testing.T) {
	model := New(99)
	model.focused = PanelMessages
	model.state.Messages = append(model.state.Messages, game.Message{
		ID:       "runner-check",
		Turn:     model.state.Turn,
		From:     "mira",
		Subject:  "runner check",
		Body:     "Need a call.",
		Audience: game.MessageAudienceRunner,
		Status:   game.MessageOpen,
	})

	updated, _ := model.Update(keyPress("down"))
	model = updated.(Model)
	if model.selectedMessage != 1 {
		t.Fatalf("selected message = %d, want 1", model.selectedMessage)
	}

	updated, _ = model.Update(keyPress("up"))
	model = updated.(Model)
	if model.selectedMessage != 0 {
		t.Fatalf("selected message after up = %d, want 0", model.selectedMessage)
	}

	updated, _ = model.Update(keyPress("up"))
	model = updated.(Model)
	if model.selectedMessage != 1 {
		t.Fatalf("selected message should wrap, got %d", model.selectedMessage)
	}
}

func TestModelRespondsToSelectedMessageFromDashboard(t *testing.T) {
	model := New(99)
	model.focused = PanelMessages
	model.state.Messages = []game.Message{{
		ID:       "runner-check",
		Turn:     model.state.Turn,
		From:     "mira",
		Subject:  "runner check",
		Body:     "Need a call.",
		Audience: game.MessageAudienceRunner,
		Status:   game.MessageOpen,
		Responses: []game.MessageResponseAction{
			mustResponseAction(t, game.ResponseAskMoreInfo),
			mustResponseAction(t, game.ResponseReassure),
		},
		TargetRunnerID: model.state.Runners[0].ID,
	}}

	updated, _ := model.Update(keyPress("r"))
	model = updated.(Model)
	if model.selectedResponse != 1 {
		t.Fatalf("selected response = %d, want 1", model.selectedResponse)
	}

	updated, _ = model.Update(keyPress("enter"))
	model = updated.(Model)

	if got, want := model.state.Messages[0].Status, game.MessageResolved; got != want {
		t.Fatalf("message status = %q, want %q", got, want)
	}
	if got, want := model.state.Messages[0].ResolvedBy, game.ResponseReassure; got != want {
		t.Fatalf("resolved by = %q, want %q", got, want)
	}
	if model.notice != "Sent response: Reassure." {
		t.Fatalf("notice = %q", model.notice)
	}
	if got, want := model.state.Runners[0].Loyalty, 5; got != want {
		t.Fatalf("runner loyalty = %d, want %d", got, want)
	}
	if got, want := model.state.Messages[len(model.state.Messages)-1].Subject, "response resolved"; got != want {
		t.Fatalf("last message subject = %q, want %q", got, want)
	}
}

func TestModelAcceptsAndAssignsJobFromDashboard(t *testing.T) {
	model := New(99)
	model.focused = PanelJobs
	jobID := model.state.AvailableJobs[0].ID
	runnerID := model.state.Runners[0].ID
	routeID := model.state.AvailableJobs[0].Routes[1].ID

	updated, _ := model.Update(keyPress("enter"))
	model = updated.(Model)
	if !containsAcceptedJob(model.state.AcceptedJobs, jobID) {
		t.Fatalf("accepted jobs missing %s", jobID)
	}
	if model.focused != PanelRunners {
		t.Fatalf("focused = %d, want %d after accepting job", model.focused, PanelRunners)
	}

	updated, _ = model.Update(keyPress("r"))
	model = updated.(Model)
	updated, _ = model.Update(keyPress("enter"))
	model = updated.(Model)

	if got, want := len(model.state.ActiveJobs), 1; got != want {
		t.Fatalf("active jobs = %d, want %d", got, want)
	}
	active := model.state.ActiveJobs[0]
	if active.JobID != jobID || active.RunnerID != runnerID || active.RouteID != routeID {
		t.Fatalf("active job = %+v, want job %s runner %s route %s", active, jobID, runnerID, routeID)
	}
	if model.state.Runners[0].State != game.RunnerOnJob {
		t.Fatalf("runner state = %q, want %q", model.state.Runners[0].State, game.RunnerOnJob)
	}
}

func TestModelResolvesActiveJobsFromDashboard(t *testing.T) {
	model := New(99)
	model.focused = PanelJobs

	updated, _ := model.Update(keyPress("enter"))
	model = updated.(Model)
	model.focused = PanelRunners
	updated, _ = model.Update(keyPress("enter"))
	model = updated.(Model)
	if got, want := len(model.state.ActiveJobs), 1; got != want {
		t.Fatalf("active jobs before resolve = %d, want %d", got, want)
	}

	updated, _ = model.Update(keyPress(" "))
	model = updated.(Model)

	if got, want := len(model.state.ActiveJobs), 0; got != want {
		t.Fatalf("active jobs after resolve = %d, want %d", got, want)
	}
	if got, want := len(model.state.LastResults), 1; got != want {
		t.Fatalf("last results = %d, want %d", got, want)
	}
	if model.notice == "" {
		t.Fatal("resolve should set a dashboard notice")
	}
}

func keyPress(text string) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Text: text, Code: []rune(text)[0]})
}

func containsAcceptedJob(jobs []game.Job, jobID string) bool {
	for _, job := range jobs {
		if job.ID == jobID {
			return true
		}
	}
	return false
}

func mustResponseAction(t *testing.T, actionID game.MessageResponseActionID) game.MessageResponseAction {
	t.Helper()
	action, ok := game.MessageResponseActionFor(actionID)
	if !ok {
		t.Fatalf("missing response action %s", actionID)
	}
	return action
}
