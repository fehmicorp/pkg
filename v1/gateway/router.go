package gateway

import (
	"fmt"
	"strconv"

	"github.com/fehmicorp/pkg/v1/config"
)

func StartServer(config *config.Config) {
	fmt.Println("------------------ Application ----------------------------")
	fmt.Printf("Name: %s\n", config.App.Name)
	fmt.Printf("Environment: %s\n", config.App.Environment)
	fmt.Printf("Version: %s\n", config.App.Version)
	fmt.Println("------------------ HTTP Server ----------------------------")
	fmt.Printf("Mode: %s\n", config.Server.Mode)
	fmt.Printf("Host: %s\n", config.Server.Host)
	fmt.Printf("Port: %s\n", strconv.Itoa(config.Server.Port))
	fmt.Printf("FQDN: %s\n", config.Server.FQDN)
	fmt.Println("------------------ Redis Client ----------------------------")
	fmt.Printf("Host: %s\n", config.Redis.Host)
	fmt.Printf("Port: %s\n", strconv.Itoa(config.Redis.Port))
	fmt.Printf("User: %s\n", config.Redis.User)
	fmt.Printf("Password: %s\n", config.Redis.Password)
	fmt.Printf("DB Id: %s\n", strconv.Itoa(config.Redis.DB))
}
