package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fehmicorp/go/pkg/v1/win/notify"
	"github.com/gorilla/websocket"
	"golang.org/x/sys/windows"
	"golang.zx2c4.com/wintun"
)

type Config struct {
	IP      string
	Netmask string
	Gateway string
	DNS     string
}

func runTunnel(ctx context.Context) {
	if runtime.GOOS == "windows" {
		runWindowsClient(ctx)
	} else {
		runLinuxClient(ctx)
	}
}

func runWindowsClient(ctx context.Context) {
	// 1. Ensure wintun.dll exists before attempting to load driver
	if err := ensureWintunDLL(); err != nil {
		log.Printf("[CLIENT WARN] Could not download wintun.dll automatically: %v", err)
	}

	workDir, err := os.Getwd()
	if err == nil {
		_, _ = windows.LoadLibrary(filepath.Join(workDir, "wintun.dll"))
	}

	// 2. Create or open Wintun adapter
	adapter, err := wintun.CreateAdapter("P2PTunnel", "Wintun", nil)
	if err != nil {
		log.Fatalf("[CLIENT ERR] Failed to create Wintun adapter (Run as Administrator!): %v", err)
	}
	defer adapter.Close()

	log.Printf("[CLIENT] Wintun interface created successfully")
	notify.SimpleAlert(Conf.AppName, "⚠️ SSL Connection Established..!")
	session, err := adapter.StartSession(0x400000) // 4MB buffer
	if err != nil {
		log.Fatalf("[CLIENT ERR] Failed to start Wintun session: %v", err)
	}
	defer session.End()

	// 3. Connect to Cloudflare Tunnel Edge
	u := url.URL{Scheme: "wss", Host: ServerURL, Path: "/"}
	log.Printf("[CLIENT] Connecting to %s...", u.String())

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		log.Fatalf("[CLIENT ERR] WebSocket connection failed: %v", err)
	}
	defer conn.Close()
	log.Println("[CLIENT] Connected successfully to Cloudflare edge server!")

	// 4. Perform Handshake
	cfg, err := performHandshake(conn)
	if err != nil {
		log.Fatalf("[CLIENT ERR] Failed IP configuration handshake: %v", err)
	}

	log.Printf("[CLIENT] Received Config: IP=%s | Mask=%s | GW=%s | DNS=%s",
		cfg.IP, cfg.Netmask, cfg.Gateway, cfg.DNS)

	// 5. Apply assigned IP address, Gateway, and DNS via netsh
	cmdIP := exec.Command("netsh", "interface", "ipv4", "set", "address",
		"name=P2PTunnel", "static", cfg.IP, cfg.Netmask, cfg.Gateway)
	if err := cmdIP.Run(); err != nil {
		log.Printf("[CLIENT WARN] Failed to set interface IP/Gateway via netsh: %v", err)
	} else {
		log.Printf("[CLIENT] Network adapter configured with IP: %s, GW: %s", cfg.IP, cfg.Gateway)
	}

	cmdDNS := exec.Command("netsh", "interface", "ipv4", "set", "dnsservers",
		"name=P2PTunnel", "static", cfg.DNS, "primary")
	if err := cmdDNS.Run(); err != nil {
		log.Printf("[CLIENT WARN] Failed to set primary DNS via netsh: %v", err)
	} else {
		log.Printf("[CLIENT] Network adapter configured with DNS: %s", cfg.DNS)
	}

	// 6. Configure Split-DNS (NRPT) for .shalimarcorp.org
	setupSplitDNS(cfg.DNS)

	// 7. WebSocket Heartbeat Loop (Prevents Cloudflare 1006 connection drops)
	go func() {
		ticker := time.NewTicker(25 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(5*time.Second)); err != nil {
					return
				}
			}
		}
	}()

	// 8. Forward Wintun Packets -> WebSocket
	go func() {
		for {
			packet, err := session.ReceivePacket()
			if err != nil {
				if err == windows.ERROR_HANDLE_EOF {
					return
				}
				continue
			}

			payload := make([]byte, len(packet))
			copy(payload, packet)
			session.ReleaseReceivePacket(packet)

			err = conn.WriteMessage(websocket.BinaryMessage, wrapFrame(payload))
			if err != nil {
				log.Printf("[CLIENT ERR] WebSocket write failed: %v", err)
				return
			}
		}
	}()

	// 9. Forward WebSocket Packets -> Wintun
	go func() {
		for {
			messageType, raw, err := conn.ReadMessage()
			if err != nil {
				log.Printf("[CLIENT ERR] WebSocket read failed: %v", err)
				return
			}

			if messageType == websocket.BinaryMessage {
				packet, err := unwrapFrame(raw)
				if err != nil {
					continue
				}

				if strings.HasPrefix(string(packet), "CONFIG") {
					continue
				}

				allocatedPacket, err := session.AllocateSendPacket(len(packet))
				if err == nil {
					copy(allocatedPacket, packet)
					session.SendPacket(allocatedPacket)
				}
			}
		}
	}()

	<-ctx.Done()
	log.Println("[CLIENT] Shutting down client...")
	cleanupSplitDNS()
}

