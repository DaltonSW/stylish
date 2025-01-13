package main

import (
	"os/exec"

	"github.com/charmbracelet/log"
	"go.dalton.dog/stylish/cmd"
)

// TODO: Look into `goreleaser`

// TODO: Finish README
// TODO: Update helptext
// TODO: Clean up and comment code

// TODO: `c` for `copy`

func main() {
	// Check that `dircolors` is installed
	checkCommand("dircolors")

	// Check that `ls` is installed
	checkCommand("ls")

	// If dependencies are fulfilled, kick it off to Cobra
	cmd.Execute()
}

func checkCommand(command string) {
	_, err := exec.LookPath(command)
	if err != nil {
		log.Fatal("Package " + command + " not found on PATH. Please install it before continuing.")
	}
}
