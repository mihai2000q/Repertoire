package user

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
)

type UserDeletedHandler struct {
	name                    string
	topic                   topics.Topic
	searchEngineService     service.SearchEngineService
	messagePublisherService service.MessagePublisherService
}

func NewUserDeletedHandler(
	searchEngineService service.SearchEngineService,
	messagePublisherService service.MessagePublisherService,
) UserDeletedHandler {
	return UserDeletedHandler{
		name:                    "user_deleted_handler",
		topic:                   topics.UserDeletedTopic,
		searchEngineService:     searchEngineService,
		messagePublisherService: messagePublisherService,
	}
}

func (u UserDeletedHandler) Handle(msg *message.Message) error {
	var userId uuid.UUID
	err := json.Unmarshal(msg.Payload, &userId)
	if err != nil {
		return err
	}

	documents, err := u.searchEngineService.GetDocuments("userId = " + userId.String())
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

	err = u.messagePublisherService.Publish(topics.DeleteFromSearchEngineTopic, idsToDelete)
	if err != nil {
		return err
	}

	return nil
}

func (u UserDeletedHandler) GetName() string {
	return u.name
}

func (u UserDeletedHandler) GetTopic() topics.Topic {
	return u.topic
}
