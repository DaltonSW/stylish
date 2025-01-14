package styling

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const FancyTitle = `      _         _ _     _
  ___| |_ _   _| (_)___| |__
 / __| __| | | | | / __| '_ \
 \__ \ |_| |_| | | \__ \ | | |
 |___/\__|\__, |_|_|___/_| |_|
          |___/               `

// Colors
var (
	// Shoutouts to Catppuccin Latte for the reference "Light" hexcodes
	Red    = lipgloss.AdaptiveColor{Dark: "#FF6666", Light: "#D24950"}
	Orange = lipgloss.AdaptiveColor{Dark: "#FFBB66", Light: "#FE640B"}
	Yellow = lipgloss.AdaptiveColor{Dark: "#FFFF66", Light: "#DF8E1D"}
	Green  = lipgloss.AdaptiveColor{Dark: "#8CFF66", Light: "#40A02B"}
	Blue   = lipgloss.AdaptiveColor{Dark: "#66B3FF", Light: "#04A5E5"}
	Purple = lipgloss.AdaptiveColor{Dark: "#D966FF", Light: "#7287FD"}
)

var HeaderStyle = lipgloss.NewStyle()

func GetColoredTitle() string {
	titleSlice := strings.Split(FancyTitle, "\n")
	red := HeaderStyle.Foreground(Red).Render(titleSlice[0])
	orange := HeaderStyle.Foreground(Orange).Render(titleSlice[1])
	yellow := HeaderStyle.Foreground(Yellow).Render(titleSlice[2])
	green := HeaderStyle.Foreground(Green).Render(titleSlice[3])
	blue := HeaderStyle.Foreground(Blue).Render(titleSlice[4])
	purple := HeaderStyle.Foreground(Purple).Render(titleSlice[5])

	return fmt.Sprintf("%v\n%v\n%v\n%v\n%v\n%v", red, orange, yellow, green, blue, purple)
}

const ConstWidth = 35
const ConstHeight = 27
