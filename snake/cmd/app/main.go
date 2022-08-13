package main

import (
	"log"
	//_ "snake_service/docs"
	"snake_service/internal/app"
	"snake_service/internal/config"
	"snake_service/pkg/logging"
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
