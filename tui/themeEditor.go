package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type ThemeModel struct {
	Theme     Theme
	StyleList list.Model

	NameInput  textinput.Model
	nameActive bool

	ColorInput textinput.Model
	foreActive bool
	backActive bool

	FilesInput  textarea.Model
	filesActive bool

	InputActive bool

	keys themeKeymap
	help help.Model
}

func NewThemeModel(theme Theme) ThemeModel {
	newHelp := help.New()
	newHelp.ShowAll = true
	newHelp.Width = ConstWidth
	var styles []list.Item
	for _, style := range theme.Styles {
		log.Debug(style)
		styles = append(styles, list.Item(&style))
	}
	del := GetItemDelgate()
	list := list.New(styles, del, ConstWidth, ConstHeight)
	list.Title = "Manage Styles for " + theme.Name
	list.SetStatusBarItemName("style", "styles")
	list.SetShowStatusBar(false)
	list.SetShowHelp(false)
	list.SetShowTitle(false)

	nameInput := textinput.New()
	nameInput.Placeholder = "New Style Name"

	colorInput := textinput.New()
	colorInput.Placeholder = "FFFFFF"
	colorInput.CharLimit = 6
	colorInput.Prompt = "#"
	colorInput.Validate = ValidHexCode

	fileArea := textarea.New()
	fileArea.Placeholder = ".mp3\n.ogg\n.wav\n.txt"

	return ThemeModel{
		Theme:      theme,
		ColorInput: colorInput,
		NameInput:  nameInput,
		FilesInput: fileArea,
		StyleList:  list,
		help:       newHelp,
		keys:       newThemeKeymap(),
	}

}

// (1) Bold  x | (f) Fore: #FFFFFF
// (2) Under x | (b) Back: #000000
// (3) Blink x | (t) Types: 18

func (m ThemeModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, textarea.Blink)
}

func (m ThemeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			if !m.InputActive {
				m.InputActive = true
				return m, m.NameInput.Focus()
			}
		case "1":
			if !m.isAnythingActive() {
				m.StyleList.SelectedItem().(*Style).ToggleBold()
			}
		case "2":
			if !m.isAnythingActive() {
				m.StyleList.SelectedItem().(*Style).ToggleUnder()
			}
		case "3":
			if !m.isAnythingActive() {
				m.StyleList.SelectedItem().(*Style).ToggleBlink()
			}
		case "f":
			if !m.isAnythingActive() {
				m.foreActive = true
				return m, m.ColorInput.Focus()
			}
		case "b":
			if !m.isAnythingActive() {
				m.backActive = true
				return m, m.ColorInput.Focus()
			}
		case "t":
			if !m.isAnythingActive() {
				m.filesActive = true
				return m, m.FilesInput.Focus()
			}
		case "esc":
			if !m.isAnythingActive() {
				m.Theme.GenerateDirColors()
				return NewLandingModel(), nil
			} else {
				m.deactivateInputs()
				return m, nil
			}
		case "ctrl+s":
			if m.isAnythingActive() {
				if m.backActive {
					m.StyleList.SelectedItem().(*Style).SetBack(m.ColorInput.Value())
				} else if m.foreActive {
					m.StyleList.SelectedItem().(*Style).SetFore(m.ColorInput.Value())
				} else if m.filesActive {
					m.StyleList.SelectedItem().(*Style).SetFiles(m.FilesInput.Value())
				}
			}
		}
	}
	var cmd tea.Cmd
	if m.isAnythingActive() {
		m.NameInput, cmd = m.NameInput.Update(msg)
		m.FilesInput, cmd = m.FilesInput.Update(msg)
		m.ColorInput, cmd = m.ColorInput.Update(msg)
	} else {
		m.StyleList, cmd = m.StyleList.Update(msg)
	}
	return m, cmd
}

func (m ThemeModel) View() string {
	if !m.isAnythingActive() {
		listHeader := CenterHorz(TitleStyle.Render("Theme Styles") + "\n" + SubtitleStyle.Render("Theme: "+m.Theme.Name))
		return RenderModel(listHeader+"\n"+m.StyleList.View(), m.help.View(m.keys))
	}

	if m.foreActive {
		footerString := ""
		if m.ColorInput.Err != nil {
			footerString = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF1155")).Render(m.ColorInput.Err.Error())
		} else {
			footerString = lipgloss.NewStyle().Foreground(lipgloss.Color("#" + m.ColorInput.Value())).Render(strings.Repeat("█", 18))
		}

		outStr := fmt.Sprintf("%v\n%v",
			CenterHorz(TitleStyle.Render("Foreground Color")),
			CenterHorz(m.ColorInput.View()))

		return RenderModel(outStr, footerString)
	} else if m.backActive {
		footerString := ""
		if m.ColorInput.Err != nil {
			footerString = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF1155")).Render(m.ColorInput.Err.Error())
		} else {
			footerString = lipgloss.NewStyle().Foreground(lipgloss.Color("#" + m.ColorInput.Value())).Render(strings.Repeat("█", 18))
		}

		outStr := fmt.Sprintf("%v\n%v",
			CenterHorz(TitleStyle.Render("Background Color")),
			CenterHorz(m.ColorInput.View()))
		return RenderModel(outStr, footerString)
	} else if m.filesActive {
		return RenderModel(fmt.Sprintf("%v\n%v", "Filetypes", m.FilesInput.View()), "")
	}

	return ""
}

func (m ThemeModel) isAnythingActive() bool {
	return m.backActive || m.foreActive || m.filesActive || m.nameActive
}

func (m *ThemeModel) deactivateInputs() {
	m.backActive = false
	m.foreActive = false
	m.nameActive = false
	m.filesActive = false

	m.ColorInput.Blur()
	m.FilesInput.Blur()
	m.NameInput.Blur()
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
