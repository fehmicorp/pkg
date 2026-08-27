package http

import (
	"net/http"

	gateway_config "github.com/fehmicorp/pkg/v1/config"
)

type Server struct {
	Config *gateway_config.Config
}

func NewServer() *http.Server {

}
