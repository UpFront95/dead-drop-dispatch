package game_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"dead-drop-dispatch/internal/content"
	"dead-drop-dispatch/internal/game"
)

func TestFirstUpgradeDefinitionsCoverMVPShop(t *testing.T) {
	definitions := game.FirstUpgradeDefinitions()
	if got, want := len(definitions), len(mvpUpgradeIDs); got != want {
		t.Fatalf("upgrade count = %d, want %d", got, want)
	}

	gotIDs := make([]game.UpgradeID, 0, len(definitions))
	seen := map[game.UpgradeID]bool{}
	for _, definition := range definitions {
		gotIDs = append(gotIDs, definition.ID)
		if definition.ID == "" {
			t.Fatal("upgrade id should not be empty")
		}
		if seen[definition.ID] {
			t.Fatalf("duplicate upgrade id %s", definition.ID)
		}
		seen[definition.ID] = true
		if definition.Name == "" {
			t.Fatalf("%s name should not be empty", definition.ID)
		}
		if definition.Description == "" {
			t.Fatalf("%s description should not be empty", definition.ID)
		}
		if definition.Cost <= 0 {
			t.Fatalf("%s cost = %d, want positive", definition.ID, definition.Cost)
		}
		if len(definition.Effects) == 0 {
			t.Fatalf("%s should list future effects", definition.ID)
		}
	}

	if !reflect.DeepEqual(gotIDs, mvpUpgradeIDs) {
		t.Fatalf("upgrade ids = %#v, want %#v", gotIDs, mvpUpgradeIDs)
	}
}

func TestUpgradeDefinitionForFindsKnownDefinition(t *testing.T) {
	definition, ok := game.UpgradeDefinitionFor(game.UpgradeSafehouse)
	if !ok {
		t.Fatal("UpgradeDefinitionFor returned ok=false")
	}
	if definition.Name != "Safehouse" {
		t.Fatalf("name = %q, want Safehouse", definition.Name)
	}
}

func TestUpgradeDefinitionForRejectsUnknownDefinition(t *testing.T) {
	_, ok := game.UpgradeDefinitionFor(game.UpgradeID("unknown"))
	if ok {
		t.Fatal("UpgradeDefinitionFor returned ok=true for unknown upgrade")
	}
}

func TestFirstUpgradeDefinitionsReturnsCopies(t *testing.T) {
	definitions := game.FirstUpgradeDefinitions()
	definitions[0].Name = "mutated"
	definitions[0].Effects[0] = "mutated effect"

	definition, ok := game.UpgradeDefinitionFor(game.UpgradeSignalRelay)
	if !ok {
		t.Fatal("missing signal relay")
	}
	if definition.Name != "Signal relay" {
		t.Fatalf("name = %q, want Signal relay", definition.Name)
	}
	if definition.Effects[0] != "route intel" {
		t.Fatalf("effect = %q, want route intel", definition.Effects[0])
	}
}

func TestAvailableUpgradesExcludesPurchasedUpgrades(t *testing.T) {
	state := content.InitialGameState(42)
	state.PurchasedUpgrades = []game.UpgradeID{game.UpgradeSafehouse, game.UpgradeScrambler}

	available := game.AvailableUpgrades(state)
	gotIDs := upgradeIDs(available)

	for _, purchased := range state.PurchasedUpgrades {
		if containsUpgradeID(gotIDs, purchased) {
			t.Fatalf("available upgrades included purchased upgrade %s", purchased)
		}
	}
	if got, want := len(available), len(mvpUpgradeIDs)-len(state.PurchasedUpgrades); got != want {
		t.Fatalf("available count = %d, want %d", got, want)
	}
}

