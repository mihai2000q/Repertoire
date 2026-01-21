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

func TestSongUpdatedHandler_WhenGetSongsFails_ShouldReturnError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewSongsUpdatedHandler(songRepository, nil)

	songIDs := []uuid.UUID{uuid.New()}

	internalError := errors.New("internal error")
	songRepository.On("GetAllByIDsWithArtistAndAlbum", new([]model.Song), songIDs).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(songIDs)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	songRepository.AssertExpectations(t)
}

func TestSongUpdatedHandler_WhenPublishFails_ShouldReturnError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewSongsUpdatedHandler(songRepository, messagePublisherService)

	songIDs := []uuid.UUID{uuid.New()}

	songs := []model.Song{{ID: uuid.New()}}
	songRepository.On("GetAllByIDsWithArtistAndAlbum", new([]model.Song), songIDs).
		Return(nil, &songs).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.UpdateFromSearchEngineTopic, mock.IsType([]any{})).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(songIDs)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	songRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestSongUpdatedHandler_WhenThereAreNoSongs_ShouldReturnNoError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewSongsUpdatedHandler(songRepository, nil)

	songIDs := []uuid.UUID{uuid.New()}

	songRepository.On("GetAllByIDsWithArtistAndAlbum", new([]model.Song), songIDs).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(songIDs)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	songRepository.AssertExpectations(t)
}

func TestSongUpdatedHandler_WhenSuccessful_ShouldPublishMessageToUpdateFromSearchEngine(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewSongsUpdatedHandler(songRepository, messagePublisherService)

	songIDs := []uuid.UUID{uuid.New(), uuid.New()}

	songs := []model.Song{{ID: uuid.New()}}
	songRepository.On("GetAllByIDsWithArtistAndAlbum", new([]model.Song), songIDs).
		Return(nil, &songs).
		Once()

	messagePublisherService.On("Publish", topics.UpdateFromSearchEngineTopic, mock.IsType([]any{})).
		Run(func(args mock.Arguments) {
			searches := args.Get(1).([]any)
			assert.Len(t, searches, len(songs))
			for i := range searches {
				assert.Contains(t, searches[i].(model.SongSearch).ID, songs[i].ID.String())
			}
		}).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(songIDs)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	songRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}
