package cf

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
)

// EnsureCloudflared Checks if cloudflared exists, downloading/installing it if missing.
func EnsureCloudflared(ctx context.Context) error {
	_, err := exec.LookPath("cloudflared")
	if err == nil {
		slog.Info("cloudflared binary found in PATH")
		return nil
	}

	slog.Warn("cloudflared binary not found, attempting automatic installation...", slog.String("os", runtime.GOOS))

	switch runtime.GOOS {
	case "windows":
		return installWindows(ctx)
	case "darwin":
		return installMacOS(ctx)
	case "linux":
		return installLinux(ctx)
	default:
		return fmt.Errorf("unsupported OS platform for automatic installation: %s", runtime.GOOS)
	}
}

// StartTunnelProcess launches cloudflared tunnel run in background
func StartTunnelProcess(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("cloudflared tunnel token is empty")
	}

	// Auto-check and install binary before running
	if err := EnsureCloudflared(ctx); err != nil {
		return fmt.Errorf("cloudflared setup failed: %w", err)
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

// Windows Installation: Uses winget, falls back to direct executable download
func installWindows(ctx context.Context) error {
	if _, err := exec.LookPath("winget"); err == nil {
		cmd := exec.CommandContext(ctx, "winget", "install", "--id", "Cloudflare.cloudflared", "-e", "--accept-source-agreements", "--accept-package-agreements")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			slog.Info("cloudflared installed successfully via winget")
			return nil
		}
	}

	// Fallback to powershell direct binary download
	psCmd := `[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12; ` +
		`Invoke-WebRequest -Uri "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-windows-amd64.exe" -OutFile "$env:SystemRoot\system32\cloudflared.exe"`

	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to download cloudflared on Windows via PowerShell: %w", err)
	}

	slog.Info("cloudflared binary downloaded into System32 successfully")
	return nil
}

// macOS Installation: Uses Homebrew, falls back to pkg download
func installMacOS(ctx context.Context) error {
	if _, err := exec.LookPath("brew"); err == nil {
		cmd := exec.CommandContext(ctx, "brew", "install", "cloudflared")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			slog.Info("cloudflared installed successfully via Homebrew")
			return nil
		}
	}

	// Fallback to direct pkg installer download
	curlCmd := exec.CommandContext(ctx, "curl", "-L", "-o", "/tmp/cloudflared.pkg", "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-darwin-amd64.pkg")
	if err := curlCmd.Run(); err != nil {
		return fmt.Errorf("failed to download cloudflared macOS package: %w", err)
	}

	installCmd := exec.CommandContext(ctx, "sudo", "installer", "-pkg", "/tmp/cloudflared.pkg", "-target", "/")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("failed to run macOS package installer: %w", err)
	}

	slog.Info("cloudflared installed successfully on macOS")
	return nil
}

// Linux Installation: Automatically detects if sudo is present (host vs root containers)
func installLinux(ctx context.Context) error {
	arch := "amd64"
	if runtime.GOARCH == "arm64" {
		arch = "arm64"
	}

	targetURL := fmt.Sprintf("https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-%s", arch)

	// Check if sudo exists in $PATH
	useSudo := false
	if _, err := exec.LookPath("sudo"); err == nil && os.Geteuid() != 0 {
		useSudo = true
	}

	// Prepare curl command
	var downloadCmd *exec.Cmd
	if useSudo {
		downloadCmd = exec.CommandContext(ctx, "sudo", "curl", "-L", targetURL, "-o", "/usr/local/bin/cloudflared")
	} else {
		downloadCmd = exec.CommandContext(ctx, "curl", "-L", targetURL, "-o", "/usr/local/bin/cloudflared")
	}

	downloadCmd.Stdout = os.Stdout
	downloadCmd.Stderr = os.Stderr
	if err := downloadCmd.Run(); err != nil {
		return fmt.Errorf("failed to download cloudflared binary on Linux: %w", err)
	}

	// Prepare chmod command
	var chmodCmd *exec.Cmd
	if useSudo {
		chmodCmd = exec.CommandContext(ctx, "sudo", "chmod", "+x", "/usr/local/bin/cloudflared")
	} else {
		chmodCmd = exec.CommandContext(ctx, "chmod", "+x", "/usr/local/bin/cloudflared")
	}

	if err := chmodCmd.Run(); err != nil {
		return fmt.Errorf("failed to set execution permissions on /usr/local/bin/cloudflared: %w", err)
	}

	slog.Info("cloudflared binary installed successfully to /usr/local/bin/cloudflared")
	return nil
}
