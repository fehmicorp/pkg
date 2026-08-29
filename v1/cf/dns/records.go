package dns

import (
	"context"
	"fmt"
	"net/http"
)

// ListDNSRecords fetches all DNS records within a zone
func (c *Client) ListDNSRecords(ctx context.Context, zoneID string) ([]DNSRecord, error) {
	path := fmt.Sprintf("/zones/%s/dns_records", zoneID)
	var resp APIResponse[[]DNSRecord]

	err := c.DoRequest(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("cloudflare error: %v", resp.Errors)
	}

	return resp.Result, nil
}

// GetDNSRecordByID retrieves a specific record
func (c *Client) GetDNSRecordByID(ctx context.Context, zoneID, recordID string) (*DNSRecord, error) {
	path := fmt.Sprintf("/zones/%s/dns_records/%s", zoneID, recordID)
	var resp APIResponse[DNSRecord]

	err := c.DoRequest(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("cloudflare error: %v", resp.Errors)
	}

	return &resp.Result, nil
}

// CreateDNSRecord inserts a new DNS record
func (c *Client) CreateDNSRecord(ctx context.Context, zoneID string, params CreateRecordParams) (*DNSRecord, error) {
	path := fmt.Sprintf("/zones/%s/dns_records", zoneID)
	var resp APIResponse[DNSRecord]

	err := c.DoRequest(ctx, http.MethodPost, path, params, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("cloudflare error: %v", resp.Errors)
	}

	return &resp.Result, nil
}

// PatchDNSRecord updates selected fields of a record
func (c *Client) PatchDNSRecord(ctx context.Context, zoneID, recordID string, params UpdateRecordParams) (*DNSRecord, error) {
	path := fmt.Sprintf("/zones/%s/dns_records/%s", zoneID, recordID)
	var resp APIResponse[DNSRecord]

	err := c.DoRequest(ctx, http.MethodPatch, path, params, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("cloudflare error: %v", resp.Errors)
	}

	return &resp.Result, nil
}

// OverwriteDNSRecord replaces a complete record using HTTP PUT
func (c *Client) OverwriteDNSRecord(ctx context.Context, zoneID, recordID string, params CreateRecordParams) (*DNSRecord, error) {
	path := fmt.Sprintf("/zones/%s/dns_records/%s", zoneID, recordID)
	var resp APIResponse[DNSRecord]

	err := c.DoRequest(ctx, http.MethodPut, path, params, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("cloudflare error: %v", resp.Errors)
	}

	return &resp.Result, nil
}

// DeleteDNSRecord removes a record
func (c *Client) DeleteDNSRecord(ctx context.Context, zoneID, recordID string) error {
	path := fmt.Sprintf("/zones/%s/dns_records/%s", zoneID, recordID)
	var resp APIResponse[DNSRecord]

	err := c.DoRequest(ctx, http.MethodDelete, path, nil, &resp)
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("cloudflare error: %v", resp.Errors)
	}

	return nil
}

// BatchEditDNSRecords handles multiple record operations in one request
func (c *Client) BatchEditDNSRecords(ctx context.Context, zoneID string, payload map[string]interface{}) (*BatchRecordResponse, error) {
	path := fmt.Sprintf("/zones/%s/dns_records/batch", zoneID)
	var resp APIResponse[BatchRecordResponse]

	err := c.DoRequest(ctx, http.MethodPost, path, payload, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("cloudflare error: %v", resp.Errors)
	}

	return &resp.Result, nil
}
