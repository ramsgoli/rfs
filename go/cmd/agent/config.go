package main

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	ServerHostname string `env:"SERVER_HOSTNAME" envDefault:"localhost"`
	ServerHttpPort int    `env:"SERVER_HTTP_PORT" envDefault:"8080"`
	ServerDataPort int    `env:"SERVER_DATA_PORT" envDefault:"8000"`
	DataPort       int    `env:"DATA_PORT" envDefault:"6000"`
	IpAddress      string `env:"IP_ADDRESS" envDefault:"localhost"`
	DataDir        string `env:"DATA_DIR" envDefault:"/tmp/data"`
}

func parseConfig() (*Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
