package game_test

import (
	"errors"
	"testing"

	"dead-drop-dispatch/internal/content"
	"dead-drop-dispatch/internal/game"
)

func TestApplyPayoutAddsCreditsAndLogsTransaction(t *testing.T) {
	state := content.InitialGameState(42)
	startCredits := state.Credits
	startLogs := len(state.EventLog)

	transaction, err := game.ApplyPayout(&state, 120, "clinic paid")
	if err != nil {
		t.Fatalf("ApplyPayout returned error: %v", err)
	}

	if got, want := state.Credits, startCredits+120; got != want {
		t.Fatalf("credits = %d, want %d", got, want)
	}
	if transaction.Kind != game.TransactionPayout {
		t.Fatalf("kind = %q, want %q", transaction.Kind, game.TransactionPayout)
	}
	if transaction.CreditsBefore != startCredits || transaction.CreditsAfter != state.Credits {
		t.Fatalf("transaction credits = %+v, want before %d after %d", transaction, startCredits, state.Credits)
	}
	if got, want := len(state.EventLog), startLogs+1; got != want {
		t.Fatalf("event log count = %d, want %d", got, want)
	}
}

func TestRecordNonpaymentLeavesCreditsUnchanged(t *testing.T) {
	state := content.InitialGameState(42)
	startCredits := state.Credits

	transaction, err := game.RecordNonpayment(&state, 180, "client ghosted")
	if err != nil {
		t.Fatalf("RecordNonpayment returned error: %v", err)
	}

	if state.Credits != startCredits {
		t.Fatalf("credits = %d, want %d", state.Credits, startCredits)
	}
	if transaction.Kind != game.TransactionNonpayment {
		t.Fatalf("kind = %q, want %q", transaction.Kind, game.TransactionNonpayment)
	}
	if transaction.Amount != 180 {
		t.Fatalf("amount = %d, want 180", transaction.Amount)
	}
}

func TestSpendBribeSubtractsCredits(t *testing.T) {
	state := content.InitialGameState(42)
	startCredits := state.Credits

	transaction, err := game.SpendBribe(&state, 75, "checkpoint")
	if err != nil {
		t.Fatalf("SpendBribe returned error: %v", err)
	}

	if got, want := state.Credits, startCredits-75; got != want {
		t.Fatalf("credits = %d, want %d", got, want)
	}
	if transaction.Kind != game.TransactionBribe {
		t.Fatalf("kind = %q, want %q", transaction.Kind, game.TransactionBribe)
	}
}

func TestSpendTreatmentRejectsInsufficientCredits(t *testing.T) {
	state := content.InitialGameState(42)
	state.Credits = 20
	startLogs := len(state.EventLog)

	_, err := game.SpendTreatment(&state, 50, "patch runner")
	if !errors.Is(err, game.ErrInsufficientCredits) {
		t.Fatalf("SpendTreatment error = %v, want %v", err, game.ErrInsufficientCredits)
	}
	if state.Credits != 20 {
		t.Fatalf("credits = %d, want 20", state.Credits)
	}
	if len(state.EventLog) != startLogs {
		t.Fatalf("event log changed on rejected spend")
	}
}

func TestApplyOperatingCostCanCreateDebt(t *testing.T) {
	state := content.InitialGameState(42)
	state.Credits = 25

	transaction, err := game.ApplyOperatingCost(&state, 40, "rent")
	if err != nil {
		t.Fatalf("ApplyOperatingCost returned error: %v", err)
	}

	if state.Credits != -15 {
		t.Fatalf("credits = %d, want -15", state.Credits)
	}
	if transaction.Kind != game.TransactionOperatingCost {
		t.Fatalf("kind = %q, want %q", transaction.Kind, game.TransactionOperatingCost)
	}
	status := game.EvaluateRunStatus(state)
	if status.Reason != game.RunEndBankrupt {
		t.Fatalf("run end reason = %q, want %q", status.Reason, game.RunEndBankrupt)
	}
}

func TestEconomyTransactionsRejectNegativeAmounts(t *testing.T) {
	state := content.InitialGameState(42)

	tests := []struct {
		name string
		fn   func() error
	}{
		{name: "payout", fn: func() error {
			_, err := game.ApplyPayout(&state, -1, "bad")
			return err
		}},
		{name: "nonpayment", fn: func() error {
			_, err := game.RecordNonpayment(&state, -1, "bad")
			return err
		}},
		{name: "bribe", fn: func() error {
			_, err := game.SpendBribe(&state, -1, "bad")
			return err
		}},
		{name: "treatment", fn: func() error {
			_, err := game.SpendTreatment(&state, -1, "bad")
			return err
		}},
		{name: "operating cost", fn: func() error {
			_, err := game.ApplyOperatingCost(&state, -1, "bad")
			return err
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if !errors.Is(err, game.ErrInvalidTransactionAmount) {
				t.Fatalf("error = %v, want %v", err, game.ErrInvalidTransactionAmount)
			}
		})
	}
}
