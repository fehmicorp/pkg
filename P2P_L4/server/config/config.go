package config

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App AppConfig `yaml:"app" json:"app"`
	TCP TCPConfig `yaml:"tcp,omitempty" json:"tcp,omitempty"`
	DNS DNSConfig `yaml:"dns,omitempty" json:"dns,omitempty"`
}
type HandshakeConfig struct {
	IP      string `json:"ip"`
	Netmask string `json:"netmask"`
}
type AppConfig struct {
	Name        string `yaml:"name" env:"APP_NAME" env-default:"gateway" json:"name"`
	Version     string `yaml:"version" env:"APP_VERSION" env-default:"1.0.0" json:"version"`
	Environment string `yaml:"environment" env:"APP_ENV" env-default:"development" json:"environment"`
}

type TCPConfig struct {
	Mode        string `yaml:"mode" env:"TCP_MODE" env-default:"tunnel" json:"mode"`
	NetworkPool string `yaml:"netpool" env:"TCP_NETPOOL" env-default:"10.8.0.0/24" json:"netpool"`
	Gateway     string `yaml:"gateway" env:"TCP_GATEWAY" env-default:"10.8.0.1" json:"gateway"`
	Port        int    `yaml:"port" env:"TCP_PORT" env-default:"8443" json:"port"`
	CFTunnel    string `yaml:"cftoken,omitempty" env:"CF_TUNNEL_TOKEN" json:"cftoken,omitempty"`
}

type DNSConfig struct {
	IP      string   `yaml:"ip" json:"ip"`
	Server  string   `yaml:"server" json:"server"`
	Zones   []string `yaml:"zones" json:"zones"`
	Refresh int      `yaml:"refresh" json:"refresh"`
}

var Conf *Config

func Init() *Config {
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

	var cfg Config

	// Read environment variables or fallback to YAML configuration file
	if _, err := os.Stat(configPath); err == nil {
		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			log.Fatalf("cannot read config file: %s", err.Error())
		}
	} else {
		// Read solely from environment variables & env-default tags if config.yaml is missing
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Fatalf("cannot read environment configuration: %s", err.Error())
		}
	}
	return &cfg
	// SetEnv(Conf)
}
