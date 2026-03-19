package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "poddington",
	Short: "A CLI tool to manage Podman containers",
	Long:  "An All-in-One CLI tool for managing Podman containers, it allows the user to orchestrate multiple servers and manage their containers from a single interface.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Poddington! Use 'poddington --help' to see available commands.")
	},
}

func init() {
	rootCmd.AddCommand(hubCmd)
	rootCmd.AddCommand(nodeCmd)
}
