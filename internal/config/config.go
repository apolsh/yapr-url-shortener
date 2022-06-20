package config

import (
	"github.com/caarlos0/env"
	"log"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS,required"`
	BaseURL       string `env:"BASE_URL,required"`
}

func Load() Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Println("Failed to load environment variables, will use default.")
		cfg = Config{ServerAddress: "http://localhost:8080", BaseURL: "localhost:8080"}
	}
	return cfg
}
