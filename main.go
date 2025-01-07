package main

import (
	"log"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"go.dalton.dog/colorgen/tui"
)

func main() {
	// Check that dircolors command is installed
	_, err := exec.LookPath("dircolors")
	if err != nil {
		log.Fatal("Package `dircolors` not found on PATH. Please install it before continuing.")
	}

	// Clear screen -- Not needed if tea.WithAltScreen() is used
	// cmd := exec.Command("clear")
	// cmd.Stdout = os.Stdout
	// cmd.Run()

	program := tea.NewProgram(tui.NewLandingModel(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		log.Fatalf("Error starting program: %v\n", err)
	}
}
