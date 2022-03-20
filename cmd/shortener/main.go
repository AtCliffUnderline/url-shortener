package main

import (
	"github.com/AtCliffUnderline/url-shortener/internal/app"
)

func main() {
	config := app.CreateConfig()
	db := app.BaseDB{}
	db.SetupConnection(config.DatabaseDSN)
	defer db.CloseConnection()
	app.StartServer(config, db)
}
