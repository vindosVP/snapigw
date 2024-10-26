package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

type Config struct {
	Port        int      `env:"HTTP_PORT" json:"port"`
	ENV         string   `env:"LOG_ENV" envDefault:"dev" json:"env"`
	Services    Services `json:"services"`
	TokenSecret string   `env:"TOKEN_SECRET" json:"-"`
	ServiceName string   `env:"SERVICE_NAME" envDefault:"apigw-ext" json:"serviceName"`
}

type Services struct {
	AuthAddr string `env:"AUTH_ADDR" json:"authAddr"`
}

func MustParse() *Config {
	cfg := &Config{}
	err := env.Parse(cfg, env.Options{RequiredIfNoDef: true})
	if err != nil {
		panic(errors.Wrap(err, "filed to parse config"))
	}
	return cfg
}
