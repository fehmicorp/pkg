package redis

type Config struct {
	Host     string `yaml:"host" env:"REDIS_HOST" env-default:"redis"`
	Port     int    `yaml:"port" env:"REDIS_PORT" env-default:"6379"`
	User     string `yaml:"user" env:"REDIS_USER" env-default:""`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db" env:"REDIS_DB" env-default:"0"`
}
