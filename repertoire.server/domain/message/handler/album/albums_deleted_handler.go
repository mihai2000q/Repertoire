package album

import (
	"encoding/json"
	"fmt"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"slices"
	"strings"

	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
)

type AlbumsDeletedHandler struct {
	name                    string
	topic                   topics.Topic
	messagePublisherService service.MessagePublisherService
	searchEngineService     service.SearchEngineService
	storageFilePathProvider provider.StorageFilePathProvider
}

func NewAlbumsDeletedHandler(
	messagePublisherService service.MessagePublisherService,
	searchEngineService service.SearchEngineService,
	storageFilePathProvider provider.StorageFilePathProvider,
) AlbumsDeletedHandler {
	return AlbumsDeletedHandler{
		name:                    "albums_deleted_handler",
		topic:                   topics.AlbumsDeletedTopic,
		messagePublisherService: messagePublisherService,
		searchEngineService:     searchEngineService,
		storageFilePathProvider: storageFilePathProvider,
	}
}

func (a AlbumsDeletedHandler) Handle(msg *watermillMessage.Message) error {
	var albums []model.Album
	err := json.Unmarshal(msg.Payload, &albums)
	if err != nil {
		return err
	}

	err = a.syncWithSearchEngine(albums)
	if err != nil {
		return err
	}

	return a.cleanupStorage(albums)
}

func (a AlbumsDeletedHandler) syncWithSearchEngine(albums []model.Album) error {
	// previously in delete album, the album was populated with songs, only if they have to be deleted too
	var albumIDs []string
	var albumSearchIDs []string
	var songSearchIDs []string
	for _, album := range albums {
		albumIDs = append(albumIDs, album.ID.String())
		albumSearchIDs = append(albumSearchIDs, album.ToSearch().ID)
		for _, song := range album.Songs {
			songSearchIDs = append(songSearchIDs, song.ToSearch().ID)
		}
	}

	ids := slices.Concat(albumSearchIDs, songSearchIDs)
	err := a.messagePublisherService.Publish(topics.DeleteFromSearchEngineTopic, ids)
	if err != nil {
		return err
	}

	// if the album already has songs populated, there is no need to update the songs, as they will be deleted
	if len(songSearchIDs) > 0 {
		return nil
	}

	// get the songs based on the album and delete their album (nullify)
	filter := fmt.Sprintf("type = %s AND album.id IN [%s]", enums.Song, strings.Join(albumIDs, ", "))
	meiliSongs, err := a.searchEngineService.GetDocuments(filter)
	if err != nil {
		return err
	}
	if len(meiliSongs) == 0 {
		return nil
	}

	var songsToUpdate []any
	for _, s := range meiliSongs {
		var song model.SongSearch
		jsonSong, _ := json.Marshal(s)
		_ = json.Unmarshal(jsonSong, &song)
		song.Album = nil
		songsToUpdate = append(songsToUpdate, song)
	}

	return a.messagePublisherService.Publish(topics.UpdateFromSearchEngineTopic, songsToUpdate)
}

func (a AlbumsDeletedHandler) cleanupStorage(albums []model.Album) error {
	var directoryPaths []string

	for _, album := range albums {
		directoryPaths = append(directoryPaths, a.storageFilePathProvider.GetAlbumDirectoryPath(album))
		// previously in delete album, the album was populated with songs, only if they have to be deleted too
		for _, song := range album.Songs {
			directoryPaths = append(directoryPaths, a.storageFilePathProvider.GetSongDirectoryPath(song))
		}
	}

	return a.messagePublisherService.Publish(topics.DeleteDirectoriesStorageTopic, directoryPaths)
}

func (a AlbumsDeletedHandler) GetName() string {
	return a.name
}

func (a AlbumsDeletedHandler) GetTopic() topics.Topic {
	return a.topic
}
