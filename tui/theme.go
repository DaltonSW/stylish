package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/charmbracelet/log"
)

var UserConfig, _ = os.UserConfigDir()
var ThemeConfigFolder = filepath.Join(UserConfig, "stylish")

// Theme represents a collection of Styles
type Theme struct {
	Name   string
	Path   string
	Styles []Style
}

// These functions fullfil the tea.DefaultItemValue interface
func (t Theme) FilterValue() string { return t.Name }
func (t Theme) Title() string       { return t.Name }
func (t Theme) Description() string { return fmt.Sprintf("Styles loaded: %v", len(t.Styles)) }

// GetAllThemes will return a slice containing all Themes
func GetAllThemes() []Theme {
	log.Debug("Trying to get all themes\n")
	var outThemes []Theme

	dir, _ := os.ReadDir(ThemeConfigFolder)

	for _, thing := range dir {
		if thing.IsDir() {
			log.Debugf("Dir found %v\n", thing.Name())
			outThemes = append(outThemes, GetTheme(thing.Name()))
		}
	}

	return outThemes
}

// GetTheme will get the theme of a given name. If the provided name
// doesn't exist, a folder for that theme will be created.
func GetTheme(name string) Theme {
	if name == "" {
		log.Fatal("Tried to create a theme with an empty name.")
	}

	outTheme := Theme{
		Name: name,
		Path: filepath.Join(ThemeConfigFolder, name),
	}

	if _, err := os.Stat(outTheme.Path); os.IsNotExist(err) {
		os.Mkdir(outTheme.Path, 0755)
	}

	log.Debugf("Theme created: %v", name)

	outTheme.Styles = outTheme.LoadStyles()

	return outTheme
}

// LoadStyles will load all of the styles for a given theme
func (t Theme) LoadStyles() []Style {
	var outStyles []Style

	log.Debugf("Trying to load styles for %v", t.Name)
	dir, err := os.ReadDir(t.Path)

	if err != nil {
		log.Fatal(err)
	}

	for _, thing := range dir {
		log.Debugf("- Thing found: %v", thing.Name())
		if !thing.IsDir() && strings.HasSuffix(thing.Name(), ".yaml") {
			styleFile, err := os.Open(filepath.Join(t.Path, thing.Name()))
			if err != nil {
				log.Fatal(err)
			}
			defer styleFile.Close()

			var style Style
			name := strings.TrimSuffix(thing.Name(), ".yaml")

			fileStat, _ := styleFile.Stat()
			if fileStat.Size() == 0 {
				style = NewStyle(t.Name, name)
			} else {
				var outStyle *Style
				decoder := yaml.NewDecoder(styleFile)
				if err := decoder.Decode(&outStyle); err != nil {
					log.Fatal(err)
				}
				if outStyle == nil {
					style = NewStyle(t.Name, name)
				} else {
					style = *outStyle
				}
			}

			outStyles = append(outStyles, style)
		}
	}

	return outStyles
}

// RemoveStyle will remove the style with a given name from both the theme's list and from the file system
func (t *Theme) RemoveStyle(styleName string) {
	newStyles := make([]Style, 0)
	for _, s := range t.Styles {
		if s.Name != styleName {
			newStyles = append(newStyles, s)
		}
	}
	t.Styles = newStyles

	path := filepath.Join(ThemeConfigFolder, t.Name, styleName+".yaml")
	os.Remove(path)
}

// GenerateDirColors will convert all of a theme's styles into an output file
func (t Theme) GenerateDirColors() error {

	path := filepath.Join(ThemeConfigFolder, t.Name)
	file, err := os.Create(filepath.Join(path, ".dircolors"))
	if err != nil {
		return err
	}
	defer file.Close()

	for _, style := range t.Styles {
		file.WriteString(style.GetDirColorBlock())
	}

	return nil
}
