package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"
)

type Style struct {
	Theme string `yaml:"theme"`
	Name  string `yaml:"name"`

	Bold  bool `yaml:"bold"`
	Under bool `yaml:"under"`
	Blink bool `yaml:"blink"`

	Fore string `yaml:"fore"`
	Back string `yaml:"back"`

	FileTypes []string `yaml:"filetypes"`
}

// These functions fullfil the tea.DefaultItemValue interface
func (s Style) Title() string {
	return lipgloss.PlaceHorizontal(lipgloss.Width(s.Description()), lipgloss.Center, s.getPreview(s.Name))
}

// 3 Row description
func (s Style) Description() string {
	// return s.threeRowDesc()
	return s.twoColDesc()
}

func (s Style) threeRowDesc() string {
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

func (s Style) twoColDesc() string {
	boxes := s.getCheckboxes()
	var fore string
	var back string
	if s.Fore == "" {
		fore = "DEFAULT"
	} else {
		fore = "#" + s.Fore
	}
	if s.Back == "" {
		back = "DEFAULT"
	} else {
		back = "#" + s.Back
	}
	topLine := fmt.Sprintf("(1) %v | (f) Fore: %v ", boxes["Bold"], fore)
	midLine := fmt.Sprintf("(2) %v | (b) Back: %v ", boxes["Under"], back)
	botLine := fmt.Sprintf("(3) %v | (t) Filetypes: %v", boxes["Blink"], len(s.FileTypes))
	outStr := fmt.Sprintf("%v\n%v\n%v\n", topLine, midLine, botLine)
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
		outStr["Bold"] = "✓ Bold "
	} else {
		outStr["Bold"] = "  Bold "
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
	var foreColor lipgloss.Color
	if s.Back == "" {
		backColor = lipgloss.Color("")
	} else {
		backColor = lipgloss.Color("#" + s.Back)
	}
	if s.Fore == "" {
		foreColor = lipgloss.Color("")
	} else {
		foreColor = lipgloss.Color("#" + s.Fore)
	}

	previewColor := lipgloss.NewStyle().Foreground(foreColor).Background(backColor).
		Bold(s.Bold).Underline(s.Under).Blink(s.Blink)

	return previewColor.Render(msg)
}

func (s *Style) ToggleBold() {
	s.Bold = !s.Bold
}

func (s *Style) ToggleUnder() {
	s.Under = !s.Under
}

func (s *Style) ToggleBlink() {
	s.Blink = !s.Blink
}

func (s *Style) SetFore(fore string) {
	s.Fore = fore
}

func (s *Style) SetBack(back string) {
	s.Back = back
}

func (s *Style) SetFiles(files string) {
	if files == "" {
		s.FileTypes = make([]string, 0)
	} else {
		s.FileTypes = strings.Split(files, "\n")
	}
}

func NewStyle(themeName, styleName string) Style {
	return Style{
		Theme:     themeName,
		Name:      styleName,
		FileTypes: make([]string, 0),
	}
}

func CopyStyle(style Style, newName string) Style {
	newStyle := NewStyle(style.Theme, newName)

	newStyle.Bold = style.Bold
	newStyle.Blink = style.Blink
	newStyle.Under = style.Under
	newStyle.Fore = style.Fore
	newStyle.Back = style.Back
	newStyle.FileTypes = append(newStyle.FileTypes, style.FileTypes...)

	return newStyle
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
	if len(s.FileTypes) < 1 {
		return ""
	}

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

	if s.Fore != "" {
		var fore termenv.Color
		if EightBitMode {
			fore = HexToEightBit(s.Fore)
		} else {
			fore = HexToRGB(s.Fore)
		}
		styleStr += fore.Sequence(false) + ";"
	}

	if s.Back != "" {
		var back termenv.Color
		if EightBitMode {
			back = HexToEightBit(s.Back)
		} else {
			back = HexToRGB(s.Back)
		}
		styleStr += back.Sequence(true) + ";"
	}

	styleStr = strings.TrimSuffix(styleStr, ";")

	for _, file := range s.FileTypes {
		outStr += fmt.Sprintf("%v %v\n", file, styleStr)
	}

	return outStr + "\n"
}
