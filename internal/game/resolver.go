package game

import (
	"fmt"
	"math/rand"
	"strings"
)

func ResolveActiveJobs(state *GameState) []JobResult {
	if len(state.ActiveJobs) == 0 {
		state.LastResults = nil
		return nil
	}

	rng := rand.New(rand.NewSource(resolveSeed(*state)))
	bundleLoads := bundleLoadMap(state.Bundles)
	results := make([]JobResult, 0, len(state.ActiveJobs))
	for _, active := range state.ActiveJobs {
		result := resolveActiveJob(state, active, bundleLoads[active.RunnerID], rng)
		results = append(results, result)
		applyJobResult(state, result)
	}

	resetActiveWork(state)
	state.LastResults = results
	appendResolutionMessage(state, results)
	return results
}

func resolveSeed(state GameState) int64 {
	return state.RandomSeed + int64(state.Night*10000) + int64(state.Turn*271)
}

func resolveActiveJob(state *GameState, active ActiveJob, bundleLoad int, rng *rand.Rand) JobResult {
	job := active.Job
	route := active.Route
	runnerIndex := findRunnerIndex(state.Runners, active.RunnerID)
	runner := Runner{ID: active.RunnerID, Name: string(active.RunnerID)}
	if runnerIndex >= 0 {
		runner = state.Runners[runnerIndex]
	}

	detectionRisk := detectionRisk(*state, job, route, runner, bundleLoad)
	delayRisk := delayRisk(*state, job, route, runner, bundleLoad)
	injuryRisk := injuryRisk(*state, job, route, runner, bundleLoad)
	damageRisk := cargoDamageRisk(job, route, bundleLoad)
	interceptionRisk := interceptionRisk(*state, job, route, runner, bundleLoad)

	detection := rollRisk(rng, detectionRisk)
	delay := rollRisk(rng, delayRisk)
	injury := rollRisk(rng, injuryRisk)
	cargoDamage := rollRisk(rng, damageRisk)
	interception := rollRisk(rng, interceptionRisk)
	complicationChance := complicationRisk(*state, job, route, runner, bundleLoad, detection, delay)
	complication := rollRisk(rng, complicationChance)
	complicationType := ComplicationNone

	outcome := OutcomeSuccess
	payout := job.Payout
	if delay || cargoDamage {
		outcome = OutcomePartial
		payout = job.Payout / 2
	}
	if injury && runner.Nerve <= 2 {
		outcome = OutcomeFailed
		payout = 0
	}
	if interception {
		outcome = OutcomeIntercepted
		payout = 0
		detection = true
	}
	if complication {
		complicationType = chooseComplicationType(job, route, runner, detection, delay, cargoDamage)
	}

	heatGain := resultHeatGain(detection, interception)
	stressGain := resultStressGain(delay, injury, bundleLoad)
	cargoIntegrity := resultCargoIntegrity(outcome, cargoDamage)
	cargoIntegrityLoss := DefaultCargoIntegrity - cargoIntegrity
	dispatchIntegrityLoss := resultDispatchIntegrityLoss(outcome, cargoDamage)

	factors := resultFactors(job, route, detection, delay, complication, complicationType, injury, cargoDamage, interception, bundleLoad)
	injuryDetail := injuryDetailForResult(*state, job, route, outcome, injury, detection, complicationType)
	result := JobResult{
		JobID:                 job.ID,
		JobTitle:              job.Title,
		RunnerID:              active.RunnerID,
		RunnerName:            runner.Name,
		FactionID:             job.Faction,
		Outcome:               outcome,
		Payout:                payout,
		HeatGain:              heatGain,
		StressGain:            stressGain,
		CargoIntegrity:        cargoIntegrity,
		CargoIntegrityLoss:    cargoIntegrityLoss,
		DispatchIntegrityLoss: dispatchIntegrityLoss,
		Delay:                 delay,
		Detection:             detection,
		Complication:          complication,
		ComplicationType:      complicationType,
		Injury:                injury,
		InjuryDetail:          injuryDetail,
		CargoDamage:           cargoDamage,
		Interception:          interception,
		Factors:               factors,
	}
	result.Summary = resultSummary(result)
	return result
}

