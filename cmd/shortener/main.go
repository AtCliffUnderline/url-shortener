package main

import (
	"flag"
	"github.com/AtCliffUnderline/url-shortener/internal/app"
)

var config = app.ApplicationConfig{}

func init() {
	flag.StringVar(&config.ServerAddress, "a", "", "Server address to run on")
	flag.StringVar(&config.BaseURL, "b", "", "Base URL for shortened links")
	flag.StringVar(&config.BaseURL, "f", "", "File storage path")
}

func main() {
	flag.Parse()
	app.StartServer(config)
}
