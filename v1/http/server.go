package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/fehmicorp/pkg/v1/config"
	"github.com/fehmicorp/pkg/v1/redis"
	"github.com/gin-gonic/gin"
)

func NewServer(
	redisCli *redis.Client,
) *http.Server {
	if config.Conf.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	app := &Server{
		Config: config.Conf,
		RouteRepo: NewRouteRepository(
			&redis.Client{
				Client: redisCli.Client,
			},
		),
	}
	router := gin.New()
	router.Use(
		gin.Logger(),
		gin.Recovery(),
	)
	router.GET(
		"/ping",
		app.HealthHandler,
	)
	return &http.Server{
		Addr: config.Conf.Server.Host +
			":" +
			strconv.Itoa(config.Conf.Server.Port),

		Handler: router,

		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}

func (s *Server) HealthHandler(
	c *gin.Context,
) {
	responseJSON(
		c,
		http.StatusOK,
		gin.H{
			"status":      true,
			"service":     s.Config.App.Name,
			"version":     s.Config.App.Version,
			"environment": s.Config.App.Environment,
		},
	)
}

func responseJSON(
	c *gin.Context,
	statusCode int,
	payload any,
) {
	c.JSON(statusCode, payload)
}
