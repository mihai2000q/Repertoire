package playlist

import (
	"encoding/json"
	"errors"
	"repertoire/server/domain/message/handler/playlist"
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

func TestPlaylistsDeletedHandler_WhenPublishDeleteFromSearchEngineFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewPlaylistsDeletedHandler(nil, messagePublisherService)

	mockPlaylists := []model.Playlist{{ID: uuid.New()}}

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockPlaylists)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
}

func TestPlaylistsDeletedHandler_WhenPublishDeleteStorageFails_ShouldReturnError(t *testing.T) {
	// given
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewPlaylistsDeletedHandler(storageFilePathProvider, messagePublisherService)

	mockPlaylists := []model.Playlist{{ID: uuid.New()}}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(nil).
		Once()

	directoryPath := "some_path"
	storageFilePathProvider.On("GetPlaylistDirectoryPath", mock.IsType(mockPlaylists[0])).
		Return(directoryPath).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, mock.IsType([]string{})).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockPlaylists)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	storageFilePathProvider.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestPlaylistsDeletedHandler_WhenSuccessful_ShouldPublishMessages(t *testing.T) {
	// given
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewPlaylistsDeletedHandler(storageFilePathProvider, messagePublisherService)

	mockPlaylists := []model.Playlist{
		{ID: uuid.New()},
		{ID: uuid.New()},
		{ID: uuid.New()},
	}
	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Run(func(args mock.Arguments) {
			ids := args.Get(1).([]string)
			assert.Len(t, ids, len(mockPlaylists))
			for i := range ids {
				assert.Equal(t, mockPlaylists[i].ToSearch().ID, ids[i])
			}
		}).
		Return(nil).
		Once()

	var directoryPaths []string
	for _, p := range mockPlaylists {
		directoryPath := "some_path" + p.ID.String()
		directoryPaths = append(directoryPaths, directoryPath)
		storageFilePathProvider.On("GetPlaylistDirectoryPath", p).Return(directoryPath).Once()
	}

	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, directoryPaths).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(mockPlaylists)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	storageFilePathProvider.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}
