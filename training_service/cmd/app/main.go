package main

import (
	"log"
	_ "training_service/docs"
	"training_service/internal/app"
	"training_service/internal/config"
	"training_service/pkg/logging"
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
