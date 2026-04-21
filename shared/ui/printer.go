package ui

import "fmt"

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
)

func PrintTitle(title string) {
	fmt.Printf("\n%s=== %s ===%s\n\n", Green, title, Reset)
}

func PrintSuccess(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s✔ SUCCESS:%s %s\n", Green, Reset, msg)
}

func PrintWarning(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("\n%s⚠ WARNING:%s %s\n", Yellow, Reset, msg)
}

func PrintKeyValue(key, value string) {
	fmt.Printf("  %s%-10s%s %s\n", Cyan, key+":", Reset, value)
}

func WrapError(format string, a ...any) error {
	return fmt.Errorf("\n"+Red+"✖ ERROR:"+Reset+" "+format, a...)
}