func TestPurchaseUpgradeSpendsCreditsAndRecordsOwnership(t *testing.T) {
	state := content.InitialGameState(42)
	state.Credits = 300
	state.EventLog = nil

	definition, transaction, err := game.PurchaseUpgrade(&state, game.UpgradeSafehouse)
	if err != nil {
		t.Fatalf("PurchaseUpgrade returned error: %v", err)
	}

	if got, want := definition.ID, game.UpgradeSafehouse; got != want {
		t.Fatalf("definition id = %q, want %q", got, want)
	}
	if got, want := state.Credits, 300-definition.Cost; got != want {
		t.Fatalf("credits = %d, want %d", got, want)
	}
	if got, want := transaction.Kind, game.TransactionUpgrade; got != want {
		t.Fatalf("transaction kind = %q, want %q", got, want)
	}
	if got, want := transaction.Amount, definition.Cost; got != want {
		t.Fatalf("transaction amount = %d, want %d", got, want)
	}
	if !game.HasUpgrade(state, game.UpgradeSafehouse) {
		t.Fatal("state should own safehouse")
	}
	if got, want := state.PurchasedUpgrades, []game.UpgradeID{game.UpgradeSafehouse}; !reflect.DeepEqual(got, want) {
		t.Fatalf("purchased upgrades = %#v, want %#v", got, want)
	}
	if got, want := len(state.EventLog), 2; got != want {
		t.Fatalf("event log count = %d, want %d", got, want)
	}
	if got := state.EventLog[0].Text; !strings.Contains(got, "Credits -260") {
		t.Fatalf("economy log = %q, want upgrade spend", got)
	}
	if got := state.EventLog[1].Text; !strings.Contains(got, "Upgrade installed: Safehouse.") {
		t.Fatalf("upgrade log = %q, want install log", got)
	}
}

func TestPurchaseUpgradeRejectsUnknownUpgrade(t *testing.T) {
	state := content.InitialGameState(42)
	startCredits := state.Credits
	startLogs := len(state.EventLog)

	_, _, err := game.PurchaseUpgrade(&state, game.UpgradeID("unknown"))

	if !errors.Is(err, game.ErrUpgradeNotFound) {
		t.Fatalf("error = %v, want %v", err, game.ErrUpgradeNotFound)
	}
	if state.Credits != startCredits {
		t.Fatalf("credits = %d, want unchanged %d", state.Credits, startCredits)
	}
	if len(state.EventLog) != startLogs {
		t.Fatalf("event log count = %d, want %d", len(state.EventLog), startLogs)
	}
}

func TestPurchaseUpgradeRejectsAlreadyPurchasedUpgrade(t *testing.T) {
	state := content.InitialGameState(42)
	state.PurchasedUpgrades = []game.UpgradeID{game.UpgradeScrambler}
	startCredits := state.Credits

	_, _, err := game.PurchaseUpgrade(&state, game.UpgradeScrambler)

	if !errors.Is(err, game.ErrUpgradeAlreadyPurchased) {
		t.Fatalf("error = %v, want %v", err, game.ErrUpgradeAlreadyPurchased)
	}
	if state.Credits != startCredits {
		t.Fatalf("credits = %d, want unchanged %d", state.Credits, startCredits)
	}
}

func TestPurchaseUpgradeRejectsInsufficientCredits(t *testing.T) {
	state := content.InitialGameState(42)
	state.Credits = 10
	state.EventLog = nil

	_, _, err := game.PurchaseUpgrade(&state, game.UpgradeDeadDropLocker)

	if !errors.Is(err, game.ErrInsufficientCredits) {
		t.Fatalf("error = %v, want %v", err, game.ErrInsufficientCredits)
	}
	if game.HasUpgrade(state, game.UpgradeDeadDropLocker) {
		t.Fatal("state should not own rejected upgrade")
	}
	if len(state.EventLog) != 0 {
		t.Fatalf("event log count = %d, want 0", len(state.EventLog))
	}
}

var mvpUpgradeIDs = []game.UpgradeID{
	game.UpgradeSignalRelay,
	game.UpgradeSafehouse,
	game.UpgradeFakeCredentialPrinter,
	game.UpgradeClinicFavor,
	game.UpgradeDeadDropLocker,
	game.UpgradeScrambler,
}

func upgradeIDs(definitions []game.UpgradeDefinition) []game.UpgradeID {
	ids := make([]game.UpgradeID, 0, len(definitions))
	for _, definition := range definitions {
		ids = append(ids, definition.ID)
	}
	return ids
}

func containsUpgradeID(ids []game.UpgradeID, want game.UpgradeID) bool {
	for _, id := range ids {
		if id == want {
			return true
		}
	}
	return false
}
