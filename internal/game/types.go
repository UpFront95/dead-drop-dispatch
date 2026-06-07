package game

type DistrictID string
type RunnerID string
type FactionID string
type CargoType string
type RouteType string
type RunnerState string
type GamePhase string
type JobOutcome string
type ComplicationID string
type ComplicationChoiceID string
type ComplicationType string
type ComplicationStatus string
type IntelSource string
type IntelConfidence string
type IntelRiskTag string
type UpgradeID string
type MessageID string
type MessageAudience string
type MessageResponseActionID string
type MessageStatus string
type RunState string
type RunEndReason string

const (
	RunnerReady   RunnerState = "ready"
	RunnerOnJob   RunnerState = "on_job"
	RunnerInjured RunnerState = "injured"
	RunnerBurned  RunnerState = "burned"
	RunnerMissing RunnerState = "missing"
)

const (
	PhaseMessages      GamePhase = "messages"
	PhaseJobs          GamePhase = "jobs"
	PhaseDispatch      GamePhase = "dispatch"
	PhaseResolve       GamePhase = "resolve"
	PhaseComplications GamePhase = "complications"
	PhaseReports       GamePhase = "reports"
	PhaseCityUpdate    GamePhase = "city_update"
	PhaseGameOver      GamePhase = "game_over"
)

const (
	RunInProgress RunState = "in_progress"
	RunWon        RunState = "won"
	RunLost       RunState = "lost"
)

const (
	RunEndNone           RunEndReason = "none"
	RunEndVictory        RunEndReason = "victory"
	RunEndShortfall      RunEndReason = "shortfall"
	RunEndBankrupt       RunEndReason = "bankrupt"
	RunEndBurned         RunEndReason = "burned"
	RunEndCollapse       RunEndReason = "collapse"
	RunEndRosterLoss     RunEndReason = "roster_loss"
	RunEndFactionLockout RunEndReason = "faction_lockout"
)

const (
	OutcomeSuccess     JobOutcome = "success"
	OutcomePartial     JobOutcome = "partial"
	OutcomeFailed      JobOutcome = "failed"
	OutcomeIntercepted JobOutcome = "intercepted"
)

const (
	ComplicationNone           ComplicationType = ""
	ComplicationCheckpoint     ComplicationType = "checkpoint"
	ComplicationScannerSweep   ComplicationType = "scanner_sweep"
	ComplicationRunnerPanic    ComplicationType = "runner_panic"
	ComplicationCargoLeak      ComplicationType = "cargo_leak"
	ComplicationSignalLoss     ComplicationType = "signal_loss"
	ComplicationGangToll       ComplicationType = "gang_toll"
	ComplicationClientTerms    ComplicationType = "client_terms"
	ComplicationDroneTail      ComplicationType = "drone_tail"
	ComplicationWitnessRefuses ComplicationType = "witness_refuses"
	ComplicationDataTrace      ComplicationType = "data_trace"
	ComplicationCurfewDrop     ComplicationType = "curfew_drop"
	ComplicationRivalCourier   ComplicationType = "rival_courier"
)

const (
	ComplicationPending  ComplicationStatus = "pending"
	ComplicationResolved ComplicationStatus = "resolved"
)

const (
	IntelSourceClient  IntelSource = "client"
	IntelSourceRunner  IntelSource = "runner"
	IntelSourceFaction IntelSource = "faction"
	IntelSourceSensor  IntelSource = "sensor"
	IntelSourceStreet  IntelSource = "street"
	IntelSourceArchive IntelSource = "archive"
)

const (
	IntelConfidenceLow    IntelConfidence = "low"
	IntelConfidenceMedium IntelConfidence = "medium"
	IntelConfidenceHigh   IntelConfidence = "high"
)

const (
	IntelRiskFalse      IntelRiskTag = "false_possible"
	IntelRiskIncomplete IntelRiskTag = "incomplete"
	IntelRiskStale      IntelRiskTag = "stale"
	IntelRiskBiased     IntelRiskTag = "biased"
)

const (
	UpgradeSignalRelay           UpgradeID = "signal_relay"
	UpgradeSafehouse             UpgradeID = "safehouse"
	UpgradeFakeCredentialPrinter UpgradeID = "fake_credential_printer"
	UpgradeClinicFavor           UpgradeID = "clinic_favor"
	UpgradeDeadDropLocker        UpgradeID = "dead_drop_locker"
	UpgradeScrambler             UpgradeID = "scrambler"
)

const (
	MessageAudienceClient  MessageAudience = "client"
	MessageAudienceRunner  MessageAudience = "runner"
	MessageAudienceFaction MessageAudience = "faction"
)

