package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fehmicorp/pkg/v1/utils/os"
)

const defaultBaseURL = "https://api.cloudflare.com/client/v4"

type Client struct {
	baseURL    string
	apiToken   string
	accountID  string
	httpClient *http.Client
}

func New(apiToken, accountID string) *Client {
	return &Client{
		baseURL:   defaultBaseURL,
		apiToken:  apiToken,
		accountID: accountID,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// DoRequest handles native HTTP calls, headers, marshaling, and error checking
func (c *Client) DoRequest(ctx context.Context, method, path string, body interface{}, out interface{}) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("dns: failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("dns: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("dns: http request failed: %w", err)
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("dns: failed to read response body: %w", err)
	}

	if out != nil {
		if err := json.Unmarshal(respData, out); err != nil {
			return fmt.Errorf("dns: failed to unmarshal response: %w (status: %d)", err, resp.StatusCode)
		}
	}

	return nil
}

func Fetch(args ...string) *Client {
	var token, id string
	if len(args) > 0 {
		token = args[0]
	}
	if len(args) > 1 {
		id = args[1]
	}
	_, apiToken := os.GetEnvV2("CF_API_TOKEN", token)
	_, accountID := os.GetEnvV2("CF_ACCOUNT_ID", id)
	return New(apiToken, accountID)
}
