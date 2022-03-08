package app

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type ApplicationConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	StoragePath   string `env:"FILE_STORAGE_PATH" envDefault:"./asd.txt"`
}

func getConfig() ApplicationConfig {
	var cfg ApplicationConfig
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
