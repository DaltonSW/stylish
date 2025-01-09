package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

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
func (s Style) Title() string {
	return lipgloss.PlaceHorizontal(lipgloss.Width(s.Description()), lipgloss.Center, s.getPreview(s.Name))
}

func (s Style) Description() string {
	boxes := s.getCheckboxes()
	topLine := fmt.Sprintf("%v | %v | %v", boxes["Bold"], boxes["Under"], boxes["Blink"])
	// midLine := fmt.Sprintf("Fore: #%v | Back: #%v", s.Fore, s.Back)
	midLine := "Fore: #456123 | Back: #789789"
	botLine := fmt.Sprintf("Filetypes: %v", len(s.FileTypes))
	// return fmt.Sprintf(outStr, checkboxes["Bold"], checkboxes["Under"], checkboxes["Blink"], s.Fore, s.Back, s.getPreview("preview.txt"))
	w := lipgloss.Width(midLine)
	outStr := fmt.Sprintf("%v\n%v\n%v\n", center(topLine, w), center(midLine, w), center(botLine, w))
	// return lipgloss.PlaceHorizontal(lipgloss.Width(midLine), lipgloss.Center, outStr)
	return outStr
}

func center(s string, w int) string {
	return lipgloss.PlaceHorizontal(w, lipgloss.Center, s)
}

func (s Style) FilterValue() string { return s.Name }

func (s Style) getCheckboxes() map[string]string {
	outStr := make(map[string]string)
	if s.Bold {
		outStr["Bold"] = "✓ Bold"
	} else {
		outStr["Bold"] = "  Bold"
	}

	if s.Under {
		outStr["Under"] = "✓ Under"
	} else {
		outStr["Under"] = "  Under"
	}

	if s.Blink {
		outStr["Blink"] = "✓ Blink"
	} else {
		outStr["Blink"] = "  Blink"
	}

	return outStr
}

func (s Style) getPreview(msg string) string {
	var backColor lipgloss.Color
	if s.Back == -1 {
		backColor = lipgloss.Color("")
	} else {
		backColor = lipgloss.Color(strconv.Itoa(s.Back))
	}
	previewColor := lipgloss.NewStyle().
		Foreground(lipgloss.Color(strconv.Itoa(s.Fore))).
		Background(backColor).
		Bold(s.Bold).Underline(s.Under).Blink(s.Blink)

	return previewColor.Render(msg)
}

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
