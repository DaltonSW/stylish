package tui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type ThemeModel struct {
	ThemeName string
	StyleList list.Model
	Viewport  viewport.Model

	NameInput   textinput.Model
	InputActive bool
}

type StyleItem string

func (s StyleItem) FilterValue() string { return string(s) }
func (s StyleItem) Title() string       { return string(s) }
func (s StyleItem) Description() string { return "" }

func NewThemeModel(theme string, listWidth, listHeight int) ThemeModel {
	model := ThemeModel{
		ThemeName: theme,
	}
	styleFiles := model.getStyles()

	var styles []list.Item
	for _, file := range styleFiles {
		styles = append(styles, StyleItem(file))
	}
	list := list.New(styles, list.NewDefaultDelegate(), listWidth, listHeight)
	list.Title = "Manage Styles for " + theme

	input := textinput.New()
	input.Placeholder = "New Style Name"

	model.NameInput = input
	model.StyleList = list

	return model

}

func (m ThemeModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ThemeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.InputActive {
				if m.NameInput.Value() == "" {
					break
				}
				m.NameInput.Blur()
				return NewStyleEditModel(m.ThemeName, m.NameInput.Value(), m.StyleList.Width(), m.StyleList.Height()), nil

			} else {

				return NewStyleEditModel(m.ThemeName, m.StyleList.SelectedItem().(list.DefaultItem).Title(), m.StyleList.Width(), m.StyleList.Height()), nil
			}
		case "n":
			if !m.InputActive {
				m.InputActive = true
				return m, m.NameInput.Focus()
			}
		case "esc", "q":
			return NewLandingModel(m.StyleList.Width(), m.StyleList.Height()), nil
		}
	case tea.WindowSizeMsg:
		m.StyleList.SetWidth(msg.Width - 4)
		m.StyleList.SetHeight(msg.Height - 4)
		ViewportBorder.Width(m.StyleList.Width() + 2).Height(m.StyleList.Height() + 2)
	}
	var cmd tea.Cmd
	if m.InputActive {
		m.NameInput, cmd = m.NameInput.Update(msg)
	} else {

		m.StyleList, cmd = m.StyleList.Update(msg)
	}
	return m, cmd
}

func (m ThemeModel) View() string {
	if m.InputActive {
		return Center(ViewportBorder.Render(m.NameInput.View()))
	} else {

		return Center(ViewportBorder.Render(m.StyleList.View()))
	}
}

func (m ThemeModel) GenerateDirColors() error {

	path := filepath.Join(ColorGenConfig, m.ThemeName)
	file, err := os.Create(filepath.Join(path, ".dircolors"))
	if err != nil {
		return err
	}
	defer file.Close()

	for _, styleName := range m.getStyles() {
		model := NewStyleEditModel(m.ThemeName, styleName, 0, 0)
		file.WriteString(model.GetDirColorBlock())
	}

	return nil

}

func (m ThemeModel) getStyles() []string {
	var outFiles []string

	themeDir := filepath.Join(ColorGenConfig, m.ThemeName)

	os.MkdirAll(themeDir, 0755)

	dir, _ := os.ReadDir(themeDir)

	for _, thing := range dir {
		if strings.HasSuffix(thing.Name(), ".yaml") || strings.HasSuffix(thing.Name(), ".yml") {
			outFiles = append(outFiles, strings.Replace(thing.Name(), ".yaml", "", -1))
		}
	}

	return outFiles
}
