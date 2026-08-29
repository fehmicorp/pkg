package dns

import "time"

// Standard Cloudflare API Response Wrapper
type APIResponse[T any] struct {
	Success  bool       `json:"success"`
	Errors   []APIError `json:"errors"`
	Messages []string   `json:"messages"`
	Result   T          `json:"result"`
}

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Zone represents a Cloudflare DNS Zone
type Zone struct {
	ID                  string    `json:"id,omitempty"`
	Name                string    `json:"name"`
	Status              string    `json:"status,omitempty"`
	Paused              bool      `json:"paused,omitempty"`
	Type                string    `json:"type,omitempty"`
	DevelopmentMode     int       `json:"development_mode,omitempty"`
	NameServers         []string  `json:"name_servers,omitempty"`
	OriginalNameServers []string  `json:"original_name_servers,omitempty"`
	ModifiedOn          time.Time `json:"modified_on,omitempty"`
	CreatedOn           time.Time `json:"created_on,omitempty"`
	Account             Account   `json:"account"`
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type CreateZoneParams struct {
	Name      string  `json:"name"`
	Account   Account `json:"account"`
	Type      string  `json:"type,omitempty"` // "full" or "partial"
	JumpStart bool    `json:"jump_start,omitempty"`
}

type UpdateZoneParams struct {
	Paused bool   `json:"paused,omitempty"`
	Type   string `json:"type,omitempty"`
}

type PurgeCacheParams struct {
	PurgeEverything bool     `json:"purge_everything,omitempty"`
	Files           []string `json:"files,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	Hosts           []string `json:"hosts,omitempty"`
}

type PurgeCacheResponse struct {
	ID string `json:"id"`
}

// DNSRecord represents a Cloudflare DNS Record
type DNSRecord struct {
	ID         string    `json:"id,omitempty"`
	Type       string    `json:"type"`    // "A", "AAAA", "CNAME", "TXT", "MX", "NS", "SRV", "CAA"
	Name       string    `json:"name"`    // "subdomain.example.com" or "example.com"
	Content    string    `json:"content"` // IP address or target domain
	Proxiable  bool      `json:"proxiable,omitempty"`
	Proxied    *bool     `json:"proxied,omitempty"` // true routes through Cloudflare proxy
	TTL        int       `json:"ttl,omitempty"`     // 1 for Auto, or seconds (e.g., 3600)
	Locked     bool      `json:"locked,omitempty"`
	ZoneID     string    `json:"zone_id,omitempty"`
	ZoneName   string    `json:"zone_name,omitempty"`
	CreatedOn  time.Time `json:"created_on,omitempty"`
	ModifiedOn time.Time `json:"modified_on,omitempty"`
	Priority   *uint16   `json:"priority,omitempty"` // Optional priority for MX/SRV records
}

type CreateRecordParams struct {
	Type     string  `json:"type"`
	Name     string  `json:"name"`
	Content  string  `json:"content"`
	TTL      int     `json:"ttl,omitempty"`
	Proxied  *bool   `json:"proxied,omitempty"`
	Priority *uint16 `json:"priority,omitempty"`
}

type UpdateRecordParams struct {
	Type     string  `json:"type,omitempty"`
	Name     string  `json:"name,omitempty"`
	Content  string  `json:"content,omitempty"`
	TTL      int     `json:"ttl,omitempty"`
	Proxied  *bool   `json:"proxied,omitempty"`
	Priority *uint16 `json:"priority,omitempty"`
}

type BatchRecordResponse struct {
	Posts   []DNSRecord `json:"posts,omitempty"`
	Puts    []DNSRecord `json:"puts,omitempty"`
	Patches []DNSRecord `json:"patches,omitempty"`
	Deletes []DNSRecord `json:"deletes,omitempty"`
}
