package http

import "github.com/fehmicorp/pkg/v1/redis"

type RouteRepository struct {
	client *redis.Client
}

func NewRouteRepository(client *redis.Client) *RouteRepository {
	return &RouteRepository{
		client: client,
	}
}
