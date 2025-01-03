package tui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ThemeModel struct {
	ThemeName string
	StyleList list.Model
}

type StyleItem string

func (s StyleItem) FilterValue() string { return string(s) }
func (s StyleItem) Title() string       { return string(s) }
func (s StyleItem) Description() string { return "" }

func NewThemeModel(theme string, listWidth, listHeight int) ThemeModel {
	styles := getThemeStyles(theme)
	list := list.New(styles, list.NewDefaultDelegate(), listWidth, listHeight)
	list.Title = "Manage Styles"
	return ThemeModel{ThemeName: theme, StyleList: list}

}

func (m ThemeModel) Init() tea.Cmd {
	return nil
}

func (m ThemeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected := m.StyleList.SelectedItem().(list.DefaultItem).Title()
			return NewStyleEditModel(m.ThemeName, selected, m.StyleList.Width(), m.StyleList.Height()), nil
		case "n":
			return NewStyleCreateModel(m.ThemeName, m.StyleList.Width(), m.StyleList.Height()), nil
		case "esc":
			return NewLandingModel(), nil
		}
	case tea.WindowSizeMsg:
		m.StyleList.SetWidth(msg.Width)
		m.StyleList.SetHeight(msg.Height)
	}
	var cmd tea.Cmd
	m.StyleList, cmd = m.StyleList.Update(msg)
	return m, cmd
}

func (m ThemeModel) View() string {
	return m.StyleList.View()
}

func getThemeStyles(theme string) []list.Item {
	var outFiles []list.Item

	themeDir := filepath.Join(ColorGenConfig, theme)

	os.MkdirAll(themeDir, 0755)

	dir, _ := os.ReadDir(themeDir)

	for _, thing := range dir {
		if strings.HasSuffix(thing.Name(), ".yaml") || strings.HasSuffix(thing.Name(), ".yml") {
			outFiles = append(outFiles, StyleItem(thing.Name()))
		}
	}

	return outFiles
}
