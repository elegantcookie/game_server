package main

import (
	"log"
	_ "prize_service/docs"
	"prize_service/internal/app"
	"prize_service/internal/config"
	"prize_service/pkg/logging"
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
