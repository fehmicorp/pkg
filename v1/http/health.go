package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
