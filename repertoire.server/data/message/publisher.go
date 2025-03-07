package message

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"repertoire/server/internal/message/topics"
)

type Publisher interface {
	PublishOnTopic(topic topics.Topic, messagePayload any) error
	Publish(topic string, messages ...*message.Message) error
	Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
	Close() error
}

type publisher struct {
	client *gochannel.GoChannel
}

func NewPublisher() Publisher {
	logger := watermill.NewStdLogger(true, true)
	return publisher{
		client: gochannel.NewGoChannel(gochannel.Config{}, logger),
	}
}

func (p publisher) PublishOnTopic(topic topics.Topic, data any) error {
	artistBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	msg := message.NewMessage(watermill.NewUUID(), artistBytes)
	msg.Metadata.Set("topic", string(topic))
	queue := string(topics.TopicToQueueMap[topic])
	return p.client.Publish(queue, msg)
}

func (p publisher) Publish(topic string, messages ...*message.Message) error {
	return p.client.Publish(topic, messages...)
}

func (p publisher) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	return p.client.Subscribe(ctx, topic)
}

func (p publisher) Close() error {
	return p.client.Close()
}
