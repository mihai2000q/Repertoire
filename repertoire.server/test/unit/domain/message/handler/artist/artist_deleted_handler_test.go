package artist

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"repertoire/server/domain/message/handler/artist"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/service"
	"repertoire/server/test/unit/domain/provider"
	"testing"
)

func TestArtistDeletedHandler_WhenPublishDeleteFromSearchEngineFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := artist.NewArtistDeletedHandler(messagePublisherService, nil, nil)

	mockArtist := model.Artist{ID: uuid.New()}

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockArtist)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
}

func TestArtistDeletedHandler_WhenGetDocumentsFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := artist.NewArtistDeletedHandler(messagePublisherService, searchEngineService, nil)

	mockArtist := model.Artist{ID: uuid.New()}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	searchEngineService.On("GetDocuments", mock.IsType("")).
		Return([]map[string]any{}, internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockArtist)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
}

func TestArtistDeletedHandler_WhenPublishUpdateFromSearchEngineFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := artist.NewArtistDeletedHandler(messagePublisherService, searchEngineService, nil)

	mockArtist := model.Artist{ID: uuid.New()}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(nil).
		Once()

	documentsToUpdate := []map[string]any{
		{"artist": nil},
	}
	searchEngineService.On("GetDocuments", mock.IsType("")).
		Return(documentsToUpdate, nil).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.UpdateFromSearchEngineTopic, mock.IsType([]any{})).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockArtist)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
}

func TestArtistDeletedHandler_WhenPublishDeleteStorageFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := artist.NewArtistDeletedHandler(messagePublisherService, searchEngineService, storageFilePathProvider)

	mockArtist := model.Artist{ID: uuid.New()}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(nil).
		Once()

	searchEngineService.On("GetDocuments", mock.IsType("")).
		Return([]map[string]any{}, nil).
		Once()

	directoryPath := "some_path"
	storageFilePathProvider.On("GetArtistDirectoryPath", mockArtist).Return(directoryPath).Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, []string{directoryPath}).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockArtist)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestArtistDeletedHandler_WhenDocumentsLengthIs0_ShouldNotReturnAnyError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := artist.NewArtistDeletedHandler(messagePublisherService, searchEngineService, storageFilePathProvider)

	mockArtist := model.Artist{ID: uuid.New()}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(nil).
		Once()

	searchEngineService.On("GetDocuments", mock.IsType("")).
		Return([]map[string]any{}, nil).
		Once()

	directoryPath := "some_path"
	storageFilePathProvider.On("GetArtistDirectoryPath", mockArtist).Return(directoryPath).Once()

	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, []string{directoryPath}).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(mockArtist)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestArtistDeletedHandler_WhenWithoutSongsOrAlbums_ShouldPublishMessagesToDeleteAndUpdateFromSearchEngine(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := artist.NewArtistDeletedHandler(messagePublisherService, searchEngineService, storageFilePathProvider)

	mockArtist := model.Artist{ID: uuid.New()}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Run(func(args mock.Arguments) {
			ids := args.Get(1).([]string)
			assert.Len(t, ids, len(mockArtist.Songs)+1)

			assert.Contains(t, ids[0], mockArtist.ID.String())
			for i, song := range mockArtist.Songs {
				assert.Contains(t, ids[1+i], song.ID.String())
			}
		}).
		Return(nil).
		Once()

	documentsToUpdate := []map[string]any{
		{
			"id":     uuid.New().String(),
			"title":  "Song 1",
			"artist": model.SongArtistSearch{ID: mockArtist.ID},
		},
		{
			"id":     uuid.New().String(),
			"title":  "Song 2",
			"artist": model.SongArtistSearch{ID: mockArtist.ID},
		},
		{
			"id":     uuid.New().String(),
			"title":  "Album 1",
			"artist": model.AlbumArtistSearch{ID: mockArtist.ID},
		},
	}

	filter := fmt.Sprintf("(type = %s OR type = %s) AND artist.id = %s", enums.Song, enums.Album, mockArtist.ID)
	searchEngineService.On("GetDocuments", filter).
		Return(documentsToUpdate, nil).
		Once()

	messagePublisherService.On("Publish", topics.UpdateFromSearchEngineTopic, mock.IsType([]any{})).
		Run(func(args mock.Arguments) {
			newDocuments := args.Get(1).([]any)
			assert.Len(t, newDocuments, len(documentsToUpdate))
			for _, song := range newDocuments {
				assert.Nil(t, song.(map[string]any)["artist"])
			}
		}).
		Return(nil).
		Once()

	directoryPath := "some_path"
	storageFilePathProvider.On("GetArtistDirectoryPath", mockArtist).Return(directoryPath).Once()

	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, []string{directoryPath}).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(mockArtist)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestArtistDeletedHandler_WhenWithSongsOrAlbums_ShouldPublishMessageToDeleteFromSearchEngine(t *testing.T) {
	tests := []struct {
		name   string
		artist model.Artist
	}{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			messagePublisherService := new(service.MessagePublisherServiceMock)
			searchEngineService := new(service.SearchEngineServiceMock)
			storageFilePathProvider := new(provider.StorageFilePathProviderMock)
			_uut := artist.NewArtistDeletedHandler(messagePublisherService, searchEngineService, storageFilePathProvider)

			messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
				Run(func(args mock.Arguments) {
					ids := args.Get(1).([]string)
					assert.Len(t, ids, len(tt.artist.Songs)+len(tt.artist.Albums)+1)

					assert.Contains(t, ids[0], tt.artist.ID.String())
					for i, song := range tt.artist.Songs {
						assert.Contains(t, ids[1+i], song.ID.String())
					}
					for i, album := range tt.artist.Albums {
						assert.Contains(t, ids[1+len(tt.artist.Songs)+i], album.ID.String())
					}
				}).
				Return(nil).
				Once()

			documentsToUpdate := []map[string]any{
				{
					"id":     uuid.New().String(),
					"title":  "Song 1",
					"artist": model.SongArtistSearch{ID: tt.artist.ID},
				},
				{
					"id":     uuid.New().String(),
					"title":  "Song 2",
					"artist": model.SongArtistSearch{ID: tt.artist.ID},
				},
				{
					"id":     uuid.New().String(),
					"title":  "Album 1",
					"artist": model.AlbumArtistSearch{ID: tt.artist.ID},
				},
			}

			filter := fmt.Sprintf("(type = %s OR type = %s) AND artist.id = %s", enums.Song, enums.Album, tt.artist.ID)
			searchEngineService.On("GetDocuments", filter).
				Return(documentsToUpdate, nil).
				Once()

			messagePublisherService.On("Publish", topics.UpdateFromSearchEngineTopic, mock.IsType([]any{})).
				Run(func(args mock.Arguments) {
					newDocuments := args.Get(1).([]any)
					assert.Len(t, newDocuments, len(documentsToUpdate))
					for _, song := range newDocuments {
						assert.Nil(t, song.(map[string]any)["artist"])
					}
				}).
				Return(nil).
				Once()

			// Cleanup Storage

			directoryPath := "some_path"
			storageFilePathProvider.On("GetArtistDirectoryPath", tt.artist).Return(directoryPath).Once()

			songDirectoryPath := "song_path"
			for _, song := range tt.artist.Songs {
				storageFilePathProvider.On("GetSongDirectoryPath", song).Return(songDirectoryPath).Once()
			}
			albumDirectoryPath := "album_path"
			for _, album := range tt.artist.Albums {
				storageFilePathProvider.On("GetAlbumDirectoryPath", album).Return(albumDirectoryPath).Once()
			}

			messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, mock.IsType([]string{})).
				Run(func(args mock.Arguments) {
					paths := args.Get(1).([]string)
					assert.Len(t, paths, len(tt.artist.Songs)+len(tt.artist.Albums)+1)
					assert.Equal(t, directoryPath, paths[0])
					for i := range tt.artist.Songs {
						assert.Equal(t, songDirectoryPath, paths[1+i])
					}
					for i := range tt.artist.Albums {
						assert.Equal(t, albumDirectoryPath, paths[1+i+len(tt.artist.Songs)])
					}
				}).
				Return(nil).
				Once()

			// when
			payload, _ := json.Marshal(tt.artist)
			msg := message.NewMessage("1", payload)
			err := _uut.Handle(msg)

			// then
			assert.NoError(t, err)

			searchEngineService.AssertExpectations(t)
			messagePublisherService.AssertExpectations(t)
			storageFilePathProvider.AssertExpectations(t)
		})
	}
}

