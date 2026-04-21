package hub

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/jefjesuswt/poddington/internal/fleet"
	podHTTP "github.com/jefjesuswt/poddington/shared/http"
	"github.com/jefjesuswt/poddington/shared/ui"
	"github.com/spf13/cobra"
	"tailscale.com/tsnet"
)

var hubConfigPath string

var Cmd = &cobra.Command{
	Use:   "hub",
	Short: "Starts Poddington in Hub mode (controller)",
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := cmd.Context()
		ui.PrintTitle("Poddington Hub Initialization")

		svc, err := InitFleetService(ctx)
		if err != nil {
			return fmt.Errorf("failed to start database. %w", err)
		}

		router := podHTTP.NewRouter()

		router.Use(withLogger, withRecoverer)

		handler := fleet.NewHandler(svc)
		handler.RegisterRoutes(router)

		tsAuthKey := os.Getenv("TS_AUTH_KEY")
		if tsAuthKey == "" {
			return ui.WrapError("TS_AUTH_KEY environment variable is required")
		}

		slog.Info("Loading configuration...", "path", hubConfigPath)
		slog.Info("Hub starting...", "db", "SQLite WAL")
		slog.Info("Connecting to Tailscale network...")

		tsServer := &tsnet.Server{
			Hostname: "poddington-hub",
			AuthKey:  tsAuthKey,
			Logf:     func(format string, args ...any) {},
		}
		defer tsServer.Close()

		ln, err := tsServer.Listen("tcp", ":80")
		if err != nil {
			return ui.WrapError("failed to create tailscale listener: %w", err)
		}

		slog.Info("Hub online and using secure mesh.", "hostname", "poddington-hub")

		if err := http.Serve(ln, router); err != nil {
			return ui.WrapError("server crashed: %w", err)
		}

		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&hubConfigPath, "config", "c", "~/.config/poddington/hub.yaml", "Path to config file")

	Cmd.AddCommand(addCommand)
}

func withLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		slog.Info("HTTP Request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", duration,
		)
	})
}

func withRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("PANIC in HTTP", "error", err, "trace", string(debug.Stack()))
				podHTTP.ErrorJSON(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
