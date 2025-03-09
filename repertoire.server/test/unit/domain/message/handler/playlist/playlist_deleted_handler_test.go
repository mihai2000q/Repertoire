package playlist

import (
	"encoding/json"
	"errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"repertoire/server/domain/message/handler/playlist"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/service"
	"testing"
)

func TestPlaylistDeletedHandler_WhenPublishFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewPlaylistDeletedHandler(messagePublisherService)

	mockPlaylist := model.Playlist{ID: uuid.New()}

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockPlaylist)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
}

func TestPlaylistDeletedHandler_WhenSuccessful_ShouldPublishMessageToDeleteFromSearchEngine(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewPlaylistDeletedHandler(messagePublisherService)

	mockPlaylist := model.Playlist{ID: uuid.New()}
	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Run(func(args mock.Arguments) {
			ids := args.Get(1).([]string)
			assert.Len(t, ids, 1)
			assert.Contains(t, ids[0], mockPlaylist.ID.String())
		}).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(mockPlaylist)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	messagePublisherService.AssertExpectations(t)
}
