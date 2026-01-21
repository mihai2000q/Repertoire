package album

import (
	"encoding/json"
	"errors"
	"fmt"
	"repertoire/server/domain/message/handler/album"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/service"
	"repertoire/server/test/unit/domain/provider"
	"strings"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAlbumsDeletedHandler_WhenPublishDeleteFromSearchEngineFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := album.NewAlbumsDeletedHandler(messagePublisherService, nil, nil)

	mockAlbums := []model.Album{{ID: uuid.New()}}

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockAlbums)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
}

func TestAlbumsDeletedHandler_WhenGetDocumentsFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := album.NewAlbumsDeletedHandler(messagePublisherService, searchEngineService, nil)

	mockAlbums := []model.Album{{ID: uuid.New()}}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	searchEngineService.On("GetDocuments", mock.IsType("")).
		Return([]map[string]any{}, internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockAlbums)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
}

func TestAlbumsDeletedHandler_WhenPublishUpdateFromSearchEngineFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := album.NewAlbumsDeletedHandler(messagePublisherService, searchEngineService, nil)

	mockAlbums := []model.Album{{ID: uuid.New()}}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(nil).
		Once()

	songsToUpdate := []map[string]any{
		{"id": uuid.New().String()},
	}
	searchEngineService.On("GetDocuments", mock.IsType("")).
		Return(songsToUpdate, nil).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.UpdateFromSearchEngineTopic, mock.IsType([]any{})).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockAlbums)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
}

func TestAlbumsDeletedHandler_WhenPublishDeleteDirectoriesStorageFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := album.NewAlbumsDeletedHandler(messagePublisherService, searchEngineService, storageFilePathProvider)

	mockAlbums := []model.Album{{ID: uuid.New()}}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(nil).
		Once()

	songsToUpdate := []map[string]any{
		{"id": uuid.New().String()},
	}
	searchEngineService.On("GetDocuments", mock.IsType("")).
		Return(songsToUpdate, nil).
		Once()

	messagePublisherService.On("Publish", topics.UpdateFromSearchEngineTopic, mock.IsType([]any{})).
		Return(nil).
		Once()

	directoryPath := "some_path"
	storageFilePathProvider.On("GetAlbumDirectoryPath", mock.IsType(mockAlbums[0])).
		Return(directoryPath).
		Times(len(mockAlbums))

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, []string{directoryPath}).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockAlbums)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestAlbumsDeletedHandler_WhenDocumentsLengthIs0_ShouldNotReturnAnyError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := album.NewAlbumsDeletedHandler(messagePublisherService, searchEngineService, storageFilePathProvider)

	mockAlbums := []model.Album{
		{ID: uuid.New()},
		{ID: uuid.New()},
	}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(nil).
		Once()

	searchEngineService.On("GetDocuments", mock.IsType("")).
		Return([]map[string]any{}, nil).
		Once()

	directoryPath := "some_path"
	storageFilePathProvider.On("GetAlbumDirectoryPath", mock.IsType(mockAlbums[0])).
		Return(directoryPath).
		Times(len(mockAlbums))

	var directoryPaths []string
	for range mockAlbums {
		directoryPaths = append(directoryPaths, directoryPath)
	}
	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, directoryPaths).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(mockAlbums)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestAlbumsDeletedHandler_WhenWithoutSongs_ShouldPublishMessageToDeleteFromSearchEngine(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := album.NewAlbumsDeletedHandler(messagePublisherService, searchEngineService, storageFilePathProvider)

	mockAlbums := []model.Album{
		{ID: uuid.New()},
		{ID: uuid.New()},
	}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Run(func(args mock.Arguments) {
			ids := args.Get(1).([]string)
			assert.Len(t, ids, len(mockAlbums))
			for i, alb := range mockAlbums {
				assert.Equal(t, alb.ToSearch().ID, ids[i])
			}
		}).
		Return(nil).
		Once()

	songsToUpdate := []map[string]any{
		{
			"id":    uuid.New().String(),
			"title": "Song 1",
			"album": model.SongAlbumSearch{ID: mockAlbums[0].ID},
		},
		{
			"id":    uuid.New().String(),
			"title": "Song 2",
			"album": model.SongAlbumSearch{ID: mockAlbums[1].ID},
		},
	}

	var albumIDs []string
	for _, alb := range mockAlbums {
		albumIDs = append(albumIDs, alb.ID.String())
	}
	filter := fmt.Sprintf("type = %s AND album.id IN [%s]", enums.Song, strings.Join(albumIDs, ", "))
	searchEngineService.On("GetDocuments", filter).
		Return(songsToUpdate, nil).
		Once()

	messagePublisherService.On("Publish", topics.UpdateFromSearchEngineTopic, mock.IsType([]any{})).
		Run(func(args mock.Arguments) {
			newSongs := args.Get(1).([]any)
			assert.Len(t, newSongs, len(songsToUpdate))
			for _, song := range newSongs {
				assert.Nil(t, song.(model.SongSearch).Album)
			}
		}).
		Return(nil).
		Once()

	var directoryPaths []string
	for _, alb := range mockAlbums {
		directoryPath := "some_path" + alb.ID.String()
		directoryPaths = append(directoryPaths, directoryPath)
		storageFilePathProvider.On("GetAlbumDirectoryPath", alb).Return(directoryPath).Once()
	}

	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, directoryPaths).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(mockAlbums)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestAlbumsDeletedHandler_WhenWithSongs_ShouldPublishMessageToDeleteFromSearchEngine(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := album.NewAlbumsDeletedHandler(messagePublisherService, searchEngineService, storageFilePathProvider)

	mockAlbums := []model.Album{
		{
			ID:    uuid.New(),
			Title: "With Songs",
			Songs: []model.Song{
				{ID: uuid.New()},
				{ID: uuid.New()},
				{ID: uuid.New()},
			},
		},
		{
			ID:    uuid.New(),
			Title: "Without Songs",
		},
	}

	var songSearchIDs []string
	for _, alb := range mockAlbums {
		for _, song := range alb.Songs {
			songSearchIDs = append(songSearchIDs, song.ToSearch().ID)
		}
	}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Run(func(args mock.Arguments) {
			ids := args.Get(1).([]string)
			assert.Len(t, ids, len(mockAlbums)+len(songSearchIDs))

			for i, alb := range mockAlbums {
				assert.Equal(t, alb.ToSearch().ID, ids[i])
			}
			for i, id := range songSearchIDs {
				assert.Equal(t, ids[i+len(mockAlbums)], id)
			}
		}).
		Return(nil).
		Once()

	var directoryPaths []string
	for _, alb := range mockAlbums {
		directoryPath := "some_path" + alb.ID.String()
		directoryPaths = append(directoryPaths, directoryPath)
		storageFilePathProvider.On("GetAlbumDirectoryPath", alb).Return(directoryPath).Once()

		for _, song := range alb.Songs {
			songDirectoryPath := "song_directory_path" + song.ID.String()
			directoryPaths = append(directoryPaths, songDirectoryPath)
			storageFilePathProvider.On("GetSongDirectoryPath", song).Return(songDirectoryPath).Once()
		}
	}

	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, directoryPaths).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(mockAlbums)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}
