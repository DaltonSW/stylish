package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type LandingModel struct {
	Title     string
	Subtitle  string
	ThemeList list.Model
}

type ThemeItem string

func (t ThemeItem) FilterValue() string { return string(t) }
func (t ThemeItem) Title() string       { return string(t) }
func (t ThemeItem) Description() string { return "" }

var UserConfig, _ = os.UserConfigDir()

var ColorGenConfig = fmt.Sprintf("%v/colorgen", UserConfig)

func NewLandingModel() LandingModel {
	themes := getThemes()
	l := list.New(themes, list.NewDefaultDelegate(), 50, 20)
	l.Title = "Select Theme To Edit"
	return LandingModel{ThemeList: l}

}

func (m LandingModel) Init() tea.Cmd {
	return nil
}

func (m LandingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected := m.ThemeList.SelectedItem().(list.DefaultItem).Title()
			return NewThemeModel(selected, m.ThemeList.Width(), m.ThemeList.Height()), nil
		case "n":
			// TODO: New folder option
		}
	case tea.WindowSizeMsg:
		m.ThemeList.SetWidth(msg.Width)
		m.ThemeList.SetHeight(msg.Height)
	}
	var cmd tea.Cmd
	m.ThemeList, cmd = m.ThemeList.Update(msg)
	return m, cmd
}

func (m LandingModel) View() string {
	// return fmt.Sprintf("%v\n\n%v\n%v\n", m.Title, m.Subtitle, m.ThemeList.View())
	return m.ThemeList.View()
}

func getThemes() []list.Item {
	var themes []list.Item

	os.MkdirAll(ColorGenConfig+"/default", 0755)

	dir, _ := os.ReadDir(ColorGenConfig)

	for _, thing := range dir {
		if thing.IsDir() {
			themes = append(themes, ThemeItem(thing.Name()))
		}
	}

	return themes
}
