package song

import (
	"encoding/json"
	"errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"repertoire/server/domain/message/handler/song"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"
)

func TestSongUpdatedHandler_WhenGetSongFails_ShouldReturnError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewSongUpdatedHandler(songRepository, messagePublisherService)

	mockSong := model.Song{ID: uuid.New()}

	internalError := errors.New("internal error")
	songRepository.On("GetWithArtistAndAlbum", &mockSong, mockSong.ID).
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

func TestSongUpdatedHandler_WhenPublishFails_ShouldReturnError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewSongUpdatedHandler(songRepository, messagePublisherService)

	mockSong := model.Song{ID: uuid.New()}

	songRepository.On("GetWithArtistAndAlbum", &mockSong, mockSong.ID).
		Return(nil, &mockSong).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.UpdateFromSearchEngineTopic, mock.IsType([]any{})).
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

func TestSongUpdatedHandler_WhenSuccessful_ShouldPublishMessageToDeleteFromSearchEngine(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewSongUpdatedHandler(songRepository, messagePublisherService)

	mockSong := model.Song{ID: uuid.New()}
	songRepository.On("GetWithArtistAndAlbum", &mockSong, mockSong.ID).
		Return(nil, &mockSong).
		Once()

	messagePublisherService.On("Publish", topics.UpdateFromSearchEngineTopic, mock.IsType([]any{})).
		Run(func(args mock.Arguments) {
			searches := args.Get(1).([]any)
			assert.Len(t, searches, 1)
			assert.Contains(t, searches[0].(model.SongSearch).ID, mockSong.ID.String())
		}).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(mockSong)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	messagePublisherService.AssertExpectations(t)
}
