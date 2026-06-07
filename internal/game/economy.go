package game

import (
	"errors"
	"fmt"
)

type EconomyTransactionKind string

const (
	TransactionPayout        EconomyTransactionKind = "payout"
	TransactionBribe         EconomyTransactionKind = "bribe"
	TransactionTreatment     EconomyTransactionKind = "treatment"
	TransactionUpgrade       EconomyTransactionKind = "upgrade"
	TransactionOperatingCost EconomyTransactionKind = "operating_cost"
	TransactionNonpayment    EconomyTransactionKind = "nonpayment"
)

var (
	ErrInvalidTransactionAmount = errors.New("invalid transaction amount")
	ErrInsufficientCredits      = errors.New("insufficient credits")
)

type EconomyTransaction struct {
	Turn          int                    `json:"turn"`
	Kind          EconomyTransactionKind `json:"kind"`
	Amount        int                    `json:"amount"`
	CreditsBefore int                    `json:"credits_before"`
	CreditsAfter  int                    `json:"credits_after"`
	Note          string                 `json:"note,omitempty"`
}

func ApplyPayout(state *GameState, amount int, note string) (EconomyTransaction, error) {
	if amount < 0 {
		return EconomyTransaction{}, ErrInvalidTransactionAmount
	}
	return applyCreditChange(state, TransactionPayout, amount, note)
}

func RecordNonpayment(state *GameState, expectedAmount int, note string) (EconomyTransaction, error) {
	if expectedAmount < 0 {
		return EconomyTransaction{}, ErrInvalidTransactionAmount
	}
	transaction := EconomyTransaction{
		Turn:          state.Turn,
		Kind:          TransactionNonpayment,
		Amount:        expectedAmount,
		CreditsBefore: state.Credits,
		CreditsAfter:  state.Credits,
		Note:          note,
	}
	appendEconomyLog(state, transaction)
	return transaction, nil
}

func SpendBribe(state *GameState, amount int, note string) (EconomyTransaction, error) {
	return spendCredits(state, TransactionBribe, amount, note, false)
}

func SpendTreatment(state *GameState, amount int, note string) (EconomyTransaction, error) {
	return spendCredits(state, TransactionTreatment, amount, note, false)
}

func SpendUpgrade(state *GameState, amount int, note string) (EconomyTransaction, error) {
	return spendCredits(state, TransactionUpgrade, amount, note, false)
}

func ApplyOperatingCost(state *GameState, amount int, note string) (EconomyTransaction, error) {
	return spendCredits(state, TransactionOperatingCost, amount, note, true)
}

func spendCredits(state *GameState, kind EconomyTransactionKind, amount int, note string, allowDebt bool) (EconomyTransaction, error) {
	if amount < 0 {
		return EconomyTransaction{}, ErrInvalidTransactionAmount
	}
	if !allowDebt && state.Credits < amount {
		return EconomyTransaction{}, ErrInsufficientCredits
	}
	return applyCreditChange(state, kind, -amount, note)
}

func applyCreditChange(state *GameState, kind EconomyTransactionKind, delta int, note string) (EconomyTransaction, error) {
	before := state.Credits
	state.Credits += delta
	transaction := EconomyTransaction{
		Turn:          state.Turn,
		Kind:          kind,
		Amount:        abs(delta),
		CreditsBefore: before,
		CreditsAfter:  state.Credits,
		Note:          note,
	}
	appendEconomyLog(state, transaction)
	return transaction, nil
}

func appendEconomyLog(state *GameState, transaction EconomyTransaction) {
	state.EventLog = append(state.EventLog, LogEntry{
		Turn: state.Turn,
		Text: economyLogText(transaction),
	})
}

func economyLogText(transaction EconomyTransaction) string {
	note := transaction.Note
	if note == "" {
		note = string(transaction.Kind)
	}
	switch transaction.Kind {
	case TransactionPayout:
		return fmt.Sprintf("Credits +%d: %s.", transaction.Amount, note)
	case TransactionNonpayment:
		return fmt.Sprintf("Expected payout %d not received: %s.", transaction.Amount, note)
	default:
		return fmt.Sprintf("Credits -%d: %s.", transaction.Amount, note)
	}
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
