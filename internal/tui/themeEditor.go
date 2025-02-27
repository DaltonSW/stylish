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
	ColorInput textinput.Model
	FilesInput textarea.Model

	showSystemFileTypes bool

	foreActive   bool
	backActive   bool
	filesActive  bool
	inputActive  bool
	deleteActive bool
	nameActive   bool
	isCopying    bool

	help help.Model
}

func NewThemeModel(theme Theme) ThemeModel {
	newHelp := help.New()
	newHelp.ShowAll = true
	newHelp.Width = ConstWidth - 3
	var styles []list.Item
	for _, style := range theme.Styles {
		log.Debug(style)
		styles = append(styles, list.Item(&style))
	}
	del := GetItemDelgate()
	list := list.New(styles, del, ConstWidth, ConstHeight)
	list.Title = "Manage Styles for " + theme.Name
	list.SetShowStatusBar(false)
	list.SetShowHelp(false)
	list.SetShowTitle(false)
	list.InfiniteScrolling = true

	nameInput := textinput.New()
	nameInput.Placeholder = "New Style Name"

	colorInput := textinput.New()
	colorInput.CharLimit = 6
	colorInput.Prompt = "#"
	colorInput.Validate = ValidHexCode

	fileArea := textarea.New()
	fileArea.Placeholder = ".mp3\n.ogg\n.wav\n.txt"
	fileArea.SetWidth(ConstWidth - 8)
	fileArea.SetHeight(ConstHeight - 10)

	return ThemeModel{
		Theme:      theme,
		ColorInput: colorInput,
		NameInput:  nameInput,
		FilesInput: fileArea,
		StyleList:  list,
		help:       newHelp,
	}

}

func (m ThemeModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, textarea.Blink)
}

