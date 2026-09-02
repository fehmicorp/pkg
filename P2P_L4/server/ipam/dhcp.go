package ipam

import (
	"fmt"
	"net"
	"sync"
)

type Pool struct {
	mu       sync.Mutex
	subnet   *net.IPNet
	assigned map[string]bool
	gateway  net.IP
}

// GetGatewayIP returns the gateway IP as a string (e.g., "10.0.0.1")
func (p *Pool) GetGatewayIP() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.gateway.String()
}

// GetCIDRMask returns the CIDR prefix length as an integer (e.g., 24 for /24)
func (p *Pool) GetCIDRMask() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	ones, _ := p.subnet.Mask.Size()
	return ones
}

func NewPool(cidr string) (*Pool, error) {
	_, subnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR: %w", err)
	}

	gateway := incrementIP(subnet.IP) // 10.0.0.1 reserved for server
	assigned := make(map[string]bool)
	assigned[gateway.String()] = true

	return &Pool{
		subnet:   subnet,
		assigned: assigned,
		gateway:  gateway,
	}, nil
}

func (p *Pool) Allocate() (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	ip := make(net.IP, len(p.gateway))
	copy(ip, p.gateway)

	for {
		ip = incrementIP(ip)
		if !p.subnet.Contains(ip) {
			return "", fmt.Errorf("ip pool exhausted")
		}

		// Skip network IP, gateway, and broadcast IP
		if ip[len(ip)-1] == 255 {
			continue
		}

		ipStr := ip.String()
		if !p.assigned[ipStr] {
			p.assigned[ipStr] = true
			return ipStr, nil
		}
	}
}

func (p *Pool) Release(ipStr string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.assigned, ipStr)
}

func (p *Pool) GetNetmask() string {
	mask := p.subnet.Mask
	return fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])
}

func incrementIP(ip net.IP) net.IP {
	next := make(net.IP, len(ip))
	copy(next, ip)
	for i := len(next) - 1; i >= 0; i-- {
		next[i]++
		if next[i] > 0 {
			break
		}
	}
	return next
}
