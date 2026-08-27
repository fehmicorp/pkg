package http

import "github.com/fehmicorp/pkg/v1/config"

type Server struct {
	Config    *config.Config
	RouteRepo *RouteRepository
}
