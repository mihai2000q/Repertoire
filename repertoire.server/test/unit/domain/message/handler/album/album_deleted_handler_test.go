package album

import (
	"encoding/json"
	"errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"repertoire/server/domain/message/handler/album"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/service"
	"testing"
)

func TestAlbumDeletedHandler_WhenPublishFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := album.NewAlbumDeletedHandler(messagePublisherService)

	mockAlbum := model.Album{ID: uuid.New()}

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockAlbum)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
}

func TestAlbumDeletedHandler_WhenSuccessful_ShouldPublishMessageToDeleteFromSearchEngine(t *testing.T) {
	tests := []struct {
		name  string
		album model.Album
	}{
		{
			"without Songs",
			model.Album{ID: uuid.New()},
		},
		{
			"with Songs",
			model.Album{
				ID: uuid.New(),
				Songs: []model.Song{
					{ID: uuid.New()},
					{ID: uuid.New()},
					{ID: uuid.New()},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			messagePublisherService := new(service.MessagePublisherServiceMock)
			_uut := album.NewAlbumDeletedHandler(messagePublisherService)

			messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
				Run(func(args mock.Arguments) {
					ids := args.Get(1).([]string)
					assert.Len(t, ids, len(tt.album.Songs)+1)

					assert.Contains(t, ids[0], tt.album.ID.String())
					for i, song := range tt.album.Songs {
						assert.Contains(t, ids[1+i], song.ID.String())
					}
				}).
				Return(nil).
				Once()

			// when
			payload, _ := json.Marshal(tt.album)
			msg := message.NewMessage("1", payload)
			err := _uut.Handle(msg)

			// then
			assert.NoError(t, err)

			messagePublisherService.AssertExpectations(t)
		})
	}
}
