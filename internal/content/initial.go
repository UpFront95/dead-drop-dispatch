package content

import "dead-drop-dispatch/internal/game"

const (
	FactionHelix    game.FactionID = "helix_municipal_security"
	FactionUnion    game.FactionID = "kestrel_dock_union"
	FactionClinic   game.FactionID = "saint_orison_clinic_network"
	FactionAsterion game.FactionID = "asterion_systems"
)

func InitialGameState(seed int64) game.GameState {
	state := game.NewGameState(game.InitialContent{
		Districts: Districts(),
		Runners:   Runners(),
		Factions:  Factions(),
		Messages: []game.Message{{
			Turn:    game.FirstTurn,
			From:    "switchboard",
			Subject: "line check",
			Body:    "Desk is live. City is listening.",
		}},
		EventLog: []game.LogEntry{{
			Turn: game.FirstTurn,
			Text: "Dispatch console initialized.",
		}},
	}, seed)
	state.AvailableJobs = game.GenerateJobs(state, JobTemplates(), game.DefaultJobsPerTurn)
	return state
}

func Districts() []game.District {
	return []game.District{
		{
			ID:             "northline",
			Name:           "Northline",
			Description:    "Corporate transit spine and checkpoint corridor.",
			Surveillance:   5,
			Traffic:        2,
			FactionControl: FactionHelix,
			Danger:         3,
			SignalQuality:  4,
		},
		{
			ID:             "floodglass",
			Name:           "Floodglass",
			Description:    "Low streets, tunnels, pump stations, and illegal clinics.",
			Surveillance:   2,
			Traffic:        4,
			FactionControl: FactionClinic,
			Danger:         3,
			SignalQuality:  2,
		},
		{
			ID:             "saint_orison_market",
			Name:           "Saint Orison Market",
			Description:    "Dense bazaar, shrines, noodle stalls, and black-market kiosks.",
			Surveillance:   2,
			Traffic:        4,
			FactionControl: FactionClinic,
			Danger:         2,
			SignalQuality:  3,
		},
		{
			ID:             "port_kestrel",
			Name:           "Port Kestrel",
			Description:    "Cargo yards, drone cranes, and container stacks.",
			Surveillance:   3,
			Traffic:        3,
			FactionControl: FactionUnion,
			Danger:         4,
			SignalQuality:  3,
		},
		{
			ID:             "crown_verge",
			Name:           "Crown Verge",
			Description:    "Luxury towers and private security zones.",
			Surveillance:   5,
			Traffic:        2,
			FactionControl: FactionAsterion,
			Danger:         4,
			SignalQuality:  5,
		},
	}
}

func Runners() []game.Runner {
	return []game.Runner{
		{
			ID:       "mira_vale",
			Name:     "Mira Vale",
			Style:    "Bike courier",
			Strength: "Fast through crowded streets.",
			Weakness: "Gains stress under surveillance.",
			Trait:    "Knows every service ramp that still opens.",
			Speed:    5,
			Stealth:  3,
			Nerve:    3,
			Talk:     2,
			Loyalty:  4,
			State:    game.RunnerReady,
		},
		{
			ID:       "kaito_senn",
			Name:     "Kaito Senn",
			Style:    "Tunnel runner",
			Strength: "Safer in Floodglass and service corridors.",
			Weakness: "Slow through corporate districts.",
			Trait:    "Can navigate by pump noise and bad wiring.",
			Speed:    3,
			Stealth:  5,
			Nerve:    4,
			Talk:     2,
			Loyalty:  4,
			State:    game.RunnerReady,
		},
		{
			ID:       "vex_calder",
			Name:     "Vex Calder",
			Style:    "Social operator",
			Strength: "Better at checkpoints and negotiations.",
			Weakness: "Side-deals become likely if loyalty drops.",
			Trait:    "Smiles like they already sold the room.",
			Speed:    3,
			Stealth:  3,
			Nerve:    4,
			Talk:     5,
			Loyalty:  3,
			State:    game.RunnerReady,
		},
	}
}

func Factions() []game.Faction {
	return []game.Faction{
		{
			ID:          FactionHelix,
			Name:        "Helix Municipal Security",
			Description: "Police-adjacent surveillance and enforcement body.",
			Reputation:  0,
			Suspicion:   1,
		},
		{
			ID:          FactionUnion,
			Name:        "Kestrel Dock Union",
			Description: "Smugglers, workers, cargo handlers, and strike captains.",
			Reputation:  1,
			Suspicion:   0,
		},
		{
			ID:          FactionClinic,
			Name:        "Saint Orison Clinic Network",
			Description: "Underground medics and biological cargo brokers.",
			Reputation:  1,
			Suspicion:   0,
		},
		{
			ID:          FactionAsterion,
			Name:        "Asterion Systems",
			Description: "Corporate client, data broker, and security contractor.",
			Reputation:  0,
			Suspicion:   1,
		},
	}
}
