package internal

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Env struct {
	ApplicationHost string
	ApplicationPort string
	Environment     string

	DatabaseHost     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabasePort     string
	DatabaseSSLMode  string

	LogOutput string
	LogLevel  string

	JwtIssuer    string
	JwtAudience  string
	JwtSecretKey string

	StorageUrl          string
	StorageClientID     string
	StorageClientSecret string
}

func NewEnv() Env {
	if path := os.Getenv("INTEGRATION_TESTING_ENVIRONMENT_FILE_PATH"); path != "" {
		_ = godotenv.Load(path)
	} else if os.Getenv("IS_RUNNING_IN_CONTAINER") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file:%v", err)
		}
	}

	env := Env{
		ApplicationHost: os.Getenv("SERVER_HOST"),
		ApplicationPort: os.Getenv("SERVER_PORT"),
		Environment:     os.Getenv("ENV"),

		DatabaseHost:     os.Getenv("DB_HOST"),
		DatabaseUser:     os.Getenv("DB_USER"),
		DatabasePassword: os.Getenv("DB_PASSWORD"),
		DatabaseName:     os.Getenv("DB_NAME"),
		DatabasePort:     os.Getenv("DB_PORT"),
		DatabaseSSLMode:  os.Getenv("DB_SSL_MODE"),

		LogOutput: os.Getenv("LOG_OUTPUT"),
		LogLevel:  os.Getenv("LOG_LEVEL"),

		JwtIssuer:    os.Getenv("JWT_ISSUER"),
		JwtAudience:  os.Getenv("JWT_AUDIENCE"),
		JwtSecretKey: os.Getenv("JWT_SECRET_KEY"),

		StorageUrl:          os.Getenv("UPLOAD_STORAGE_URL"),
		StorageClientID:     os.Getenv("STORAGE_CLIENT_ID"),
		StorageClientSecret: os.Getenv("STORAGE_CLIENT_SECRET"),
	}
	return env
}

var DevelopmentEnvironment = "development"
