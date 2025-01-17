package cmd

import (
	"go.dalton.dog/stylish/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stylish",
	Short: "stylish is a simple and intuitive path to a prettier ls experience",
	Run: func(cmd *cobra.Command, args []string) {
		program := tea.NewProgram(tui.NewLandingModel(), tea.WithAltScreen())
		if _, err := program.Run(); err != nil {
			log.Fatalf("Error running program: %v\n", err)
		}
	},
}

func Execute() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
