package main

import (
	"log"
	_ "qualifications_service/docs"
	"qualifications_service/internal/app"
	"qualifications_service/internal/config"
	"qualifications_service/pkg/logging"
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