const (
	MessageOpen     MessageStatus = "open"
	MessageResolved MessageStatus = "resolved"
)

const (
	ResponseRefuse      MessageResponseActionID = "refuse"
	ResponseAskMorePay  MessageResponseActionID = "ask_more_pay"
	ResponseAskMoreInfo MessageResponseActionID = "ask_more_info"
	ResponseThreaten    MessageResponseActionID = "threaten"
	ResponseReassure    MessageResponseActionID = "reassure"
	ResponseStall       MessageResponseActionID = "stall"
	ResponseDeceive     MessageResponseActionID = "deceive"
	ResponseCancel      MessageResponseActionID = "cancel"
	ResponseAccept      MessageResponseActionID = "accept"
)

const (
	ChoiceTalkThrough   ComplicationChoiceID = "talk_through"
	ChoiceReroute       ComplicationChoiceID = "reroute"
	ChoiceBribe         ComplicationChoiceID = "bribe"
	ChoiceAbandon       ComplicationChoiceID = "abandon"
	ChoiceHideCargo     ComplicationChoiceID = "hide_cargo"
	ChoiceRushThrough   ComplicationChoiceID = "rush_through"
	ChoiceSpoofTag      ComplicationChoiceID = "spoof_tag"
	ChoiceReassure      ComplicationChoiceID = "reassure"
	ChoiceOrderForward  ComplicationChoiceID = "order_forward"
	ChoiceAbort         ComplicationChoiceID = "abort"
	ChoiceContinue      ComplicationChoiceID = "continue"
	ChoiceSeekClinic    ComplicationChoiceID = "seek_clinic"
	ChoiceDumpCargo     ComplicationChoiceID = "dump_cargo"
	ChoiceTrustRoute    ComplicationChoiceID = "trust_route"
	ChoiceWait          ComplicationChoiceID = "wait"
	ChoicePayToll       ComplicationChoiceID = "pay_toll"
	ChoiceThreaten      ComplicationChoiceID = "threaten"
	ChoiceCallFavor     ComplicationChoiceID = "call_favor"
	ChoiceAcceptTerms   ComplicationChoiceID = "accept_terms"
	ChoiceRenegotiate   ComplicationChoiceID = "renegotiate"
	ChoiceRefuseTerms   ComplicationChoiceID = "refuse_terms"
	ChoiceLoseTail      ComplicationChoiceID = "lose_tail"
	ChoiceJamSignal     ComplicationChoiceID = "jam_signal"
	ChoiceShelter       ComplicationChoiceID = "shelter"
	ChoiceCoaxWitness   ComplicationChoiceID = "coax_witness"
	ChoiceSedateWitness ComplicationChoiceID = "sedate_witness"
	ChoiceScrubTrace    ComplicationChoiceID = "scrub_trace"
	ChoiceBurnNode      ComplicationChoiceID = "burn_node"
	ChoiceDecoyPacket   ComplicationChoiceID = "decoy_packet"
	ChoiceUseSafehouse  ComplicationChoiceID = "use_safehouse"
	ChoiceBreakCurfew   ComplicationChoiceID = "break_curfew"
	ChoiceRaceCourier   ComplicationChoiceID = "race_courier"
	ChoiceBlockCourier  ComplicationChoiceID = "block_courier"
	ChoiceShareRoute    ComplicationChoiceID = "share_route"
)

const (
	CargoDataShard          CargoType = "data_shard"
	CargoMedicalCooler      CargoType = "medical_cooler"
	CargoWitness            CargoType = "witness"
	CargoContrabandPackage  CargoType = "contraband_package"
	CargoCorporatePrototype CargoType = "corporate_prototype"
)

const (
	RouteMainArtery     RouteType = "main_artery"
	RouteServiceTunnels RouteType = "service_tunnels"
	RouteMarketWeave    RouteType = "market_weave"
	RouteDroneCorridor  RouteType = "drone_corridor"
	RouteFloodline      RouteType = "floodline"
)

