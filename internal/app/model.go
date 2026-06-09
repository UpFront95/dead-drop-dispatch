package app

import (
	tea "charm.land/bubbletea/v2"

	"dead-drop-dispatch/internal/content"
	"dead-drop-dispatch/internal/game"
	"dead-drop-dispatch/internal/tui"
)

type Panel int

const (
	PanelCity Panel = iota
	PanelRunners
	PanelJobs
	PanelMessages
	PanelDetail
	panelCount
)

type ScreenTab int

const (
	ScreenDashboard ScreenTab = iota
	ScreenJobs
	ScreenRouting
	ScreenRunners
	ScreenEquipment
	ScreenHelp
	screenTabCount
)

type Model struct {
	state              game.GameState
	width              int
	height             int
	focused            Panel
	tab                ScreenTab
	selectedDistrict   int
	selectedJobIndex   int
	selectedRunner     int
	selectedRouteIndex int
	selectedMessage    int
	selectedResponse   int
	messageScroll      int
	notice             string
	showDistrictBrief  bool
	showHelp           bool
	styles             tui.Styles
}

func New(seed int64) Model {
	return Model{
		state:  content.InitialGameState(seed),
		width:  tui.TargetWidth,
		height: tui.TargetHeight,
		styles: tui.NewStyles(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.RequestBackgroundColor
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		m.styles = tui.NewStylesForBackground(msg.IsDark())
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			m.moveSelection(-1)
		case "down", "j":
			m.moveSelection(1)
		case "enter":
			m.confirmSelection()
		case "esc":
			m.back()
		case "r":
			m.cycleRoute()
		case " ", "space":
			m.advanceTurnPhase()
		case "tab":
			m.focused = (m.focused + 1) % panelCount
		case "shift+tab":
			m.focused = (m.focused + panelCount - 1) % panelCount
		case "[":
			m.tab = (m.tab + screenTabCount - 1) % screenTabCount
		case "]":
			m.tab = (m.tab + 1) % screenTabCount
		case "1":
			m.tab = ScreenDashboard
		case "2":
			m.tab = ScreenJobs
		case "3":
			m.tab = ScreenRouting
		case "4":
			m.tab = ScreenRunners
		case "5":
			m.tab = ScreenEquipment
		case "6":
			m.tab = ScreenHelp
		case "?":
			m.showHelp = !m.showHelp
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	return tui.RenderDashboard(tui.DashboardView{
		State:              m.state,
		Width:              m.width,
		Height:             m.height,
		Focused:            int(m.focused),
		ActiveTab:          int(m.tab),
		SelectedDistrict:   m.selectedDistrict,
		SelectedJobIndex:   m.selectedJobIndex,
		SelectedRunner:     m.selectedRunner,
		SelectedRouteIndex: m.selectedRouteIndex,
		SelectedMessage:    m.selectedMessage,
		SelectedResponse:   m.selectedResponse,
		MessageScroll:      m.messageScroll,
		Notice:             m.notice,
		ShowDistrictBrief:  m.showDistrictBrief,
		ShowHelp:           m.showHelp,
		Styles:             m.styles,
	})
}

func (m *Model) moveSelection(delta int) {
	m.notice = ""
	switch m.focused {
	case PanelCity:
		m.selectedDistrict = wrapIndex(m.selectedDistrict+delta, len(m.state.Districts))
	case PanelJobs:
		m.selectedJobIndex = wrapIndex(m.selectedJobIndex+delta, len(m.state.AvailableJobs))
	case PanelRunners:
		m.selectedRunner = wrapIndex(m.selectedRunner+delta, len(m.state.Runners))
	case PanelDetail:
		m.selectedRouteIndex = wrapIndex(m.selectedRouteIndex+delta, pendingRouteCount(m.state))
	case PanelMessages:
		m.selectedMessage = wrapIndex(m.selectedMessage+delta, len(m.state.Messages))
		m.messageScroll = maxInt(0, m.selectedMessage-3)
		m.selectedResponse = 0
	}
}

func (m *Model) confirmSelection() {
	switch m.focused {
	case PanelCity:
		m.openDistrictBriefing()
	case PanelJobs:
		m.acceptSelectedJob()
	case PanelRunners:
		m.assignPendingJob()
	case PanelMessages:
		m.respondToSelectedMessage()
	default:
		m.notice = "Focus JOB BOARD to accept, RUNNERS to assign."
	}
}

func (m *Model) back() {
	if m.showDistrictBrief {
		m.showDistrictBrief = false
		m.notice = ""
		return
	}
	if m.showHelp {
		m.showHelp = false
		return
	}
	if len(m.state.AcceptedJobs) > 0 {
		job := m.state.AcceptedJobs[0]
		if err := game.CancelAcceptedJob(&m.state, job.ID); err != nil {
			m.notice = err.Error()
			return
		}
		m.selectedRouteIndex = 0
		m.selectedJobIndex = wrapIndex(m.selectedJobIndex, len(m.state.AvailableJobs))
		m.focused = PanelJobs
		m.notice = "Canceled " + job.Title + ". Contract returned to board."
		return
	}
	m.notice = ""
}

func (m *Model) openDistrictBriefing() {
	if len(m.state.Districts) == 0 {
		m.notice = "No district data available."
		return
	}
	m.showDistrictBrief = true
	m.notice = ""
}

func (m *Model) acceptSelectedJob() {
	if len(m.state.AvailableJobs) == 0 {
		m.notice = "No posted jobs to accept."
		return
	}
	job := m.state.AvailableJobs[wrapIndex(m.selectedJobIndex, len(m.state.AvailableJobs))]
	if err := game.AcceptJob(&m.state, job.ID); err != nil {
		m.notice = err.Error()
		return
	}
	m.selectedJobIndex = wrapIndex(m.selectedJobIndex, len(m.state.AvailableJobs))
	m.selectedRouteIndex = 0
	m.focused = PanelRunners
	m.notice = "Accepted " + job.Title + ". Select a runner."
}

func (m *Model) assignPendingJob() {
	if len(m.state.AcceptedJobs) == 0 {
		m.notice = "Accept a job first."
		return
	}
	if len(m.state.Runners) == 0 {
		m.notice = "No runners available."
		return
	}
	job := m.state.AcceptedJobs[0]
	runner := m.state.Runners[wrapIndex(m.selectedRunner, len(m.state.Runners))]
	if len(job.Routes) == 0 {
		m.notice = "Accepted job has no routes."
		return
	}
	route := job.Routes[wrapIndex(m.selectedRouteIndex, len(job.Routes))]
	bundling := runnerLoad(m.state, runner.ID) > 0
	if err := game.AssignAcceptedJob(&m.state, job.ID, runner.ID, route.ID); err != nil {
		m.notice = err.Error()
		return
	}
	m.selectedRouteIndex = 0
	if bundling {
		m.notice = "Bundled " + job.Title + " with " + runner.Name + "."
	} else {
		m.notice = "Assigned " + job.Title + " to " + runner.Name + "."
	}
}

func (m *Model) cycleRoute() {
	if m.focused == PanelMessages {
		m.selectedResponse = wrapIndex(m.selectedResponse+1, selectedMessageResponseCount(m.state, m.selectedMessage))
		m.notice = ""
		return
	}
	m.selectedRouteIndex = wrapIndex(m.selectedRouteIndex+1, pendingRouteCount(m.state))
	m.notice = ""
}

func (m *Model) advanceTurnPhase() {
	advance := game.AdvanceTurnPhase(&m.state)
	m.selectedRouteIndex = 0
	m.selectedMessage = wrapIndex(m.selectedMessage, len(m.state.Messages))
	m.selectedResponse = 0
	results := advance.Results
	if len(results) == 1 {
		m.notice = "Resolved " + results[0].JobTitle + ": " + string(results[0].Outcome) + "."
		return
	}
	m.notice = advance.Summary
}

func (m *Model) respondToSelectedMessage() {
	if len(m.state.Messages) == 0 {
		m.notice = "No messages to answer."
		return
	}
	messageIndex := wrapIndex(m.selectedMessage, len(m.state.Messages))
	message := m.state.Messages[messageIndex]
	responses := messageResponseOptions(message)
	if len(responses) == 0 {
		m.notice = "Selected message has no open response."
		return
	}
	response := responses[wrapIndex(m.selectedResponse, len(responses))]
	if _, err := game.ResolveMessageResponse(&m.state, message.ID, response.ID); err != nil {
		m.notice = err.Error()
		return
	}
	m.selectedMessage = messageIndex
	m.selectedResponse = 0
	m.notice = "Sent response: " + response.Label + "."
}

func pendingRouteCount(state game.GameState) int {
	if len(state.AcceptedJobs) == 0 {
		return 0
	}
	return len(state.AcceptedJobs[0].Routes)
}

func selectedMessageResponseCount(state game.GameState, selected int) int {
	if len(state.Messages) == 0 {
		return 0
	}
	return len(messageResponseOptions(state.Messages[wrapIndex(selected, len(state.Messages))]))
}

func messageResponseOptions(message game.Message) []game.MessageResponseAction {
	if message.Status == game.MessageResolved || message.Audience == "" {
		return nil
	}
	if len(message.Responses) > 0 {
		return append([]game.MessageResponseAction(nil), message.Responses...)
	}
	return game.MessageResponseActionsFor(message.Audience)
}

func wrapIndex(index int, length int) int {
	if length <= 0 {
		return 0
	}
	index %= length
	if index < 0 {
		index += length
	}
	return index
}

func runnerLoad(state game.GameState, runnerID game.RunnerID) int {
	load := 0
	for _, active := range state.ActiveJobs {
		if active.RunnerID == runnerID {
			load++
		}
	}
	return load
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
