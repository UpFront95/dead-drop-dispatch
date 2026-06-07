package game

const (
	SignalRelayIntelClaimBonus              = 1
	SignalRelaySignalRiskReduction          = 8
	SafehouseBaseStressRecovery             = 1
	SafehouseUpgradeStressRecoveryBonus     = 1
	FakeCredentialCheckpointStressReduction = 1
	FakeCredentialBribeDiscount             = 25
	ClinicFavorInjuryRecoveryTurns          = 1
	DefaultInjuryRecoveryTurns              = 2
	DeadDropLockerContrabandRiskReduction   = 8
	ScramblerDataTraceRiskReduction         = 10
)

type RunnerRecoveryReport struct {
	StressRecovered int `json:"stress_recovered"`
	InjuryTicks     int `json:"injury_ticks"`
	RunnersReadied  int `json:"runners_readied"`
}

func intelClaimsPerJob(state GameState) int {
	claims := DefaultIntelClaimsPerJob
	if HasUpgrade(state, UpgradeSignalRelay) {
		claims += SignalRelayIntelClaimBonus
	}
	return claims
}

func signalRelayRiskReduction(state GameState, job Job, route Route) int {
	if !HasUpgrade(state, UpgradeSignalRelay) {
		return 0
	}
	if containsString(job.RiskFactors, "weak signal") || containsString(route.Traits, "signal drops") {
		return SignalRelaySignalRiskReduction
	}
	return 0
}

func contrabandRiskReduction(state GameState, job Job) int {
	if HasUpgrade(state, UpgradeDeadDropLocker) && job.Cargo == CargoContrabandPackage {
		return DeadDropLockerContrabandRiskReduction
	}
	return 0
}

func dataTraceRiskReduction(state GameState, job Job) int {
	if HasUpgrade(state, UpgradeScrambler) && job.Cargo == CargoDataShard {
		return ScramblerDataTraceRiskReduction
	}
	return 0
}

func injuryRecoveryTurns(state GameState) int {
	if HasUpgrade(state, UpgradeClinicFavor) {
		return ClinicFavorInjuryRecoveryTurns
	}
	return DefaultInjuryRecoveryTurns
}

func checkpointStressGain(state GameState, baseGain int) int {
	if HasUpgrade(state, UpgradeFakeCredentialPrinter) {
		return maxInt(baseGain-FakeCredentialCheckpointStressReduction, 0)
	}
	return baseGain
}

func checkpointBribeCost(state GameState) int {
	if HasUpgrade(state, UpgradeFakeCredentialPrinter) {
		return maxInt(ComplicationBribeCost-FakeCredentialBribeDiscount, 0)
	}
	return ComplicationBribeCost
}

// RecoverRunners applies between-cycle runner recovery. Injuries tick down first;
// ready runners use the safehouse bonus to shed stress faster.
func RecoverRunners(state *GameState) RunnerRecoveryReport {
	stressRecovery := SafehouseBaseStressRecovery
	if HasUpgrade(*state, UpgradeSafehouse) {
		stressRecovery += SafehouseUpgradeStressRecoveryBonus
	}

	report := RunnerRecoveryReport{}
	for i := range state.Runners {
		if state.Runners[i].State == RunnerInjured {
			before := state.Runners[i].Recovery
			state.Runners[i].Recovery = maxInt(before-1, 0)
			if state.Runners[i].Recovery < before {
				report.InjuryTicks++
			}
			if state.Runners[i].Recovery == 0 {
				state.Runners[i].State = RunnerReady
				report.RunnersReadied++
			}
			continue
		}
		if state.Runners[i].State == RunnerReady {
			before := state.Runners[i].Stress
			state.Runners[i].Stress = maxInt(before-stressRecovery, 0)
			report.StressRecovered += before - state.Runners[i].Stress
		}
	}
	return report
}
