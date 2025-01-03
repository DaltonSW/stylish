package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"go.dalton.dog/colorgen/tui"
)

func main() {
	if _, err := tea.NewProgram(tui.NewLandingModel(), tea.WithAltScreen()).Run(); err != nil {
		log.Fatalf("Error starting program: %v\n", err)
	}
}