func (m ThemeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var style *Style
	if len(m.StyleList.Items()) > 0 {
		style = m.StyleList.SelectedItem().(*Style)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n", "c": // New style
			if !m.isAnythingActive() {
				if msg.String() == "c" {
					m.isCopying = true
				}
				m.NameInput.SetValue("")
				m.nameActive = true
				return m, m.NameInput.Focus()
			} else if m.deleteActive {
				m.deleteActive = false
				return m, nil
			}
		case "d": // Delete style
			if !m.isAnythingActive() {
				m.deleteActive = true
				return m, nil
			}
		case "1": // Toggle Bold
			if !m.isAnythingActive() {
				style.ToggleBold()
				style.SaveStyle()
			}
		case "2": // Toggle Underline
			if !m.isAnythingActive() {
				style.ToggleUnder()
				style.SaveStyle()
			}
		case "3": // Toggle Blinking
			if !m.isAnythingActive() {
				m.StyleList.SelectedItem().(*Style).ToggleBlink()
				style.SaveStyle()
			}
		case "f": // Edit Foreground
			if !m.isAnythingActive() {
				m.ColorInput.SetValue(style.Fore)
				m.foreActive = true
				return m, m.ColorInput.Focus()
			}
		case "b": // Edit Background
			if !m.isAnythingActive() {
				m.ColorInput.SetValue(style.Back)
				m.backActive = true
				return m, m.ColorInput.Focus()
			}
		case "t": // Edit filetypes
			if !m.isAnythingActive() {
				m.FilesInput.SetValue(strings.Join(style.FileTypes, "\n"))
				m.filesActive = true
				return m, m.FilesInput.Focus()
			}
		case "y":
			if m.deleteActive {
				m.StyleList.RemoveItem(m.StyleList.Index())
				m.Theme.RemoveStyle(style.Name)
				m.deleteActive = false
				return m, nil
			}
		case "esc": // Close theme editor
			if !m.isAnythingActive() {
				m.Theme.GenerateDirColors()
				model := NewLandingModel()
				return model, model.Init()
			}
		case "ctrl+h": // Show detailed system filetypes helptext
			if m.filesActive {
				m.showSystemFileTypes = !m.showSystemFileTypes
				return m, textarea.Blink
			}
		case "ctrl+s": // Save and close
			if m.isAnythingActive() {
				if m.nameActive {
					val := m.NameInput.Value()
					if m.Theme.DoesStyleExist(val) {
						m.deactivateInputs()
						m.NameInput.SetValue("")
						return m, nil
					}
					var newStyle Style
					if m.isCopying {
						newStyle = CopyStyle(*style, val)
						newStyle.SaveStyle()

					} else {
						newStyle = NewStyle(m.Theme.Name, val)

					}
					m.Theme.Styles = append(m.Theme.Styles, newStyle)
					m.StyleList.InsertItem(len(m.StyleList.Items()), &newStyle)
					m.StyleList.CursorDown()
					var cmd tea.Cmd
					m.StyleList, cmd = m.StyleList.Update(msg)
					m.deactivateInputs()
					m.NameInput.SetValue("")
					return m, cmd

				}
				if m.backActive {
					style.SetBack(m.ColorInput.Value())
				} else if m.foreActive {
					style.SetFore(m.ColorInput.Value())
				} else if m.filesActive {
					style.SetFiles(m.FilesInput.Value())
				}
				style.SaveStyle()
				m.deactivateInputs()

				return m, nil
			}
		case "ctrl+c": // Cancel and close
			if m.isAnythingActive() {
				m.deactivateInputs()
				return m, nil
			} else {
				m.Theme.GenerateDirColors()
				model := NewLandingModel()
				return model, model.Init()
			}
		case "ctrl+q": // Clear value to default
			if m.foreActive {
				style.SetFore("")
			} else if m.backActive {
				style.SetBack("")
			} else if m.filesActive {
				style.SetFiles("")
			}
			m.deactivateInputs()
			style.SaveStyle()
			return m, nil
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
		return RenderModel(listHeader+"\n"+m.StyleList.View(), m.help.View(themeKeys))
	} else if m.deleteActive {
		return RenderModel(Center(TitleStyle.Render("Delete this style? (y/n)")), "")
	} else if m.foreActive || m.backActive {
		return m.getColorModel()
	} else if m.filesActive {
		return RenderModel(fmt.Sprintf("%v\n\n%v",
			CenterHorz(TitleStyle.Render("Filetypes")), CenterHorz(m.FilesInput.View())), m.getFileAreaHelpText())
	} else if m.nameActive {
		return RenderModel(fmt.Sprintf("%v\n\n%v\n",
			CenterHorz(TitleStyle.Render("New Style")), CenterHorz(m.NameInput.View())), m.getEditHelpTextNoClear())
	}

	return ""
}

func (m ThemeModel) getColorModel() string {
	var titleStr string
	if m.foreActive {
		titleStr = TitleStyle.Render("Foreground Color")
	} else {
		titleStr = TitleStyle.Render("Background Color")
	}

	footerString := ""
	if m.ColorInput.Err != nil {
		footerString = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF1155")).Render(m.ColorInput.Err.Error())
	} else {
		footerString = lipgloss.NewStyle().Foreground(lipgloss.Color("#" + m.ColorInput.Value())).Render(strings.Repeat("█", 18))
	}

	outStr := fmt.Sprintf("%v\n\n%v\n\n%v\n\n%v",
		CenterHorz(titleStr),
		CenterHorz(m.ColorInput.View()),
		CenterHorz(footerString),
		m.getEditHelpText())

	return RenderModel(outStr, "")
}

func (m ThemeModel) isAnythingActive() bool {
	return m.backActive || m.foreActive || m.filesActive || m.nameActive || m.deleteActive
}

func (m *ThemeModel) deactivateInputs() {
	m.backActive = false
	m.foreActive = false
	m.nameActive = false
	m.isCopying = false
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
	New    key.Binding
	Copy   key.Binding
	Filter key.Binding
}

func (k themeKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Quit, k.New, k.Delete}
}

func (k themeKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Quit},
		{k.New, k.Delete, k.Copy, k.Filter},
	}
}

var themeKeys = themeKeymap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "Up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "Down"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("ctrl+c", "Quit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d", "x"),
		key.WithHelp("d", "Delete Style"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "New Style"),
	),
	Copy: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "Copy Style"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "Filter"),
	),
}