func applyJobResult(state *GameState, result JobResult) {
	applyResultCredits(state, result)
	applyResultHeat(state, result)
	applyResultRunner(state, result)
	applyResultFaction(state, result)
	applyResultDispatchIntegrity(state, result)
	applyResultComplication(state, result)
	appendResultLog(state, result)
}

func applyResultCredits(state *GameState, result JobResult) {
	if result.Payout > 0 {
		_, _ = ApplyPayout(state, result.Payout, fmt.Sprintf("%s completed", result.JobTitle))
	}
}

func applyResultHeat(state *GameState, result JobResult) {
	state.Heat = clampInt(state.Heat+result.HeatGain, StartingHeat, MaximumHeat)
}

func applyResultRunner(state *GameState, result JobResult) {
	if runnerIndex := findRunnerIndex(state.Runners, result.RunnerID); runnerIndex >= 0 {
		state.Runners[runnerIndex].Stress = clampInt(state.Runners[runnerIndex].Stress+result.StressGain, 0, MaxRunnerStress)
		if result.Injury {
			state.Runners[runnerIndex].State = RunnerInjured
			state.Runners[runnerIndex].Recovery = injuryRecoveryTurns(*state)
		} else if state.Runners[runnerIndex].State == RunnerOnJob {
			state.Runners[runnerIndex].State = RunnerReady
		}
		if result.Outcome == OutcomeFailed || result.Outcome == OutcomeIntercepted {
			state.Runners[runnerIndex].Loyalty--
			state.Runners[runnerIndex].Loyalty = maxInt(state.Runners[runnerIndex].Loyalty, MinRunnerLoyalty)
		}
	}
}

func applyResultFaction(state *GameState, result JobResult) {
	if factionIndex := findFactionIndex(state.Factions, resultJobFaction(state.ActiveJobs, result.JobID)); factionIndex >= 0 {
		switch result.Outcome {
		case OutcomeSuccess:
			state.Factions[factionIndex].Reputation++
		case OutcomeFailed, OutcomeIntercepted:
			state.Factions[factionIndex].Suspicion++
		}
		if result.Detection {
			state.Factions[factionIndex].Suspicion++
		}
	}
}

func applyResultDispatchIntegrity(state *GameState, result JobResult) {
	state.DispatchIntegrity = maxInt(state.DispatchIntegrity-result.DispatchIntegrityLoss, DispatchIntegrityFailure)
}

func applyResultComplication(state *GameState, result JobResult) {
	if complication, ok := QueueComplicationFromResult(state, result); ok {
		state.EventLog = append(state.EventLog, LogEntry{
			Turn: state.Turn,
			Text: complication.Summary,
		})
		appendComplicationOpenedMessage(state, complication)
	}
}

func appendResultLog(state *GameState, result JobResult) {
	state.EventLog = append(state.EventLog, LogEntry{
		Turn: state.Turn,
		Text: result.Summary,
	})
}

func resetActiveWork(state *GameState) {
	for i := range state.Runners {
		if state.Runners[i].State == RunnerOnJob {
			state.Runners[i].State = RunnerReady
		}
	}
	state.ActiveJobs = nil
	state.Bundles = nil
}

func appendResolutionMessage(state *GameState, results []JobResult) {
	if len(results) == 0 {
		return
	}
	lines := make([]string, 0, len(results))
	for _, result := range results {
		lines = append(lines, result.Summary)
	}
	state.LastResults = results
	state.Messages = append(state.Messages, Message{
		Turn:    state.Turn,
		From:    "after-action",
		Subject: "run report",
		Body:    strings.Join(lines, " / "),
	})
}

