package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/apolsh/yapr-url-shortener/internal/logger"
	"github.com/caarlos0/env"
)

var log = logger.LoggerOfComponent("config")

// Config конфигурационные данные приложения
type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080" json:"server_address"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080" json:"base_url"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	AuthSecretKey   string `env:"AUTH_SECRET_KEY" envDefault:"very_secret_key"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	HTTPSEnabled    bool   `env:"ENABLE_HTTPS" json:"enable_https"`
	ConfigFilePath  string `env:"CONFIG"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"info"`
}

func (c *Config) populateEmptyFields(another Config) {
	if c.ServerAddress == "" && another.ServerAddress != "" {
		c.ServerAddress = another.ServerAddress
	}
	if c.BaseURL == "" && another.BaseURL != "" {
		c.BaseURL = another.BaseURL
	}
	if c.FileStoragePath == "" && another.FileStoragePath != "" {
		c.FileStoragePath = another.FileStoragePath
	}
	if c.AuthSecretKey == "" && another.AuthSecretKey != "" {
		c.AuthSecretKey = another.AuthSecretKey
	}
	if c.DatabaseDSN == "" && another.DatabaseDSN != "" {
		c.DatabaseDSN = another.DatabaseDSN
	}
	if c.ConfigFilePath == "" && another.ConfigFilePath != "" {
		c.ConfigFilePath = another.ConfigFilePath
	}
	if !c.HTTPSEnabled && another.HTTPSEnabled {
		c.HTTPSEnabled = another.HTTPSEnabled
	}
}

// Load считывает переменные окружения и флаги, приоритет отдается флагам
func Load() Config {

	var mainConfig Config

	flag.StringVar(&mainConfig.ServerAddress, "a", "", "HTTP server start address")
	flag.StringVar(&mainConfig.BaseURL, "b", "", "base address of the resulting shortened URL")
	flag.StringVar(&mainConfig.FileStoragePath, "f", "", "path to file with abbreviated URLs")
	flag.StringVar(&mainConfig.DatabaseDSN, "d", "", "database DSN")
	flag.BoolVar(&mainConfig.HTTPSEnabled, "s", false, "enable HTTPS with self signed certificate")
	flag.StringVar(&mainConfig.ConfigFilePath, "c", "", "config file path")
	if mainConfig.ConfigFilePath == "" {
		flag.StringVar(&mainConfig.ConfigFilePath, "config", "", "config file path")
	}

	flag.Parse()

	var envsConfig Config
	if err := env.Parse(&envsConfig); err != nil {
		panic(err)
	}

	mainConfig.populateEmptyFields(envsConfig)

	if envsConfig.ConfigFilePath != "" {
		f, err := os.Open(envsConfig.ConfigFilePath)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to read file: %s, cause: %w", envsConfig.ConfigFilePath, err))
		}

		var configFile Config
		err = json.NewDecoder(f).Decode(&configFile)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to parse file: %s, cause: %w", envsConfig.ConfigFilePath, err))
		}
		mainConfig.populateEmptyFields(configFile)
	}

	return mainConfig
}
