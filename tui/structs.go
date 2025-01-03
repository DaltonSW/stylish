package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Theme struct {
	Name   string
	Path   string
	Styles []*Style
}

func (t Theme) FilterValue() string { return t.Name }
func (t Theme) Title() string       { return t.Name }
func (t Theme) Description() string { return "" }

func GetAllThemeNames() []string {
	var outNames []string

	dir, _ := os.ReadDir(ColorGenConfig)

	for _, thing := range dir {
		if thing.IsDir() {
			outNames = append(outNames, thing.Name())
		}
	}

	return outNames

}

func GetTheme(name string) *Theme {
	if name == "" {
		return nil
	}

	outTheme := &Theme{
		Name: name,
		Path: filepath.Join(ColorGenConfig, name),
	}

	if _, err := os.Stat(outTheme.Path); os.IsNotExist(err) {
		os.Mkdir(outTheme.Path, 0755)
	}

	outTheme.LoadStyles()

	return outTheme
}

func (t *Theme) LoadStyles() {
	dir, _ := os.ReadDir(t.Path)

	for _, thing := range dir {
		if thing.IsDir() {
			t.Styles = append(t.Styles, GetStyle(t.Name, thing.Name()))
		}
	}
}

func (t *Theme) GenerateTheme() error {

	path := filepath.Join(ColorGenConfig, t.Name)
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

func GetStyle(themeName, styleName string) *Style {
	style := loadStyle(themeName, styleName)
	if style == nil {
		style = &Style{
			Theme:     themeName,
			Name:      styleName,
			Fore:      -1,
			Back:      -1,
			FileTypes: make([]string, 0),
		}
	}

	return style
}

func loadStyle(theme, styleName string) *Style {
	file, err := os.Open(filepath.Join(ColorGenConfig, theme, styleName+".yaml"))
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

func SaveStyle(style Style) error {
	path := filepath.Join(ColorGenConfig, style.Theme)
	file, err := os.Create(filepath.Join(path, style.Name+".yaml"))
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := yaml.NewEncoder(file)
	return encoder.Encode(style)
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
