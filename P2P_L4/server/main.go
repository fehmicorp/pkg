package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"p2pl4/server/cf"
	"p2pl4/server/config"
	"p2pl4/server/dns"
	"p2pl4/server/ipam"
	ssltun "p2pl4/server/tunnel"
)

func main() {
	conf := config.Init()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	serverIP := conf.DNS.IP
	GatewayIp := conf.TCP.Gateway
	if conf.TCP.Mode == "tunnel" {
		cidr := conf.TCP.NetworkPool // "10.8.0.0/24"
		port := fmt.Sprintf(":%d", conf.TCP.Port)

		pool, err := ipam.NewPool(cidr)
		if err != nil {
			log.Fatalf("[SERVER ERR] Failed to initialize IPAM pool: %v", err)
		}

		// 1. Initialize TUN interface and set IP 10.8.0.1/24
		go ssltun.SetupTunnel(pool, port, GatewayIp)

		// 2. Start Cloudflare Tunnel Daemon
		token := conf.TCP.CFTunnel
		if token != "" {
			go func() {
				if err := cf.StartTunnelProcess(ctx, token); err != nil {
					slog.Error("Failed to launch cloudflared process", slog.String("error", err.Error()))
				}
			}()
		}

		// 3. BLOCK until tun0 interface is created and active with 10.8.0.1
		waitForInterfaceIP(serverIP, 10*time.Second)
	}

	// 4. Start DNS listener directly on 10.8.0.1:53 NOW that interface exists
	dnsServer := dns.NewDNSServer(conf.DNS)
	if err := dnsServer.Start(serverIP); err != nil {
		log.Fatalf("[SERVER ERR] Failed to start DNS server: %v", err)
	}
	defer dnsServer.Stop()

	slog.Info("Server running in isolated mode over CF tunnel. Press Ctrl+C to stop.")
	<-ctx.Done()
}

// Blocks execution until target IP is actively bound to a network interface
func waitForInterfaceIP(targetIP string, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		ifaces, err := net.Interfaces()
		if err == nil {
			for _, iface := range ifaces {
				addrs, err := iface.Addrs()
				if err != nil {
					continue
				}
				for _, addr := range addrs {
					if ipnet, ok := addr.(*net.IPNet); ok {
						if ipnet.IP.String() == targetIP {
							log.Printf("[TUNNEL] Interface %s active with IP %s", iface.Name, targetIP)
							return
						}
					}
				}
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	log.Printf("[TUNNEL WARN] Timeout waiting for IP %s on TUN interface; proceeding with bind", targetIP)
}
