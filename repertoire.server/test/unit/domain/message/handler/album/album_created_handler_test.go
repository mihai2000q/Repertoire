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

func TestAlbumCreatedHandler_WhenGetArtistFails_ShouldReturnError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := album.NewAlbumCreatedHandler(artistRepository, nil)

	mockAlbum := model.Album{Title: "Something", ArtistID: &[]uuid.UUID{uuid.New()}[0]}

	internalError := errors.New("internal error")
	artistRepository.On("Get", new(model.Artist), *mockAlbum.ArtistID).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockAlbum)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	artistRepository.AssertExpectations(t)
}

func TestAlbumCreatedHandler_WhenPublishFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := album.NewAlbumCreatedHandler(nil, messagePublisherService)

	mockAlbum := model.Album{ID: uuid.New()}

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.AddToSearchEngineTopic, mock.IsType([]any{})).
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

func TestAlbumCreatedHandler_WhenSuccessful_ShouldPublishMessageToAddAlbumToSearchEngine(t *testing.T) {
	tests := []struct {
		name  string
		album model.Album
	}{
		{
			"without Artist",
			model.Album{ID: uuid.New()},
		},
		{
			"with Artist",
			model.Album{
				ID:     uuid.New(),
				Artist: &model.Artist{ID: uuid.New()},
			},
		},
		{
			"with Artist ID",
			model.Album{
				ID:       uuid.New(),
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			artistRepository := new(repository.ArtistRepositoryMock)
			messagePublisherService := new(service.MessagePublisherServiceMock)
			_uut := album.NewAlbumCreatedHandler(artistRepository, messagePublisherService)

			var artist model.Artist
			if tt.album.ArtistID != nil {
				artist = model.Artist{ID: *tt.album.ArtistID}
				artistRepository.On("Get", new(model.Artist), *tt.album.ArtistID).
					Return(nil, &artist).
					Once()
			}

			messagePublisherService.On("Publish", topics.AddToSearchEngineTopic, mock.IsType([]any{})).
				Run(func(args mock.Arguments) {
					searches := args.Get(1).([]any)
					if tt.album.Artist != nil {
						assert.Len(t, searches, 2)
						assert.Contains(t, searches[1].(model.ArtistSearch).ID, tt.album.Artist.ID.String())
						assert.Equal(t, searches[0].(model.AlbumSearch).Artist.ID, tt.album.Artist.ID)
					} else {
						assert.Len(t, searches, 1)
					}

					assert.Contains(t, searches[0].(model.AlbumSearch).ID, tt.album.ID.String())
					if tt.album.ArtistID != nil {
						assert.Equal(t, searches[0].(model.AlbumSearch).Artist.ID, *tt.album.ArtistID)
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

			artistRepository.AssertExpectations(t)
			messagePublisherService.AssertExpectations(t)
		})
	}
}
