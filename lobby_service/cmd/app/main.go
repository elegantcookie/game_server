package main

import (
	//_ "lobby_service/docs"
	"lobby_service/internal/app"
	"lobby_service/internal/config"
	"lobby_service/pkg/logging"
	"log"
)

func main() {
	log.Print("config initialization")
	cfg := config.GetConfig()

	log.Printf("logging initialized.")

	logger := logging.GetLogger(cfg.AppConfig.LogLevel)

	a, err := app.NewApp(cfg, &logger)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("running Application")
	a.Run()
}
