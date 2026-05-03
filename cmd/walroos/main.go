package main

import (
	"fmt"
	"os"

	"github.com/jefjesuswt/walroos/cmd/walroos/hub"
	"github.com/jefjesuswt/walroos/cmd/walroos/node"
	"github.com/jefjesuswt/walroos/shared/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "walroos",
	Short: "A CLI tool to manage Podman containers",
	Long:  "An All-in-One CLI tool for managing Podman containers, it allows the user to orchestrate multiple servers and manage their containers from a single interface.",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintTitle("🐳 WALROOS CORE CLI")
		fmt.Println("  Use 'walroos --help' to see available commands.")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.AddCommand(hub.Cmd)
	rootCmd.AddCommand(node.Cmd)
}
