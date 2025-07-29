package artist

import (
	"encoding/json"
	"errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"repertoire/server/domain/message/handler/artist"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"
)

func TestArtistUpdatedHandler_WhenGetArtistFails_ShouldReturnError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewArtistUpdatedHandler(artistRepository, nil)

	artistID := uuid.New()

	internalError := errors.New("internal error")
	artistRepository.On("GetWithSongsOrAlbums", new(model.Artist), artistID, true, true).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(artistID)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, internalError, err)

	artistRepository.AssertExpectations(t)
}

func TestArtistUpdatedHandler_WhenPublishFails_ShouldReturnError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := artist.NewArtistUpdatedHandler(artistRepository, messagePublisherService)

	mockArtist := model.Artist{ID: uuid.New()}
	artistRepository.On("GetWithSongsOrAlbums", new(model.Artist), mockArtist.ID, true, true).
		Return(nil, &mockArtist).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.UpdateFromSearchEngineTopic, mock.IsType([]any{})).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockArtist.ID)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, internalError, err)

	artistRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestArtistUpdatedHandler_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name   string
		artist model.Artist
	}{
		{
			"without Songs or Albums",
			model.Artist{ID: uuid.New()},
		},
		{
			"with Songs",
			model.Artist{
				ID: uuid.New(),
				Songs: []model.Song{
					{ID: uuid.New()},
					{ID: uuid.New()},
					{ID: uuid.New()},
				},
			},
		},
		{
			"with Albums",
			model.Artist{
				ID: uuid.New(),
				Albums: []model.Album{
					{ID: uuid.New()},
					{ID: uuid.New()},
					{ID: uuid.New()},
				},
			},
		},
		{
			"with Albums and Songs",
			model.Artist{
				ID: uuid.New(),
				Songs: []model.Song{
					{ID: uuid.New()},
					{ID: uuid.New()},
					{ID: uuid.New()},
				},
				Albums: []model.Album{
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
			artistRepository := new(repository.ArtistRepositoryMock)
			messagePublisherService := new(service.MessagePublisherServiceMock)
			_uut := artist.NewArtistUpdatedHandler(artistRepository, messagePublisherService)

			artistRepository.On("GetWithSongsOrAlbums", new(model.Artist), tt.artist.ID, true, true).
				Return(nil, &tt.artist).
				Once()

			messagePublisherService.On("Publish", topics.UpdateFromSearchEngineTopic, mock.IsType([]any{})).
				Run(func(args mock.Arguments) {
					searches := args.Get(1).([]any)
					assert.Len(t, searches, len(tt.artist.Songs)+len(tt.artist.Albums)+1)

					assert.Contains(t, searches[0].(model.ArtistSearch).ID, tt.artist.ID.String())
					for i, song := range tt.artist.Songs {
						assert.Contains(t, searches[1+i].(model.SongSearch).ID, song.ID.String())
						assert.Equal(t, searches[1+i].(model.SongSearch).Artist.ID, tt.artist.ID)
					}
					for i, album := range tt.artist.Albums {
						assert.Contains(t, searches[1+len(tt.artist.Songs)+i].(model.AlbumSearch).ID, album.ID.String())
						assert.Equal(t, searches[1+len(tt.artist.Songs)+i].(model.AlbumSearch).Artist.ID, tt.artist.ID)
					}
				}).
				Return(nil).
				Once()

			// when
			payload, _ := json.Marshal(tt.artist.ID)
			msg := message.NewMessage("1", payload)
			err := _uut.Handle(msg)

			// then
			assert.NoError(t, err)

			artistRepository.AssertExpectations(t)
			messagePublisherService.AssertExpectations(t)
		})
	}
}
