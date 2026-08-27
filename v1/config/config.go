package config

import (
	"github.com/fehmicorp/pkg/v1/gateway"
	"github.com/fehmicorp/pkg/v1/http"
	"github.com/fehmicorp/pkg/v1/redis"
)

type Config struct {
	App     gateway.Config   `yaml:"app"`
	Server  http.Config      `yaml:"server"`
	Redis   redis.Config     `yaml:"redis"`
	Monitor MonitoringConfig `yaml:"monitor"`
	Logging LoggingConfig    `yaml:"logging"`
}

type LoggingConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	Dir   string `yaml:"dir" env:"LOG_DIR" env-default:"./logs"`
}

type MonitoringConfig struct {
	PrometheusEnabled bool `yaml:"prometheus_enabled" env:"PROM_ENABLED" env-default:"true"`
}
