package util

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	PinkColor  = "#FFFDF5"
	GreenColor = "#FFDDEE"
)

var (
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	StatusStyle = lipgloss.NewStyle().
			Inherit(StatusBarStyle).
			Foreground(lipgloss.Color(PinkColor)).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	statusText = lipgloss.NewStyle().Inherit(StatusBarStyle)

	docStyle = lipgloss.NewStyle()
)

func StyleText(color string, title string, text string) string {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	physicalWidth = physicalWidth - 1
	w := lipgloss.Width

	doc := strings.Builder{}

	statusKey := StatusStyle.Foreground(lipgloss.Color(color)).Render(title)
	statusVal := statusText.Copy().
		Width(physicalWidth - w(statusKey)).
		Render(text)

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		statusVal,
	)

	doc.WriteString(StatusBarStyle.Width(physicalWidth).Render(bar))
	docStyle = docStyle.MaxWidth(physicalWidth)

	return docStyle.Render(doc.String())
}
