package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func Load() Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	return cfg
}
