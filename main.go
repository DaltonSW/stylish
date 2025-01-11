package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"go.dalton.dog/stylish/tui"
)

// TODO: Swap to cobra for commands
// TODO: Move stuff into `cmd` and `internal` to publish on pkg.go.dev
// TODO: Look into `goreleaser`
// TODO: Finish README
// TODO: Update helptext
// TODO: Clean up and comment code

func main() {
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
	switch os.Args[1] {
	case "help":
		log.Printf("Help command")

	case "preview":
		output := doApply()
		output = strings.TrimPrefix(output, "LS_COLORS='")
		output = strings.TrimSpace(output)
		output = strings.TrimSuffix(output, "';\nexport LS_COLORS")

		os.Setenv("LS_COLORS", output)

		binary, pathErr := exec.LookPath("ls")
		if pathErr != nil {
			log.Fatal(errors.New("Path err: " + pathErr.Error()))
		}

		execErr := syscall.Exec(binary, []string{"ls", "--color=auto"}, os.Environ())
		if execErr != nil {
			log.Fatal(errors.New("Exec err: " + execErr.Error()))
		}

	case "apply":
		fmt.Print(doApply())
	case "apply-8bit":
		tui.EightBitMode = true
		fmt.Print(doApply())

	default:
		log.Fatal("Command not found")
	}

}

func doApply() string {
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
	return string(cmdOut)
}
