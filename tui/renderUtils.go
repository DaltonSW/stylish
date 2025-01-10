package tui

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const Title = "stylish"
const Subtitle = "~ Feel good in your shell ~"

const HexCodePattern = "[0-9a-fA-F]{6}"

const ConstWidth = 35
const ConstHeight = 27

const TrueColorFore = "38;2;%d;%d;%d"
const TrueColorBack = "48;2;%d;%d;%d"

const EightBitFore = "38;5;%d"
const EightBitBack = "48;5;%d"

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

func ValidHexCode(input string) error {
	match, err := regexp.MatchString(HexCodePattern, input)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("Enter a valid hex code")
	}

	return nil

}

func HexToRGBA(hex string) color.RGBA {
	values, _ := strconv.ParseUint(string(hex), 16, 32)
	return color.RGBA{R: uint8(values >> 16), G: uint8((values >> 8) & 0xFF), B: uint8(values & 0xFF), A: 255}
}

func HexToEightBit(hex string) uint8 {
	color := HexToRGBA(hex)

	return (color.R*7/255)<<5 + (color.G*7/255)<<2 + (color.B * 3 / 255)
}
