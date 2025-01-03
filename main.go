package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"go.dalton.dog/colorgen/tui"
)

func main() {
	w, h := tui.GetTermSize()
	program := tea.NewProgram(tui.NewLandingModel(w, h), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		log.Fatalf("Error starting program: %v\n", err)
	}
}