func detectionRisk(state GameState, job Job, route Route, runner Runner, bundleLoad int) int {
	risk := 25 + routeModifier(route.Type, map[RouteType]int{
		RouteMainArtery:     18,
		RouteDroneCorridor:  16,
		RouteMarketWeave:    5,
		RouteFloodline:      -4,
		RouteServiceTunnels: -8,
	})
	risk += districtSurveillance(state, job.Destination) * 7
	risk += factionSuspicion(state, job.Faction) * 5
	risk += cargoModifier(job.Cargo, map[CargoType]int{
		CargoDataShard:          8,
		CargoContrabandPackage:  12,
		CargoCorporatePrototype: 16,
	})
	risk -= signalRelayRiskReduction(state, job, route)
	risk -= contrabandRiskReduction(state, job)
	risk -= dataTraceRiskReduction(state, job)
	risk -= runner.Stealth * 6
	risk -= runner.Talk * 2
	risk += (bundleLoad - 1) * 10
	return clampRisk(risk)
}

func delayRisk(state GameState, job Job, route Route, runner Runner, bundleLoad int) int {
	risk := 15 + route.TimeCost*12 - runner.Speed*5
	if route.TimeCost > job.DeadlineTurns {
		risk += 18
	}
	if route.Type == RouteFloodline || route.Type == RouteServiceTunnels {
		risk += 8
	}
	risk -= signalRelayRiskReduction(state, job, route)
	risk += (bundleLoad - 1) * 12
	return clampRisk(risk)
}

func injuryRisk(state GameState, job Job, route Route, runner Runner, bundleLoad int) int {
	risk := 8 + districtDanger(state, job.Destination)*5 - runner.Nerve*4
	if route.Type == RouteServiceTunnels || route.Type == RouteFloodline {
		risk += 7
	}
	risk += runner.Stress * 3
	risk += (bundleLoad - 1) * 5
	return clampRisk(risk)
}

func cargoDamageRisk(job Job, route Route, bundleLoad int) int {
	risk := cargoModifier(job.Cargo, map[CargoType]int{
		CargoMedicalCooler:      18,
		CargoWitness:            10,
		CargoCorporatePrototype: 12,
	})
	if route.Type == RouteFloodline {
		risk += 10
	}
	risk += (bundleLoad - 1) * 8
	return clampRisk(risk)
}

func interceptionRisk(state GameState, job Job, route Route, runner Runner, bundleLoad int) int {
	risk := cargoModifier(job.Cargo, map[CargoType]int{
		CargoCorporatePrototype: 12,
		CargoContrabandPackage:  7,
	})
	if containsString(job.RiskFactors, "betrayal risk") || containsString(job.RiskFactors, "intercept interest") {
		risk += 8
	}
	if route.Type == RouteDroneCorridor {
		risk += 4
	}
	risk -= contrabandRiskReduction(state, job)
	risk -= runner.Loyalty * 2
	risk += (bundleLoad - 1) * 4
	return clampRisk(risk)
}

func complicationRisk(state GameState, job Job, route Route, runner Runner, bundleLoad int, detection bool, delay bool) int {
	risk := 8 + routeModifier(route.Type, map[RouteType]int{
		RouteMainArtery:     7,
		RouteDroneCorridor:  9,
		RouteMarketWeave:    5,
		RouteFloodline:      6,
		RouteServiceTunnels: 4,
	})
	risk += districtDanger(state, job.Destination) * 3
	risk += cargoModifier(job.Cargo, map[CargoType]int{
		CargoMedicalCooler:      8,
		CargoWitness:            10,
		CargoContrabandPackage:  7,
		CargoCorporatePrototype: 7,
	})
	risk -= signalRelayRiskReduction(state, job, route)
	risk -= contrabandRiskReduction(state, job)
	risk -= dataTraceRiskReduction(state, job)
	risk -= runner.Nerve * 3
	if detection {
		risk += 10
	}
	if delay {
		risk += 7
	}
	risk += (bundleLoad - 1) * 6
	return clampRisk(risk)
}

func chooseComplicationType(job Job, route Route, runner Runner, detection bool, delay bool, cargoDamage bool) ComplicationType {
	if cargoDamage && job.Cargo == CargoMedicalCooler {
		return ComplicationCargoLeak
	}
	if detection && job.Cargo == CargoDataShard {
		return ComplicationDataTrace
	}
	if detection && (route.Type == RouteMainArtery || route.Type == RouteDroneCorridor) {
		return ComplicationScannerSweep
	}
	if runner.Stress >= 6 || job.Cargo == CargoWitness {
		return ComplicationRunnerPanic
	}
	if delay || route.Type == RouteFloodline || route.Type == RouteServiceTunnels {
		return ComplicationSignalLoss
	}
	return ComplicationCheckpoint
}

