package main

import (
	"context"
	"log"
	"user_service/internal/app"
	"user_service/internal/config"
	"user_service/internal/user"
	"user_service/internal/user/db"
	"user_service/pkg/client/mongodb"
	"user_service/pkg/logging"
)

func main() {
	log.Print("config initialization")
	cfg := config.GetConfig()

	mongodbClient, err := mongodb.NewClient(context.Background(), cfg.MongoDB.Host, cfg.MongoDB.Port,
		cfg.MongoDB.Username, cfg.MongoDB.Password, cfg.MongoDB.Database, cfg.MongoDB.AuthDB)
	if err != nil {
		panic(err)
	}

	log.Printf("logging initialized.")

	logger := logging.GetLogger(cfg.AppConfig.LogLevel)

	user1 := user.User{
		ID:           "",
		Username:     "aboba",
		PasswordHash: "12asfj123",
	}

	storage := db.NewStorage(mongodbClient, "users", &logger)
	user1ID, err := storage.Create(context.Background(), user1)
	if err != nil {
		panic(err)
	}
	logger.Info(user1ID)

	a, err := app.NewApp(cfg, &logger)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("running Application")
	a.Run()
}
