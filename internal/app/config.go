package app

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type ApplicationConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func getConfig() ApplicationConfig {
	var cfg ApplicationConfig
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localhost:8080"
	}
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = ":8080"
	}

	return cfg
}
