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
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"
)

func TestAlbumUpdatedHandler_WhenGetAlbumFails_ShouldReturnError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewAlbumUpdatedHandler(albumRepository, nil)

	albumID := uuid.New()

	internalError := errors.New("internal error")
	albumRepository.On("GetWithSongsAndArtist", new(model.Album), albumID).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(albumID)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, internalError, err)

	albumRepository.AssertExpectations(t)
}

func TestAlbumUpdatedHandler_WhenPublishFails_ShouldReturnError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := album.NewAlbumUpdatedHandler(albumRepository, messagePublisherService)

	mockAlbum := model.Album{ID: uuid.New()}
	albumRepository.On("GetWithSongsAndArtist", new(model.Album), mockAlbum.ID).
		Return(nil, &mockAlbum).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.UpdateFromSearchEngineTopic, mock.IsType([]any{})).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockAlbum.ID)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, internalError, err)

	albumRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestAlbumUpdatedHandler_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
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
			albumRepository := new(repository.AlbumRepositoryMock)
			messagePublisherService := new(service.MessagePublisherServiceMock)
			_uut := album.NewAlbumUpdatedHandler(albumRepository, messagePublisherService)

			albumRepository.On("GetWithSongsAndArtist", new(model.Album), tt.album.ID).
				Return(nil, &tt.album).
				Once()

			messagePublisherService.On("Publish", topics.UpdateFromSearchEngineTopic, mock.IsType([]any{})).
				Run(func(args mock.Arguments) {
					searches := args.Get(1).([]any)
					assert.Len(t, searches, len(tt.album.Songs)+1)

					assert.Contains(t, searches[0].(model.AlbumSearch).ID, tt.album.ID.String())
					for i, song := range tt.album.Songs {
						assert.Contains(t, searches[1+i].(model.SongSearch).ID, song.ID.String())
						assert.Equal(t, searches[1+i].(model.SongSearch).Album.ID, tt.album.ID)
					}
				}).
				Return(nil).
				Once()

			// when
			payload, _ := json.Marshal(tt.album.ID)
			msg := message.NewMessage("1", payload)
			err := _uut.Handle(msg)

			// then
			assert.NoError(t, err)

			albumRepository.AssertExpectations(t)
			messagePublisherService.AssertExpectations(t)
		})
	}
}
