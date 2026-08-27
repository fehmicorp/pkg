package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	Client *goredis.Client
}

func NewClient(
	host string,
	port int,
	user string,
	password string,
	db int,
) (*Client, error) {

	rdb := goredis.NewClient(&goredis.Options{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Password:     password,
		Username:     user,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{
		Client: rdb,
	}, nil
}

// Ping checks the connection status to the Redis server.
func (c *Client) Ping(ctx context.Context) *goredis.StatusCmd {
	return c.Client.Ping(ctx)
}

// Set writes a key-value pair with an expiration duration.
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *goredis.StatusCmd {
	return c.Client.Set(ctx, key, value, expiration)
}

// Get fetches the value for a given key.
func (c *Client) Get(ctx context.Context, key string) *goredis.StringCmd {
	return c.Client.Get(ctx, key)
}

// Close gracefully closes the client connection.
func (c *Client) Close() error {
	return c.Client.Close()
}
