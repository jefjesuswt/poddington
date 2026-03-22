package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	PodIndigo   = lipgloss.Color("#25283d") // Space Indigo
	PodGrape    = lipgloss.Color("#8f3985") // Grape Soda
	PodFrosted  = lipgloss.Color("#98dfea") // Frosted Blue
	PodSeaGreen = lipgloss.Color("#07beb8") // Light Sea Green
	PodPowder   = lipgloss.Color("#efd9ce") // Powder Petal
)

var (
	ColorSuccess = PodSeaGreen
	ColorWarning = lipgloss.Color("#e0af68")
	ColorError   = lipgloss.Color("#f7768e")
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(PodGrape).
			PaddingLeft(1).
			MarginBottom(1).
			MarginTop(1)

	ColID    = lipgloss.NewStyle().Foreground(PodFrosted).Width(15)
	ColName  = lipgloss.NewStyle().Foreground(PodPowder).Bold(true).Width(25)
	ColState = lipgloss.NewStyle().Width(12)
	ColImage = lipgloss.NewStyle().Foreground(PodGrape).Italic(true)

	InspectBlock = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(PodIndigo).
			PaddingLeft(2)

	KeyStyle   = lipgloss.NewStyle().Foreground(PodFrosted).Width(10).Bold(true)
	ValueStyle = lipgloss.NewStyle().Foreground(PodPowder)
)

var LogBlock = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), false, false, false, true).
	BorderForeground(PodIndigo).
	PaddingLeft(2).
	Foreground(lipgloss.AdaptiveColor{Light: string(PodIndigo), Dark: string(PodPowder)})

func hexToRGB(hex string) (r, g, b float64) {
	hex = strings.TrimPrefix(hex, "#")
	parsed, _ := strconv.ParseUint(hex, 16, 32)
	r = float64(parsed >> 16)
	g = float64((parsed >> 8) & 0xFF)
	b = float64(parsed & 0xFF)
	return
}

func rgbToHex(r, g, b float64) string {
	return fmt.Sprintf("#%02x%02x%02x", int(r), int(g), int(b))
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func ApplyGradient(text string) string {
	colors := []string{"#25283d", "#8f3985", "#98dfea", "#07beb8", "#efd9ce"}

	var builder strings.Builder
	builder.Grow(len(text) * 15)

	length := len(text)
	if length <= 1 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(colors[0])).Render(text)
	}

	for i, char := range text {
		t := float64(i) / float64(length-1)
		segmentFloat := t * float64(len(colors)-1)
		segment := int(segmentFloat)
		if segment >= len(colors)-1 {
			segment = len(colors) - 2
		}

		fraction := segmentFloat - float64(segment)

		r1, g1, b1 := hexToRGB(colors[segment])
		r2, g2, b2 := hexToRGB(colors[segment+1])

		r := lerp(r1, r2, fraction)
		g := lerp(g1, g2, fraction)
		b := lerp(b1, b2, fraction)

		smoothHexColor := rgbToHex(r, g, b)
		builder.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(smoothHexColor)).Render(string(char)))
	}

	return builder.String()
}

func PrintGradientTitle(text string) {
	gradText := ApplyGradient(text)
	fmt.Printf("\n  %s\n\n", lipgloss.NewStyle().Bold(true).Render(gradText))
}

func PrintTitle(format string, a ...any) {
	text := fmt.Sprintf(format, a...)
	gradText := ApplyGradient(text)
	fmt.Println(TitleStyle.Render(gradText))
}

func PrintListHeader() {
	fmt.Printf("  %s%s%s%s\n",
		ColID.Foreground(PodGrape).Render("CONTAINER ID"),
		ColName.Foreground(PodGrape).Render("NAME"),
		ColState.Foreground(PodGrape).Render("STATE"),
		ColImage.Foreground(PodGrape).Render("IMAGE"),
	)
}

func PrintContainerRow(id, name, state, image string) {
	st := lipgloss.NewStyle().Foreground(PodIndigo).Render(state)
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

func PrintKeyValue(key, value string) {
	if value == "" {
		value = "N/A"
	}
	fmt.Printf("%s %s\n", KeyStyle.Render(key+":"), ValueStyle.Render(value))
}

func PrintList(key string, items []string) {
	if len(items) == 0 {
		PrintKeyValue(key, "None")
		return
	}

	fmt.Printf("%s\n", KeyStyle.Render(key+":"))
	for _, item := range items {
		dot := lipgloss.NewStyle().Foreground(PodGrape).Render("•")
		fmt.Printf("  %s %s\n", dot, ValueStyle.Render(item))
	}
}

func PrintSuccess(format string, a ...any) {
	text := fmt.Sprintf(format, a...)
	tag := lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render("✔ SUCCESS:")
	fmt.Printf("\n  %s %s\n", tag, text)
}

func PrintInfo(format string, a ...any) {
	text := fmt.Sprintf(format, a...)
	tag := lipgloss.NewStyle().Foreground(PodFrosted).Bold(true).Render("ℹ INFO:")
	fmt.Printf("\n  %s %s\n", tag, text)
}

func PrintWarning(format string, a ...any) {
	text := fmt.Sprintf(format, a...)
	tag := lipgloss.NewStyle().Foreground(ColorWarning).Bold(true).Render("⚠ WARNING:")
	fmt.Printf("\n  %s %s\n", tag, text)
}

func PrintError(format string, a ...any) error {
	text := fmt.Errorf(format, a...).Error()
	errTag := lipgloss.NewStyle().Foreground(ColorError).Bold(true).Render("✖ ERROR:")
	return fmt.Errorf("\n  %s %s\n", errTag, text)
}

func PrintLogs(target, logs string) {
	PrintTitle("Logs: %s", target)

	cleanLogs := strings.TrimSpace(logs)
	if cleanLogs == "" {
		fmt.Println(LogBlock.Render("No logs available for this container."))
	} else {
		fmt.Println(LogBlock.Render(cleanLogs))
	}
	fmt.Println()
}
