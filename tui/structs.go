package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

func (t *Theme) ReplaceStyle(style Style) {
	for i, s := range t.Styles {
		if s.Name == style.Name {
			t.Styles[i] = style
		}
	}
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

type Style struct {
	Theme string `yaml:"theme"`
	Name  string `yaml:"name"`

	Bold  bool `yaml:"bold"`
	Under bool `yaml:"under"`
	Blink bool `yaml:"blink"`

	Fore int `yaml:"fore"`
	Back int `yaml:"back"`

	FileTypes []string `yaml:"filetypes"`
}

// These functions fullfil the tea.DefaultItemValue interface
func (s Style) FilterValue() string { return s.Name }
func (s Style) Title() string       { return s.Name }
func (s Style) Description() string { return "" }

func NewStyle(themeName, styleName string) Style {
	return Style{
		Theme:     themeName,
		Name:      styleName,
		Fore:      -1,
		Back:      -1,
		FileTypes: make([]string, 0),
	}
}

func GetStyle(theme, styleName string) Style {
	themePath := filepath.Join(ThemeConfigFolder, theme)
	var outStyle Style
	style := loadStyle(themePath, styleName)
	if style == nil {
		outStyle = NewStyle(theme, styleName)
	} else {
		outStyle = *style
	}

	return outStyle
}

func loadStyle(theme, styleName string) *Style {
	file, err := os.Open(filepath.Join(ThemeConfigFolder, theme, styleName+".yaml"))
	if err != nil {
		return nil
	}
	defer file.Close()

	var style *Style
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&style); err != nil {
		return nil
	}

	return style
}

func (s Style) SaveStyle() {
	path := filepath.Join(ThemeConfigFolder, s.Theme)
	file, err := os.Create(filepath.Join(path, s.Name+".yaml"))
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	encoder := yaml.NewEncoder(file)
	err = encoder.Encode(s)
	if err != nil {
		log.Fatal(err)
	}
}

func (s Style) GetDirColorBlock() string {
	outStr := " # " + s.Name + "\n\n"

	styleStr := ""

	if s.Bold {
		styleStr += "1;"
	}

	if s.Under {
		styleStr += "4;"
	}

	if s.Blink {
		styleStr += "5;"
	}

	if s.Fore != -1 {
		styleStr += fmt.Sprintf("38;5;%v;", strconv.Itoa(s.Fore))
	}

	if s.Back != -1 {
		styleStr += fmt.Sprintf("48;5;%v;", strconv.Itoa(s.Back))
	}

	styleStr = strings.TrimSuffix(styleStr, ";")

	for _, file := range s.FileTypes {
		outStr += fmt.Sprintf("%v %v\n", file, styleStr)
	}

	return outStr
}
