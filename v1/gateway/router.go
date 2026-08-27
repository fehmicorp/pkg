package gateway

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/fehmicorp/pkg/v1/config"
	"github.com/fehmicorp/pkg/v1/redis"
)

func StartServer() {
	TestRedis()
}

func TestRedis() {
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		slog.Error("redis ping failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("redis connection verified via PING")
}

func Print() {
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
