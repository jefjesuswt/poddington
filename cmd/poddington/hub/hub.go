package hub

import (
	"github.com/charmbracelet/log"
	"github.com/jefjesuswt/poddington/internal/ui"
	"github.com/spf13/cobra"
)

var hubConfigPath string

var Cmd = &cobra.Command{
	Use:   "hub",
	Short: "Starts Poddington in Hub mode (controller)",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintTitle("Poddington Hub Initialization")
		log.Info("Loading configuration...", "path", hubConfigPath)
		log.Info("Listening for incoming node connections...", "port", "8443")
	},
}

func init() {
	Cmd.Flags().StringVarP(&hubConfigPath, "config", "c", "~/.config/poddington/hub.yaml", "Path to config file")
}
