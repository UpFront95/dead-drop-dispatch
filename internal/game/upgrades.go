package game

import (
	"errors"
	"fmt"
)

var (
	ErrUpgradeNotFound         = errors.New("upgrade not found")
	ErrUpgradeAlreadyPurchased = errors.New("upgrade already purchased")
)

func FirstUpgradeDefinitions() []UpgradeDefinition {
	return copyUpgradeDefinitions([]UpgradeDefinition{
		upgrade(
			UpgradeSignalRelay,
			"Signal relay",
			"Improves route intel and comms around signal-poor districts.",
			220,
			"route intel",
			"signal recovery",
		),
		upgrade(
			UpgradeSafehouse,
			"Safehouse",
			"Gives runners a place to cool down between turns or nights.",
			260,
			"stress recovery",
			"curfew shelter",
		),
		upgrade(
			UpgradeFakeCredentialPrinter,
			"Fake credential printer",
			"Produces throwaway credentials for checkpoint pressure.",
			300,
			"checkpoint outcomes",
			"authority pressure",
		),
		upgrade(
			UpgradeClinicFavor,
			"Clinic favor",
			"Buys medical goodwill before the desk needs emergency care.",
			240,
			"injury recovery",
			"treatment discount",
		),
		upgrade(
			UpgradeDeadDropLocker,
			"Dead-drop locker",
			"Adds protected storage for contraband and awkward package timing.",
			180,
			"contraband risk",
			"cargo storage",
		),
		upgrade(
			UpgradeScrambler,
			"Scrambler",
			"Disrupts trace pressure on data shard routes.",
			320,
			"data trace risk",
			"scanner exposure",
		),
	})
}

func UpgradeDefinitionFor(upgradeID UpgradeID) (UpgradeDefinition, bool) {
	for _, definition := range FirstUpgradeDefinitions() {
		if definition.ID == upgradeID {
			return definition, true
		}
	}
	return UpgradeDefinition{}, false
}

func AvailableUpgrades(state GameState) []UpgradeDefinition {
	definitions := FirstUpgradeDefinitions()
	available := make([]UpgradeDefinition, 0, len(definitions))
	for _, definition := range definitions {
		if !HasUpgrade(state, definition.ID) {
			available = append(available, definition)
		}
	}
	return available
}

func HasUpgrade(state GameState, upgradeID UpgradeID) bool {
	for _, purchased := range state.PurchasedUpgrades {
		if purchased == upgradeID {
			return true
		}
	}
	return false
}

func PurchaseUpgrade(state *GameState, upgradeID UpgradeID) (UpgradeDefinition, EconomyTransaction, error) {
	definition, ok := UpgradeDefinitionFor(upgradeID)
	if !ok {
		return UpgradeDefinition{}, EconomyTransaction{}, ErrUpgradeNotFound
	}
	if HasUpgrade(*state, upgradeID) {
		return UpgradeDefinition{}, EconomyTransaction{}, ErrUpgradeAlreadyPurchased
	}

	transaction, err := SpendUpgrade(state, definition.Cost, fmt.Sprintf("%s upgrade", definition.Name))
	if err != nil {
		return UpgradeDefinition{}, EconomyTransaction{}, err
	}
	state.PurchasedUpgrades = append(state.PurchasedUpgrades, upgradeID)
	state.EventLog = append(state.EventLog, LogEntry{
		Turn: state.Turn,
		Text: fmt.Sprintf("Upgrade installed: %s.", definition.Name),
	})
	return definition, transaction, nil
}

func upgrade(id UpgradeID, name string, description string, cost int, effects ...string) UpgradeDefinition {
	return UpgradeDefinition{
		ID:          id,
		Name:        name,
		Description: description,
		Cost:        cost,
		Effects:     append([]string(nil), effects...),
	}
}

func copyUpgradeDefinitions(definitions []UpgradeDefinition) []UpgradeDefinition {
	copied := make([]UpgradeDefinition, len(definitions))
	for i, definition := range definitions {
		copied[i] = definition
		copied[i].Effects = append([]string(nil), definition.Effects...)
	}
	return copied
}
