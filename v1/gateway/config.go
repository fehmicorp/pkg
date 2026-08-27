package gateway

type Config struct {
	Name        string `yaml:"name" env:"APP_NAME" env-default:"gateway"`
	Version     string `yaml:"version" env:"APP_VERSION" env-default:"1.0.0"`
	Environment string `yaml:"environment" env:"APP_ENV" env-default:"development"`
}
