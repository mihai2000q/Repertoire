package service

import (
	"github.com/stretchr/testify/mock"
	"repertoire/server/data/message"
	"repertoire/server/internal/message/topics"
)

type MessagePublisherServiceMock struct {
	mock.Mock
}

func (m *MessagePublisherServiceMock) GetClient() message.Publisher {
	args := m.Called()
	return args.Get(0).(message.Publisher)
}

func (m *MessagePublisherServiceMock) Publish(topic topics.Topic, messagePayload any) error {
	args := m.Called(topic, messagePayload)
	return args.Error(0)
}
