package app

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type ApplicationConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	StoragePath   string `env:"FILE_STORAGE_PATH"`
}

func getMergedConfig(cliConfig ApplicationConfig) ApplicationConfig {
	var cfg ApplicationConfig
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	if cliConfig.ServerAddress != "" {
		cfg.ServerAddress = cliConfig.ServerAddress
	}
	if cliConfig.BaseURL != "" {
		cfg.BaseURL = cliConfig.BaseURL
	}
	if cliConfig.StoragePath != "" {
		cfg.StoragePath = cliConfig.StoragePath
	}

	return cfg
}