func resultFactors(job Job, route Route, detection bool, delay bool, complication bool, complicationType ComplicationType, injury bool, cargoDamage bool, interception bool, bundleLoad int) []string {
	factors := append([]string(nil), job.RiskFactors...)
	factors = append(factors, route.Traits...)
	if bundleLoad > 1 {
		factors = append(factors, "bundle pressure")
	}
	if detection {
		factors = append(factors, "detected")
	}
	if delay {
		factors = append(factors, "delayed")
	}
	if complication {
		factors = append(factors, "complication: "+string(complicationType))
	}
	if injury {
		factors = append(factors, "runner hurt")
	}
	if cargoDamage {
		factors = append(factors, "cargo damage")
	}
	if interception {
		factors = append(factors, "interception")
	}
	return uniqueStrings(factors)
}

func resultHeatGain(detection bool, interception bool) int {
	heatGain := 0
	if detection {
		heatGain += DetectionHeatGain
	}
	if interception {
		heatGain += InterceptionHeatGain
	}
	return heatGain
}

func resultStressGain(delay bool, injury bool, bundleLoad int) int {
	stressGain := BaseJobStressGain
	if delay {
		stressGain += DelayStressGain
	}
	if injury {
		stressGain += InjuryStressGain
	}
	if bundleLoad > 1 {
		stressGain += BundleStressGain
	}
	return stressGain
}

func resultCargoIntegrity(outcome JobOutcome, cargoDamage bool) int {
	switch outcome {
	case OutcomeFailed, OutcomeIntercepted:
		return DefaultCargoIntegrity - CargoFailureIntegrityLoss
	default:
		if cargoDamage {
			return DefaultCargoIntegrity - CargoDamageIntegrityLoss
		}
		return DefaultCargoIntegrity
	}
}

func resultDispatchIntegrityLoss(outcome JobOutcome, cargoDamage bool) int {
	loss := 0
	if outcome == OutcomeFailed || outcome == OutcomeIntercepted {
		loss += FailedJobDispatchIntegrityLoss
	}
	if cargoDamage {
		loss += CargoDamageDispatchIntegrityLoss
	}
	return loss
}

func resultSummary(result JobResult) string {
	var summary string
	causes := resultCauseClauses(result)
	causeText := strings.Join(causes, "; ")

	switch result.Outcome {
	case OutcomeSuccess:
		summary = fmt.Sprintf("%s completed %s cleanly. The delivery landed with cargo intact.", result.RunnerName, result.JobTitle)
	case OutcomePartial:
		if causeText == "" {
			causeText = "the client accepted the drop under protest"
		}
		summary = fmt.Sprintf("%s completed %s, but the delivery came in rough: %s.", result.RunnerName, result.JobTitle, causeText)
	case OutcomeFailed:
		if causeText == "" {
			causeText = "the handoff collapsed before contact"
		}
		summary = fmt.Sprintf("%s could not complete %s: %s. No payout cleared.", result.RunnerName, result.JobTitle, causeText)
	case OutcomeIntercepted:
		if causeText == "" {
			causeText = "security caught the line before the drop"
		}
		summary = fmt.Sprintf("%s lost %s to an intercept: %s. Cargo is gone and heat climbs.", result.RunnerName, result.JobTitle, causeText)
	default:
		summary = fmt.Sprintf("%s returned from %s: %s.", result.RunnerName, result.JobTitle, result.Outcome)
	}

	if len(result.Factors) > 0 {
		limit := min(len(result.Factors), 3)
		summary += " Factors: " + strings.Join(result.Factors[:limit], ", ") + "."
	}
	return summary
}

