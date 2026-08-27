package http

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/fehmicorp/pkg/v1/config"
	"github.com/fehmicorp/pkg/v1/redis"
	"github.com/gin-gonic/gin"
)

type RouteRepository struct {
	client *redis.Client
}

type Server struct {
	Config    *config.Config
	RouteRepo *RouteRepository
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

func SetupRouter(app *Server, prefixes ...string) (r *gin.Engine) {
	router := gin.New()
	router.Use(
		gin.Logger(),
		gin.Recovery(),
	)
	var prefix string
	if len(prefixes) > 0 && prefixes[0] != "" {
		prefix = "/" + strings.Trim(prefixes[0], "/")
	}
	baseGroup := router.Group(prefix)
	{
		baseGroup.GET("/health", app.HealthHandler)
	}
	router.NoRoute(func(c *gin.Context) {
		if prefix != "" && !strings.HasPrefix(c.Request.URL.Path, prefix) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		app.RoutesHandler(c)
	})
	return router
}
