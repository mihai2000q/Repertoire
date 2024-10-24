package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
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
	JwtIssuer        string
	JwtAudience      string
	JwtSecretKey     string
}

func NewEnv() Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file:%v", err)
	}

	env := Env{
		ApplicationPort:  os.Getenv("SERVER_PORT"),
		Environment:      os.Getenv("ENV"),
		DatabaseHost:     os.Getenv("DB_HOST"),
		DatabaseUser:     os.Getenv("DB_USER"),
		DatabasePassword: os.Getenv("DB_PASSWORD"),
		DatabaseName:     os.Getenv("DB_NAME"),
		DatabasePort:     os.Getenv("DB_PORT"),
		DatabaseSSLMode:  os.Getenv("DB_SSL_MODE"),
		LogOutput:        os.Getenv("LOG_OUTPUT"),
		LogLevel:         os.Getenv("LOG_LEVEL"),
		JwtIssuer:        os.Getenv("JWT_ISSUER"),
		JwtAudience:      os.Getenv("JWT_AUDIENCE"),
		JwtSecretKey:     os.Getenv("JWT_SECRET_KEY"),
	}
	return env
}

var DevelopmentEnvironment = "development"
