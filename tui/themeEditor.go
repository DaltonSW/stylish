package tui

import (
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

type ThemeModel struct {
	Theme     Theme
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

func NewThemeModel(theme Theme) ThemeModel {
	newHelp := help.New()
	newHelp.ShowAll = true
	newHelp.Width = ConstWidth
	var styles []list.Item
	for _, style := range theme.Styles {
		styles = append(styles, list.Item(style))
	}
	list := list.New(styles, list.NewDefaultDelegate(), ConstWidth, ConstHeight)
	list.Title = "Manage Styles for " + theme.Name
	list.SetStatusBarItemName("style", "styles")
	list.SetShowHelp(false)
	list.SetShowTitle(false)

	input := textinput.New()
	input.Placeholder = "New Style Name"

	return ThemeModel{
		Theme:     theme,
		NameInput: input,
		StyleList: list,
		help:      newHelp,
		keys:      newThemeKeymap(),
	}

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
				return NewStyleEditModel(m.Theme.Name, m.NameInput.Value()), nil

			} else {
				return NewStyleEditModel(m.Theme.Name, m.StyleList.SelectedItem().(list.DefaultItem).Title()), nil
			}
		case "n":
			if !m.InputActive {
				m.InputActive = true
				return m, m.NameInput.Focus()
			}
		case "esc", "q":
			return NewLandingModel(), nil
		}
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
		return RenderModel(m.NameInput.View(), "")
	} else {
		listHeader := CenterHorz(TitleStyle.Render("Theme Styles") + "\n" + SubtitleStyle.Render("Theme: "+m.Theme.Name))
		return RenderModel(listHeader+"\n"+m.StyleList.View(), m.help.View(m.keys))
	}
}

func (m ThemeModel) GenerateDirColors() error {
	path := filepath.Join(ThemeConfigFolder, m.Theme.Name)
	file, err := os.Create(filepath.Join(path, ".dircolors"))
	if err != nil {
		return err
	}
	defer file.Close()

	for _, styleName := range m.getStyles() {
		model := NewStyleEditModel(m.Theme.Name, styleName)
		file.WriteString(model.GetDirColorBlock())
	}

	return nil

}

func (m ThemeModel) getStyles() []string {
	var outFiles []string

	themeDir := filepath.Join(ThemeConfigFolder, m.Theme.Name)

	os.MkdirAll(themeDir, 0755)

	dir, _ := os.ReadDir(themeDir)

	for _, thing := range dir {
		if strings.HasSuffix(thing.Name(), ".yaml") || strings.HasSuffix(thing.Name(), ".yml") {
			outFiles = append(outFiles, strings.Replace(thing.Name(), ".yaml", "", -1))
		}
	}

	return outFiles
}

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
