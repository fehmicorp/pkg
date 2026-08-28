package gateway

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/cloudflare/cloudflared/connection"
	"github.com/cloudflare/cloudflared/tunnelrpc/pogs"
	"github.com/fehmicorp/pkg/v1/config"
	"github.com/google/uuid"
)

type Connector struct {
	config config.ClientConfig
}

func NewConnector(cfg config.ClientConfig) *Connector {
	return &Connector{
		config: cfg,
	}
}

// Connect starts the client-side bridge proxying local requests to the remote private endpoint via Cloudflare Tunnel.
func (c *Connector) Connect(ctx context.Context) error {
	creds := c.config.Tunnel
	if creds.TunnelID == "" || creds.TunnelSecret == "" || creds.AccountTag == "" {
		return fmt.Errorf("missing required cloudflare tunnel credentials")
	}

	tunnelID, err := uuid.Parse(creds.TunnelID)
	if err != nil {
		return fmt.Errorf("invalid tunnel ID format: %w", err)
	}

	// 1. Configure Edge Transport Config
	connConfig := &connection.Config{
		TunnelID:   tunnelID,
		Secret:     []byte(creds.TunnelSecret),
		AccountTag: creds.AccountTag,
		TLSConfig:  &tls.Config{MinVersion: tls.VersionTLS12},
		MaxRetries: 5,
		ClientTags: []pogs.Tag{},
	}

	slog.Info("Establishing outbound Cloudflare Edge tunnel transport...", slog.String("tunnel_id", tunnelID.String()))
	observer := connection.NewObserver(slog.Default())

	edgeConn, err := connection.NewEdgeConn(connConfig, observer, 0)
	if err != nil {
		return fmt.Errorf("failed to open edge tunnel transport: %w", err)
	}

	// 2. Start edge connection worker
	go func() {
		if err := edgeConn.Serve(ctx); err != nil {
			slog.Error("Edge transport error", slog.String("error", err.Error()))
		}
	}()

	// 3. Bind local listener for client application traffic
	localAddr := creds.LocalBindAddr
	if localAddr == "" {
		localAddr = "127.0.0.1:9090"
	}

	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		edgeConn.Close()
		return fmt.Errorf("failed to bind local listener %s: %w", localAddr, err)
	}
	defer listener.Close()

	slog.Info("Client proxy listener ready",
		slog.String("local_bind", localAddr),
		slog.String("target_endpoint", creds.TargetEndpoint),
	)

	// Shutdown listener when context cancels
	go func() {
		<-ctx.Done()
		listener.Close()
		edgeConn.Close()
	}()

	// 4. Accept local connections and bridge to remote private server via Cloudflare
	for {
		localConn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				slog.Error("Failed to accept local connection", slog.String("error", err.Error()))
				continue
			}
		}

		go c.bridgeConnection(ctx, localConn, creds.TargetEndpoint)
	}
}

// DialTunnel establishes a direct HTTP/2 or QUIC stream over Cloudflare to the target endpoint.
func (c *Connector) GetHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				// Route outbound HTTP connections through local tunnel listener
				localAddr := c.config.Tunnel.LocalBindAddr
				if localAddr == "" {
					localAddr = "127.0.0.1:9090"
				}
				var dialer net.Dialer
				return dialer.DialContext(ctx, "tcp", localAddr)
			},
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Connector) bridgeConnection(ctx context.Context, localConn net.Conn, targetEndpoint string) {
	defer localConn.Close()

	// Dial remote endpoint via Cloudflare Tunnel virtual network / stream
	var dialer net.Dialer
	remoteConn, err := dialer.DialContext(ctx, "tcp", targetEndpoint)
	if err != nil {
		slog.Error("Failed to connect to target private endpoint", slog.String("target", targetEndpoint), slog.String("error", err.Error()))
		return
	}
	defer remoteConn.Close()

	// Bi-directional copy between local client socket and Cloudflare tunnel target connection
	errChan := make(chan error, 2)
	go func() {
		_, err := io.Copy(remoteConn, localConn)
		errChan <- err
	}()
	go func() {
		_, err := io.Copy(localConn, remoteConn)
		errChan <- err
	}()

	<-errChan
}
