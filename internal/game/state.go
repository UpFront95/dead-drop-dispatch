package game

const (
	StartingCredits                 = 800
	StartingHeat                    = 0
	StartingDispatchIntegrity       = 100
	FirstNight                      = 1
	FirstTurn                       = 1
	DefaultTurnsPerNight            = 6
	DefaultRunNights                = 7
	DefaultRandomSeed         int64 = 1701
	MaxJobsPerRunner                = 2
)

type GameState struct {
	Turn              int         `json:"turn"`
	Night             int         `json:"night"`
	TurnsPerNight     int         `json:"turns_per_night"`
	RunNights         int         `json:"run_nights"`
	Credits           int         `json:"credits"`
	Heat              int         `json:"heat"`
	DispatchIntegrity int         `json:"dispatch_integrity"`
	Phase             GamePhase   `json:"phase"`
	Districts         []District  `json:"districts"`
	Runners           []Runner    `json:"runners"`
	Factions          []Faction   `json:"factions"`
	AvailableJobs     []Job       `json:"available_jobs"`
	ActiveJobs        []ActiveJob `json:"active_jobs"`
	Bundles           []Bundle    `json:"bundles"`
	Messages          []Message   `json:"messages"`
	EventLog          []LogEntry  `json:"event_log"`
	RandomSeed        int64       `json:"random_seed"`
}

type InitialContent struct {
	Districts []District
	Runners   []Runner
	Factions  []Faction
	Messages  []Message
	EventLog  []LogEntry
}

func NewGameState(content InitialContent, seed int64) GameState {
	if seed == 0 {
		seed = DefaultRandomSeed
	}

	return GameState{
		Turn:              FirstTurn,
		Night:             FirstNight,
		TurnsPerNight:     DefaultTurnsPerNight,
		RunNights:         DefaultRunNights,
		Credits:           StartingCredits,
		Heat:              StartingHeat,
		DispatchIntegrity: StartingDispatchIntegrity,
		Phase:             PhaseDispatch,
		Districts:         append([]District(nil), content.Districts...),
		Runners:           append([]Runner(nil), content.Runners...),
		Factions:          append([]Faction(nil), content.Factions...),
		Messages:          append([]Message(nil), content.Messages...),
		EventLog:          append([]LogEntry(nil), content.EventLog...),
		RandomSeed:        seed,
	}
}
