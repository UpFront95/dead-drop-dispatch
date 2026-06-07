package game

import "errors"

const DefaultIntelClaimsPerJob = 2

var (
	ErrIntelSourceNotSet     = errors.New("intel source not set")
	ErrUnsupportedIntelClaim = errors.New("unsupported intel claim")
)

func NewIntelReport(state GameState, source IntelSource, claims []string, omittedTags []string, riskTags []IntelRiskTag) (IntelReport, error) {
	if source == "" {
		return IntelReport{}, ErrIntelSourceNotSet
	}
	claims = uniqueStrings(claims)
	omittedTags = uniqueStrings(omittedTags)
	riskTags = uniqueIntelRiskTags(riskTags)
	if len(omittedTags) > 0 && !intelRiskTagsContain(riskTags, IntelRiskIncomplete) {
		riskTags = append(riskTags, IntelRiskIncomplete)
	}

	return IntelReport{
		Source:      source,
		Night:       state.Night,
		Turn:        state.Turn,
		Staleness:   0,
		Confidence:  intelConfidenceForRiskTags(riskTags),
		Claims:      append([]string(nil), claims...),
		OmittedTags: append([]string(nil), omittedTags...),
		RiskTags:    append([]IntelRiskTag(nil), riskTags...),
	}, nil
}

func JobIntelReport(state GameState, riskFactors []string) IntelReport {
	claims := uniqueStrings(riskFactors)
	omittedTags := []string{}
	claimLimit := intelClaimsPerJob(state)
	if len(claims) > claimLimit {
		omittedTags = append([]string(nil), claims[claimLimit:]...)
		claims = append([]string(nil), claims[:claimLimit]...)
	}

	report, err := NewIntelReport(state, IntelSourceClient, claims, omittedTags, nil)
	if err != nil {
		return IntelReport{}
	}
	return report
}

func AgeIntelReport(report IntelReport, state GameState) IntelReport {
	aged := copyIntelReport(report)
	turnsPerNight := state.TurnsPerNight
	if turnsPerNight <= 0 {
		turnsPerNight = DefaultTurnsPerNight
	}
	aged.Staleness = maxInt(0, (state.Night-report.Night)*turnsPerNight+(state.Turn-report.Turn))
	if aged.Staleness > 0 && !intelRiskTagsContain(aged.RiskTags, IntelRiskStale) {
		aged.RiskTags = append(aged.RiskTags, IntelRiskStale)
	}
	aged.Confidence = intelConfidenceForRiskTags(aged.RiskTags)
	return aged
}

func IntelReportSupportsClaim(report IntelReport, claim string) bool {
	for _, supported := range report.Claims {
		if supported == claim {
			return true
		}
	}
	return false
}

func ValidateIntelClaim(report IntelReport, claim string) error {
	if !IntelReportSupportsClaim(report, claim) {
		return ErrUnsupportedIntelClaim
	}
	return nil
}

func copyIntelReport(report IntelReport) IntelReport {
	copied := report
	copied.Claims = append([]string(nil), report.Claims...)
	copied.OmittedTags = append([]string(nil), report.OmittedTags...)
	copied.RiskTags = append([]IntelRiskTag(nil), report.RiskTags...)
	return copied
}

func intelConfidenceForRiskTags(riskTags []IntelRiskTag) IntelConfidence {
	if intelRiskTagsContain(riskTags, IntelRiskFalse) {
		return IntelConfidenceLow
	}
	if intelRiskTagsContain(riskTags, IntelRiskIncomplete) || intelRiskTagsContain(riskTags, IntelRiskStale) || intelRiskTagsContain(riskTags, IntelRiskBiased) {
		return IntelConfidenceMedium
	}
	return IntelConfidenceHigh
}

func intelRiskTagsContain(riskTags []IntelRiskTag, want IntelRiskTag) bool {
	for _, riskTag := range riskTags {
		if riskTag == want {
			return true
		}
	}
	return false
}

func uniqueIntelRiskTags(values []IntelRiskTag) []IntelRiskTag {
	seen := map[IntelRiskTag]bool{}
	result := make([]IntelRiskTag, 0, len(values))
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	return result
}
