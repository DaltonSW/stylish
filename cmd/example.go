package cmd

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"go.dalton.dog/stylish/internal/tui"
)

func init() {
	rootCmd.AddCommand(exampleCmd)
}

var exampleCmd = &cobra.Command{
	Use:     "example",
	Short:   "Generates an example directory with dummy files to showcase your theme.",
	Example: "stylish example <theme>",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		binary, pathErr := exec.LookPath("tree")
		if pathErr != nil {
			log.Fatal(errors.New("Path err: " + pathErr.Error()))
		}

		oldColors := os.Getenv("LS_COLORS")
		defer os.Setenv("LS_COLORS", oldColors)

		var themeName string
		if len(os.Args) < 3 {
			themeName = "default"
		} else {
			themeName = os.Args[2]
		}

		theme := tui.GetTheme(themeName)

		createThemeExampleDir(theme)

		output := doApply(themeName)
		output = strings.TrimPrefix(output, "LS_COLORS='")
		output = strings.TrimSpace(output)
		output = strings.TrimSuffix(output, "';\nexport LS_COLORS")

		os.Setenv("LS_COLORS", output)

		execErr := syscall.Exec(binary, []string{"tree", filepath.Join(theme.Path, "example")}, os.Environ())
		if execErr != nil {
			log.Fatal(errors.New("Exec err: " + execErr.Error()))
		}
	},
}

func createThemeExampleDir(theme tui.Theme) {
	outputDir := filepath.Join(theme.Path, "example")
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	for _, style := range theme.Styles {
		styleDir := filepath.Join(outputDir, style.Name)
		err := os.MkdirAll(styleDir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		for i, fileType := range style.FileTypes {
			if i >= 3 {
				break
			}
			if fileType[0] != '.' {
				continue
			}
			// Check if a file with this filetype already exists
			exists, err := fileTypeExistsInDir(fileType, styleDir)
			if err != nil {
				log.Fatal(err)
			}
			if exists {
				continue
			}
			randFileName := fmt.Sprintf("%v%v", randomFileName(), fileType)
			fileName := filepath.Join(styleDir, randFileName)
			file, err := os.Create(fileName)
			if err != nil {
				log.Fatal(err)
			}
			file.Close()
		}
	}
}

func fileTypeExistsInDir(fileType, dirPath string) (bool, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			if filepath.Ext(entry.Name()) == fileType {
				return true, nil
			}
		}
	}
	return false, nil
}

func randomFileName() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	n := rand.Intn(16) + 1 // Random length between 1 and 16
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
