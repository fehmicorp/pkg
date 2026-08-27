package config

import (
	"github.com/fehmicorp/pkg/v1/redis"
)

type Config struct {
	App    AppConfig    `yaml:"app"`
	Server HttpConfig   `yaml:"server"`
	Redis  redis.Config `yaml:"redis"`
	// Monitor MonitoringConfig `yaml:"monitor"`
	// Logging LoggingConfig    `yaml:"logging"`
}

type AppConfig struct {
	Name        string `yaml:"name" env:"APP_NAME" env-default:"gateway"`
	Version     string `yaml:"version" env:"APP_VERSION" env-default:"1.0.0"`
	Environment string `yaml:"environment" env:"APP_ENV" env-default:"development"`
}

type HttpConfig struct {
	Mode            string `yaml:"mode" env:"SERVER_MODE" env-default:"http"`
	Host            string `yaml:"host" env:"SERVER_HOST" env-default:"localhost"`
	Port            int    `yaml:"port" env:"SERVER_PORT" env-default:"8050"`
	FQDN            string `yaml:"fqdn" env:"SERVER_FQDN" env-default:"https://api.fehmicorp.in/"`
	ReadTimeoutSec  int    `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT" env-default:"10"`
	WriteTimeoutSec int    `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT" env-default:"10"`
}

type LoggingConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	Dir   string `yaml:"dir" env:"LOG_DIR" env-default:"./logs"`
}

type MonitoringConfig struct {
	PrometheusEnabled bool `yaml:"prometheus_enabled" env:"PROM_ENABLED" env-default:"true"`
}
