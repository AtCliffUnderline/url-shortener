package app

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/imdario/mergo"
	"log"
)

type ApplicationConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	StoragePath   string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN   string `env:"DATABASE_DSN" envDefault:"postgres://postgres:root@localhost:5432/golang"`
}

func CreateConfig() ApplicationConfig {
	var config = ApplicationConfig{}
	var flagConfig = ApplicationConfig{}

	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&flagConfig.ServerAddress, "a", "", "Server address to run on")
	flag.StringVar(&flagConfig.BaseURL, "b", "", "Base URL for shortened links")
	flag.StringVar(&flagConfig.StoragePath, "f", "", "File storage path")
	flag.StringVar(&flagConfig.StoragePath, "d", "", "Database DSN")
	flag.StringVar(&flagConfig.StoragePath, "database-dsn", "", "Database DSN")
	flag.Parse()

	if err := mergo.Merge(&config, flagConfig, mergo.WithOverride); err != nil {
		panic(err)
	}

	return config
}
