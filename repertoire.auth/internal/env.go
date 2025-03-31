package internal

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Env struct {
	ApplicationHost string
	ApplicationPort string
	Environment     string
	LogLevel        string

	DatabaseHost     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabasePort     string
	DatabaseSSLMode  string

	ClientID     string
	ClientSecret string

	JwtPrivateKey     string
	JwtPublicKey      string
	JwtIssuer         string
	JwtAudience       string
	JwtExpirationTime string

	StorageJwtSecretKey      string
	StorageJwtAudience       string
	StorageJwtExpirationTime string

	CentrifugoJwtSecretKey      string
	CentrifugoJwtAudience       string
	CentrifugoJwtExpirationTime string
}

func NewEnv() Env {
	if os.Getenv("IS_RUNNING_IN_CONTAINER") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file:%v", err)
		}
	}

	env := Env{
		ApplicationHost: os.Getenv("SERVER_HOST"),
		ApplicationPort: os.Getenv("SERVER_PORT"),
		Environment:     os.Getenv("ENV"),
		LogLevel:        os.Getenv("LOG_LEVEL"),

		DatabaseHost:     os.Getenv("DB_HOST"),
		DatabaseUser:     os.Getenv("DB_USER"),
		DatabasePassword: os.Getenv("DB_PASSWORD"),
		DatabaseName:     os.Getenv("DB_NAME"),
		DatabasePort:     os.Getenv("DB_PORT"),
		DatabaseSSLMode:  os.Getenv("DB_SSL_MODE"),

		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),

		JwtPrivateKey:     os.Getenv("JWT_PRIVATE_KEY"),
		JwtPublicKey:      os.Getenv("JWT_PUBLIC_KEY"),
		JwtIssuer:         os.Getenv("JWT_ISSUER"),
		JwtExpirationTime: os.Getenv("JWT_EXPIRATION_TIME"),

		StorageJwtSecretKey:      os.Getenv("STORAGE_JWT_SECRET_KEY"),
		StorageJwtAudience:       os.Getenv("STORAGE_JWT_AUDIENCE"),
		StorageJwtExpirationTime: os.Getenv("STORAGE_JWT_EXPIRATION_TIME"),

		CentrifugoJwtSecretKey:      os.Getenv("CENTRIFUGO_JWT_SECRET_KEY"),
		CentrifugoJwtAudience:       os.Getenv("CENTRIFUGO_JWT_AUDIENCE"),
		CentrifugoJwtExpirationTime: os.Getenv("CENTRIFUGO_JWT_EXPIRATION_TIME"),
	}
	env.JwtPrivateKey = strings.Replace(env.JwtPrivateKey, "\\n", "\n", -1)
	env.JwtPublicKey = strings.Replace(env.JwtPublicKey, "\\n", "\n", -1)
	return env
}

var DevelopmentEnvironment = "development"
