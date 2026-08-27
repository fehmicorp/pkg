package main

import (
	"fmt"
	"strconv"

	"github.com/fehmicorp/pkg/v1/config"
)

func main() {
	config.Load()
	fmt.Println("------------------ Application ----------------------------")
	fmt.Printf("Name: %s\n", config.Conf.App.Name)
	fmt.Printf("Environment: %s\n", config.Conf.App.Environment)
	fmt.Printf("Version: %s\n", config.Conf.App.Version)
	fmt.Println("------------------ HTTP Server ----------------------------")
	fmt.Printf("Mode: %s\n", config.Conf.Server.Mode)
	fmt.Printf("Host: %s\n", config.Conf.Server.Host)
	fmt.Printf("Port: %s\n", strconv.Itoa(config.Conf.Server.Port))
	fmt.Printf("FQDN: %s\n", config.Conf.Server.FQDN)
	fmt.Println("------------------ Redis Client ----------------------------")
	fmt.Printf("Host: %s\n", config.Conf.Redis.Host)
	fmt.Printf("Port: %s\n", strconv.Itoa(config.Conf.Redis.Port))
	fmt.Printf("User: %s\n", config.Conf.Redis.User)
	fmt.Printf("Password: %s\n", config.Conf.Redis.Password)
	fmt.Printf("DB Id: %s\n", strconv.Itoa(config.Conf.Redis.DB))
}
