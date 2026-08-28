package gateway

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fehmicorp/pkg/v1/config"
	httpHandler "github.com/fehmicorp/pkg/v1/http"
)

func StartServer(prefixes ...string) {
	var targetPrefix string
	if len(prefixes) > 0 && strings.TrimSpace(prefixes[0]) != "" {
		targetPrefix = prefixes[0]
	} else {
		targetPrefix = os.Getenv("PREFIX")
	}

	server := httpHandler.NewServer(targetPrefix)

	slog.Info(
		"server started",
		slog.String("address", server.Addr),
		slog.String("prefix", targetPrefix),
		slog.String("environment", config.Conf.App.Environment),
	)

	// Create context bound to OS interruption signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Connect to Cloudflare if token is available, otherwise run normally
	if config.Conf.Server.CFTunnel != "" {
		slog.Info("Cloudflare Tunnel token detected; initializing tunnel daemon")
		if err := StartTunnelProcess(ctx, config.Conf.Server.CFTunnel); err != nil {
			slog.Error("Failed to initiate Cloudflare tunnel", slog.String("error", err.Error()))
		}
	} else {
		slog.Info("No Cloudflare Tunnel token configured; starting server normally")
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(
				"server failed",
				slog.String("error", err.Error()),
			)
			os.Exit(1)
		}
	}()

	// Wait for termination signal
	<-ctx.Done()
	slog.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error(
			"failed to shutdown server",
			slog.String("error", err.Error()),
		)
	}

	slog.Info("server shutdown successfully")
}
