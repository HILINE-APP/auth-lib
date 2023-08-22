package auth

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AUTH_URL string `envconfig:"AUTH_URL" default:"http://localhost:8080"`
}

func Get() Config {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)
	return cfg
}
