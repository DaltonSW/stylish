package tui

import (
	"fmt"
	"strconv"

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

type StyleModel struct {
	Theme string

	StyleName string
	NameInput textinput.Model

	Bold  bool
	Under bool
	Blink bool

	ForeColor int
	BackColor int

	FileArea textarea.Model

	Focused int

	ViewWidth  int
	ViewHeight int
}

func NewStyleCreateModel(theme string, viewWidth, viewHeight int) StyleModel {
	nameInput := textinput.New()
	nameInput.Focus()
	nameInput.Placeholder = "Style Name"
	nameInput.PromptStyle = focusedStyle
	nameInput.TextStyle = focusedStyle

	foreSlider := 128
	backSlider := 0

	fileArea := textarea.New()
	fileArea.Placeholder = ".mp3\n.gif\n.docx\n..."
	fileArea.Blur()

	return StyleModel{
		Theme:      theme,
		NameInput:  nameInput,
		ForeColor:  foreSlider,
		BackColor:  backSlider,
		ViewWidth:  viewWidth,
		ViewHeight: viewHeight,
		FileArea:   fileArea,
	}
}

func NewStyleEditModel(theme, style string, viewWidth, viewHeight int) StyleModel {
	nameInput := textinput.New()
	// nameInput.Focus()
	nameInput.Placeholder = "Style Name"
	nameInput.SetValue(style)
	nameInput.PromptStyle = blurredStyle
	nameInput.TextStyle = blurredStyle

	foreSlider := 128
	backSlider := 0

	fileArea := textarea.New()
	fileArea.Placeholder = ".mp3\n.gif\n.docx\n..."
	fileArea.Blur()

	return StyleModel{
		Theme:      theme,
		NameInput:  nameInput,
		ForeColor:  foreSlider,
		BackColor:  backSlider,
		ViewWidth:  viewWidth,
		ViewHeight: viewHeight,
		Focused:    1,
		FileArea:   fileArea,
	}
}

func (m StyleModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m StyleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var outCmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return NewThemeModel(m.Theme, m.ViewWidth, m.ViewHeight), nil
		case "tab", "down", "shift+tab", "up":
			if msg.String() == "tab" || msg.String() == "down" {

				m.Focused++
				if m.Focused > 6 {
					m.Focused = 0
				}
			} else {
				m.Focused--
				if m.Focused < 0 {
					m.Focused = 6
				}
			}
			if m.Focused == 0 {
				outCmds = append(outCmds, m.NameInput.Focus())
				m.NameInput.PromptStyle = focusedStyle
				m.NameInput.TextStyle = focusedStyle
			} else {
				m.NameInput.Blur()
				m.NameInput.PromptStyle = blurredStyle
				m.NameInput.TextStyle = blurredStyle
			}

			if m.Focused == 6 {
				outCmds = append(outCmds, m.FileArea.Focus())
			} else {
				m.FileArea.Blur()
			}
		case "space", "enter":
			// Toggle checkboxes if focused
			switch m.Focused {
			case 1:
				m.Bold = !m.Bold
			case 2:
				m.Under = !m.Under
			case 3:
				m.Blink = !m.Blink
			}
		case "left", "right":
			// Adjust sliders
			if m.Focused == 4 {
				m.ForeColor = clamp(m.ForeColor+sliderAdjustment(msg.String()), 0, 255)
			} else if m.Focused == 5 {
				m.BackColor = clamp(m.BackColor+sliderAdjustment(msg.String()), 0, 255)
			}
		case "s":
			// TODO: Save the style
			// saveStyle(m.folder, m.fileName, m)
		}
	case tea.WindowSizeMsg:
		m.ViewWidth = msg.Width
		m.ViewHeight = msg.Height

	}

	var cmd tea.Cmd
	m.NameInput, cmd = m.NameInput.Update(msg)
	outCmds = append(outCmds, cmd)

	return m, tea.Batch(outCmds...)
}

func (m StyleModel) View() string {

	return Center(fmt.Sprintf("%v\n\n%v\n\n%v\n\n%v\n\n%v", m.renderName(), m.renderStyles(), m.renderSliders(), m.renderPreview(), m.renderFileArea()))

}

func (m StyleModel) renderName() string {
	return CenterHorz(fmt.Sprintf("%v\n\n%v", TitleStyle.Render("Name"), m.NameInput.View()))
}

func (m StyleModel) renderStyles() string {
	sOut := TitleStyle.Render("Styles") + "\n\n" +
		"%s\n" +
		"%s\n" +
		"%s\n\n"

	return CenterHorz(fmt.Sprintf(sOut,
		checkboxView(m.Bold, "Bold     ", m.Focused == 1),
		checkboxView(m.Under, "Underline", m.Focused == 2),
		checkboxView(m.Blink, "Blink    ", m.Focused == 3)))
}

func (m StyleModel) renderSliders() string {
	var foreStr string
	var backStr string
	if m.Focused == 4 {
		foreStr = focusedStyle.Render(" > Foreground:") + " %v [%03d]"
		backStr = "   Background: %v [%03d]"
	} else if m.Focused == 5 {
		foreStr = "   Foreground: %v [%03d]"
		backStr = focusedStyle.Render(" > Background:") + " %v [%03d]"
	} else {
		foreStr = "   Foreground: %v [%03d]"
		backStr = "   Background: %v [%03d]"
	}
	foreStr = fmt.Sprintf(foreStr, renderSlider(m.ForeColor, m.ViewWidth/2), m.ForeColor)
	backStr = fmt.Sprintf(backStr, renderSlider(m.BackColor, m.ViewWidth/2), m.BackColor)

	return CenterHorz(fmt.Sprintf(TitleStyle.Render("Colors")+"\n\n%v\n%v", foreStr, backStr))
}

func (m StyleModel) renderPreview() string {
	previewColor := lipgloss.NewStyle().
		Foreground(lipgloss.Color(strconv.Itoa(m.ForeColor))).
		Background(lipgloss.Color(strconv.Itoa(m.BackColor))).
		Bold(m.Bold).Underline(m.Under).Blink(m.Blink)

	return CenterHorz(TitleStyle.Render("Preview") + "\n\n" + previewColor.Render("file.example"))
}

func (m StyleModel) renderFileArea() string {
	return CenterHorz(fmt.Sprintf("%v\n%v", TitleStyle.Render("File Types"), m.FileArea.View()))
}

func renderSlider(value, width int) string {
	bar := ""
	totalBlocks := width
	position := value * totalBlocks / 255
	for i := 0; i < totalBlocks; i++ {
		if i <= position {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
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
