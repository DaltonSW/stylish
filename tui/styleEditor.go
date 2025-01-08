package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle
)

var ControlOrder []string = []string{
	"Bold",
	"Under",
	"Blink",
	"Fore",
	"Back",
	"Files",
	"Save",
	"Discard",
}

type StyleModel struct {
	theme Theme
	style Style

	NameInput textinput.Model

	foreSlider progress.Model
	backSlider progress.Model

	FileArea textarea.Model

	Focused int

	keymap styleKeymap
	help   help.Model
}

func NewStyleEditModel(theme Theme, style Style) StyleModel {
	helpModel := help.New()
	helpModel.ShowAll = true

	foreSlider := progress.New()
	foreSlider.ShowPercentage = false
	foreSlider.Width = ConstWidth - 6
	backSlider := progress.New()
	backSlider.ShowPercentage = false
	backSlider.Width = ConstWidth - 6

	nameInput := textinput.New()
	nameInput.Placeholder = "Style Name"
	nameInput.SetValue(style.Name)
	nameInput.PromptStyle = focusedStyle
	nameInput.TextStyle = focusedStyle
	nameInput.Prompt = ""

	fileArea := textarea.New()
	if len(style.FileTypes) > 1 {
		fileArea.SetValue(strings.Join(style.FileTypes, "\n"))
	} else {
		fileArea.Placeholder = ".mp3\n.gif\n.docx\n..."
	}
	fileArea.Blur()

	return StyleModel{
		theme:      theme,
		style:      style,
		foreSlider: foreSlider,
		backSlider: backSlider,
		FileArea:   fileArea,
		NameInput:  nameInput,
		keymap:     newStyleKeymap(),
		help:       helpModel,
	}
}

func (m StyleModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, textarea.Blink)
}

func (m StyleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var outCmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return NewThemeModel(m.theme), nil
		case "tab", "down", "shift+tab", "up":
			if msg.String() == "tab" || msg.String() == "down" {

				m.Focused++
				if m.Focused >= len(ControlOrder) {
					m.Focused = 0
				}
			} else {
				m.Focused--
				if m.Focused < 0 {
					m.Focused = len(ControlOrder) - 1
				}
			}

			if ControlOrder[m.Focused] == "Files" {
				outCmds = append(outCmds, m.FileArea.Focus())
			} else {
				m.FileArea.Blur()
			}
		case "space", "enter":
			// Toggle checkboxes if focused
			switch ControlOrder[m.Focused] {
			case "Bold":
				m.style.Bold = !m.style.Bold
			case "Under":
				m.style.Under = !m.style.Under
			case "Blink":
				m.style.Blink = !m.style.Blink
			case "Save":
				m.style.SaveStyle()
				m.theme.ReplaceStyle(m.style)
				return NewThemeModel(m.theme), nil
			case "Discard":
				return NewThemeModel(m.theme), nil
			}
		case "left", "right":
			// Adjust sliders
			if ControlOrder[m.Focused] == "Fore" {
				m.style.Fore = clamp(m.style.Fore+sliderAdjustment(msg.String()), -1, 255)
			} else if ControlOrder[m.Focused] == "Back" {
				m.style.Back = clamp(m.style.Back+sliderAdjustment(msg.String()), -1, 255)
			}
		}
	}

	var cmd tea.Cmd
	m.NameInput, cmd = m.NameInput.Update(msg)
	outCmds = append(outCmds, cmd)

	m.FileArea, cmd = m.FileArea.Update(msg)
	outCmds = append(outCmds, cmd)

	return m, tea.Batch(outCmds...)
}

func (m StyleModel) View() string {
	msgBody :=
		fmt.Sprintf(
			"%v\n\n"+ // Style name
				"%v\n\n"+ // File preview
				"%v\n\n"+ // Style options
				"%v\n\n"+ // Color sliders
				"%v\n\n"+ // Filetypes
				"%v\n\n", // Selection Buttons
			// "%v\n\n", // Help keymap
			m.renderName(),
			m.renderPreview(),
			m.renderStyles(),
			m.renderSliders(),
			m.renderFileArea(),
			m.renderButtons(),
			// CenterHorz(m.help.View(m.keymap)),
		)

	return RenderModel(msgBody, m.help.View(m.keymap))
}

func (m StyleModel) renderName() string {
	return CenterHorz(fmt.Sprintf("%v: %v", TitleStyle.Render("Style"), m.NameInput.View()))
}

func (m StyleModel) renderStyles() string {
	sOut := TitleStyle.Render("Styles") + "\n\n" +
		"%s\n" +
		"%s\n" +
		"%s"

	return CenterHorz(fmt.Sprintf(sOut,
		checkboxView(m.style.Bold, "Bold     ", ControlOrder[m.Focused] == "Bold"),
		checkboxView(m.style.Under, "Underline", ControlOrder[m.Focused] == "Under"),
		checkboxView(m.style.Blink, "Blink    ", ControlOrder[m.Focused] == "Blink")))
}

