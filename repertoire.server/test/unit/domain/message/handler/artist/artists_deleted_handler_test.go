package artist

import (
	"encoding/json"
	"errors"
	"fmt"
	"repertoire/server/domain/message/handler/artist"
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

func TestArtistsDeletedHandler_WhenPublishDeleteFromSearchEngineFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := artist.NewArtistsDeletedHandler(messagePublisherService, nil, nil)

	mockArtists := []model.Artist{{ID: uuid.New()}}

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockArtists)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
}

func TestArtistsDeletedHandler_WhenGetDocumentsFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := artist.NewArtistsDeletedHandler(messagePublisherService, searchEngineService, nil)

	mockArtists := []model.Artist{{ID: uuid.New()}}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	searchEngineService.On("GetDocuments", mock.IsType("")).
		Return([]map[string]any{}, internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockArtists)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
}

func TestArtistsDeletedHandler_WhenPublishUpdateFromSearchEngineFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := artist.NewArtistsDeletedHandler(messagePublisherService, searchEngineService, nil)

	mockArtists := []model.Artist{{ID: uuid.New()}}

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
	payload, _ := json.Marshal(mockArtists)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
}

func TestArtistsDeletedHandler_WhenPublishDeleteStorageFails_ShouldReturnError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := artist.NewArtistsDeletedHandler(messagePublisherService, searchEngineService, storageFilePathProvider)

	mockArtists := []model.Artist{{ID: uuid.New()}}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(nil).
		Once()

	searchEngineService.On("GetDocuments", mock.IsType("")).
		Return([]map[string]any{}, nil).
		Once()

	directoryPath := "some_path"
	storageFilePathProvider.On("GetArtistDirectoryPath", mock.IsType(mockArtists[0])).
		Return(directoryPath).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, []string{directoryPath}).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(mockArtists)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestArtistsDeletedHandler_WhenDocumentsLengthIs0_ShouldNotReturnAnyError(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := artist.NewArtistsDeletedHandler(messagePublisherService, searchEngineService, storageFilePathProvider)

	mockArtists := []model.Artist{{ID: uuid.New()}}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(nil).
		Once()

	searchEngineService.On("GetDocuments", mock.IsType("")).
		Return([]map[string]any{}, nil).
		Once()

	directoryPath := "some_path"
	storageFilePathProvider.On("GetArtistDirectoryPath", mock.IsType(mockArtists[0])).Return(directoryPath).Once()

	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, []string{directoryPath}).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(mockArtists)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestArtistsDeletedHandler_WhenWithoutSongsOrAlbums_ShouldPublishMessagesToDeleteAndUpdateFromSearchEngine(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := artist.NewArtistsDeletedHandler(messagePublisherService, searchEngineService, storageFilePathProvider)

	mockArtists := []model.Artist{{ID: uuid.New()}, {ID: uuid.New()}}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Run(func(args mock.Arguments) {
			ids := args.Get(1).([]string)
			assert.Len(t, ids, len(mockArtists))
			for i, art := range mockArtists {
				assert.Equal(t, art.ToSearch().ID, ids[i])
			}
		}).
		Return(nil).
		Once()

	documentsToUpdate := []map[string]any{
		{
			"id":     uuid.New().String(),
			"title":  "Song 1",
			"artist": model.SongArtistSearch{ID: mockArtists[0].ID},
		},
		{
			"id":     uuid.New().String(),
			"title":  "Song 2",
			"artist": model.SongArtistSearch{ID: mockArtists[1].ID},
		},
		{
			"id":     uuid.New().String(),
			"title":  "Album 1",
			"artist": model.AlbumArtistSearch{ID: mockArtists[0].ID},
		},
	}

	var artistIDs []string
	for _, art := range mockArtists {
		artistIDs = append(artistIDs, art.ID.String())
	}
	filter := fmt.Sprintf(
		"(type = %s OR type = %s) AND artist.id IN [%s]",
		enums.Song,
		enums.Album,
		strings.Join(artistIDs, ","),
	)
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

	var directoryPaths []string
	for _, art := range mockArtists {
		directoryPath := "some_path" + art.ID.String()
		storageFilePathProvider.On("GetArtistDirectoryPath", art).Return(directoryPath).Once()
		directoryPaths = append(directoryPaths, directoryPath)
	}

	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, directoryPaths).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(mockArtists)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	messagePublisherService.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestArtistsDeletedHandler_WhenWithSongsAndOrAlbums_ShouldPublishMessageToDeleteFromSearchEngine(t *testing.T) {
	// given
	messagePublisherService := new(service.MessagePublisherServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := artist.NewArtistsDeletedHandler(messagePublisherService, searchEngineService, storageFilePathProvider)

	mockArtists := []model.Artist{
		{
			ID:   uuid.New(),
			Name: "With Albums",
			Albums: []model.Album{
				{ID: uuid.New()},
				{ID: uuid.New()},
				{ID: uuid.New()},
			},
		},
		{
			ID:   uuid.New(),
			Name: "Empty",
		},
		{
			ID:   uuid.New(),
			Name: "With Songs",
			Songs: []model.Song{
				{ID: uuid.New()},
				{ID: uuid.New()},
			},
		},
		{
			ID:   uuid.New(),
			Name: "With Albums and Songs",
			Albums: []model.Album{
				{ID: uuid.New()},
			},
			Songs: []model.Song{
				{ID: uuid.New()},
			},
		},
	}

	var albumSearchIDs []string
	var songSearchIDs []string
	for _, art := range mockArtists {
		for _, album := range art.Albums {
			albumSearchIDs = append(albumSearchIDs, album.ToSearch().ID)
		}
		for _, song := range art.Songs {
			songSearchIDs = append(songSearchIDs, song.ToSearch().ID)
		}
	}

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Run(func(args mock.Arguments) {
			ids := args.Get(1).([]string)
			assert.Len(t, ids, len(mockArtists)+len(albumSearchIDs)+len(songSearchIDs))

			// assert IDs one by one
			for i, art := range mockArtists {
				assert.Equal(t, art.ToSearch().ID, ids[i])
			}
			for i, albumID := range albumSearchIDs {
				assert.Equal(t, albumID, ids[len(mockArtists)+i])
			}
			for i, songID := range songSearchIDs {
				assert.Equal(t, songID, ids[len(mockArtists)+len(albumSearchIDs)+i])
			}
		}).
		Return(nil).
		Once()

	// Cleanup Storage
	artistDirectoryPath := "some_artist_path"
	albumDirectoryPath := "album_path"
	songDirectoryPath := "song_path"

	for _, art := range mockArtists {
		storageFilePathProvider.On("GetArtistDirectoryPath", art).Return(artistDirectoryPath).Once()
		for _, album := range art.Albums {
			storageFilePathProvider.On("GetAlbumDirectoryPath", album).Return(albumDirectoryPath).Once()
		}
		for _, song := range art.Songs {
			storageFilePathProvider.On("GetSongDirectoryPath", song).Return(songDirectoryPath).Once()
		}
	}

	messagePublisherService.On("Publish", topics.DeleteDirectoriesStorageTopic, mock.IsType([]string{})).
		Run(func(args mock.Arguments) {
			paths := args.Get(1).([]string)
			assert.Len(t, paths, len(mockArtists)+len(albumSearchIDs)+len(songSearchIDs))
			// assert directory paths one by one
			index := 0
			for _, art := range mockArtists {
				assert.Equal(t, artistDirectoryPath, paths[index])
				index++
				for range art.Albums {
					assert.Equal(t, albumDirectoryPath, paths[index])
					index++
				}
				for range art.Songs {
					assert.Equal(t, songDirectoryPath, paths[index])
					index++
				}
			}
		}).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(mockArtists)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	searchEngineService.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}
