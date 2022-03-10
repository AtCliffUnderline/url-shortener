package main

import (
	"github.com/AtCliffUnderline/url-shortener/internal/app"
)

func main() {
	config := app.CreateConfig()
	app.StartServer(config)
}
