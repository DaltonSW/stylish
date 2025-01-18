package main

// TODO: Look into `goreleaser`
// TODO: Set up a brew tap

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"go.dalton.dog/stylish/cmd"
)

//go:embed themes/default
var DefaultTheme embed.FS

func main() {
	// Check that `dircolors` is installed
	checkCommand("dircolors")

	// Check that `ls` is installed
	checkCommand("ls")

	// Check that stylish dir exists, creating it and copying the default theme into it if needed
	checkConfigDir()

	// If dependencies are fulfilled, kick it off to Cobra
	cmd.Execute()
}

func checkCommand(command string) {
	_, err := exec.LookPath(command)
	if err != nil {
		log.Fatal("Package " + command + " not found on PATH. Please install it before continuing.")
	}
}

func checkConfigDir() {
	userCfgDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Unable to determine user's config directory.")
	}

	stylishDir := filepath.Join(userCfgDir, "stylish")

	_, err = os.Stat(stylishDir)

	if os.IsNotExist(err) { // If the directory doesn't exist...
		err := os.MkdirAll(stylishDir, 0755)
		if err != nil {
			log.Fatal("Unable to create the stylish dir in the user config directory.")
		}

		copyDefaultTheme(filepath.Join(stylishDir, "default"))
	}
}

func copyDefaultTheme(targetDir string) {
	err := fs.WalkDir(DefaultTheme, "themes/default", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(path, "themes/default/")
		targetPath := filepath.Join(targetDir, relPath)

		if d.IsDir() {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
		} else {
			file, err := DefaultTheme.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open embedded file %s: %w", path, err)
			}
			defer file.Close()

			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", targetPath, err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, file); err != nil {
				return fmt.Errorf("failed to copy file %s: %w", path, err)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