// -----------------------------------------------------------------------------
// Handshake & Dependency Helpers
// -----------------------------------------------------------------------------

func performHandshake(conn *websocket.Conn) (*Config, error) {
	messageType, raw, err := conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("failed reading handshake message: %w", err)
	}

	if messageType != websocket.BinaryMessage {
		return nil, fmt.Errorf("expected binary frame during handshake, got type %d", messageType)
	}

	payload, err := unwrapFrame(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid frame format during handshake: %w", err)
	}

	msg := string(payload)
	if !strings.HasPrefix(msg, "CONFIG|") {
		return nil, fmt.Errorf("unexpected handshake frame payload: %s", msg)
	}

	cfg := &Config{}
	parts := strings.Split(msg, "|")
	for _, part := range parts[1:] {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "IP":
			cfg.IP = kv[1]
		case "NETMASK":
			cfg.Netmask = kv[1]
		case "GATEWAY":
			cfg.Gateway = kv[1]
		case "DNS":
			cfg.DNS = kv[1]
		}
	}

	if cfg.IP == "" || cfg.Netmask == "" {
		return nil, fmt.Errorf("incomplete CONFIG payload received: %s", msg)
	}

	return cfg, nil
}

func setupSplitDNS(dnsIP string) {
	psScript := fmt.Sprintf(`
		Get-DnsClientNrptRule | Where-Namespace -eq ".shalimarcorp.org" | Remove-DnsClientNrptRule -Force -ErrorAction SilentlyContinue;
		Add-DnsClientNrptRule -Namespace ".shalimarcorp.org" -NameServers "%s" -ErrorAction SilentlyContinue;
		Clear-DnsClientCache;
	`, dnsIP)

	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psScript)
	if err := cmd.Run(); err != nil {
		log.Printf("[CLIENT WARN] Failed to add NRPT Split-DNS rule: %v", err)
	} else {
		log.Printf("[CLIENT] NRPT Split-DNS rule registered for .shalimarcorp.org -> %s", dnsIP)
	}
}

func cleanupSplitDNS() {
	psScript := `
		Get-DnsClientNrptRule | Where-Namespace -eq ".shalimarcorp.org" | Remove-DnsClientNrptRule -Force -ErrorAction SilentlyContinue;
		Clear-DnsClientCache;
	`
	_ = exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psScript).Run()
}

func ensureWintunDLL() error {
	dllName := "wintun.dll"
	if _, err := os.Stat(dllName); err == nil {
		return nil
	}

	arch := "amd64"
	if runtime.GOARCH == "arm64" {
		arch = "arm64"
	} else if runtime.GOARCH == "386" {
		arch = "x86"
	}

	log.Println("[CLIENT] wintun.dll not found. Downloading official Wintun release package...")

	psCommand := fmt.Sprintf(`
		$ProgressPreference = 'SilentlyContinue';
		[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12;
		Invoke-WebRequest -Uri "https://www.wintun.net/builds/wintun-0.14.1.zip" -OutFile "wintun.zip";
		Expand-Archive -Path "wintun.zip" -DestinationPath "wintun_temp" -Force;
		Copy-Item -Path "wintun_temp\wintun\bin\%s\wintun.dll" -Destination ".\wintun.dll" -Force;
		Remove-Item -Path "wintun.zip", "wintun_temp" -Recurse -Force;
	`, arch)

	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to download and extract wintun.dll via PowerShell: %w", err)
	}

	log.Println("[CLIENT] Successfully downloaded and placed wintun.dll!")
	return nil
}

// -----------------------------------------------------------------------------
// Linux/macOS Fallback
// -----------------------------------------------------------------------------

func runLinuxClient(ctx context.Context) {
	log.Println("[CLIENT] Linux/macOS execution mode engaged")
}

// Framing Helpers
func wrapFrame(payload []byte) []byte {
	length := uint16(len(payload))
	frame := make([]byte, 2+len(payload))
	binary.BigEndian.PutUint16(frame[0:2], length)
	copy(frame[2:], payload)
	return frame
}

func unwrapFrame(frame []byte) ([]byte, error) {
	if len(frame) < 2 {
		return nil, io.ErrUnexpectedEOF
	}
	length := binary.BigEndian.Uint16(frame[0:2])
	if int(length) > len(frame)-2 {
		return nil, io.ErrShortBuffer
	}
	return frame[2 : 2+length], nil
}