type District struct {
	ID             DistrictID `json:"id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Surveillance   int        `json:"surveillance"`
	Traffic        int        `json:"traffic"`
	FactionControl FactionID  `json:"faction_control"`
	Danger         int        `json:"danger"`
	SignalQuality  int        `json:"signal_quality"`
}

type Runner struct {
	ID          RunnerID     `json:"id"`
	Name        string       `json:"name"`
	Style       string       `json:"style"`
	Strength    string       `json:"strength"`
	Weakness    string       `json:"weakness"`
	Trait       string       `json:"trait"`
	Speed       int          `json:"speed"`
	Stealth     int          `json:"stealth"`
	Nerve       int          `json:"nerve"`
	Talk        int          `json:"talk"`
	Loyalty     int          `json:"loyalty"`
	Stress      int          `json:"stress"`
	State       RunnerState  `json:"state"`
	Recovery    int          `json:"recovery"`
	BurnedZones []DistrictID `json:"burned_zones,omitempty"`
}

type Faction struct {
	ID          FactionID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Reputation  int       `json:"reputation"`
	Suspicion   int       `json:"suspicion"`
}

type Job struct {
	ID            string        `json:"id"`
	TemplateID    string        `json:"template_id"`
	Title         string        `json:"title"`
	ClientMessage string        `json:"client_message"`
	Cargo         CargoType     `json:"cargo"`
	Origin        DistrictID    `json:"origin"`
	Destination   DistrictID    `json:"destination"`
	DeadlineTurns int           `json:"deadline_turns"`
	Payout        int           `json:"payout"`
	Faction       FactionID     `json:"faction"`
	Modifiers     []string      `json:"modifiers,omitempty"`
	RiskFactors   []string      `json:"risk_factors,omitempty"`
	Intel         []IntelReport `json:"intel,omitempty"`
	Routes        []Route       `json:"routes,omitempty"`
}

type JobTemplate struct {
	ID              string      `json:"id"`
	Title           string      `json:"title"`
	ClientMessage   string      `json:"client_message"`
	Cargo           CargoType   `json:"cargo"`
	Faction         FactionID   `json:"faction"`
	BasePayout      int         `json:"base_payout"`
	BaseDeadline    int         `json:"base_deadline"`
	Modifiers       []string    `json:"modifiers,omitempty"`
	RiskFactors     []string    `json:"risk_factors,omitempty"`
	PreferredRoutes []RouteType `json:"preferred_routes,omitempty"`
}

type Route struct {
	ID         string       `json:"id"`
	Type       RouteType    `json:"type"`
	Name       string       `json:"name"`
	Districts  []DistrictID `json:"districts"`
	Traits     []string     `json:"traits,omitempty"`
	TimeCost   int          `json:"time_cost"`
	StressHint string       `json:"stress_hint,omitempty"`
}

type ActiveJob struct {
	JobID    string   `json:"job_id"`
	RunnerID RunnerID `json:"runner_id"`
	RouteID  string   `json:"route_id"`
	Job      Job      `json:"job"`
	Route    Route    `json:"route"`
}

type Bundle struct {
	RunnerID  RunnerID    `json:"runner_id"`
	Jobs      []ActiveJob `json:"jobs"`
	Penalties []string    `json:"penalties,omitempty"`
}

type UpgradeDefinition struct {
	ID          UpgradeID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Cost        int       `json:"cost"`
	Effects     []string  `json:"effects,omitempty"`
}

type JobResult struct {
	JobID                 string           `json:"job_id"`
	JobTitle              string           `json:"job_title"`
	RunnerID              RunnerID         `json:"runner_id"`
	RunnerName            string           `json:"runner_name"`
	FactionID             FactionID        `json:"faction_id,omitempty"`
	Outcome               JobOutcome       `json:"outcome"`
	Payout                int              `json:"payout"`
	HeatGain              int              `json:"heat_gain"`
	StressGain            int              `json:"stress_gain"`
	CargoIntegrity        int              `json:"cargo_integrity"`
	CargoIntegrityLoss    int              `json:"cargo_integrity_loss"`
	DispatchIntegrityLoss int              `json:"dispatch_integrity_loss"`
	Delay                 bool             `json:"delay"`
	Detection             bool             `json:"detection"`
	Complication          bool             `json:"complication"`
	ComplicationType      ComplicationType `json:"complication_type,omitempty"`
	Injury                bool             `json:"injury"`
	CargoDamage           bool             `json:"cargo_damage"`
	Interception          bool             `json:"interception"`
	Factors               []string         `json:"factors,omitempty"`
	Summary               string           `json:"summary"`
}

type Complication struct {
	ID                 ComplicationID       `json:"id"`
	Type               ComplicationType     `json:"type"`
	Status             ComplicationStatus   `json:"status"`
	Title              string               `json:"title"`
	Prompt             string               `json:"prompt"`
	Choices            []ComplicationChoice `json:"choices,omitempty"`
	Turn               int                  `json:"turn"`
	Night              int                  `json:"night"`
	JobID              string               `json:"job_id"`
	JobTitle           string               `json:"job_title"`
	RunnerID           RunnerID             `json:"runner_id"`
	RunnerName         string               `json:"runner_name"`
	FactionID          FactionID            `json:"faction_id,omitempty"`
	Outcome            JobOutcome           `json:"outcome"`
	CargoIntegrity     int                  `json:"cargo_integrity"`
	CargoIntegrityLoss int                  `json:"cargo_integrity_loss"`
	DelayTurns         int                  `json:"delay_turns"`
	ResolvedBy         ComplicationChoiceID `json:"resolved_by,omitempty"`
	ResolutionEffects  []string             `json:"resolution_effects,omitempty"`
	Factors            []string             `json:"factors,omitempty"`
	Summary            string               `json:"summary"`
}

type ComplicationDefinition struct {
	Type        ComplicationType     `json:"type"`
	Title       string               `json:"title"`
	Prompt      string               `json:"prompt"`
	RiskTag     string               `json:"risk_tag"`
	Description string               `json:"description"`
	Choices     []ComplicationChoice `json:"choices"`
}

type ComplicationChoice struct {
	ID          ComplicationChoiceID `json:"id"`
	Label       string               `json:"label"`
	Description string               `json:"description"`
}

type ComplicationResolution struct {
	ComplicationID         ComplicationID       `json:"complication_id"`
	ChoiceID               ComplicationChoiceID `json:"choice_id"`
	Summary                string               `json:"summary"`
	CreditsDelta           int                  `json:"credits_delta,omitempty"`
	HeatDelta              int                  `json:"heat_delta,omitempty"`
	RunnerStressDelta      int                  `json:"runner_stress_delta,omitempty"`
	RunnerLoyaltyDelta     int                  `json:"runner_loyalty_delta,omitempty"`
	CargoIntegrityDelta    int                  `json:"cargo_integrity_delta,omitempty"`
	DispatchIntegrityDelta int                  `json:"dispatch_integrity_delta,omitempty"`
	DelayTurns             int                  `json:"delay_turns,omitempty"`
	FactionReputationDelta int                  `json:"faction_reputation_delta,omitempty"`
	FactionSuspicionDelta  int                  `json:"faction_suspicion_delta,omitempty"`
	Effects                []string             `json:"effects,omitempty"`
}

type RunStatus struct {
	State   RunState     `json:"state"`
	Reason  RunEndReason `json:"reason"`
	Summary string       `json:"summary"`
}

type IntelReport struct {
	Source      IntelSource     `json:"source"`
	Night       int             `json:"night"`
	Turn        int             `json:"turn"`
	Staleness   int             `json:"staleness"`
	Confidence  IntelConfidence `json:"confidence"`
	Claims      []string        `json:"claims,omitempty"`
	OmittedTags []string        `json:"omitted_tags,omitempty"`
	RiskTags    []IntelRiskTag  `json:"risk_tags,omitempty"`
}

type Message struct {
	ID                MessageID               `json:"id,omitempty"`
	Turn              int                     `json:"turn"`
	From              string                  `json:"from"`
	Subject           string                  `json:"subject"`
	Body              string                  `json:"body"`
	Audience          MessageAudience         `json:"audience,omitempty"`
	Status            MessageStatus           `json:"status,omitempty"`
	Responses         []MessageResponseAction `json:"responses,omitempty"`
	TargetRunnerID    RunnerID                `json:"target_runner_id,omitempty"`
	TargetFactionID   FactionID               `json:"target_faction_id,omitempty"`
	ResolvedBy        MessageResponseActionID `json:"resolved_by,omitempty"`
	ResolutionEffects []string                `json:"resolution_effects,omitempty"`
	Summary           string                  `json:"summary,omitempty"`
}

type MessageResponseAction struct {
	ID          MessageResponseActionID `json:"id"`
	Label       string                  `json:"label"`
	Description string                  `json:"description"`
	Audiences   []MessageAudience       `json:"audiences"`
}

type MessageResponseResolution struct {
	MessageID              MessageID               `json:"message_id"`
	ActionID               MessageResponseActionID `json:"action_id"`
	Summary                string                  `json:"summary"`
	CreditsDelta           int                     `json:"credits_delta,omitempty"`
	HeatDelta              int                     `json:"heat_delta,omitempty"`
	DispatchIntegrityDelta int                     `json:"dispatch_integrity_delta,omitempty"`
	RunnerStressDelta      int                     `json:"runner_stress_delta,omitempty"`
	RunnerLoyaltyDelta     int                     `json:"runner_loyalty_delta,omitempty"`
	FactionReputationDelta int                     `json:"faction_reputation_delta,omitempty"`
	FactionSuspicionDelta  int                     `json:"faction_suspicion_delta,omitempty"`
	Effects                []string                `json:"effects,omitempty"`
}

type LogEntry struct {
	Turn int    `json:"turn"`
	Text string `json:"text"`
}
