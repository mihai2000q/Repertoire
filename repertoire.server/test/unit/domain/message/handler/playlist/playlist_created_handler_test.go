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

func TestPlaylistCreatedHandler_WhenPublishFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewPlaylistCreatedHandler(messagePublisherService)

	mockPlaylist := model.Playlist{ID: uuid.New()}

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.AddToSearchEngineTopic, mock.IsType([]any{})).
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

func TestPlaylistCreatedHandler_WhenSuccessful_ShouldPublishMessageToAddToSearchEngine(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewPlaylistCreatedHandler(messagePublisherService)

	mockPlaylist := model.Playlist{ID: uuid.New()}

	messagePublisherService.On("Publish", topics.AddToSearchEngineTopic, mock.IsType([]any{})).
		Run(func(args mock.Arguments) {
			searches := args.Get(1).([]any)
			assert.Len(t, searches, 1)
			assert.Contains(t, searches[0].(model.PlaylistSearch).ID, mockPlaylist.ID.String())
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
