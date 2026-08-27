package http

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/fehmicorp/pkg/v1/config"
	"github.com/fehmicorp/pkg/v1/http/utils/header"
	"github.com/gin-gonic/gin"
)

func (s *Server) RoutesHandler(c *gin.Context) {
	ctx := c.Request.Context()
	path := c.Request.URL.Path
	route, err := s.RouteRepo.GetRouteByPath(ctx, path)
	if err != nil {
		responseJSON(c, http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		return
	}
}

func (r *RouteRepository) GetRouteByPath(ctx context.Context, path string) (*header.RouteConfig, error) {
	if path == "" {
		return nil, errors.New("path is empty")
	}
	initialRedisKey := config.GetEnvOrDefault("INITKEY", "gateway:api:")
	key := initialRedisKey + strings.ReplaceAll(strings.TrimPrefix(path, "/"), "/", ":")
	payload, err := r.client.Client.Get(ctx, key).Result()
}