func (m ThemeModel) getEditHelpTextNoClear() string {
	keyStyle := m.help.Styles.FullKey
	descStyle := m.help.Styles.FullDesc
	return fmt.Sprintf("%v\n%v",
		CenterHorz(keyStyle.Render("ctrl+s")+descStyle.Render(" [Save]   ")),
		CenterHorz(keyStyle.Render("ctrl+c")+descStyle.Render(" [Discard]")))
}

func (m ThemeModel) getEditHelpText() string {
	keyStyle := m.help.Styles.FullKey
	descStyle := m.help.Styles.FullDesc
	return fmt.Sprintf("%v\n%v\n%v",
		CenterHorz(keyStyle.Render("ctrl+s")+descStyle.Render(" [Save]   ")),
		CenterHorz(keyStyle.Render("ctrl+c")+descStyle.Render(" [Discard]")),
		CenterHorz(keyStyle.Render("ctrl+q")+descStyle.Render(" [Clear]  ")))
}

func (m ThemeModel) getFileAreaHelpText() string {
	keyStyle := m.help.Styles.FullKey
	descStyle := m.help.Styles.FullDesc
	if !m.showSystemFileTypes {
		return fmt.Sprintf("%v\n%v\n%v\n\n%v",
			CenterHorz(keyStyle.Render("ctrl+s")+descStyle.Render(" [Save]   ")),
			CenterHorz(keyStyle.Render("ctrl+c")+descStyle.Render(" [Discard]")),
			CenterHorz(keyStyle.Render("ctrl+q")+descStyle.Render(" [Clear]  ")),
			CenterHorz(keyStyle.Render("ctrl+h")+descStyle.Render(" [Show/Hide System Types]")))
	} else {
		return fmt.Sprintf("%v\n\n%v",
			CenterHorz(keyStyle.Render("ctrl+h")+descStyle.Render(" [Show/Hide System Types]")),
			CenterHorz(m.getInitStrings()))
	}
}

func (m ThemeModel) getInitStrings() string {
	keyStyle := m.help.Styles.FullKey
	descStyle := m.help.Styles.FullDesc
	outStr := CenterHorz(TitleStyle.Render("System File Types")) + "\n"
	outStr += CenterHorz(keyStyle.Render("  FILE")+descStyle.Render(" Normal File             ")) + "\n"
	outStr += CenterHorz(keyStyle.Render("   DIR")+descStyle.Render(" Normal Directory        ")) + "\n"
	outStr += CenterHorz(keyStyle.Render("  LINK")+descStyle.Render(" Symbolic Link           ")) + "\n"
	outStr += CenterHorz(keyStyle.Render("  FIFO")+descStyle.Render(" Pipe                    ")) + "\n"
	outStr += CenterHorz(keyStyle.Render("  EXEC")+descStyle.Render(" Execute permissions     ")) + "\n"
	outStr += CenterHorz(keyStyle.Render("ORPHAN")+descStyle.Render(" Sym Link to Missing File")) + "\n"
	outStr += CenterHorz(keyStyle.Render("SETUID")+descStyle.Render(" File w/ u+s             ")) + "\n"
	outStr += CenterHorz(keyStyle.Render("SETGID")+descStyle.Render(" File w/ g+s             ")) + "\n"
	outStr += CenterHorz(keyStyle.Render("STICKY")+descStyle.Render(" Dir w/ +t, no o or w    ")) + "\n"
	// outStr += CenterHorz(keyStyle.Render("MISSING")+descStyle.Render(" Missing Files")) + "\n"
	// outStr += CenterHorz(keyStyle.Render("CAPABILITY")+descStyle.Render(" File w/ capability")) + "\n"
	// outStr += CenterHorz(keyStyle.Render("MULTIHARDLINK")+descStyle.Render(" File w/ >1 Link")) + "\n"
	// outStr += CenterHorz(keyStyle.Render("STICKY_OTHER_WRITABLE")+descStyle.Render(" Dir w/ +t,o+w")) + "\n"
	// outStr += CenterHorz(keyStyle.Render("OTHER_WRITABLE")+descStyle.Render(" Dir w/ o+w, no t")) + "\n"

	return outStr
}
