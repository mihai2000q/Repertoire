package song

import (
	"encoding/json"
	"errors"
	"repertoire/server/domain/message/handler/song"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/service"
	"repertoire/server/test/unit/domain/provider"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSongsDeletedHandler_WhenPublishDeleteFromSearchEngineFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewSongsDeletedHandler(nil, messagePublisherService)

	mockSongs := []model.Song{{ID: uuid.New()}}

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockSongs)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
}

func TestSongsDeletedHandler_WhenPublishDeleteStorageFails_ShouldReturnError(t *testing.T) {
	// given
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewSongsDeletedHandler(storageFilePathProvider, messagePublisherService)

	mockSongs := []model.Song{{ID: uuid.New()}}
	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(nil).
		Once()

	directoryPath := "some_path"
	storageFilePathProvider.On("GetSongDirectoryPath", mock.IsType(model.Song{})).
		Return(directoryPath).
		Times(len(mockSongs))

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, mock.IsType([]string{})).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockSongs)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	storageFilePathProvider.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestSongsDeletedHandler_WhenSuccessful_ShouldPublishMessageToDeleteFromSearchEngine(t *testing.T) {
	// given
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewSongsDeletedHandler(storageFilePathProvider, messagePublisherService)

	mockSongs := []model.Song{
		{ID: uuid.New()},
		{ID: uuid.New()},
	}
	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Run(func(args mock.Arguments) {
			ids := args.Get(1).([]string)
			assert.Len(t, ids, len(mockSongs))
			for i := range ids {
				assert.Equal(t, mockSongs[i].ToSearch().ID, ids[i])
			}
		}).
		Return(nil).
		Once()

	var directoryPaths []string
	for _, s := range mockSongs {
		directoryPath := "some_path" + s.ID.String()
		directoryPaths = append(directoryPaths, directoryPath)
		storageFilePathProvider.On("GetSongDirectoryPath", s).Return(directoryPath).Once()
	}

	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, directoryPaths).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(mockSongs)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	storageFilePathProvider.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}
