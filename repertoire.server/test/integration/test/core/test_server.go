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
	"regexp"
	"repertoire/server/api"
	"repertoire/server/data"
	"repertoire/server/data/cache"
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
var MeiliCache cache.MeiliCache
var httpServer *http.Server

type TestServer struct {
	WithMeili           bool
	WithStorage         bool
	WithCentrifugo      bool
	EnvPath             string
	app                 *fx.App
	dbContainer         *postgresTest.PostgresContainer
	storageServer       *httptest.Server
	centrifugoContainer testcontainers.Container
	meiliContainer      *meilisearchTest.MeilisearchContainer
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
	if ts.WithCentrifugo {
		ts.setupCentrifugoContainer()
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
		fx.Populate(&MeiliCache),
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
		log.Printf("failed to terminate postgres db container: %s", err)
	}
	if ts.WithMeili {
		if err := testcontainers.TerminateContainer(ts.meiliContainer); err != nil {
			log.Printf("failed to terminate meiliearch container: %s", err)
		}
	}
	if ts.WithStorage {
		ts.storageServer.Close()
	}
	if ts.WithCentrifugo {
		if err := testcontainers.TerminateContainer(ts.centrifugoContainer); err != nil {
			log.Printf("failed to terminate centrifugo container: %s", err)
		}
	}
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
				ExpiresIn string
			}{
				"some token",
				"1h",
			}
			bytes, _ := json.Marshal(response)
			_, _ = w.Write(bytes)
		}
	}))
	_ = os.Setenv("AUTH_STORAGE_URL", ts.storageServer.URL)
	_ = os.Setenv("STORAGE_UPLOAD_URL", ts.storageServer.URL)
}

func (ts *TestServer) setupCentrifugoContainer() {
	containerRequest := testcontainers.ContainerRequest{
		Image:        "centrifugo/centrifugo:v6",
		ExposedPorts: []string{"8000/tcp"},
		Cmd:          []string{"centrifugo", "-c", "/centrifugo/config.json"},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      "../../../centrifugo-config.json",
				ContainerFilePath: "/centrifugo/config.json",
				FileMode:          0644,
			},
		},
		WaitingFor: wait.ForLog("serving websocket"),
	}
	ts.centrifugoContainer, _ = testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerRequest,
			Started:          true,
		})
	// Get Random Port and set it to the environment variable
	port, _ := ts.centrifugoContainer.MappedPort(context.Background(), "8000/tcp")
	regex := regexp.MustCompile(`localhost:\d{4}`)
	newUrl := regex.ReplaceAllString(os.Getenv("CENTRIFUGO_URL"), "localhost:"+port.Port())
	_ = os.Setenv("CENTRIFUGO_URL", newUrl)
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
	regex := regexp.MustCompile(`localhost:\d{4}`)
	newUrl := regex.ReplaceAllString(os.Getenv("MEILI_URL"), "localhost:"+port.Port())
	_ = os.Setenv("MEILI_URL", newUrl)

	// Initialize Indexes and Filterable Attributes
	meiliClient := meilisearch.New(env.MeiliUrl, meilisearch.WithAPIKey(env.MeiliMasterKey))

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
