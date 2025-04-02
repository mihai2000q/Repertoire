package core

import (
	"context"
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

var (
	Dsn           string
	MessageBroker message.Publisher
	MeiliCache    cache.MeiliCache
	httpServer    *http.Server
)

var jwtInfo = struct {
	Issuer     string
	Audience   string
	PrivateKey string
	PublicKey  string
}{
	Issuer:   "Repertoire",
	Audience: "Repertoire",
	PrivateKey: `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`,
	PublicKey: `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAhN+jlepyRESen3yU/Nui
3tb5/4j1ckKEiQnGRm3oQtXNBpUJkiP21h/zJIv9n5NfdP6TAsgoCGfhEHPezQ55
d5B8qHRaZQ7nbXlZDq4mdSI4OFC8U5RsepQX3MEvNtk0P4BfNMXbeFDrAfu8WeAm
Y2U5tcqbwKCNPUSJBQvLJvfFHVR87h+9ovJK35d5qYaQZ75EHphJd+A3K93Uasng
Lw9Uf65PCm3f2Snj4yunhi853jPR6HZ8jftO3YraSD01odqyP5OgYaH3LUlPPexy
XiOCnuDtdGGv+IYsQNCpDOhAP0Z/FRk6NRVBC+ORVGqEeMo7/PEDpG3YMuYljI6P
FiJjOh/DkNUnZVSFHtaa9OPeAtzp0DOqvwEJSlfAx74E7XcrXIJOU7S13Rl/Qqr8
A8qGyrQNLyWixS8wm28iQT0nl7qCwVvh1ZuvJ1xUpfs2rfh2HAJ1PSz5hLY7xT5P
Dl399rQURyk5+k9hcmxr9Oe7cmjLbnq+2HdE1zNcGmSUZ8+pDbSrVhNCFOBzxypV
b0P2ZpMt6xrlcB/w3BU8/CIpbLnr+GJVHGFUB+CSdCv82R4979dZIqzOaz0ZZ+CC
DYxvFS9YDlIdyR8/ohisnBjqADjQefBJy/WY7fGWBU2sItZjMEm1ohqj60hjaN2s
fZIWqtSmnVopY/qDGKBzQjUCAwEAAQ==
-----END PUBLIC KEY-----`,
}

// CentrifugoJwtInfo - Extracted info from local config.json
var CentrifugoJwtInfo = struct {
	SecretKey string
	Issuer    string
	Audience  string
}{
	"This-is-a-very-super-duper-secret-key-and-it-shall-stay-like-this",
	"http://localhost:8030/auth",
	"Centrifugo",
}

type TestServer struct {
	WithMeili           bool
	WithStorage         bool
	WithCentrifugo      bool
	EnvPath             string
	app                 *fx.App
	dbContainer         *postgresTest.PostgresContainer
	storageServer       *httptest.Server
	authServer          *httptest.Server
	centrifugoContainer testcontainers.Container
	meiliContainer      *meilisearchTest.MeilisearchContainer
}

func (ts *TestServer) Start() {
	// Setup environment public key for custom tokens
	_ = os.Setenv("JWT_PUBLIC_KEY", jwtInfo.PublicKey)

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
	if ts.WithStorage || ts.WithCentrifugo {
		ts.setupAuthServer()
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
	if ts.WithStorage || ts.WithCentrifugo {
		ts.authServer.Close()
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
	}))
	_ = os.Setenv("STORAGE_UPLOAD_URL", ts.storageServer.URL)
}

func (ts *TestServer) setupAuthServer() {
	testServer := testAuthServer{}
	ts.authServer = httptest.NewServer(testServer.handle())
	_ = os.Setenv("AUTH_URL", ts.authServer.URL)
}

func (ts *TestServer) setupCentrifugoContainer() {
	containerRequest := testcontainers.ContainerRequest{
		Image:        "centrifugo/centrifugo:v6",
		ExposedPorts: []string{"8003/tcp"},
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
	port, _ := ts.centrifugoContainer.MappedPort(context.Background(), "8003/tcp")
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
	meiliClient := meilisearch.New(newUrl, meilisearch.WithAPIKey(env.MeiliMasterKey))

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
		"title", "name", "updatedAt", "createdAt", "album", "album.title", "artist", "artist.name",
	})
	if err != nil {
		log.Println(err)
	}

	meiliClient.Close()
}
