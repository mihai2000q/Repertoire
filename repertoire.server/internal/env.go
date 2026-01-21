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

	DatabaseHost     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabasePort     string
	DatabaseSSLMode  string

	LogOutput string
	LogLevel  string

	JwtPublicKey string

	AuthUrl          string
	AuthClientID     string
	AuthClientSecret string

	StorageUrl string

	MeiliUrl       string
	MeiliMasterKey string
	MeiliAuthKey   string

	CentrifugoUrl string
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

		JwtPublicKey: os.Getenv("JWT_PUBLIC_KEY"),

		AuthUrl:          os.Getenv("AUTH_URL"),
		AuthClientID:     os.Getenv("AUTH_CLIENT_ID"),
		AuthClientSecret: os.Getenv("AUTH_CLIENT_SECRET"),

		StorageUrl: os.Getenv("STORAGE_UPLOAD_URL"),

		MeiliUrl:       os.Getenv("MEILI_URL"),
		MeiliMasterKey: os.Getenv("MEILI_MASTER_KEY"),
		MeiliAuthKey:   os.Getenv("MEILI_WEBHOOK_AUTHORIZATION_KEY"),

		CentrifugoUrl: os.Getenv("CENTRIFUGO_URL"),
	}
	env.JwtPublicKey = strings.Replace(env.JwtPublicKey, "\\n", "\n", -1)
	return env
}

var DevelopmentEnvironment = "development"
var DebugLogLevel = "DEBUG"
