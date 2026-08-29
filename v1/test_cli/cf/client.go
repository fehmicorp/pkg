package warp

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const warpMSIURL = "https://1111-releases.cloudflareclient.com/windows/Cloudflare_WARP_Release-x64.msi"

type ConnectorConfig struct {
	Token string
}

type Connector struct {
	cfg      ConnectorConfig
	execPath string
}

func NewConnector(cfg ConnectorConfig) *Connector {
	return &Connector{cfg: cfg}
}

func (c *Connector) FindWarpBinary() (string, error) {
	if path, err := exec.LookPath("warp-cli"); err == nil {
		return path, nil
	}

	var knownPaths []string
	switch runtime.GOOS {
	case "windows":
		programFiles := os.Getenv("ProgramFiles")
		if programFiles == "" {
			programFiles = `C:\Program Files`
		}
		knownPaths = []string{
			filepath.Join(programFiles, "Cloudflare", "Cloudflare WARP", "warp-cli.exe"),
		}
	case "darwin":
		knownPaths = []string{
			"/usr/local/bin/warp-cli",
			"/Applications/Cloudflare WARP.app/Contents/Resources/warp-cli",
		}
	case "linux":
		knownPaths = []string{
			"/usr/bin/warp-cli",
			"/usr/local/bin/warp-cli",
		}
	}

	for _, p := range knownPaths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("warp-cli executable not found")
}

func downloadFile(ctx context.Context, url string, destPath string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil // Follow redirects
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	_, copyErr := io.Copy(out, resp.Body)
	closeErr := out.Close()

	if copyErr != nil {
		return fmt.Errorf("failed to write file: %w", copyErr)
	}
	if closeErr != nil {
		return fmt.Errorf("failed to flush/close file: %w", closeErr)
	}

	fi, err := os.Stat(destPath)
	if err != nil || fi.Size() < 1000000 {
		return fmt.Errorf("downloaded file is incomplete or invalid")
	}

	return nil
}

func (c *Connector) InstallWarpWindows(ctx context.Context) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("automatic installation is only supported on Windows")
	}

	slog.Info("Cloudflare WARP not found. Starting automatic download...")

	msiPath := filepath.Join(os.TempDir(), "Cloudflare_WARP_Release-x64.msi")
	_ = os.Remove(msiPath)

	if err := downloadFile(ctx, warpMSIURL, msiPath); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	absMsiPath, err := filepath.Abs(msiPath)
	if err != nil {
		absMsiPath = msiPath
	}

	slog.Info("Download completed. Triggering elevated installer (UAC prompt)...", slog.String("msi_path", absMsiPath))

	psArgs := fmt.Sprintf(`Start-Process msiexec.exe -ArgumentList '/i "%s" /qn /norestart' -Verb RunAs -Wait`, absMsiPath)
	cmd := exec.CommandContext(ctx, "powershell.exe", "-NoProfile", "-NonInteractive", "-Command", psArgs)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("elevated installation failed: %w", err)
	}

	slog.Info("Installation finished. Waiting for Cloudflare WARP service initialization...")
	time.Sleep(5 * time.Second)

	return nil
}

func (c *Connector) EnsureBinary(ctx context.Context) error {
	path, err := c.FindWarpBinary()
	if err == nil {
		c.execPath = path
		slog.Info("Found warp-cli binary", slog.String("path", c.execPath))
		return nil
	}

	if err := c.InstallWarpWindows(ctx); err != nil {
		return fmt.Errorf("auto-installation failed: %w", err)
	}

	path, err = c.FindWarpBinary()
	if err != nil {
		return fmt.Errorf("warp-cli still missing after installation: %w", err)
	}

	c.execPath = path
	slog.Info("warp-cli successfully installed and located", slog.String("path", c.execPath))
	return nil
}

func (c *Connector) ResetConfiguration(ctx context.Context) {
	slog.Info("Clearing existing WARP registration and configuration...")

	delCmd := exec.CommandContext(ctx, c.execPath, "registration", "delete")
	_ = delCmd.Run()

	resetCmd := exec.CommandContext(ctx, c.execPath, "settings", "reset")
	_ = resetCmd.Run()
}
func (c *Connector) RegisterConnector(ctx context.Context) error {
	if strings.TrimSpace(c.cfg.Token) == "" {
		return fmt.Errorf("warp registration token cannot be empty")
	}

	if err := c.EnsureBinary(ctx); err != nil {
		return err
	}

	c.ResetConfiguration(ctx)

	slog.Info("Applying MDM deployment token...")

	// Use 'mdm set-config' with --token flag
	cmd := exec.CommandContext(ctx, c.execPath, "mdm", "set-config", "--token", c.cfg.Token)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Fallback to 'registration token' if MDM token mode requires direct registration token CLI format
		slog.Warn("mdm set-config failed, attempting registration token fallback...")
		fallbackCmd := exec.CommandContext(ctx, c.execPath, "registration", "token", c.cfg.Token)
		fallbackCmd.Stdout = os.Stdout
		fallbackCmd.Stderr = os.Stderr

		if fallbackErr := fallbackCmd.Run(); fallbackErr != nil {
			return fmt.Errorf("failed to register warp token via mdm or registration token: %w", err)
		}
	}

	slog.Info("WARP MDM registration successfully completed")
	return nil
}

func (c *Connector) Connect(ctx context.Context) error {
	if err := c.EnsureBinary(ctx); err != nil {
		return err
	}

	slog.Info("Establishing WARP connection...")

	cmd := exec.CommandContext(ctx, c.execPath, "connect")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute warp-cli connect: %w", err)
	}

	slog.Info("WARP connect command executed successfully")
	return nil
}

func (c *Connector) Disconnect(ctx context.Context) error {
	if c.execPath == "" {
		path, err := c.FindWarpBinary()
		if err != nil {
			return err
		}
		c.execPath = path
	}

	cmd := exec.CommandContext(ctx, c.execPath, "disconnect")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
