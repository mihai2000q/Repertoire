package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/meilisearch/meilisearch-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"repertoire/server/api"
	"repertoire/server/data"
	"repertoire/server/data/message"
	"repertoire/server/domain"
	"repertoire/server/internal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pressly/goose"
	"github.com/testcontainers/testcontainers-go"
	meilisearchTest "github.com/testcontainers/testcontainers-go/modules/meilisearch"
	postgresTest "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/fx"
)

var Dsn string
var MessageBroker message.Publisher
var httpServer *http.Server

type TestServer struct {
	WithMeili      bool
	WithStorage    bool
	EnvPath        string
	app            *fx.App
	dbContainer    *postgresTest.PostgresContainer
	storageServer  *httptest.Server
	meiliContainer *meilisearchTest.MeilisearchContainer
}

func (ts *TestServer) Start() {
	// Setup Environment Variable to anything, so that it checks the right path of .env
	relativePath := "../../../"
	if ts.EnvPath != "" {
		relativePath = ts.EnvPath
	}
	_ = os.Setenv("INTEGRATION_TESTING_ENVIRONMENT_FILE_PATH", relativePath+".env")

	env := internal.NewEnv()

	ts.setupPostgresContainer(env, relativePath)
	if ts.WithMeili {
		ts.setupMeiliContainer(env)
	}
	if ts.WithStorage {
		ts.setupStorageServer()
	}

	// Setup application modules and populate the router
	// Implicitly, the application will connect to the database
	gin.SetMode(gin.TestMode)
	ts.app = fx.New(
		internal.Module,
		data.Module,
		domain.Module,
		api.Module,
		fx.Populate(&httpServer),
		fx.Populate(&MessageBroker),
	)

	// Start application
	if err := ts.app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func (ts *TestServer) Stop() {
	if err := ts.app.Stop(context.Background()); err != nil {
		log.Fatal(err)
	}
	if err := testcontainers.TerminateContainer(ts.dbContainer); err != nil {
		log.Printf("failed to terminate postgres dbContainer: %s", err)
	}
	if ts.WithMeili {
		if err := testcontainers.TerminateContainer(ts.meiliContainer); err != nil {
			log.Printf("failed to terminate meiliearch dbContainer: %s", err)
		}
	}
	if ts.WithStorage {
		ts.storageServer.Close()
	}
}

func getHttpServer() *http.Server {
	return httpServer
}

func (ts *TestServer) setupPostgresContainer(env internal.Env, relativePath string) {
	// Setup Container
	container, err := postgresTest.Run(context.Background(),
		"postgres:17",
		postgresTest.WithDatabase(env.DatabaseName),
		postgresTest.WithUsername(env.DatabaseUser),
		postgresTest.WithPassword(env.DatabasePassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	ts.dbContainer = container
	if err != nil {
		log.Fatal(err)
	}

	// Get Random Port and set the environment variable
	port, _ := ts.dbContainer.MappedPort(context.Background(), "5432/tcp")
	_ = os.Setenv("DB_PORT", port.Port())
	Dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		env.DatabaseHost,
		env.DatabaseUser,
		env.DatabasePassword,
		env.DatabaseName,
		port.Port(),
		env.DatabaseSSLMode,
	)

	// apply migrations to database
	postgresDB, _ := gorm.Open(postgres.Open(Dsn))
	db, _ := postgresDB.DB()
	if err = goose.Up(db, relativePath+"migrations/database/"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	_ = db.Close()
}

func (ts *TestServer) setupStorageServer() {
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
	_ = os.Setenv("AUTH_STORAGE_URL", ts.storageServer.URL)
	_ = os.Setenv("UPLOAD_STORAGE_URL", ts.storageServer.URL)
}

func (ts *TestServer) setupMeiliContainer(env internal.Env) {
	// Setup Container
	container, err := meilisearchTest.Run(
		context.Background(),
		"getmeili/meilisearch:v1.13.0",
		meilisearchTest.WithMasterKey(env.MeiliMasterKey),
	)
	if err != nil {
		log.Fatal(err)
	}
	ts.meiliContainer = container

	// Get Random Port and set the environment variable
	port, _ := ts.meiliContainer.MappedPort(context.Background(), "7700/tcp")
	_ = os.Setenv("MEILI_PORT", port.Port())

	// Initialize Indexes and Filterable Attributes
	url := "http://" + env.MeiliHost + ":" + port.Port()
	meiliClient := meilisearch.New(url, meilisearch.WithAPIKey(env.MeiliMasterKey))

	_, err = meiliClient.CreateIndex(&meilisearch.IndexConfig{
		Uid:        "search",
		PrimaryKey: "id",
	})
	if err != nil {
		log.Println(err)
	}

	_, err = meiliClient.Index("search").UpdateFilterableAttributes(&[]string{
		"type", "userId", "album", "album.id", "artist", "artist.id",
	})
	if err != nil {
		log.Println(err)
	}

	_, err = meiliClient.Index("search").UpdateSortableAttributes(&[]string{
		"title", "name", "updatedAt", "album", "album.title", "artist", "artist.name",
	})
	if err != nil {
		log.Println(err)
	}

	meiliClient.Close()
}
