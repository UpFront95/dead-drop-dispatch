package game

import (
	"fmt"
	"math/rand"
)

const DefaultJobsPerTurn = 3

func RefreshAvailableJobs(state *GameState, count int) []Job {
	state.AvailableJobs = GenerateJobs(*state, state.JobTemplates, count)
	state.EventLog = append(state.EventLog, LogEntry{
		Turn: state.Turn,
		Text: fmt.Sprintf("Job board refreshed with %d postings.", len(state.AvailableJobs)),
	})
	return state.AvailableJobs
}

func GenerateJobs(state GameState, templates []JobTemplate, count int) []Job {
	if count <= 0 || len(templates) == 0 || len(state.Districts) < 2 {
		return nil
	}

	rng := rand.New(rand.NewSource(jobSeed(state)))
	jobs := make([]Job, 0, count)
	for i := 0; i < count; i++ {
		template := templates[(rng.Intn(len(templates))+i)%len(templates)]
		origin, destination := pickDistrictPair(rng, state.Districts)
		job := buildJob(state, template, origin, destination, i)
		jobs = append(jobs, job)
	}
	return jobs
}

func jobSeed(state GameState) int64 {
	return state.RandomSeed + int64(state.Night*1000) + int64(state.Turn*37)
}

func pickDistrictPair(rng *rand.Rand, districts []District) (District, District) {
	originIndex := rng.Intn(len(districts))
	destinationIndex := rng.Intn(len(districts) - 1)
	if destinationIndex >= originIndex {
		destinationIndex++
	}
	return districts[originIndex], districts[destinationIndex]
}

func buildJob(state GameState, template JobTemplate, origin District, destination District, index int) Job {
	jobID := fmt.Sprintf("n%02dt%02d-%02d", state.Night, state.Turn, index+1)
	deadline := template.BaseDeadline
	if deadline <= 0 {
		deadline = 2
	}

	pressure := origin.Traffic + destination.Surveillance + destination.Danger
	payout := template.BasePayout + pressure*18
	if deadline <= 1 {
		payout += 90
	}
	riskFactors := jobRiskFactors(template, origin, destination)

	return Job{
		ID:            jobID,
		TemplateID:    template.ID,
		Title:         template.Title,
		ClientMessage: template.ClientMessage,
		Cargo:         template.Cargo,
		Origin:        origin.ID,
		Destination:   destination.ID,
		DeadlineTurns: deadline,
		Payout:        payout,
		Faction:       template.Faction,
		Modifiers:     append([]string(nil), template.Modifiers...),
		RiskFactors:   riskFactors,
		Intel:         []IntelReport{JobIntelReport(state, riskFactors)},
		Routes:        generateRoutes(jobID, template, origin, destination),
	}
}

func jobRiskFactors(template JobTemplate, origin District, destination District) []string {
	factors := append([]string(nil), template.RiskFactors...)
	if destination.Surveillance >= 4 {
		factors = append(factors, "destination surveillance")
	}
	if origin.Traffic >= 4 || destination.Traffic >= 4 {
		factors = append(factors, "traffic pressure")
	}
	if origin.Danger >= 4 || destination.Danger >= 4 {
		factors = append(factors, "violent district")
	}
	if origin.SignalQuality <= 2 || destination.SignalQuality <= 2 {
		factors = append(factors, "weak signal")
	}
	return uniqueStrings(factors)
}

func generateRoutes(jobID string, template JobTemplate, origin District, destination District) []Route {
	routeTypes := template.PreferredRoutes
	if len(routeTypes) == 0 {
		routeTypes = []RouteType{RouteMainArtery, RouteServiceTunnels, RouteMarketWeave}
	}
	if len(routeTypes) > 4 {
		routeTypes = routeTypes[:4]
	}

	routes := make([]Route, 0, len(routeTypes))
	for i, routeType := range routeTypes {
		routes = append(routes, Route{
			ID:         fmt.Sprintf("%s-r%d", jobID, i+1),
			Type:       routeType,
			Name:       routeName(routeType),
			Districts:  []DistrictID{origin.ID, destination.ID},
			Traits:     routeTraits(routeType, origin, destination),
			TimeCost:   routeTimeCost(routeType, origin, destination),
			StressHint: routeStressHint(routeType),
		})
	}
	return routes
}

func routeName(routeType RouteType) string {
	switch routeType {
	case RouteMainArtery:
		return "Main artery"
	case RouteServiceTunnels:
		return "Service tunnels"
	case RouteMarketWeave:
		return "Market weave"
	case RouteDroneCorridor:
		return "Drone corridor"
	case RouteFloodline:
		return "Floodline"
	default:
		return "Unknown route"
	}
}

func routeTraits(routeType RouteType, origin District, destination District) []string {
	traits := []string{}
	switch routeType {
	case RouteMainArtery:
		traits = append(traits, "fast", "visible")
	case RouteServiceTunnels:
		traits = append(traits, "concealed", "slow")
	case RouteMarketWeave:
		traits = append(traits, "crowded", "talk helps")
	case RouteDroneCorridor:
		traits = append(traits, "watched", "precise timing")
	case RouteFloodline:
		traits = append(traits, "unstable", "signal drops")
	}
	if destination.Surveillance >= 4 {
		traits = append(traits, "watched destination")
	}
	if origin.Traffic >= 4 || destination.Traffic >= 4 {
		traits = append(traits, "traffic cover")
	}
	return uniqueStrings(traits)
}

func routeTimeCost(routeType RouteType, origin District, destination District) int {
	cost := 1
	if origin.Traffic >= 4 || destination.Traffic >= 4 {
		cost++
	}
	switch routeType {
	case RouteServiceTunnels, RouteFloodline:
		cost++
	case RouteMainArtery:
		if cost > 1 {
			cost--
		}
	}
	return cost
}

func routeStressHint(routeType RouteType) string {
	switch routeType {
	case RouteMainArtery:
		return "public exposure"
	case RouteServiceTunnels:
		return "close quarters"
	case RouteMarketWeave:
		return "crowd noise"
	case RouteDroneCorridor:
		return "scanner timing"
	case RouteFloodline:
		return "bad footing"
	default:
		return "unknown pressure"
	}
}

func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	return result
}
