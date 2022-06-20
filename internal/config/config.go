package config

import (
	"github.com/caarlos0/env"
	"log"
)

const (
	defaultServerAddress = "http://localhost:8080"
	defaultBaseUrl       = "localhost:8080"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func Load() Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Println("Failed to load environment variables, will use default.")
		cfg = Config{ServerAddress: "http://localhost:8080", BaseURL: "localhost:8080"}
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseUrl
	}
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = defaultServerAddress
	}

	return cfg
}
