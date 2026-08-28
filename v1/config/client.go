package config

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type ClientConfig struct {
	App    AppConfig   `yaml:"app" json:"app"`
	Tunnel TunnelCreds `yaml:"tunnel" json:"tunnel"`
}
type TunnelCreds struct {
	AccountTag     string `yaml:"account_tag" json:"account_tag"`
	TunnelID       string `yaml:"tunnel_id" json:"tunnel_id"`
	TunnelSecret   string `yaml:"tunnel_secret" json:"tunnel_secret"`
	TargetEndpoint string `yaml:"target_endpoint" json:"target_endpoint"` // e.g., "private-api.fehmi.cloud:443" or "10.0.0.5:8080"
	LocalBindAddr  string `yaml:"local_bind_addr" json:"local_bind_addr"` // e.g., "127.0.0.1:8080"
}

var CliConf *ClientConfig

func ClientInit() {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		workDir, err := os.Getwd()
		if err != nil {
			slog.Error("failed to get working directory", slog.String("error", err.Error()))
			os.Exit(1)
		}
		configPath = filepath.Join(workDir, "config.yaml")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg ClientConfig
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("can not read config file: %s", err.Error())
	}

	CliConf = &cfg
}
