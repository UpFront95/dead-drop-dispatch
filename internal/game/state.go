package game

const (
	StartingCredits                        = 800
	StartingHeat                           = 0
	StartingDispatchIntegrity              = 100
	DefaultCargoIntegrity                  = 100
	BankruptcyCreditFloor                  = 0
	VictoryCreditTarget                    = 2500
	MaximumHeat                            = 10
	DispatchIntegrityFailure               = 0
	HostileFactionSuspicion                = 5
	FactionLockoutCount                    = 3
	FirstNight                             = 1
	FirstTurn                              = 1
	DefaultTurnsPerNight                   = 6
	DefaultRunNights                       = 7
	DefaultRandomSeed                int64 = 1701
	MaxJobsPerRunner                       = 2
	MaxRunnerStress                        = 10
	MinRunnerLoyalty                       = 0
	DetectionHeatGain                      = 2
	InterceptionHeatGain                   = 2
	BaseJobStressGain                      = 1
	DelayStressGain                        = 1
	InjuryStressGain                       = 2
	BundleStressGain                       = 1
	FailedJobDispatchIntegrityLoss         = 5
	CargoDamageDispatchIntegrityLoss       = 2
	CargoDamageIntegrityLoss               = 35
	CargoFailureIntegrityLoss              = 100
	NightlyOperatingCost                   = 150
	NightlyHeatDecay                       = 1
)

type GameState struct {
	Turn              int            `json:"turn"`
	Night             int            `json:"night"`
	TurnsPerNight     int            `json:"turns_per_night"`
	RunNights         int            `json:"run_nights"`
	Credits           int            `json:"credits"`
	Heat              int            `json:"heat"`
	DispatchIntegrity int            `json:"dispatch_integrity"`
	Phase             GamePhase      `json:"phase"`
	Districts         []District     `json:"districts"`
	Runners           []Runner       `json:"runners"`
	Factions          []Faction      `json:"factions"`
	JobTemplates      []JobTemplate  `json:"job_templates,omitempty"`
	AvailableJobs     []Job          `json:"available_jobs"`
	AcceptedJobs      []Job          `json:"accepted_jobs"`
	ActiveJobs        []ActiveJob    `json:"active_jobs"`
	Bundles           []Bundle       `json:"bundles"`
	PurchasedUpgrades []UpgradeID    `json:"purchased_upgrades,omitempty"`
	Complications     []Complication `json:"complications"`
	LastResults       []JobResult    `json:"last_results"`
	Messages          []Message      `json:"messages"`
	EventLog          []LogEntry     `json:"event_log"`
	RandomSeed        int64          `json:"random_seed"`
}

type InitialContent struct {
	Districts    []District
	Runners      []Runner
	Factions     []Faction
	JobTemplates []JobTemplate
	Messages     []Message
	EventLog     []LogEntry
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
		JobTemplates:      append([]JobTemplate(nil), content.JobTemplates...),
		Messages:          append([]Message(nil), content.Messages...),
		EventLog:          append([]LogEntry(nil), content.EventLog...),
		RandomSeed:        seed,
	}
}

func EvaluateRunStatus(state GameState) RunStatus {
	if state.Credits < BankruptcyCreditFloor {
		return lostStatus(RunEndBankrupt, "Dispatch is bankrupt.")
	}
	if state.Heat >= MaximumHeat {
		return lostStatus(RunEndBurned, "Heat reached its limit. The desk is burned.")
	}
	if state.DispatchIntegrity <= DispatchIntegrityFailure {
		return lostStatus(RunEndCollapse, "Dispatch integrity collapsed.")
	}
	if allRunnersOutOfRun(state.Runners) {
		return lostStatus(RunEndRosterLoss, "No runners are available for work.")
	}
	if hostileFactionCount(state.Factions) >= FactionLockoutCount {
		return lostStatus(RunEndFactionLockout, "Too many factions are hostile to the desk.")
	}
	if runComplete(state) {
		if state.Credits >= VictoryCreditTarget {
			return RunStatus{
				State:   RunWon,
				Reason:  RunEndVictory,
				Summary: "Seven nights survived with enough credits to keep the operation alive.",
			}
		}
		return lostStatus(RunEndShortfall, "Seven nights survived, but the desk missed its credit target.")
	}
	return RunStatus{
		State:   RunInProgress,
		Reason:  RunEndNone,
		Summary: "Run is still active.",
	}
}

func lostStatus(reason RunEndReason, summary string) RunStatus {
	return RunStatus{
		State:   RunLost,
		Reason:  reason,
		Summary: summary,
	}
}

func runComplete(state GameState) bool {
	return state.Night > state.RunNights || (state.Night == state.RunNights && state.Turn > state.TurnsPerNight)
}

func allRunnersOutOfRun(runners []Runner) bool {
	if len(runners) == 0 {
		return true
	}
	for _, runner := range runners {
		if runner.State == RunnerReady || runner.State == RunnerOnJob {
			return false
		}
	}
	return true
}

func hostileFactionCount(factions []Faction) int {
	count := 0
	for _, faction := range factions {
		if faction.Suspicion >= HostileFactionSuspicion {
			count++
		}
	}
	return count
}
