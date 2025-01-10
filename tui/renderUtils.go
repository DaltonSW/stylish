package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const Title = "stylish"
const Subtitle = "~ Feel good in your shell ~"

const ConstWidth = 35
const ConstHeight = 27

var ViewportBorder = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#4400FF")) //.Width(ConstWidth).Height(ConstHeight)

var TitleStyle = lipgloss.NewStyle().Underline(true).Bold(true).Italic(true)
var SubtitleStyle = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#888888"))

var FocusedAreaStyle = textarea.Style{}

var BlurredAreaStyle = textarea.Style{}

func GetItemDelgate() (del list.DefaultDelegate) {
	del = list.NewDefaultDelegate()
	styles := list.NewDefaultItemStyles()

	styles.SelectedTitle = lipgloss.NewStyle().Padding(0, 0, 0, 2)
	styles.SelectedDesc = lipgloss.NewStyle().Padding(0, 0, 0, 2)

	del.Styles = styles

	del.SetHeight(4)

	return del
}

func GetTermSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	return width, height
}

func CenterHorz(msg string) string {
	return lipgloss.PlaceHorizontal(ConstWidth, lipgloss.Center, msg)

}

func Center(msg string) string {
	return lipgloss.Place(ConstWidth, ConstHeight, lipgloss.Center, lipgloss.Center, msg)
}

func ProgramHeader() string {
	return lipgloss.PlaceHorizontal(ConstWidth+2, lipgloss.Center, fmt.Sprintf("%v\n%v", TitleStyle.Render(Title), SubtitleStyle.Render(Subtitle)))
}

func RenderModel(body, footer string) string {
	return Center(fmt.Sprintf("%v\n%v", ProgramHeader(), ViewportBorder.Render(fmt.Sprintf("%v\n%v", body, CenterHorz(footer)))))
}
