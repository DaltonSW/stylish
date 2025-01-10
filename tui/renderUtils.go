package tui

import (
	"errors"
	"fmt"
	// "image/color"
	"log"
	"os"
	"regexp"
	// "strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
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

var EightBitMode = false
var DefaultTermFore lipgloss.Color
var DefaultTermBack lipgloss.Color

var ViewportBorder = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#4400FF")).Height(ConstHeight)

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

func HexToRGB(hex string) termenv.RGBColor {
	if strings.HasPrefix(hex, "#") {
		return termenv.RGBColor(hex)
	} else {
		return termenv.RGBColor("#" + hex)
	}

}

func HexToEightBit(hex string) termenv.Color {
	prof256 := termenv.ANSI256
	return prof256.Convert(HexToRGB(hex))
}
