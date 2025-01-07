package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type themeKeymap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding

	Delete key.Binding
	Create key.Binding
	Filter key.Binding
}

func (k themeKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Quit, k.Create, k.Delete}
}

func (k themeKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select, k.Quit},
		{k.Create, k.Delete, k.Filter}, //, k.Quit},
		// {k.Quit, k.Filter},
	}
}

func newThemeKeymap() themeKeymap {
	return themeKeymap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/↑", "Up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/↓", "Down"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Select"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("ctrl+c", "Quit"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d", "x"),
			key.WithHelp("d", "Delete Style"),
		),
		Create: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "New Style"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "Filter"),
		),
	}
}

type ThemeModel struct {
	ThemeName string
	StyleList list.Model
	Viewport  viewport.Model

	NameInput   textinput.Model
	InputActive bool

	keys themeKeymap
	help help.Model
}

type StyleItem string

func (s StyleItem) FilterValue() string { return string(s) }
func (s StyleItem) Title() string       { return string(s) }
func (s StyleItem) Description() string { return "" }

func NewThemeModel(theme string) ThemeModel {
	newHelp := help.New()
	newHelp.ShowAll = true
	model := ThemeModel{
		ThemeName: theme,
		help:      newHelp,
		keys:      newThemeKeymap(),
	}
	styleFiles := model.getStyles()

	var styles []list.Item
	for _, file := range styleFiles {
		styles = append(styles, StyleItem(file))
	}
	list := list.New(styles, list.NewDefaultDelegate(), ConstWidth, ConstHeight)
	list.Title = "Manage Styles for " + theme
	list.SetStatusBarItemName("style", "styles")
	list.SetShowHelp(false)

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
				return NewStyleEditModel(m.ThemeName, m.NameInput.Value()), nil

			} else {

				return NewStyleEditModel(m.ThemeName, m.StyleList.SelectedItem().(list.DefaultItem).Title()), nil
			}
		case "n":
			if !m.InputActive {
				m.InputActive = true
				return m, m.NameInput.Focus()
			}
		case "esc", "q":
			// return NewLandingModel(m.StyleList.Width(), m.StyleList.Height()), nil
			return NewLandingModel(), nil
		}
		// case tea.WindowSizeMsg:
		// m.StyleList.SetWidth(msg.Width - 4)
		// m.StyleList.SetHeight(msg.Height - 4)
		// ViewportBorder.Width(m.StyleList.Width() + 2).Height(m.StyleList.Height() + 2)
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
		return m.NameInput.View()
		// return Center(ViewportBorder.Render(m.NameInput.View()))
	} else {
		return fmt.Sprintf("%v\n%v", ProgramHeader(), ViewportBorder.Render(m.StyleList.View()+"\n"+CenterHorz(m.help.View(m.keys))))
		// return ProgramHeader() + "\n" + ViewportBorder.Render(m.StyleList.View()+"\n"+m.help.View(m.keys))
		// return Center(ViewportBorder.Render(m.StyleList.View()))
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
		model := NewStyleEditModel(m.ThemeName, styleName)
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
