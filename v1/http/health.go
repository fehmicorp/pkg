package http

import (
	"net/http"

	"github.com/fehmicorp/pkg/v1/gateway"
	"github.com/gin-gonic/gin"
)

func (s *gateway.Config) HealthHandler(
	c *gin.Context,
) {
	responseJSON(
		c,
		http.StatusOK,
		gin.H{
			"status":      true,
			"service":     s.Name,
			"version":     s.Version,
			"environment": s.Environment,
		},
	)
}
