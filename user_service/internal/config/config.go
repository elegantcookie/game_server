package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	IsDebug       bool `env:"IS_DEBUG" env-default:"false"`
	IsDevelopment bool `env:"IS_DEV" env-default:"false"`
	Listen        struct {
		SocketFile string `env:"SOCKET_FILE" env-default:"app.sock"`
		Type       string `env:"LISTEN_TYPE" env-default:"port"`
		BindIP     string `env:"BIND_IP" env-default:"0.0.0.0"`
		Port       string `env:"PORT" env-default:"10002"`
	}
	AppConfig struct {
		LogLevel string `env:"LOG_LEVEL" env-default:"trace"`
	}
	MongoDB struct {
		Host     string `env:"HOST" env-default:"localhost"`
		Port     string `env:"PORT" env-default:"27017"`
		Username string `env:"ADMIN_USERNAME"`
		Password string `env:"ADMIN_PASSWORD"`
		Database string `env:"DATABASE" env-default:"user-service"`
		AuthDB   string `env:"AUTH_DB"`
	}
	Keys struct {
		JWTSignKey string `env:"JWT_SIGN_KEY" env-default:"alsfjak12h4i1h2uas7f7241231o1u2io5u12asopua0w9812"`
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Printf("gather config")

		instance = &Config{}

		if err := cleanenv.ReadEnv(instance); err != nil {
			helpText := "An error occurred during reading config"
			help, _ := cleanenv.GetDescription(instance, &helpText)
			log.Println(help)
			log.Fatal(err)
		}
	})
	return instance
}
