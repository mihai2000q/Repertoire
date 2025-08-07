package user

import (
	"encoding/json"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal/message/topics"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

type UserDeletedHandler struct {
	name                    string
	topic                   topics.Topic
	searchEngineService     service.SearchEngineService
	storageFilePathProvider provider.StorageFilePathProvider
	messagePublisherService service.MessagePublisherService
}

func NewUserDeletedHandler(
	searchEngineService service.SearchEngineService,
	storageFilePathProvider provider.StorageFilePathProvider,
	messagePublisherService service.MessagePublisherService,
) UserDeletedHandler {
	return UserDeletedHandler{
		name:                    "user_deleted_handler",
		topic:                   topics.UserDeletedTopic,
		searchEngineService:     searchEngineService,
		storageFilePathProvider: storageFilePathProvider,
		messagePublisherService: messagePublisherService,
	}
}

func (u UserDeletedHandler) Handle(msg *message.Message) error {
	var userID uuid.UUID
	err := json.Unmarshal(msg.Payload, &userID)
	if err != nil {
		return err
	}

	err = u.syncWithSearchEngine(userID)
	if err != nil {
		return err
	}

	directoryPath := u.storageFilePathProvider.GetUserDirectoryPath(userID)
	return u.messagePublisherService.Publish(topics.DeleteDirectoriesStorageTopic, []string{directoryPath})
}

func (u UserDeletedHandler) syncWithSearchEngine(userID uuid.UUID) error {
	documents, err := u.searchEngineService.GetDocuments("userId = " + userID.String())
	if err != nil {
		return err
	}
	if len(documents) == 0 {
		return nil
	}

	var idsToDelete []string
	for _, document := range documents {
		idsToDelete = append(idsToDelete, document["id"].(string))
	}

	return u.messagePublisherService.Publish(topics.DeleteFromSearchEngineTopic, idsToDelete)
}

func (u UserDeletedHandler) GetName() string {
	return u.name
}

func (u UserDeletedHandler) GetTopic() topics.Topic {
	return u.topic
}
