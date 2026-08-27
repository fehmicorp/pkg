package gateway

type Config struct {
	Name        string `yaml:"name" env:"APP_NAME" env-default:"gateway"`
	Version     string `yaml:"version" env:"APP_VERSION" env-default:"1.0.0"`
	Environment string `yaml:"environment" env:"APP_ENV" env-default:"development"`
}

type LoggingConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	Dir   string `yaml:"dir" env:"LOG_DIR" env-default:"./logs"`
}

type MonitoringConfig struct {
	PrometheusEnabled bool `yaml:"prometheus_enabled" env:"PROM_ENABLED" env-default:"true"`
}
