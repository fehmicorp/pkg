package ssltun

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"sync"

	"p2pl4/server/ipam"

	"github.com/songgao/water"
)

var (
	clients   = make(map[string]net.Conn)
	clientsMu sync.RWMutex
	gatewayIP string
)

func SetupTunnel(pool *ipam.Pool, listenPort string, gwIp string) {
	cfg := water.Config{
		DeviceType: water.TUN,
	}

	iface, err := water.New(cfg)
	if err != nil {
		log.Fatalf("[SERVER ERR] Failed to create TUN interface: %v", err)
	}
	log.Printf("[SERVER] TUN Interface Created: %s", iface.Name())

	// Allocate the first IP (10.8.0.1) from IPAM to reserve it for the Gateway/DNS
	reservedIP, err := pool.Allocate()
	if err != nil {
		log.Printf("[IPAM WARN] Failed to allocate gateway IP from pool: %v", err)
	} else {
		log.Printf("[IPAM] Allocated Gateway/DNS IP: %s", reservedIP)
	}

	// Configure server TUN interface with Gateway IP (e.g., 10.8.0.1/24)
	gatewayIP = gwIp
	cidrMask := pool.GetCIDRMask()
	gatewayCIDR := fmt.Sprintf("%s/%d", gatewayIP, cidrMask)

	if err := exec.Command("ip", "addr", "add", gatewayCIDR, "dev", iface.Name()).Run(); err != nil {
		log.Printf("[SERVER WARN] Failed to assign IP to TUN interface: %v", err)
	}
	if err := exec.Command("ip", "link", "set", "dev", iface.Name(), "up").Run(); err != nil {
		log.Printf("[SERVER WARN] Failed to bring TUN interface up: %v", err)
	}

	log.Printf("[SERVER] Configured TUN interface with Gateway IP %s", gatewayCIDR)

	// Start packet reader loop: TUN -> Client Routing Table
	go handleTUNToClients(iface)

	listener, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatalf("[SERVER ERR] Failed to bind TCP listener on %s: %v", listenPort, err)
	}
	defer listener.Close()

	log.Printf("[SERVER] Edge Tunnel Server listening on %s", listenPort)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[SERVER ERR] Accept error: %v", err)
			continue
		}
		log.Printf("[SERVER] Client connected from %s", conn.RemoteAddr())
		go handleClient(conn, iface, pool)
	}
}

func handleClient(conn net.Conn, iface *water.Interface, pool *ipam.Pool) {
	defer conn.Close()

	// 1. Allocate dynamic IP for the connected client (will start from 10.8.0.2)
	assignedIP, err := pool.Allocate()
	if err != nil {
		log.Printf("[SERVER ERR] Dynamic IP allocation failed for %s: %v", conn.RemoteAddr(), err)
		return
	}
	defer pool.Release(assignedIP)

	log.Printf("[SERVER] Assigned Dynamic IP %s to client %s", assignedIP, conn.RemoteAddr())

	// 2. Register active connection in routing table
	registerClient(assignedIP, conn)
	defer unregisterClient(assignedIP)

	// 3. Send Handshake Config frame containing IP, Netmask, Gateway, and DNS to Client
	dnsIP := gatewayIP // DNS replica runs locally on Gateway IP (10.8.0.1)
	configMsg := fmt.Sprintf("CONFIG|IP:%s|NETMASK:%s|GATEWAY:%s|DNS:%s", assignedIP, pool.GetNetmask(), gatewayIP, dnsIP)

	if err := writeFrame(conn, []byte(configMsg)); err != nil {
		log.Printf("[SERVER ERR] Failed to send CONFIG handshake to %s: %v", assignedIP, err)
		return
	}

	log.Printf("[SERVER] Handshake sent to %s -> IP: %s | GW: %s | DNS: %s",
		conn.RemoteAddr(), assignedIP, gatewayIP, dnsIP)

	// 4. Client -> TUN Read Loop
	for {
		packet, err := readFrame(conn)
		if err != nil {
			log.Printf("[SERVER] Client disconnected (%s - %s): %v", assignedIP, conn.RemoteAddr(), err)
			return
		}

		// Write packet directly to TUN interface for processing/routing
		if _, err := iface.Write(packet); err != nil {
			log.Printf("[SERVER ERR] Error writing packet to TUN interface: %v", err)
		}
	}
}

// Reads IP packets from TUN interface and dispatches them to the targeted client connection
func handleTUNToClients(iface *water.Interface) {
	buf := make([]byte, 2048)
	for {
		n, err := iface.Read(buf)
		if err != nil {
			log.Printf("[SERVER ERR] TUN read error: %v", err)
			return
		}

		packet := buf[:n]

		// Ensure it's an IPv4 packet (version field in IP header == 4)
		if n < 20 || (packet[0]>>4) != 4 {
			continue
		}

		// Extract destination IP from IPv4 header (bytes 16-19)
		destIP := net.IP(packet[16:20]).String()

		// Route packet to specific client connection
		clientsMu.RLock()
		clientConn, exists := clients[destIP]
		clientsMu.RUnlock()

		if exists {
			if err := writeFrame(clientConn, packet); err != nil {
				log.Printf("[SERVER ERR] Failed writing packet to client %s: %v", destIP, err)
			}
		}
	}
}

// Helper methods for managing client routing table
func registerClient(ip string, conn net.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	clients[ip] = conn
}

func unregisterClient(ip string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	delete(clients, ip)
}

// Framing Helpers
func writeFrame(w io.Writer, payload []byte) error {
	length := uint16(len(payload))
	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return err
	}
	_, err := w.Write(payload)
	return err
}

func readFrame(r io.Reader) ([]byte, error) {
	var length uint16
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return nil, err
	}
	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}
	return payload, nil
}
