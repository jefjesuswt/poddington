package hub

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/jefjesuswt/walroos/internal/fleet"
	"github.com/jefjesuswt/walroos/shared/ui"
	"github.com/jefjesuswt/walroos/shared/whttp"
	"github.com/spf13/cobra"
	"tailscale.com/tsnet"
)

var hubConfigPath string

var Cmd = &cobra.Command{
	Use:   "hub",
	Short: "Starts Walroos in Hub mode (controller)",
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := cmd.Context()
		ui.PrintTitle("Walroos Hub Initialization")

		svc, cleanup, err := InitFleetService(ctx)
		if err != nil {
			return fmt.Errorf("failed to start database. %w", err)
		}

		defer func() {
			if err := cleanup(); err != nil {
				slog.Error("error closing database", "err", err)
			} else {
				slog.Info("Database closed successfully")
			}
		}()

		router := whttp.NewRouter()

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
			Hostname: "walroos-hub",
			AuthKey:  tsAuthKey,
			Logf:     func(format string, args ...any) {},
		}
		defer tsServer.Close()

		ln, err := tsServer.Listen("tcp", ":80")
		if err != nil {
			return ui.WrapError("failed to create tailscale listener: %w", err)
		}

		slog.Info("Hub online and using secure mesh.", "hostname", "walroos-hub")

		if err := http.Serve(ln, router); err != nil {
			return ui.WrapError("server crashed: %w", err)
		}

		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&hubConfigPath, "config", "c", "~/.config/walroos/hub.yaml", "Path to config file")

	Cmd.AddCommand(addCommand)
}

func withLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		duration := time.Since(start)

		slog.Info("HTTP Request",
			"method", req.Method,
			"path", req.URL.Path,
			"duration", duration,
		)
	})
}

func withRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("PANIC in HTTP", "error", err, "trace", string(debug.Stack()))
				whttp.ErrorJSON(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		next.ServeHTTP(w, req)
	})
}
