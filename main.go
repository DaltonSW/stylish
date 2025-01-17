package main

import (
	"os/exec"

	"github.com/charmbracelet/log"
	"go.dalton.dog/stylish/cmd"
)

// TODO: Look into `goreleaser`
// TODO: Set up a brew tap

// TODO: Finish README
//	- Document commands
//	- Add Quick Start
//	- Make it a bit less cluttered?
//	- Add Table of Contents

// TODO: If `stylish` folder doesn't exist, create it and also export the `default` theme
// TODO: Embed default themes
// TODO: `export` command to save any embedded themes into the `stylish` folder
// TODO: `list` command to list out embedded themes

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
