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
	var activePrefix string
	if len(prefixes) > 0 {
		activePrefix = prefixes[0]
	}

	slog.Info(
		"server started",
		slog.String("address", server.Addr),
		slog.String("prefix", activePrefix),
		slog.String("environment", config.Conf.App.Environment),
	)

	done := make(chan os.Signal, 1)
	signal.Notify(
		done,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(
				"server failed",
				slog.String("error", err.Error()),
			)
			os.Exit(1)
		}
	}()

	<-done
	slog.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error(
			"failed to shutdown server",
			slog.String("error", err.Error()),
		)
	}

	slog.Info("server shutdown successfully")
}
