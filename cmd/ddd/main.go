package main

import (
	"fmt"

	"dead-drop-dispatch/internal/content"
)

func main() {
	state := content.InitialGameState(0)

	fmt.Printf("Dead Drop Dispatch\n")
	fmt.Printf("Night %d/%d, Turn %d/%d\n", state.Night, state.RunNights, state.Turn, state.TurnsPerNight)
	fmt.Printf("Credits %d | Heat %d | Integrity %d\n", state.Credits, state.Heat, state.DispatchIntegrity)
	fmt.Printf("Districts %d | Runners %d | Factions %d\n", len(state.Districts), len(state.Runners), len(state.Factions))
	fmt.Printf("Desk: %s\n", state.Messages[0].Body)
}
