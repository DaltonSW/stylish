package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"go.dalton.dog/stylish/internal/tui"
)

func init() {
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(applyEightBitCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Turns the theme's yaml files into a dircolors compatible format",
	Long: `Takes the theme's yaml files and turns it 
	into a dircolors compatible format. Should be used
	with eval in your shell's init script.`,
	Example: "eval $(stylish apply <theme>)",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(doApply(args[0]))
	},
}
var applyEightBitCmd = &cobra.Command{

	Use:     "apply-eightbit",
	Short:   "apply command, but constrained to 8-bit colors",
	Example: "eval $(stylish apply-eightbit <theme>)",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tui.EightBitMode = true
		fmt.Print(doApply(args[0]))
	},
}

func doApply(themeName string) string {
	theme := tui.GetTheme(themeName)
	err := theme.GenerateDirColors()
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("dircolors", filepath.Join(theme.Path, ".dircolors"))
	cmdOut, cmdErr := cmd.Output()
	if cmdErr != nil {
		log.Fatal(cmdErr.Error() + string(cmdOut))
	}
	return string(cmdOut)
}
