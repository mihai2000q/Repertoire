package album

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"repertoire/server/domain/message/handler/album"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/service"
	"repertoire/server/test/unit/domain/provider"
	"testing"
)

func TestAlbumDeletedHandler_WhenPublishDeleteFromSearchEngineFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := album.NewAlbumDeletedHandler(messagePublisherService, nil, nil)

	mockAlbum := model.Album{ID: uuid.New()}

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
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

func TestAlbumDeletedHandler_WhenGetDocumentsFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := album.NewAlbumDeletedHandler(messagePublisherService, searchEngineService, nil)

	mockAlbum := model.Album{ID: uuid.New()}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	searchEngineService.On("GetDocuments", mock.IsType("")).
		Return([]map[string]any{}, internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockAlbum)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
}

func TestAlbumDeletedHandler_WhenPublishUpdateFromSearchEngineFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := album.NewAlbumDeletedHandler(messagePublisherService, searchEngineService, nil)

	mockAlbum := model.Album{ID: uuid.New()}

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
	payload, _ := json.Marshal(mockAlbum)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
}

func TestAlbumDeletedHandler_WhenPublishDeleteDirectoriesStorageFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := album.NewAlbumDeletedHandler(messagePublisherService, searchEngineService, storageFilePathProvider)

	mockAlbum := model.Album{ID: uuid.New()}

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
	storageFilePathProvider.On("GetAlbumDirectoryPath", mockAlbum).Return(directoryPath).Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, []string{directoryPath}).
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
	searchEngineService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestAlbumDeletedHandler_WhenDocumentsLengthIs0_ShouldNotReturnAnyError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := album.NewAlbumDeletedHandler(messagePublisherService, searchEngineService, storageFilePathProvider)

	mockAlbum := model.Album{ID: uuid.New()}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(nil).
		Once()

	searchEngineService.On("GetDocuments", mock.IsType("")).
		Return([]map[string]any{}, nil).
		Once()

	directoryPath := "some_path"
	storageFilePathProvider.On("GetAlbumDirectoryPath", mockAlbum).Return(directoryPath).Once()

	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, []string{directoryPath}).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(mockAlbum)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestAlbumDeletedHandler_WhenSuccessful_ShouldPublishMessageToDeleteFromSearchEngine(t *testing.T) {
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
			messagePublisherService := new(service.MessagePublisherServiceMock)
			searchEngineService := new(service.SearchEngineServiceMock)
			storageFilePathProvider := new(provider.StorageFilePathProviderMock)
			_uut := album.NewAlbumDeletedHandler(messagePublisherService, searchEngineService, storageFilePathProvider)

			messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
				Run(func(args mock.Arguments) {
					ids := args.Get(1).([]string)
					assert.Len(t, ids, len(tt.album.Songs)+1)

					assert.Contains(t, ids[0], tt.album.ID.String())
					for i, song := range tt.album.Songs {
						assert.Contains(t, ids[1+i], song.ID.String())
					}
				}).
				Return(nil).
				Once()

			if len(tt.album.Songs) == 0 {
				songsToUpdate := []map[string]any{
					{
						"id":    uuid.New().String(),
						"title": "Song 1",
						"album": model.SongAlbumSearch{ID: tt.album.ID},
					},
					{
						"id":    uuid.New().String(),
						"title": "Song 2",
						"album": model.SongAlbumSearch{ID: tt.album.ID},
					},
				}
				filter := fmt.Sprintf("type = %s AND album.id = %s", enums.Song, tt.album.ID)
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
			}

			directoryPath := "some_path"
			storageFilePathProvider.On("GetAlbumDirectoryPath", tt.album).Return(directoryPath).Once()

			songDirectoryPath := "song_directory_path"
			for _, song := range tt.album.Songs {
				storageFilePathProvider.On("GetSongDirectoryPath", song).Return(songDirectoryPath).Once()
			}

			messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, mock.IsType([]string{})).
				Run(func(args mock.Arguments) {
					directoryPaths := args.Get(1).([]string)
					assert.Len(t, directoryPaths, len(tt.album.Songs)+1)
					assert.Equal(t, directoryPaths[0], directoryPath)
					for i, path := range directoryPaths {
						if i == 0 {
							continue
						}
						assert.Equal(t, songDirectoryPath, path)
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

			messagePublisherService.AssertExpectations(t)
			searchEngineService.AssertExpectations(t)
			storageFilePathProvider.AssertExpectations(t)
		})
	}
}
