package service

import (
	"encoding/json"
	"repertoire/server/data/message"
	"repertoire/server/internal/message/topics"

	"github.com/ThreeDotsLabs/watermill"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
)

type MessagePublisherService interface {
	GetClient() message.Publisher
	Publish(topic topics.Topic, messagePayload any) error
}

type messagePublisherService struct {
	client message.Publisher
}

func NewMessagePublisherService(client message.Publisher) MessagePublisherService {
	return messagePublisherService{client: client}
}

func (m messagePublisherService) GetClient() message.Publisher {
	return m.client
}

func (m messagePublisherService) Publish(topic topics.Topic, data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	msg := watermillMessage.NewMessage(watermill.NewUUID(), bytes)
	msg.Metadata.Set("topic", string(topic))
	queue := string(topics.TopicToQueueMap[topic])
	return m.client.Publish(queue, msg)
}
