package song

import (
	"encoding/json"
	"errors"
	"repertoire/server/domain/message/handler/song"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSongCreatedHandler_WhenGetArtistFails_ShouldReturnError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := song.NewSongCreatedHandler(artistRepository, nil, nil)

	mockSong := model.Song{ID: uuid.New(), ArtistID: &[]uuid.UUID{uuid.New()}[0]}

	internalError := errors.New("internal error")
	artistRepository.On("Get", new(model.Artist), *mockSong.ArtistID).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockSong)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	artistRepository.AssertExpectations(t)
}

func TestSongCreatedHandler_WhenGetAlbumFails_ShouldReturnError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := song.NewSongCreatedHandler(nil, albumRepository, nil)

	mockSong := model.Song{ID: uuid.New(), AlbumID: &[]uuid.UUID{uuid.New()}[0]}

	internalError := errors.New("internal error")
	albumRepository.On("Get", new(model.Album), *mockSong.AlbumID).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockSong)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	albumRepository.AssertExpectations(t)
}

func TestSongCreatedHandler_WhenPublishFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewSongCreatedHandler(nil, nil, messagePublisherService)

	mockSong := model.Song{ID: uuid.New()}

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.AddToSearchEngineTopic, mock.IsType([]any{})).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockSong)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
}

func TestSongCreatedHandler_WhenSuccessful_ShouldPublishMessageToAddToSearchEngine(t *testing.T) {
	tests := []struct {
		name string
		song model.Song
	}{
		{
			"without entities",
			model.Song{ID: uuid.New()},
		},
		{
			"with Artist and album",
			model.Song{
				ID:     uuid.New(),
				Artist: &model.Artist{ID: uuid.New()},
				Album:  &model.Album{ID: uuid.New()},
			},
		},
		{
			"with Artist ID",
			model.Song{
				ID:       uuid.New(),
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
			},
		},
		{
			"with Album and Artist ID",
			model.Song{
				ID:       uuid.New(),
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
				Album:    &model.Album{ID: uuid.New()},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			artistRepository := new(repository.ArtistRepositoryMock)
			albumRepository := new(repository.AlbumRepositoryMock)
			messagePublisherService := new(service.MessagePublisherServiceMock)
			_uut := song.NewSongCreatedHandler(artistRepository, albumRepository, messagePublisherService)

			var artist model.Artist
			if tt.song.ArtistID != nil {
				artist = model.Artist{ID: *tt.song.ArtistID}
				artistRepository.On("Get", new(model.Artist), *tt.song.ArtistID).
					Return(nil, &artist).
					Once()
			}

			messagePublisherService.On("Publish", topics.AddToSearchEngineTopic, mock.IsType([]any{})).
				Run(func(args mock.Arguments) {
					searches := args.Get(1).([]any)
					var artistIndex *int
					var albumIndex *int
					if tt.song.Artist != nil && tt.song.Album != nil {
						assert.Len(t, searches, 3)
						artistIndex = &[]int{1}[0]
						albumIndex = &[]int{2}[0]
					} else if tt.song.Artist != nil {
						assert.Len(t, searches, 2)
						artistIndex = &[]int{1}[0]
					} else if tt.song.Album != nil {
						assert.Len(t, searches, 2)
						albumIndex = &[]int{1}[0]
					}

					assert.Contains(t, searches[0].(model.SongSearch).ID, tt.song.ID.String())
					if tt.song.ArtistID != nil {
						assert.Equal(t, searches[0].(model.SongSearch).Artist.ID, *tt.song.ArtistID)
					}
					if tt.song.AlbumID != nil {
						assert.Equal(t, searches[0].(model.SongSearch).Album.ID, *tt.song.AlbumID)
					}
					if artistIndex != nil {
						assert.Contains(t, searches[*artistIndex].(model.ArtistSearch).ID, tt.song.Artist.ID.String())
					}
					if albumIndex != nil {
						assert.Contains(t, searches[*albumIndex].(model.AlbumSearch).ID, tt.song.Album.ID.String())
					}
				}).
				Return(nil).
				Once()

			// when
			payload, _ := json.Marshal(tt.song)
			msg := message.NewMessage("1", payload)
			err := _uut.Handle(msg)

			// then
			assert.NoError(t, err)

			artistRepository.AssertExpectations(t)
			albumRepository.AssertExpectations(t)
			messagePublisherService.AssertExpectations(t)
		})
	}
}
