package message

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type Publisher interface {
	Publish(topic string, messages ...*message.Message) error
	Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
	Close() error
}

func NewPublisher() Publisher {
	logger := watermill.NewStdLogger(false, false)
	return gochannel.NewGoChannel(gochannel.Config{}, logger)
}
