package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fehmicorp/pkg/v1/config"
)

func main() {
	config.ClientInit()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	connector := cf.NewConnector(*config.CliConf)
	go func() {
		if err := connector.Connect(ctx); err != nil {
			slog.Error("Client Tunnel Connector failed", slog.String("error", err.Error()))
		}
	}()
	time.Sleep(500 * time.Millisecond)
	client := connector.GetHTTPClient()
	req, err := client.Get("http://" + config.CliConf.Tunnel.LocalBindAddr + "/api/v1/status")
	if err != nil {
		slog.Error("Request to private server failed", slog.String("error", err.Error()))
	} else {
		defer req.Body.Close()
		body, _ := io.ReadAll(req.Body)
		slog.Info("Response from private server over Cloudflare Tunnel", slog.String("body", string(body)))
	}

	<-ctx.Done()
	slog.Info("Shutting down client gateway...")
}
