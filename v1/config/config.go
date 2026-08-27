package config

import "github.com/fehmicorp/pkg/v1/gateway"

type Config struct {
	App      gateway.Config `yaml:"app"`
	Server   http.Config    `yaml:"server"`
	Redis    RedisConfig    `yaml:"redis"`
	Services ServicesConfig `yaml:"services"`
	Logging  LoggingConfig  `yaml:"logging"`
}
