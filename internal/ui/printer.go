package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// tokyo night inspired
const (
	ColorPrimary = lipgloss.Color("#7aa2f7")
	ColorText    = lipgloss.Color("#c0caf5")
	ColorSubtext = lipgloss.Color("#565f89")
	ColorSuccess = lipgloss.Color("#9ece6a")
	ColorError   = lipgloss.Color("#f7768e")
)

var (
	// title style (minimalist with accent on the left)
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Bold(true).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(ColorPrimary).
			PaddingLeft(1).
			MarginBottom(1).
			MarginTop(1)

	// list columns styles (fixed width for table alignment)
	ColID    = lipgloss.NewStyle().Foreground(ColorSubtext).Width(15)
	ColName  = lipgloss.NewStyle().Foreground(ColorText).Bold(true).Width(25)
	ColState = lipgloss.NewStyle().Width(12) // dynamic state color
	ColImage = lipgloss.NewStyle().Foreground(ColorSubtext).Italic(true)

	// inspect block style (vertical border and padding)
	InspectBlock = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(ColorSubtext).
			PaddingLeft(2)

	KeyStyle   = lipgloss.NewStyle().Foreground(ColorPrimary).Width(10)
	ValueStyle = lipgloss.NewStyle().Foreground(ColorText)
)

func PrintTitle(format string, a ...any) {
	text := fmt.Sprintf(format, a...)
	fmt.Println(TitleStyle.Render(text))
}

func PrintListHeader() {
	fmt.Printf("  %s%s%s%s\n",
		ColID.Foreground(ColorPrimary).Render("CONTAINER ID"),
		ColName.Foreground(ColorPrimary).Render("NAME"),
		ColState.Foreground(ColorPrimary).Render("STATE"),
		ColImage.Foreground(ColorPrimary).Render("IMAGE"),
	)
}

func PrintContainerRow(id, name, state, image string) {
	st := lipgloss.NewStyle().Foreground(ColorSubtext).Render(state)
	if state == "running" {
		st = lipgloss.NewStyle().Foreground(ColorSuccess).Render(state)
	}

	fmt.Printf("  %s%s%s%s\n",
		ColID.Render(id),
		ColName.Render(name),
		ColState.Render(st),
		ColImage.Render(image),
	)
}

func PrintError(format string, a ...any) error {
	text := fmt.Errorf(format, a...).Error()
	errTag := lipgloss.NewStyle().Foreground(ColorError).Bold(true).Render("✖ ERROR:")
	return fmt.Errorf("\n  %s %s\n", errTag, text)
}

func PrintSuccess(format string, a ...any) {
	text := fmt.Sprintf(format, a...)
	tag := lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render("✔ SUCCESS:")
	fmt.Printf("\n  %s %s\n", tag, text)
}

func PrintInfo(format string, a ...any) {
	text := fmt.Sprintf(format, a...)
	tag := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render("ℹ INFO:")
	fmt.Printf("\n  %s %s\n", tag, text)
}

func PrintList(key string, items []string) {
	if len(items) == 0 {
		PrintKeyValue(key, "None")
		return
	}

	fmt.Printf("  %s\n", KeyStyle.Render(key+":"))
	for _, item := range items {
		// Pinta un puntito estilo viñeta antes de cada item
		dot := lipgloss.NewStyle().Foreground(ColorPrimary).Render("•")
		fmt.Printf("    %s %s\n", dot, ValueStyle.Render(item))
	}
}

func PrintKeyValue(key, value string) {
	if value == "" {
		value = "N/A"
	}
	fmt.Printf("  %s %s\n", KeyStyle.Render(key+":"), ValueStyle.Render(value))
}
