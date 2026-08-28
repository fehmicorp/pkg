package config

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

var Conf *Config

func Init() {
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
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("can not read config file: %s", err.Error())
	}
	Conf = &cfg
}

func GetEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