func (m StyleModel) renderSliders() string {
	var foreStr string
	var backStr string
	if ControlOrder[m.Focused] == "Fore" {
		foreStr = focusedStyle.Render(" > Foreground") + " [%03d]\n%v"
		backStr = "   Background [%03d]\n%v"
	} else if ControlOrder[m.Focused] == "Back" {
		foreStr = "   Foreground [%03d]\n%v"
		backStr = focusedStyle.Render(" > Background") + " [%03d]\n%v"
	} else {
		foreStr = "   Foreground [%03d]\n%v"
		backStr = "   Background [%03d]\n%v"
	}
	foreStr = fmt.Sprintf(foreStr, m.style.Fore, "← "+m.foreSlider.ViewAs(((float64)(m.style.Fore+2)/257))+" →")
	backStr = fmt.Sprintf(backStr, m.style.Back, "← "+m.backSlider.ViewAs(((float64)(m.style.Back+2)/257))+" →")

	return CenterHorz(fmt.Sprintf(TitleStyle.Render("Colors")+"\n\n%v\n\n%v", foreStr, backStr))
}

func (m StyleModel) renderPreview() string {
	var backColor lipgloss.Color
	if m.style.Back == -1 {
		backColor = lipgloss.Color("")
	} else {
		backColor = lipgloss.Color(strconv.Itoa(m.style.Back))
	}
	previewColor := lipgloss.NewStyle().
		Foreground(lipgloss.Color(strconv.Itoa(m.style.Fore))).
		Background(backColor).
		Bold(m.style.Bold).Underline(m.style.Under).Blink(m.style.Blink)

	return CenterHorz(TitleStyle.Render("Preview") + ": " + previewColor.Render("file.example"))
}

func (m StyleModel) renderFileArea() string {
	return CenterHorz(fmt.Sprintf("%v\n%v", TitleStyle.Render("File Types"), m.FileArea.View()))
}

func (m StyleModel) renderButtons() string {
	var save string
	var discard string
	if ControlOrder[m.Focused] == "Save" {
		save = focusedStyle.Render("[ Save & Exit ]")
		discard = "[ " + blurredStyle.Render("Discard") + " ]"
	} else if ControlOrder[m.Focused] == "Discard" {
		save = "[ " + blurredStyle.Render("Save & Exit") + " ]"
		discard = focusedStyle.Render("[ Discard ]")
	} else {
		save = "[ " + blurredStyle.Render("Save & Exit") + " ]"
		discard = "[ " + blurredStyle.Render("Discard") + " ]"

	}
	return CenterHorz(fmt.Sprintf("%v\n%v", save, discard))
}

func (m StyleModel) GetDirColorBlock() string {
	outStr := " # " + m.style.Name + "\n\n"

	styleStr := ""

	if m.style.Bold {
		styleStr += "1;"
	}

	if m.style.Under {
		styleStr += "4;"
	}

	if m.style.Blink {
		styleStr += "5;"
	}

	if m.style.Fore != -1 {
		styleStr += fmt.Sprintf("38;5;%v;", strconv.Itoa(m.style.Fore))
	}

	if m.style.Back != -1 {
		styleStr += fmt.Sprintf("48;5;%v;", strconv.Itoa(m.style.Back))
	}

	styleStr = strings.TrimSuffix(styleStr, ";")

	for _, file := range m.style.FileTypes {
		outStr += fmt.Sprintf("%v %v\n", file, styleStr)
	}

	return outStr
}

func checkboxView(checked bool, label string, focused bool) string {
	var check string
	var focus string
	var style lipgloss.Style
	if checked {
		check = "[x]"
	} else {
		check = "[ ]"
	}

	if focused {
		focus = " >"
		style = focusedStyle
	} else {
		focus = "  "
		style = noStyle
	}

	return style.Render(fmt.Sprintf("%v %v %v", focus, check, label))
}

func sliderAdjustment(key string) int {
	if key == "right" {
		return 1
	} else if key == "left" {
		return -1
	}
	return 0
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Region: Keymap stuff

type styleKeymap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding

	Select key.Binding
	Quit   key.Binding
}

func (k styleKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Quit}
}

func (k styleKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Quit}, //, k.Select, k.Quit},
		{k.Left, k.Right, k.Select},
		// {k.Select, k.Quit},
	}
}

func newStyleKeymap() styleKeymap {
	return styleKeymap{
		Up: key.NewBinding(
			key.WithKeys("shift+tab", "up"),
			key.WithHelp("ShTab/↑", "Up"),
		),
		Down: key.NewBinding(
			key.WithKeys("tab", "down"),
			key.WithHelp("tab/↓", "Down"),
		),

		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("h/←", "-1"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("l/→", "+1"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Toggle"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("ctrl+c", "Quit"),
		),
	}
}
