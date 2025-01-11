package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"go.dalton.dog/stylish/tui"
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
	} else {
		handleArgs()
	}
}

func handleArgs() {
	// Else, evaluate what the command is asking for
	// I'm probably only going to have `help`, `generate`, and `apply` as subcommands
	switch os.Args[1] {
	case "help":
		log.Printf("Help command")

	case "generate":
		if len(os.Args) < 3 {
			log.Fatal("No theme provided")
		}
		themeName := os.Args[2]
		log.Print("Generate command for " + themeName)
		theme := tui.GetTheme(themeName)
		err := theme.GenerateDirColors()
		if err != nil {
			log.Fatal(err)
		} else {
			log.Print("Successfully generated file at " + filepath.Join(theme.Path, ".dircolors"))
		}

	case "apply":
		doApply()
	case "apply-8bit":
		tui.EightBitMode = true
		doApply()

	default:
		log.Fatal("Command not found")
	}

}

func doApply() {

	var theme string
	if len(os.Args) < 3 {
		theme = "default"
	} else {
		theme = os.Args[2]
	}

	err := tui.GetTheme(theme).GenerateDirColors()
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("dircolors", filepath.Join(tui.ThemeConfigFolder, theme, ".dircolors"))
	cmdOut, cmdErr := cmd.Output()
	if cmdErr != nil {
		log.Fatal(cmdErr.Error() + string(cmdOut))
	}
	fmt.Print(string(cmdOut))
}
