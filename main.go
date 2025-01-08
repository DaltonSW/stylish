package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"go.dalton.dog/colorgen/tui"
)

func main() {
	// log.SetLevel(log.DebugLevel)

	// Check that dircolors command is installed
	_, err := exec.LookPath("dircolors")
	if err != nil {
		log.Fatal("Package `dircolors` not found on PATH. Please install it before continuing.")
	}

	// If no arguments are provided, just run the TUI
	if len(os.Args) == 1 {
		program := tea.NewProgram(tui.NewLandingModel(), tea.WithAltScreen())
		if _, err := program.Run(); err != nil {
			log.Fatalf("Error running program: %v\n", err)
		}
		return
	}

	// Else, evaluate what the command is asking for
	// I'm probably only going to have `help`, `generate`, and `apply` as subcommands
	switch os.Args[1] {
	case "help":
		log.Printf("Help command")
	case "generate":
		if len(os.Args) < 3 {
			log.Fatal("No theme provided")
		}
		theme := os.Args[2]
		log.Print("Generate command for " + theme)
		// newModel := tui.NewThemeModel(theme)
		// err := newModel.GenerateDirColors()
		// if err != nil {
		// 	log.Fatal(err)
		// } else {
		// 	log.Print("Successfully generated file at " + filepath.Join(tui.ThemeConfigFolder, theme, ".dircolors"))
		// }
	case "apply":
		if len(os.Args) < 3 {
			log.Fatal("No theme provided")
		}
		theme := os.Args[2]
		// log.Printf("Apply command for " + theme)
		cmd := exec.Command("dircolors", filepath.Join(tui.ThemeConfigFolder, theme, ".dircolors"))
		cmdOut, cmdErr := cmd.Output()
		if cmdErr != nil {
			log.Fatal(cmdErr.Error() + string(cmdOut))
		}
		// cmdStr := string(cmdOut)
		// outStr := strings.TrimPrefix(strings.TrimSuffix(cmdStr, ";\nexport LS_COLORS\n"), "LS_COLORS=")
		fmt.Print(string(cmdOut))
	default:
		log.Fatal("Command not found")
	}

}
