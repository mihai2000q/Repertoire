package message

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"repertoire/server/data/logger"
)

type Publisher interface {
	Publish(topic string, messages ...*message.Message) error
	Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
	Close() error
}

func NewPublisher(logger *logger.WatermillLogger) Publisher {
	return gochannel.NewGoChannel(gochannel.Config{}, logger)
}
