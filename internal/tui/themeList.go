package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type LandingModel struct {
	ThemeList    list.Model
	ThemeInput   textinput.Model
	InputActive  bool
	isCopying    bool
	themeToCopy  string
	DeleteActive bool

	keys landingKeymap
	help help.Model
}

func NewLandingModel() LandingModel {
	log.Debug("Trying to create landing model")

	themes := GetAllThemes()
	var items []list.Item
	for _, t := range themes {
		items = append(items, list.Item(t))
	}

	l := list.New(items, list.NewDefaultDelegate(), ConstWidth, ConstHeight)
	l.SetStatusBarItemName("theme", "themes")
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowHelp(false)

	themeInput := textinput.New()
	themeInput.Placeholder = "Theme Name"

	newHelp := help.New()
	newHelp.ShowAll = true
	newHelp.Width = ConstWidth - 2

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
				name := m.ThemeInput.Value()
				if m.isCopying {
					srcDir := filepath.Join(ThemeConfigFolder, m.themeToCopy)
					destDir := filepath.Join(ThemeConfigFolder, name)

					err := os.CopyFS(destDir, os.DirFS(srcDir))
					if err != nil {
						log.Fatal(err)
					}
					m.isCopying = false
					m.themeToCopy = ""
				}
				m.ThemeInput.Blur()
				return NewThemeModel(GetTheme(name)), nil

			} else {
				selected := m.ThemeList.SelectedItem().(Theme)
				return NewThemeModel(selected), nil
			}
		case "d":
			if !m.InputActive && !m.DeleteActive {
				m.DeleteActive = true
			}
		case "y":
			if m.DeleteActive {
				selected := m.ThemeList.SelectedItem().(Theme)
				m.deleteTheme(selected.Name)
				m.DeleteActive = false
				m.ThemeList.RemoveItem(m.ThemeList.Index())
			}

		case "g":
			selected := m.ThemeList.SelectedItem().(Theme)
			err := selected.GenerateDirColors()
			if err != nil {
				log.Fatal(err)
			}

			return m, nil
		case "n", "c":
			if !m.InputActive && !m.DeleteActive {
				m.InputActive = true
				if msg.String() == "c" {
					m.isCopying = true
					m.themeToCopy = m.ThemeList.SelectedItem().(Theme).Name
				}
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
		return RenderModel(Center(fmt.Sprintf("%v\n%v", TitleStyle.Render("New Theme Name"), m.ThemeInput.View())), "")
	} else if m.DeleteActive {
		return RenderModel(Center(TitleStyle.Render("Delete this theme? (y/n)")), "")
	} else {
		listHeader := CenterHorz(TitleStyle.Render("Current Themes") + "\n" + SubtitleStyle.Render(ThemeConfigFolder))
		return RenderModel(listHeader+"\n"+m.ThemeList.View(), m.help.View(m.keys))
	}
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

type landingKeymap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding

	Delete key.Binding
	New    key.Binding
	Copy   key.Binding
	Filter key.Binding
}

func (k landingKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Quit, k.New, k.Delete, k.Copy}
}

func (k landingKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Quit, k.Select},
		{k.New, k.Delete, k.Copy, k.Filter},
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
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "New Theme"),
		),
		Copy: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "Copy Theme"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "Filter"),
		),
	}
}
