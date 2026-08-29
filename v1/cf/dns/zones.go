package dns

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// GetZones lists all zones accessible by the authenticated API key/token
func (c *Client) GetZones(ctx context.Context) ([]Zone, error) {
	path := "/zones"
	var resp APIResponse[[]Zone]

	err := c.DoRequest(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("cloudflare error: %v", resp.Errors)
	}

	return resp.Result, nil
}

// GetZoneByID fetches details for a single zone by ID
func (c *Client) GetZoneByID(ctx context.Context, zoneID string) (*Zone, error) {
	path := fmt.Sprintf("/zones/%s", zoneID)
	var resp APIResponse[Zone]

	err := c.DoRequest(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("cloudflare error: %v", resp.Errors)
	}

	return &resp.Result, nil
}

// CreateZone provisions a new domain inside Cloudflare
func (c *Client) CreateZone(ctx context.Context, params CreateZoneParams) (*Zone, error) {
	path := "/zones"
	var resp APIResponse[Zone]

	err := c.DoRequest(ctx, http.MethodPost, path, params, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("cloudflare error: %v", resp.Errors)
	}

	return &resp.Result, nil
}

// UpdateZone edits attributes of a zone (e.g., pause status)
func (c *Client) UpdateZone(ctx context.Context, zoneID string, params UpdateZoneParams) (*Zone, error) {
	path := fmt.Sprintf("/zones/%s", zoneID)
	var resp APIResponse[Zone]

	err := c.DoRequest(ctx, http.MethodPatch, path, params, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("cloudflare error: %v", resp.Errors)
	}

	return &resp.Result, nil
}

// DeleteZone drops a zone from Cloudflare
func (c *Client) DeleteZone(ctx context.Context, zoneID string) (string, error) {
	path := fmt.Sprintf("/zones/%s", zoneID)
	type deleteResult struct {
		ID string `json:"id"`
	}
	var resp APIResponse[deleteResult]

	err := c.DoRequest(ctx, http.MethodDelete, path, nil, &resp)
	if err != nil {
		return "", err
	}
	if !resp.Success {
		return "", fmt.Errorf("cloudflare error: %v", resp.Errors)
	}

	return resp.Result.ID, nil
}

// PurgeCache invalidates cached assets inside a specific zone
func (c *Client) PurgeCache(ctx context.Context, zoneID string, params PurgeCacheParams) (*PurgeCacheResponse, error) {
	path := fmt.Sprintf("/zones/%s/purge_cache", zoneID)
	var resp APIResponse[PurgeCacheResponse]

	err := c.DoRequest(ctx, http.MethodPost, path, params, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("cloudflare error: %v", resp.Errors)
	}

	return &resp.Result, nil
}

func GetZones(ctx context.Context, client *Client) []Zone {
	zones, err := client.GetZones(ctx)
	if err != nil {
		log.Fatalf("Error getting zones: %v", err)
	}
	return zones
}

func GetZoneIDs(ctx context.Context, client *Client) []string {
	zones, err := client.GetZones(ctx)
	if err != nil {
		log.Fatalf("Error getting zones: %v", err)
	}
	ids := make([]string, len(zones))
	for i, zone := range zones {
		ids[i] = zone.ID
	}
	return ids
}

func GetZonesNames(ctx context.Context, client *Client) []string {
	zones, err := client.GetZones(ctx)
	if err != nil {
		log.Fatalf("Error getting zones: %v", err)
	}
	names := make([]string, len(zones))
	for i, zone := range zones {
		names[i] = zone.Name
	}
	return names
}