func TestArtistDeletedHandler_WhenWithSongsAndAlbums_ShouldPublishMessageToDeleteFromSearchEngine(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := artist.NewArtistDeletedHandler(messagePublisherService, searchEngineService, storageFilePathProvider)

	mockArtist := model.Artist{
		ID: uuid.New(),
		Albums: []model.Album{
			{ID: uuid.New()},
			{ID: uuid.New()},
			{ID: uuid.New()},
		},
		Songs: []model.Song{
			{ID: uuid.New()},
			{ID: uuid.New()},
			{ID: uuid.New()},
		},
	}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Run(func(args mock.Arguments) {
			ids := args.Get(1).([]string)
			assert.Len(t, ids, len(mockArtist.Songs)+len(mockArtist.Albums)+1)

			assert.Contains(t, ids[0], mockArtist.ID.String())
			for i, song := range mockArtist.Songs {
				assert.Contains(t, ids[1+i], song.ID.String())
			}
			for i, album := range mockArtist.Albums {
				assert.Contains(t, ids[1+len(mockArtist.Songs)+i], album.ID.String())
			}
		}).
		Return(nil).
		Once()

	// Cleanup Storage

	directoryPath := "some_path"
	storageFilePathProvider.On("GetArtistDirectoryPath", mockArtist).Return(directoryPath).Once()

	songDirectoryPath := "song_path"
	for _, song := range mockArtist.Songs {
		storageFilePathProvider.On("GetSongDirectoryPath", song).Return(songDirectoryPath).Once()
	}
	albumDirectoryPath := "album_path"
	for _, album := range mockArtist.Albums {
		storageFilePathProvider.On("GetAlbumDirectoryPath", album).Return(albumDirectoryPath).Once()
	}

	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, mock.IsType([]string{})).
		Run(func(args mock.Arguments) {
			paths := args.Get(1).([]string)
			assert.Len(t, paths, len(mockArtist.Songs)+len(mockArtist.Albums)+1)
			assert.Equal(t, directoryPath, paths[0])
			for i := range mockArtist.Songs {
				assert.Equal(t, songDirectoryPath, paths[1+i])
			}
			for i := range mockArtist.Albums {
				assert.Equal(t, albumDirectoryPath, paths[1+i+len(mockArtist.Songs)])
			}
		}).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(mockArtist)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}
