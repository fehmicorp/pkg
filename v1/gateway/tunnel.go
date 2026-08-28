package gateway

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

// StartTunnelProcess launches cloudflared tunnel run in background
func StartTunnelProcess(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("cloudflared tunnel token is empty")
	}

	cmd := exec.CommandContext(ctx, "cloudflared", "tunnel", "run", "--token", token)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	slog.Info("Starting Cloudflare Tunnel daemon...")
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start cloudflared daemon: %w", err)
	}

	go func() {
		if err := cmd.Wait(); err != nil && ctx.Err() == nil {
			slog.Error("Cloudflare tunnel process exited unexpectedly", slog.String("error", err.Error()))
		}
	}()

	return nil
}
