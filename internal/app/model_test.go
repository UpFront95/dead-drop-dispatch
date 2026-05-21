package app

import (
	"testing"

	"dead-drop-dispatch/internal/tui"
)

func TestNewModelUsesInitialState(t *testing.T) {
	model := New(99)

	if model.width != tui.TargetWidth {
		t.Fatalf("width = %d, want %d", model.width, tui.TargetWidth)
	}
	if model.height != tui.TargetHeight {
		t.Fatalf("height = %d, want %d", model.height, tui.TargetHeight)
	}
	if model.state.RandomSeed != 99 {
		t.Fatalf("seed = %d, want 99", model.state.RandomSeed)
	}
	if len(model.state.Districts) != 5 {
		t.Fatalf("district count = %d, want 5", len(model.state.Districts))
	}
}
