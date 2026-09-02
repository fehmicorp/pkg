package dns

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"p2pl4/server/config"

	"github.com/miekg/dns"
)

type DNSServer struct {
	config    config.DNSConfig
	cache     map[string][]dns.RR
	cacheMu   sync.RWMutex
	udpServer *dns.Server
	tcpServer *dns.Server
}

func NewDNSServer(cfg config.DNSConfig) *DNSServer {
	return &DNSServer{
		config: cfg,
		cache:  make(map[string][]dns.RR),
	}
}

func (d *DNSServer) Start(bindIP string) error {
	upstreamAddr := d.config.Server
	if !strings.Contains(upstreamAddr, ":") {
		upstreamAddr = upstreamAddr + ":53"
	}
	d.config.Server = upstreamAddr

	// Explicitly bind to 10.8.0.1:53
	listenAddr := fmt.Sprintf("%s:53", strings.TrimSpace(bindIP))
	dns.HandleFunc(".", d.handleQuery)

	d.udpServer = &dns.Server{Addr: listenAddr, Net: "udp"}
	d.tcpServer = &dns.Server{Addr: listenAddr, Net: "tcp"}

	go func() {
		log.Printf("[DNS] Launching UDP listener on %s", listenAddr)
		if err := d.udpServer.ListenAndServe(); err != nil {
			log.Printf("[DNS ERR] UDP listener crashed: %v", err)
		}
	}()

	go func() {
		log.Printf("[DNS] Launching TCP listener on %s", listenAddr)
		if err := d.tcpServer.ListenAndServe(); err != nil {
			log.Printf("[DNS ERR] TCP listener crashed: %v", err)
		}
	}()

	go func() {
		d.syncZones()
		d.scheduleZoneSync()
	}()

	return nil
}

func (d *DNSServer) Stop() {
	if d.udpServer != nil {
		_ = d.udpServer.Shutdown()
	}
	if d.tcpServer != nil {
		_ = d.tcpServer.Shutdown()
	}
}

// handleQuery answers incoming client queries from local cache or proxies directly to upstream
func (d *DNSServer) handleQuery(w dns.ResponseWriter, req *dns.Msg) {
	if len(req.Question) == 0 {
		return
	}

	msg := new(dns.Msg)
	msg.SetReply(req)
	msg.Authoritative = true

	for _, q := range req.Question {
		qName := strings.ToLower(q.Name)

		// 1. Check local AXFR cache first
		if records, found := d.lookupCache(qName, q.Qtype); found {
			msg.Answer = append(msg.Answer, records...)
			_ = w.WriteMsg(msg)
			return
		}

		// 2. Fallback: Proxy query directly to upstream
		resp, err := d.proxyQuery(req)
		if err == nil && resp != nil {
			_ = w.WriteMsg(resp)
			return
		} else {
			log.Printf("[DNS WARN] Direct proxy fallback for %s to %s failed: %v", qName, d.config.Server, err)
		}
	}

	_ = w.WriteMsg(msg)
}

func (d *DNSServer) lookupCache(qName string, qType uint16) ([]dns.RR, bool) {
	d.cacheMu.RLock()
	defer d.cacheMu.RUnlock()

	records, exists := d.cache[qName]
	if !exists {
		return nil, false
	}

	var matched []dns.RR
	for _, rr := range records {
		if rr.Header().Rrtype == qType || qType == dns.TypeANY {
			matched = append(matched, rr)
		}
	}

	return matched, len(matched) > 0
}

func (d *DNSServer) proxyQuery(req *dns.Msg) (*dns.Msg, error) {
	c := new(dns.Client)
	c.Timeout = 3 * time.Second
	resp, _, err := c.Exchange(req, d.config.Server)
	return resp, err
}

// scheduleZoneSync periodically refreshes selected zones
func (d *DNSServer) scheduleZoneSync() {
	interval := time.Duration(d.config.Refresh) * time.Minute
	if interval <= 0 {
		interval = 60 * time.Minute
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		log.Printf("[DNS] Periodic refresh interval (%d mins) reached. Fetching zones from %s...", d.config.Refresh, d.config.Server)
		d.syncZones()
	}
}

// syncZones executes Zone Transfer (AXFR) or logs fallback state
func (d *DNSServer) syncZones() {
	newCache := make(map[string][]dns.RR)

	for _, zone := range d.config.Zones {
		zoneName := dns.Fqdn(zone)
		records, err := d.fetchAXFR(zoneName)
		if err != nil {
			log.Printf("[DNS WARN] AXFR zone transfer failed for %s (%v). Active queries will direct-proxy to %s", zoneName, err, d.config.Server)
			continue
		}

		for _, rr := range records {
			hdr := rr.Header()
			ownerKey := strings.ToLower(hdr.Name)
			newCache[ownerKey] = append(newCache[ownerKey], rr)
		}
		log.Printf("[DNS] Successfully replicated zone %s (%d records cached)", zoneName, len(records))
	}

	d.cacheMu.Lock()
	d.cache = newCache
	d.cacheMu.Unlock()
}

// fetchAXFR attempts a DNS zone transfer with a 5-second context timeout
func (d *DNSServer) fetchAXFR(zone string) ([]dns.RR, error) {
	tr := new(dns.Transfer)
	m := new(dns.Msg)
	m.SetAxfr(zone)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	channel, err := tr.In(m, d.config.Server)
	if err != nil {
		return nil, err
	}

	var records []dns.RR
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("AXFR transfer context timeout")
		case env, ok := <-channel:
			if !ok {
				if len(records) == 0 {
					return nil, fmt.Errorf("empty AXFR response received")
				}
				return records, nil
			}
			if env.Error != nil {
				return nil, env.Error
			}
			records = append(records, env.RR...)
		}
	}
}
