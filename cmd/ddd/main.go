package main

import (
	"log"

	tea "charm.land/bubbletea/v2"

	"dead-drop-dispatch/internal/app"
)

func main() {
	program := tea.NewProgram(app.New(0))
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}
