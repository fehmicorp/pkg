package config

import (
	"github.com/fehmicorp/pkg/v1/redis"
)

type Config struct {
	App    AppConfig    `yaml:"app" json:"app"`
	Server HttpConfig   `yaml:"http" json:"http"`
	Redis  redis.Config `yaml:"redis" json:"redis"`
}

type AppConfig struct {
	Name        string `yaml:"name" env:"APP_NAME" env-default:"gateway" json:"name"`
	Version     string `yaml:"version" env:"APP_VERSION" env-default:"1.0.0" json:"version"`
	Environment string `yaml:"environment" env:"APP_ENV" env-default:"development" json:"environment"`
}

type HttpConfig struct {
	Mode            string `yaml:"mode" env:"SERVER_MODE" env-default:"http" json:"mode"`
	Host            string `yaml:"host" env:"SERVER_HOST" env-default:"localhost" json:"host"`
	Port            int    `yaml:"port" env:"SERVER_PORT" env-default:"8050" json:"port"`
	CFTunnel        string `yaml:"cftoken,omitempty" env:"CF_TUNNEL_TOKEN" json:"cftoken,omitempty"`
	FQDN            string `yaml:"fqdn" env:"SERVER_FQDN" env-default:"https://api.fehmicorp.in/" json:"fqdn"`
	ReadTimeoutSec  int    `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT" env-default:"10" json:"read_timeout"`
	WriteTimeoutSec int    `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT" env-default:"10" json:"write_timeout"`
}

type LoggingConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	Dir   string `yaml:"dir" env:"LOG_DIR" env-default:"./logs"`
}

type MonitoringConfig struct {
	PrometheusEnabled bool `yaml:"prometheus_enabled" env:"PROM_ENABLED" env-default:"true"`
}
