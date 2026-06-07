package game

import (
	"errors"
	"fmt"
)

var (
	ErrJobNotAvailable    = errors.New("job not available")
	ErrJobNotAccepted     = errors.New("job not accepted")
	ErrRunnerNotFound     = errors.New("runner not found")
	ErrRunnerBusy         = errors.New("runner unavailable")
	ErrRunnerAtCapacity   = errors.New("runner at capacity")
	ErrRouteNotFound      = errors.New("route not found")
	ErrBundleIncompatible = errors.New("bundle incompatible")
)

func AcceptJob(state *GameState, jobID string) error {
	index := findJobIndex(state.AvailableJobs, jobID)
	if index < 0 {
		return ErrJobNotAvailable
	}

	job := state.AvailableJobs[index]
	state.AvailableJobs = append(state.AvailableJobs[:index], state.AvailableJobs[index+1:]...)
	state.AcceptedJobs = append(state.AcceptedJobs, job)
	state.EventLog = append(state.EventLog, LogEntry{
		Turn: state.Turn,
		Text: fmt.Sprintf("Accepted contract: %s.", job.Title),
	})
	return nil
}

func AssignAcceptedJob(state *GameState, jobID string, runnerID RunnerID, routeID string) error {
	jobIndex := findJobIndex(state.AcceptedJobs, jobID)
	if jobIndex < 0 {
		return ErrJobNotAccepted
	}

	runnerIndex := findRunnerIndex(state.Runners, runnerID)
	if runnerIndex < 0 {
		return ErrRunnerNotFound
	}
	job := state.AcceptedJobs[jobIndex]
	route, ok := findRoute(job, routeID)
	if !ok {
		return ErrRouteNotFound
	}

	bundleIndex := findBundleIndex(state.Bundles, runnerID)
	if err := canAssignToRunner(*state, runnerIndex, bundleIndex, job, route); err != nil {
		return err
	}

	active := ActiveJob{
		JobID:    job.ID,
		RunnerID: runnerID,
		RouteID:  route.ID,
		Job:      job,
		Route:    route,
	}

	state.AcceptedJobs = append(state.AcceptedJobs[:jobIndex], state.AcceptedJobs[jobIndex+1:]...)
	state.ActiveJobs = append(state.ActiveJobs, active)
	state.Runners[runnerIndex].State = RunnerOnJob
	if bundleIndex >= 0 {
		state.Bundles[bundleIndex].Jobs = append(state.Bundles[bundleIndex].Jobs, active)
		state.Bundles[bundleIndex].Penalties = BundlePenalties(state.Bundles[bundleIndex])
	} else {
		state.Bundles = append(state.Bundles, Bundle{
			RunnerID: runnerID,
			Jobs:     []ActiveJob{active},
		})
	}

	action := "Assigned"
	if bundleIndex >= 0 {
		action = "Bundled"
	}
	state.EventLog = append(state.EventLog, LogEntry{
		Turn: state.Turn,
		Text: fmt.Sprintf("%s %s to %s.", action, job.Title, state.Runners[runnerIndex].Name),
	})
	return nil
}

func BundlePenalties(bundle Bundle) []string {
	if len(bundle.Jobs) < 2 {
		return nil
	}

	penalties := []string{"stress gain"}
	if hasDifferentDestinations(bundle.Jobs) {
		penalties = append(penalties, "destination complexity")
	}
	if hasCargoConflict(bundle.Jobs) {
		penalties = append(penalties, "cargo conflict")
	}
	if hasDetectionExposure(bundle.Jobs) {
		penalties = append(penalties, "detection exposure")
	}
	if hasDelayPressure(bundle.Jobs) {
		penalties = append(penalties, "delay pressure")
	}
	return uniqueStrings(penalties)
}

func canAssignToRunner(state GameState, runnerIndex int, bundleIndex int, job Job, route Route) error {
	runner := state.Runners[runnerIndex]
	if bundleIndex < 0 {
		if runner.State != RunnerReady {
			return ErrRunnerBusy
		}
		return nil
	}

	bundle := state.Bundles[bundleIndex]
	if len(bundle.Jobs) >= MaxJobsPerRunner {
		return ErrRunnerAtCapacity
	}
	if runner.State != RunnerOnJob {
		return ErrRunnerBusy
	}
	if !BundleCompatible(bundle, job, route) {
		return ErrBundleIncompatible
	}
	return nil
}

func BundleCompatible(bundle Bundle, job Job, route Route) bool {
	if len(bundle.Jobs) == 0 {
		return true
	}
	if job.Cargo == CargoWitness {
		return false
	}
	for _, active := range bundle.Jobs {
		if active.Job.Cargo == CargoWitness {
			return false
		}
		if active.Route.Type == route.Type || routesShareDistrict(active.Route, route) {
			return true
		}
	}
	return false
}

func findJobIndex(jobs []Job, jobID string) int {
	for i, job := range jobs {
		if job.ID == jobID {
			return i
		}
	}
	return -1
}

func findRunnerIndex(runners []Runner, runnerID RunnerID) int {
	for i, runner := range runners {
		if runner.ID == runnerID {
			return i
		}
	}
	return -1
}

func findBundleIndex(bundles []Bundle, runnerID RunnerID) int {
	for i, bundle := range bundles {
		if bundle.RunnerID == runnerID {
			return i
		}
	}
	return -1
}

func findRoute(job Job, routeID string) (Route, bool) {
	for _, route := range job.Routes {
		if route.ID == routeID {
			return route, true
		}
	}
	return Route{}, false
}

func routesShareDistrict(a Route, b Route) bool {
	for _, aDistrict := range a.Districts {
		for _, bDistrict := range b.Districts {
			if aDistrict == bDistrict {
				return true
			}
		}
	}
	return false
}

func hasDifferentDestinations(jobs []ActiveJob) bool {
	if len(jobs) < 2 {
		return false
	}
	destination := jobs[0].Job.Destination
	for _, job := range jobs[1:] {
		if job.Job.Destination != destination {
			return true
		}
	}
	return false
}

func hasCargoConflict(jobs []ActiveJob) bool {
	seenPhysicalRisk := false
	for _, active := range jobs {
		switch active.Job.Cargo {
		case CargoMedicalCooler, CargoContrabandPackage, CargoCorporatePrototype:
			if seenPhysicalRisk {
				return true
			}
			seenPhysicalRisk = true
		}
	}
	return false
}

func hasDetectionExposure(jobs []ActiveJob) bool {
	for _, active := range jobs {
		if active.Job.Cargo == CargoContrabandPackage || active.Job.Cargo == CargoCorporatePrototype {
			return true
		}
		if active.Route.Type == RouteMainArtery || active.Route.Type == RouteDroneCorridor {
			return true
		}
	}
	return false
}

func hasDelayPressure(jobs []ActiveJob) bool {
	totalTime := 0
	shortestDeadline := 0
	for _, active := range jobs {
		totalTime += active.Route.TimeCost
		if shortestDeadline == 0 || active.Job.DeadlineTurns < shortestDeadline {
			shortestDeadline = active.Job.DeadlineTurns
		}
	}
	return shortestDeadline > 0 && totalTime > shortestDeadline
}
