package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const Title = "stylish"
const Subtitle = "~ Feel pretty in your shell ~"

const ConstWidth = 40
const ConstHeight = 20

// var OverallRender = lipgloss.NewStyle().Align(lipgloss.Left, lipgloss.Left)
var ViewportBorder = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#4400FF")).Width(ConstWidth).Height(ConstHeight) //.Width(50).Height(20).Align(lipgloss.Center, lipgloss.Center)

var TitleStyle = lipgloss.NewStyle().Underline(true).Bold(true).Italic(true).Foreground(lipgloss.Color("purple"))
var SubtitleStyle = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#888888"))

var FocusedAreaStyle = textarea.Style{}

var BlurredAreaStyle = textarea.Style{}

func GetTermSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	return width, height
}

func CenterHorz(msg string) string {
	return lipgloss.PlaceHorizontal(ConstWidth+1, lipgloss.Center, msg)

}

func Center(msg string) string {
	return lipgloss.Place(ConstWidth, ConstHeight, lipgloss.Center, lipgloss.Center, msg)
}

func ProgramHeader() string {
	return CenterHorz(fmt.Sprintf("%v\n%v\n", TitleStyle.Render(Title), SubtitleStyle.Render(Subtitle)))
}
