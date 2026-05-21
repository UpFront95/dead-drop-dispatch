package game

type DistrictID string
type RunnerID string
type FactionID string
type CargoType string
type RouteType string
type RunnerState string
type GamePhase string

const (
	RunnerReady   RunnerState = "ready"
	RunnerOnJob   RunnerState = "on_job"
	RunnerInjured RunnerState = "injured"
	RunnerBurned  RunnerState = "burned"
	RunnerMissing RunnerState = "missing"
)

const (
	PhaseDispatch GamePhase = "dispatch"
	PhaseResolve  GamePhase = "resolve"
	PhaseGameOver GamePhase = "game_over"
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
	ID            string     `json:"id"`
	TemplateID    string     `json:"template_id"`
	Title         string     `json:"title"`
	ClientMessage string     `json:"client_message"`
	Cargo         CargoType  `json:"cargo"`
	Origin        DistrictID `json:"origin"`
	Destination   DistrictID `json:"destination"`
	DeadlineTurns int        `json:"deadline_turns"`
	Payout        int        `json:"payout"`
	Faction       FactionID  `json:"faction"`
	Modifiers     []string   `json:"modifiers,omitempty"`
	RiskFactors   []string   `json:"risk_factors,omitempty"`
	Routes        []Route    `json:"routes,omitempty"`
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
}

type Bundle struct {
	RunnerID RunnerID    `json:"runner_id"`
	Jobs     []ActiveJob `json:"jobs"`
}

type Message struct {
	Turn    int    `json:"turn"`
	From    string `json:"from"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type LogEntry struct {
	Turn int    `json:"turn"`
	Text string `json:"text"`
}
