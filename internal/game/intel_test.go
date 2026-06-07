package game_test

import (
	"errors"
	"reflect"
	"testing"

	"dead-drop-dispatch/internal/content"
	"dead-drop-dispatch/internal/game"
)

func TestNewIntelReportRequiresSource(t *testing.T) {
	state := content.InitialGameState(42)

	_, err := game.NewIntelReport(state, "", []string{"traffic pressure"}, nil, nil)

	if !errors.Is(err, game.ErrIntelSourceNotSet) {
		t.Fatalf("error = %v, want %v", err, game.ErrIntelSourceNotSet)
	}
}

func TestNewIntelReportRecordsSourceConfidenceAndIncompleteTags(t *testing.T) {
	state := content.InitialGameState(42)

	report, err := game.NewIntelReport(
		state,
		game.IntelSourceStreet,
		[]string{"traffic pressure", "traffic pressure", "weak signal"},
		[]string{"destination surveillance"},
		[]game.IntelRiskTag{game.IntelRiskBiased, game.IntelRiskBiased},
	)
	if err != nil {
		t.Fatalf("NewIntelReport returned error: %v", err)
	}

	if got, want := report.Source, game.IntelSourceStreet; got != want {
		t.Fatalf("source = %q, want %q", got, want)
	}
	if got, want := report.Night, state.Night; got != want {
		t.Fatalf("night = %d, want %d", got, want)
	}
	if got, want := report.Turn, state.Turn; got != want {
		t.Fatalf("turn = %d, want %d", got, want)
	}
	if got, want := report.Staleness, 0; got != want {
		t.Fatalf("staleness = %d, want %d", got, want)
	}
	if got, want := report.Confidence, game.IntelConfidenceMedium; got != want {
		t.Fatalf("confidence = %q, want %q", got, want)
	}
	if got, want := report.Claims, []string{"traffic pressure", "weak signal"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("claims = %#v, want %#v", got, want)
	}
	if got, want := report.OmittedTags, []string{"destination surveillance"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("omitted tags = %#v, want %#v", got, want)
	}
	if got, want := report.RiskTags, []game.IntelRiskTag{game.IntelRiskBiased, game.IntelRiskIncomplete}; !reflect.DeepEqual(got, want) {
		t.Fatalf("risk tags = %#v, want %#v", got, want)
	}
}

func TestNewIntelReportMarksFalseRiskAsLowConfidence(t *testing.T) {
	state := content.InitialGameState(42)

	report, err := game.NewIntelReport(
		state,
		game.IntelSourceFaction,
		[]string{"union politics"},
		nil,
		[]game.IntelRiskTag{game.IntelRiskFalse},
	)
	if err != nil {
		t.Fatalf("NewIntelReport returned error: %v", err)
	}

	if got, want := report.Confidence, game.IntelConfidenceLow; got != want {
		t.Fatalf("confidence = %q, want %q", got, want)
	}
	if !reflect.DeepEqual(report.RiskTags, []game.IntelRiskTag{game.IntelRiskFalse}) {
		t.Fatalf("risk tags = %#v, want false risk tag", report.RiskTags)
	}
}

func TestJobIntelReportLimitsClaimsAndMarksOmissions(t *testing.T) {
	state := content.InitialGameState(42)

	report := game.JobIntelReport(state, []string{
		"cargo integrity",
		"client urgency",
		"weak signal",
		"traffic pressure",
	})

	if got, want := report.Source, game.IntelSourceClient; got != want {
		t.Fatalf("source = %q, want %q", got, want)
	}
	if got, want := report.Confidence, game.IntelConfidenceMedium; got != want {
		t.Fatalf("confidence = %q, want %q", got, want)
	}
	if got, want := report.Claims, []string{"cargo integrity", "client urgency"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("claims = %#v, want %#v", got, want)
	}
	if got, want := report.OmittedTags, []string{"weak signal", "traffic pressure"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("omitted tags = %#v, want %#v", got, want)
	}
	if got, want := report.RiskTags, []game.IntelRiskTag{game.IntelRiskIncomplete}; !reflect.DeepEqual(got, want) {
		t.Fatalf("risk tags = %#v, want %#v", got, want)
	}
}

func TestAgeIntelReportAddsStalenessRiskDeterministically(t *testing.T) {
	state := content.InitialGameState(42)
	report, err := game.NewIntelReport(state, game.IntelSourceSensor, []string{"drone corridor watched"}, nil, nil)
	if err != nil {
		t.Fatalf("NewIntelReport returned error: %v", err)
	}
	state.Turn += 2

	aged := game.AgeIntelReport(report, state)

	if got, want := aged.Staleness, 2; got != want {
		t.Fatalf("staleness = %d, want %d", got, want)
	}
	if got, want := aged.Confidence, game.IntelConfidenceMedium; got != want {
		t.Fatalf("confidence = %q, want %q", got, want)
	}
	if got, want := aged.RiskTags, []game.IntelRiskTag{game.IntelRiskStale}; !reflect.DeepEqual(got, want) {
		t.Fatalf("risk tags = %#v, want %#v", got, want)
	}
	if got, want := report.Staleness, 0; got != want {
		t.Fatalf("original report staleness = %d, want unchanged %d", got, want)
	}
}

func TestValidateIntelClaimRejectsUnsupportedClaim(t *testing.T) {
	state := content.InitialGameState(42)
	report, err := game.NewIntelReport(state, game.IntelSourceClient, []string{"traffic pressure"}, []string{"weak signal"}, nil)
	if err != nil {
		t.Fatalf("NewIntelReport returned error: %v", err)
	}

	if err := game.ValidateIntelClaim(report, "traffic pressure"); err != nil {
		t.Fatalf("ValidateIntelClaim returned error for supported claim: %v", err)
	}
	err = game.ValidateIntelClaim(report, "weak signal")

	if !errors.Is(err, game.ErrUnsupportedIntelClaim) {
		t.Fatalf("error = %v, want %v", err, game.ErrUnsupportedIntelClaim)
	}
}

func TestGeneratedJobsCarryIntelSnapshots(t *testing.T) {
	state := content.InitialGameState(42)

	for _, job := range state.AvailableJobs {
		if len(job.Intel) != 1 {
			t.Fatalf("job %s intel count = %d, want 1", job.ID, len(job.Intel))
		}
		report := job.Intel[0]
		if report.Source != game.IntelSourceClient {
			t.Fatalf("job %s intel source = %q, want %q", job.ID, report.Source, game.IntelSourceClient)
		}
		if report.Night != state.Night || report.Turn != state.Turn {
			t.Fatalf("job %s intel timestamp = n%d t%d, want n%d t%d", job.ID, report.Night, report.Turn, state.Night, state.Turn)
		}
		if len(report.Claims) == 0 {
			t.Fatalf("job %s intel has no claims", job.ID)
		}
		for _, claim := range report.Claims {
			if err := game.ValidateIntelClaim(report, claim); err != nil {
				t.Fatalf("job %s claim %q rejected: %v", job.ID, claim, err)
			}
		}
	}
}
