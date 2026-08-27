package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/fehmicorp/pkg/v1/config"
	"github.com/fehmicorp/pkg/v1/redis"
	"github.com/gin-gonic/gin"
)

func NewServer() *http.Server {
	redisClient := RedisConnect()
	if config.Conf.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	app := &Server{
		Config: config.Conf,
		RouteRepo: NewRouteRepository(
			&redis.Client{
				Client: redisClient.Client,
			},
		),
	}
	router := gin.New()
	router.Use(
		gin.Logger(),
		gin.Recovery(),
	)
	router.GET(
		"/health",
		app.HealthHandler,
	)
	return &http.Server{
		Addr: config.Conf.Server.Host +
			":" +
			strconv.Itoa(config.Conf.Server.Port),
		Handler:           router,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}

func responseJSON(
	c *gin.Context,
	statusCode int,
	payload any,
) {
	c.JSON(statusCode, payload)
}
