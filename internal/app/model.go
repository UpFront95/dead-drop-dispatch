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
	PanelJobs
	PanelRunners
	PanelMessages
	PanelDetail
	panelCount
)

type Model struct {
	state    game.GameState
	width    int
	height   int
	focused  Panel
	showHelp bool
	styles   tui.Styles
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
		case "tab":
			m.focused = (m.focused + 1) % panelCount
		case "shift+tab":
			m.focused = (m.focused + panelCount - 1) % panelCount
		case "?":
			m.showHelp = !m.showHelp
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	return tui.RenderDashboard(tui.DashboardView{
		State:    m.state,
		Width:    m.width,
		Height:   m.height,
		Focused:  int(m.focused),
		ShowHelp: m.showHelp,
		Styles:   m.styles,
	})
}
