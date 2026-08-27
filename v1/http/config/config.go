package httpconfig

type HttpServer struct {
	Mode            string `yaml:"host" env:"SERVER_MODE" env-default:"http"`
	Host            string `yaml:"host" env:"SERVER_HOST" env-default:"localhost"`
	Port            int    `yaml:"port" env:"SERVER_PORT" env-default:"8050"`
	FQDN            string `yaml:"host" env:"SERVER_FQDN" env-default:"https://api.fehmicorp.in/"`
	ReadTimeoutSec  int    `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT" env-default:"10"`
	WriteTimeoutSec int    `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT" env-default:"10"`
}
