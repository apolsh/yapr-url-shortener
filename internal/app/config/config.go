package config

import (
	"flag"
	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	AuthSecretKey   string `env:"AUTH_SECRET_KEY" envDefault:"very_secret_key"`
}

func Load() Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	var addressFlagValue string
	var baseURLFlagValue string
	var fileStoragePathFlagValue string

	flag.StringVar(&addressFlagValue, "a", "", "HTTP server start address")
	flag.StringVar(&baseURLFlagValue, "b", "", "base address of the resulting shortened URL")
	flag.StringVar(&fileStoragePathFlagValue, "f", "", "path to file with abbreviated URLs")

	flag.Parse()

	if addressFlagValue != "" {
		cfg.ServerAddress = addressFlagValue
	}
	if baseURLFlagValue != "" {
		cfg.BaseURL = baseURLFlagValue
	}
	if fileStoragePathFlagValue != "" {
		cfg.FileStoragePath = fileStoragePathFlagValue
	}

	return cfg
}
