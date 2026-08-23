package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	DBURL string `env:"DB_URL"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)

	return cfg, err
}