func resultCauseClauses(result JobResult) []string {
	causes := []string{}
	if result.Delay {
		causes = append(causes, "route delay forced a late handoff")
	}
	if result.CargoDamage {
		causes = append(causes, "cargo was damaged in transit")
	}
	if result.InjuryDetail != nil {
		causes = append(causes, result.InjuryDetail.Summary)
	} else if result.Injury {
		causes = append(causes, "the runner was hurt during the run")
	}
	if result.Detection {
		causes = append(causes, "the route drew heat")
	}
	if result.Complication {
		if result.ComplicationType != ComplicationNone {
			causes = append(causes, fmt.Sprintf("%s complication interrupted the line", strings.ReplaceAll(string(result.ComplicationType), "_", " ")))
		} else {
			causes = append(causes, "a complication interrupted the line")
		}
	}
	if result.Outcome == OutcomePartial && result.Payout > 0 {
		causes = append(causes, "the client cut the payout")
	}
	return causes
}

func injuryDetailForResult(state GameState, job Job, route Route, outcome JobOutcome, injury bool, detection bool, complicationType ComplicationType) *InjuryDetail {
	if !injury {
		return nil
	}

	severity := "minor"
	if outcome == OutcomeFailed || outcome == OutcomeIntercepted || detection || complicationType != ComplicationNone {
		severity = "serious"
	}
	if route.Type == RouteDroneCorridor || route.Type == RouteFloodline || job.Cargo == CargoWitness {
		severity = "serious"
	}

	cause := injuryCause(job, route, detection, complicationType)
	recovery := injuryRecoveryTurns(state)
	summary := fmt.Sprintf("%s injury: %s; recovery estimate %d turns", severity, cause, recovery)
	return &InjuryDetail{
		Severity:      severity,
		Cause:         cause,
		RecoveryTurns: recovery,
		Summary:       summary,
	}
}

func injuryCause(job Job, route Route, detection bool, complicationType ComplicationType) string {
	if complicationType != ComplicationNone {
		return strings.ReplaceAll(string(complicationType), "_", " ") + " pressure"
	}
	if detection {
		return "checkpoint contact"
	}
	if job.Cargo == CargoWitness {
		return "witness panic"
	}
	switch route.Type {
	case RouteServiceTunnels:
		return "tunnel fall"
	case RouteFloodline:
		return "floodline impact"
	case RouteDroneCorridor:
		return "drone corridor scramble"
	case RouteMarketWeave:
		return "market crowd crush"
	case RouteMainArtery:
		return "traffic breakaway"
	}
	return "route accident"
}

func bundleLoadMap(bundles []Bundle) map[RunnerID]int {
	loads := make(map[RunnerID]int, len(bundles))
	for _, bundle := range bundles {
		loads[bundle.RunnerID] = len(bundle.Jobs)
	}
	return loads
}

func rollRisk(rng *rand.Rand, risk int) bool {
	return rng.Intn(100) < clampRisk(risk)
}

func clampRisk(risk int) int {
	if risk < 0 {
		return 0
	}
	if risk > 95 {
		return 95
	}
	return risk
}

func clampInt(value int, low int, high int) int {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}

func maxInt(value int, floor int) int {
	if value < floor {
		return floor
	}
	return value
}

func routeModifier(routeType RouteType, modifiers map[RouteType]int) int {
	return modifiers[routeType]
}

func cargoModifier(cargo CargoType, modifiers map[CargoType]int) int {
	return modifiers[cargo]
}

func districtSurveillance(state GameState, districtID DistrictID) int {
	for _, district := range state.Districts {
		if district.ID == districtID {
			return district.Surveillance
		}
	}
	return 0
}

func districtDanger(state GameState, districtID DistrictID) int {
	for _, district := range state.Districts {
		if district.ID == districtID {
			return district.Danger
		}
	}
	return 0
}

func factionSuspicion(state GameState, factionID FactionID) int {
	for _, faction := range state.Factions {
		if faction.ID == factionID {
			return faction.Suspicion
		}
	}
	return 0
}

func findFactionIndex(factions []Faction, factionID FactionID) int {
	for i, faction := range factions {
		if faction.ID == factionID {
			return i
		}
	}
	return -1
}

func resultJobFaction(activeJobs []ActiveJob, jobID string) FactionID {
	for _, active := range activeJobs {
		if active.JobID == jobID {
			return active.Job.Faction
		}
	}
	return ""
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
