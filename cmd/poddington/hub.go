package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var hubConfigPath string

var hubCmd = &cobra.Command{
	Use:   "hub",
	Short: "Starts Poddington in Hub mode (controller)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Running Poddington Hub using config set on %s...\n", hubConfigPath)
		fmt.Println("Poddington Hub is running.")
	},
}

func init() {
	hubCmd.Flags().StringVar(&hubConfigPath, "config", "c", "~/.config/poddington/hub.yaml")
}
