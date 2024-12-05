package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"repertoire/server/api"
	"repertoire/server/data"
	"repertoire/server/domain"
	"repertoire/server/internal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	postgresTest "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/fx"
)

var Dsn string
var httpServer *http.Server

type TestServer struct {
	app           *fx.App
	container     *postgresTest.PostgresContainer
	storageServer *httptest.Server
}

func Start(envPath ...string) *TestServer {
	ts := &TestServer{}

	// Setup Environment Variable to anything, so that it checks the right path of .env
	if len(envPath) > 0 {
		_ = os.Setenv("INTEGRATION_TESTING_ENVIRONMENT_FILE_PATH", envPath[0])
	} else {
		_ = os.Setenv("INTEGRATION_TESTING_ENVIRONMENT_FILE_PATH", "../../../.env")
	}

	env := internal.NewEnv()

	// Setup Postgres Docker Container
	postgresContainer, err := postgresTest.Run(context.Background(),
		"postgres:17",
		postgresTest.WithDatabase(env.DatabaseName),
		postgresTest.WithUsername(env.DatabaseUser),
		postgresTest.WithPassword(env.DatabasePassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	ts.container = postgresContainer
	if err != nil {
		log.Fatal(err)
	}

	// Get Random Port and set the environment variable
	port, _ := ts.container.MappedPort(context.Background(), "5432/tcp")
	_ = os.Setenv("DB_PORT", port.Port())
	Dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		env.DatabaseHost,
		env.DatabaseUser,
		env.DatabasePassword,
		env.DatabaseName,
		port.Port(),
		env.DatabaseSSLMode,
	)

	// Start Storage Server
	ts.storageServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		// when asking for another auth token, just return something
		if r.Method == http.MethodPost {
			response := struct {
				Token     string
				TokenType string
				ExpiresIn string
			}{
				"some token",
				"Bearer",
				"1h",
			}
			bytes, _ := json.Marshal(response)
			_, _ = w.Write(bytes)
		}
	}))
	_ = os.Setenv("UPLOAD_STORAGE_URL", ts.storageServer.URL)

	// Setup application modules and populate the router
	// Implicitly, the application will connect to the database
	gin.SetMode(gin.TestMode)
	ts.app = fx.New(
		internal.Module,
		data.Module,
		domain.Module,
		api.Module,
		fx.Populate(&httpServer),
	)

	// Start application
	if err = ts.app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}

	return ts
}

func Stop(ts *TestServer) {
	ts.storageServer.Close()
	if err := ts.app.Stop(context.Background()); err != nil {
		log.Fatal(err)
	}
	if err := testcontainers.TerminateContainer(ts.container); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
}

func getHttpServer() *http.Server {
	return httpServer
}
