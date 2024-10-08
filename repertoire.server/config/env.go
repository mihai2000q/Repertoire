package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Env struct {
	ApplicationPort  string
	Environment      string
	DatabaseHost     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabasePort     string
	DatabaseSSLMode  string
	LogOutput        string
	LogLevel         string
}

func NewEnv() Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file:%v", err)
	}

	env := Env{
		ApplicationPort:  os.Getenv("SERVER_PORT"),
		Environment:      os.Getenv("ENVIRONMENT"),
		DatabaseHost:     os.Getenv("DB_HOST"),
		DatabaseUser:     os.Getenv("DB_USER"),
		DatabasePassword: os.Getenv("DB_PASSWORD"),
		DatabaseName:     os.Getenv("DB_NAME"),
		DatabasePort:     os.Getenv("DB_PORT"),
		DatabaseSSLMode:  os.Getenv("DB_SSL_MODE"),
		LogOutput:        os.Getenv("LOG_OUTPUT"),
		LogLevel:         os.Getenv("LOG_LEVEL"),
	}
	return env
}
