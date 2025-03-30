package internal

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	ApplicationHost string
	ApplicationPort string
	Environment     string

	JwtSecretKey string
	JwtIssuer    string
	JwtAudience  string

	UploadDirectory string
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

		JwtSecretKey: os.Getenv("JWT_SECRET_KEY"),
		JwtIssuer:    os.Getenv("JWT_ISSUER"),
		JwtAudience:  os.Getenv("JWT_AUDIENCE"),

		UploadDirectory: os.Getenv("UPLOAD_DIRECTORY"),
	}
	return env
}

var DevelopmentEnvironment = "development"
