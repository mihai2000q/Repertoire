package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"repertoire/server/internal"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/centrifugal/centrifuge-go"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/meilisearch/meilisearch-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Clients

func GetDatabase(t *testing.T) *gorm.DB {
	db, _ := gorm.Open(postgres.Open(core.Dsn))
	t.Cleanup(func() {
		d, _ := db.DB()
		_ = d.Close()
	})
	return db
}

func GetSearchClient(t *testing.T) meilisearch.ServiceManager {
	env := GetEnv()
	client := meilisearch.New(env.MeiliUrl, meilisearch.WithAPIKey(env.MeiliMasterKey))
	t.Cleanup(func() {
		client.Close()
	})
	return client
}

func GetCentrifugoClient(t *testing.T) *centrifuge.Client {
	env := GetEnv()
	client := centrifuge.NewJsonClient(env.CentrifugoUrl, centrifuge.Config{})
	client.SetToken(createCentrifugoToken())
	t.Cleanup(func() {
		client.Close()
	})
	return client
}

func GetEnv() internal.Env {
	return internal.NewEnv()
}

// Meilisearch

func WaitForSearchTasksToStart(client meilisearch.ServiceManager, totalTasks int64) {
	for {
		tasks, _ := client.GetTasks(nil)
		if tasks.Total != totalTasks {
			break
		}
	}
}

func WaitForAllSearchTasks(client meilisearch.ServiceManager) {
	for {
		breakOuterFor := true
		tasks, _ := client.GetTasks(nil)
		for _, taskResult := range tasks.Results {
			if taskResult.Status == meilisearch.TaskStatusEnqueued ||
				taskResult.Status == meilisearch.TaskStatusProcessing {
				breakOuterFor = false
				break
			}
		}
		if breakOuterFor {
			break
		}
	}
	fmt.Print("letsgo")
}

func createCentrifugoToken() string {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": uuid.New().String(),
		"sub": "Integration Testing",
		"iss": core.CentrifugoJwtInfo.Issuer,
		"aud": core.CentrifugoJwtInfo.Audience,
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(time.Hour).Unix(),
	})
	token, _ := claims.SignedString([]byte(core.CentrifugoJwtInfo.SecretKey))

	return token
}

// Message Handling

func PublishToTopic(topic topics.Topic, data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	msg := message.NewMessage(watermill.NewUUID(), bytes)
	msg.Metadata.Set("topic", string(topic))
	queue := string(topics.TopicToQueueMap[topic])
	return core.MessageBroker.Publish(queue, msg)
}

type SubscribedToTopic struct {
	Messages <-chan *message.Message
	Topic    topics.Topic
}

func SubscribeToTopic(topic topics.Topic) SubscribedToTopic {
	messages, _ := core.MessageBroker.Subscribe(context.Background(), string(topics.TopicToQueueMap[topic]))
	return SubscribedToTopic{
		Messages: messages,
		Topic:    topic,
	}
}

// Seeding

func SeedAndCleanupData(t *testing.T, users []model.User, seed func(*gorm.DB)) {
	db := GetDatabase(t)
	seed(db)
	t.Cleanup(func() {
		for _, user := range users {
			db.Select(clause.Associations).Delete(user)
		}
	})
}

func SeedAndCleanupSearchData(t *testing.T, items []any) {
	searchClient := GetSearchClient(t)

	_, _ = searchClient.Index("search").AddDocuments(items, nil)
	WaitForAllSearchTasks(searchClient)

	t.Cleanup(func() {
		_, _ = searchClient.Index("search").DeleteAllDocuments(nil)
	})
}

// Misc Utils

func AttachFileToMultipartBody(fileName string, formName string, multiWriter *multipart.Writer) {
	tempFile, _ := os.CreateTemp("", fileName)
	defer func(name string) {
		_ = os.Remove(name)
	}(tempFile.Name())

	fileWriter, _ := multiWriter.CreateFormFile(formName, tempFile.Name())

	file, _ := os.Open(tempFile.Name())
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	_, _ = file.WriteTo(fileWriter)
}

func UnmarshalDocument[T any](document any) T {
	bytes, _ := json.Marshal(document)
	var marshalledDocument T
	_ = json.Unmarshal(bytes, &marshalledDocument)
	return marshalledDocument
}
