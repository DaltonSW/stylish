package tui

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type landingKeymap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding

	Delete   key.Binding
	Create   key.Binding
	Generate key.Binding
	Filter   key.Binding
}

func (k landingKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Quit, k.Create, k.Delete, k.Generate}
}

func (k landingKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select, k.Quit},
		{k.Create, k.Delete, k.Generate, k.Filter},
	}
}

func newLandingKeymap() landingKeymap {
	return landingKeymap{
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
			key.WithHelp("d", "Delete Theme"),
		),
		Create: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "New Theme"),
		),
		Generate: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "Generate File"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "Filter"),
		),
	}
}

type LandingModel struct {
	Title        string
	Subtitle     string
	ThemeList    list.Model
	ThemeInput   textinput.Model
	InputActive  bool
	DeleteActive bool

	keys landingKeymap
	help help.Model
}

type ThemeItem string

func (t ThemeItem) FilterValue() string { return string(t) }
func (t ThemeItem) Title() string       { return string(t) }
func (t ThemeItem) Description() string { return "" }

var UserConfig, _ = os.UserConfigDir()

var ThemeConfigFolder = fmt.Sprintf("%v/stylish", UserConfig)

func NewLandingModel() LandingModel {
	themeNames := GetAllThemeNames()
	themeList := make([]list.Item, 0, len(themeNames))
	for _, name := range themeNames {
		if name == "" {
			continue
		}
		themeList = append(themeList, list.Item(GetTheme(name)))
	}
	l := list.New(themeList, list.NewDefaultDelegate(), ConstWidth, ConstHeight)
	l.SetStatusBarItemName("theme", "themes")
	l.Title = "Select Theme To Edit"
	l.SetShowTitle(false)
	l.SetShowHelp(false)

	themeInput := textinput.New()
	themeInput.Placeholder = "Theme Name"

	newHelp := help.New()
	newHelp.ShowAll = true
	newHelp.Width = ConstWidth

	return LandingModel{
		ThemeList:  l,
		ThemeInput: themeInput,

		keys: newLandingKeymap(),
		help: newHelp,
	}

}

func (m LandingModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m LandingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.InputActive {
				m.createTheme(m.ThemeInput.Value())
				m.ThemeInput.Blur()
				return NewThemeModel(m.ThemeInput.Value()), nil

			} else {
				selected := m.ThemeList.SelectedItem().(list.DefaultItem).Title()
				return NewThemeModel(selected), nil
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
			newModel := NewThemeModel(selected)
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
		return Center(ViewportBorder.Render("Delete this theme? (y/n)"))
	} else {
		listHeader := CenterHorz(TitleStyle.Render("Current Themes") + "\n" + SubtitleStyle.Render(ThemeConfigFolder))
		return RenderModel(listHeader+"\n"+m.ThemeList.View(), m.help.View(m.keys))
	}
}

func getThemes() []list.Item {
	var themes []list.Item

	os.MkdirAll(ThemeConfigFolder+"/default", 0755)

	dir, _ := os.ReadDir(ThemeConfigFolder)

	for _, thing := range dir {
		if thing.IsDir() {
			themes = append(themes, ThemeItem(thing.Name()))
		}
	}

	return themes
}

func (m *LandingModel) createTheme(name string) {
	path := filepath.Join(ThemeConfigFolder, name)
	if err := os.MkdirAll(path, 0755); err != nil {
		fmt.Printf("Error creating theme folder: %v\n", err)
	}
}

func (m *LandingModel) deleteTheme(name string) {
	path := filepath.Join(ThemeConfigFolder, name)
	if err := os.RemoveAll(path); err != nil {
		fmt.Printf("Error deleting theme folder: %v\n", err)
	}
}
