package cmd

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(previewCmd)
}

var previewCmd = &cobra.Command{
	Use:     "preview",
	Short:   "Shows how your current directory would look with the given theme.",
	Example: "stylish preview <theme>",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		oldColors := os.Getenv("LS_COLORS")
		defer os.Setenv("LS_COLORS", oldColors)

		output := doApply(args[0])
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
	},
}
