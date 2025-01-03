package tui

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type LandingModel struct {
	Title        string
	Subtitle     string
	ThemeList    list.Model
	ThemeInput   textinput.Model
	InputActive  bool
	DeleteActive bool
}

type ThemeItem string

func (t ThemeItem) FilterValue() string { return string(t) }
func (t ThemeItem) Title() string       { return string(t) }
func (t ThemeItem) Description() string { return "" }

var UserConfig, _ = os.UserConfigDir()

var ColorGenConfig = fmt.Sprintf("%v/colorgen", UserConfig)

func NewLandingModel(width, height int) LandingModel {
	ViewportBorder.Width(width).Height(height)
	themeNames := GetAllThemeNames()
	themeList := make([]list.Item, 0, len(themeNames))
	for _, name := range themeNames {
		if name == "" {
			continue
		}
		themeList = append(themeList, list.Item(GetTheme(name)))
	}
	l := list.New(themeList, list.NewDefaultDelegate(), width, height)
	l.Title = "Select Theme To Edit"

	themeInput := textinput.New()
	themeInput.Placeholder = "Theme Name"

	return LandingModel{
		Title:      "ColorGen",
		Subtitle:   "Put the glam in your term",
		ThemeList:  l,
		ThemeInput: themeInput,
	}

}

func (m LandingModel) Init() tea.Cmd {
	return nil
}

func (m LandingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.InputActive {
				m.createTheme(m.ThemeInput.Value())
				m.ThemeInput.Blur()
				return NewThemeModel(m.ThemeInput.Value(), m.ThemeList.Width(), m.ThemeList.Height()), nil

			} else {
				selected := m.ThemeList.SelectedItem().(list.DefaultItem).Title()
				return NewThemeModel(selected, m.ThemeList.Width(), m.ThemeList.Height()), nil
			}
		case "d":
			if !m.InputActive && !m.DeleteActive {
				m.DeleteActive = true
			}
		case "y":
			if m.DeleteActive {
				selected := m.ThemeList.SelectedItem().(ThemeItem)
				m.deleteTheme(string(selected))
				m.DeleteActive = false
				m.ThemeList.RemoveItem(m.ThemeList.Index())
			}

		case "g":
			selected := m.ThemeList.SelectedItem().(list.DefaultItem).Title()
			newModel := NewThemeModel(selected, m.ThemeList.Width(), m.ThemeList.Height())
			err := newModel.GenerateDirColors()
			if err != nil {
				log.Fatal(err)
			}

			return m, nil
		case "n":
			if !m.InputActive && !m.DeleteActive {
				m.InputActive = true
				return m, m.ThemeInput.Focus()
			}
			if m.DeleteActive {
				m.DeleteActive = false
			}
		case "esc":
			if m.InputActive {
				m.InputActive = false
				m.ThemeInput.Blur()
				m.ThemeInput.SetValue("")
			} else {
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.ThemeList.SetWidth(msg.Width / 2)
		m.ThemeList.SetHeight(msg.Height / 2)
		ViewportBorder.Width(m.ThemeList.Width() + 2).Height(m.ThemeList.Height() + 2)
	}
	var cmd tea.Cmd
	if m.InputActive {
		m.ThemeInput, cmd = m.ThemeInput.Update(msg)
	} else {
		m.ThemeList, cmd = m.ThemeList.Update(msg)
	}
	return m, cmd
}

func (m LandingModel) View() string {
	if m.InputActive {
		return Center(ViewportBorder.Render(m.ThemeInput.View()))
	} else if m.DeleteActive {
		return ViewportBorder.Render("Delete this theme? (y/n)")
	} else {
		return Center(fmt.Sprintf("%v\n%v", m.getHeader(), ViewportBorder.Render(m.ThemeList.View())))
	}
}

func (m LandingModel) getHeader() string {
	return fmt.Sprintf("%v\n\n%v\n", m.Title, m.Subtitle)
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

func (m *LandingModel) createTheme(name string) {
	path := filepath.Join(ColorGenConfig, name)
	if err := os.MkdirAll(path, 0755); err != nil {
		fmt.Printf("Error creating theme folder: %v\n", err)
	}
}

func (m *LandingModel) deleteTheme(name string) {
	path := filepath.Join(ColorGenConfig, name)
	if err := os.RemoveAll(path); err != nil {
		fmt.Printf("Error deleting theme folder: %v\n", err)
	}
}
