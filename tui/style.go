package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

	Focused int

	ViewWidth  int
	ViewHeight int
}

func NewStyleCreateModel(theme string, viewWidth, viewHeight int) StyleModel {
	nameInput := textinput.New()
	nameInput.Focus()
	nameInput.Placeholder = "Style (and File) Name"

	foreSlider := 128
	backSlider := 0

	return StyleModel{
		Theme:      theme,
		NameInput:  nameInput,
		ForeColor:  foreSlider,
		BackColor:  backSlider,
		ViewWidth:  viewWidth,
		ViewHeight: viewHeight,
	}
}

func NewStyleEditModel(theme, style string, viewWidth, viewHeight int) StyleModel {
	nameInput := textinput.New()
	// nameInput.Focus()
	nameInput.Placeholder = "Style (and File) Name"
	nameInput.SetValue(style)

	foreSlider := 128
	backSlider := 0

	return StyleModel{
		Theme:      theme,
		NameInput:  nameInput,
		ForeColor:  foreSlider,
		BackColor:  backSlider,
		ViewWidth:  viewWidth,
		ViewHeight: viewHeight,
		Focused:    1,
	}
}

func (m StyleModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m StyleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return NewThemeModel(m.Theme, m.ViewWidth, m.ViewHeight), nil
		case "tab":
			// Cycle through focusable elements
			m.Focused = (m.Focused + 1) % 7
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

	// Update name input if focused
	if m.Focused == 0 {
		var cmd tea.Cmd
		m.NameInput, cmd = m.NameInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m StyleModel) View() string {
	return fmt.Sprintf(
		"Name: %s\n\n"+
			"Styles:\n"+
			"%s Bold\n"+
			"%s Underline\n"+
			"%s Blink\n\n"+
			"Foreground: %s [%d]\n"+
			"Background: %s [%d]\n"+
			"Controls: [TAB] Cycle Focus | [SPACE] Toggle | [←/→] Adjust Slider | [S] Save | [ESC] Back",
		m.NameInput.View(),
		checkboxView(m.Bold, "Bold"),
		checkboxView(m.Under, "Underline"),
		checkboxView(m.Blink, "Blink"),
		renderSlider(m.ForeColor), m.ForeColor,
		renderSlider(m.BackColor), m.BackColor,
	)
}

func renderSlider(value int) string {
	bar := ""
	totalBlocks := 20
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

func checkboxView(checked bool, label string) string {
	if checked {
		return fmt.Sprintf("[x] %s", label)
	}
	return fmt.Sprintf("[ ] %s", label)
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
