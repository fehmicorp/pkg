package http

import (
	"log/slog"
	"os"

	"github.com/fehmicorp/pkg/v1/config"
	"github.com/fehmicorp/pkg/v1/redis"
)

type RouteRepository struct {
	client *redis.Client
}

func NewRouteRepository(client *redis.Client) *RouteRepository {
	return &RouteRepository{
		client: client,
	}
}

func RedisConnect() (client *redis.Client) {
	redisClient, err := redis.NewClient(
		config.Conf.Redis.Host,
		int(config.Conf.Redis.Port),
		config.Conf.Redis.User,
		config.Conf.Redis.Password,
		config.Conf.Redis.DB,
	)
	if err != nil {
		slog.Error(
			"failed to connect redis",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	return redisClient
}
