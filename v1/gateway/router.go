package gateway

import gwconfig "github.com/fehmicorp/pkg/v1/gateway/config"

func StartServer() {
	cfg := gwconfig.Init()

	// redisClient, err := redis.NewClient(
	// 	cfg.Redis.Host,
	// 	int(cfg.Redis.Port),
	// 	cfg.Redis.User,
	// 	cfg.Redis.Password,
	// 	cfg.Redis.DB,
	// )

	// if err != nil {
	// 	slog.Error(
	// 		"failed to connect redis",
	// 		slog.String("error", err.Error()),
	// 	)
	// 	os.Exit(1)
	// }

	// server := httphandler.NewServer(cfg, redisClient.Client)
}
