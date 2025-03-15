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
	"repertoire/server/test/unit/domain/provider"
	"testing"
)

func TestPlaylistDeletedHandler_WhenPublishDeleteFromSearchEngineFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewPlaylistDeletedHandler(nil, messagePublisherService)

	mockPlaylist := model.Playlist{ID: uuid.New()}

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
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

func TestPlaylistDeletedHandler_WhenPublishDeleteStorageFails_ShouldReturnError(t *testing.T) {
	// given
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewPlaylistDeletedHandler(storageFilePathProvider, messagePublisherService)

	mockPlaylist := model.Playlist{ID: uuid.New()}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(nil).
		Once()

	directoryPath := "some_path"
	storageFilePathProvider.On("GetPlaylistDirectoryPath", mockPlaylist).Return(directoryPath).Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, []string{directoryPath}).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockPlaylist)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	storageFilePathProvider.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestPlaylistDeletedHandler_WhenSuccessful_ShouldPublishMessages(t *testing.T) {
	// given
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewPlaylistDeletedHandler(storageFilePathProvider, messagePublisherService)

	mockPlaylist := model.Playlist{ID: uuid.New()}
	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Run(func(args mock.Arguments) {
			ids := args.Get(1).([]string)
			assert.Len(t, ids, 1)
			assert.Contains(t, ids[0], mockPlaylist.ID.String())
		}).
		Return(nil).
		Once()

	directoryPath := "some_path"
	storageFilePathProvider.On("GetPlaylistDirectoryPath", mockPlaylist).Return(directoryPath).Once()

	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, []string{directoryPath}).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(mockPlaylist)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	storageFilePathProvider.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}
