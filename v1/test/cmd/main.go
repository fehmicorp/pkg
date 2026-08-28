package main

import (
	"github.com/fehmicorp/pkg/v1/config"
	"github.com/fehmicorp/pkg/v1/gateway"
)

func main() {
	config.Init()
	gateway.StartServer()
}
