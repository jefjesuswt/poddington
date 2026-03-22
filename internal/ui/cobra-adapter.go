package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)


func SetupCobraUI(rootCmd *cobra.Command) {
	rootCmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		printHelp(c)
	})

	rootCmd.SetUsageFunc(func(c *cobra.Command) error {
		printUsage(c)
		return nil
	})

	rootCmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		return PrintError("Invalid argument: %v", err)
	})
}

func printHelp(c *cobra.Command) {
	fmt.Println()

	// title
	if c.Parent() == nil {
		PrintGradientTitle("🐳 PODDINGTON CORE CLI")
		fmt.Printf("  %s\n\n", ValueStyle.Render(c.Long))
	} else {
		fmt.Printf("  %s\n\n", ValueStyle.Render(c.Short))
	}

	// usage
	PrintTitle("USAGE")
	fmt.Printf("    %s %s\n\n", ColName.Render(c.CommandPath()), ValueStyle.Render(strings.TrimPrefix(c.UseLine(), c.CommandPath()+" ")))

	// commands
	if len(c.Commands()) > 0 {
		PrintTitle("COMMANDS")
		for _, sub := range c.Commands() {
			if sub.IsAvailableCommand() || sub.Name() == "help" {
				fmt.Printf("    %s %s\n", ColImage.Render(fmt.Sprintf("%-15s", sub.Name())), ValueStyle.Render(sub.Short))
			}
		}
		fmt.Println()
	}

	// flags
	if c.HasAvailableLocalFlags() {
		PrintTitle("FLAGS")
		flagsOut := strings.TrimRight(c.LocalFlags().FlagUsages(), "\n")
		fmt.Println(ValueStyle.Render(flagsOut))
		fmt.Println()
	}
}

func printUsage(c *cobra.Command) {
	fmt.Println()
	fmt.Printf("  %s\n", ApplyGradient("USAGE"))
	fmt.Printf("    %s %s\n\n", ColName.Render(c.CommandPath()), ValueStyle.Render(strings.TrimPrefix(c.UseLine(), c.CommandPath()+" ")))

	helpCmd := lipgloss.NewStyle().Foreground(PodGrape).Render(fmt.Sprintf("%s --help", c.CommandPath()))
	fmt.Printf("  Run '%s' for more information.\n\n", helpCmd)
}
