package album

import (
	"encoding/json"
	"fmt"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
)

type AlbumDeletedHandler struct {
	name                    string
	topic                   topics.Topic
	messagePublisherService service.MessagePublisherService
	searchEngineService     service.SearchEngineService
	storageFilePathProvider provider.StorageFilePathProvider
}

func NewAlbumDeletedHandler(
	messagePublisherService service.MessagePublisherService,
	searchEngineService service.SearchEngineService,
	storageFilePathProvider provider.StorageFilePathProvider,
) AlbumDeletedHandler {
	return AlbumDeletedHandler{
		name:                    "album_deleted_handler",
		topic:                   topics.AlbumDeletedTopic,
		messagePublisherService: messagePublisherService,
		searchEngineService:     searchEngineService,
		storageFilePathProvider: storageFilePathProvider,
	}
}

func (a AlbumDeletedHandler) Handle(msg *watermillMessage.Message) error {
	var album model.Album
	err := json.Unmarshal(msg.Payload, &album)
	if err != nil {
		return err
	}

	err = a.syncWithSearchEngine(album)
	if err != nil {
		return err
	}

	return a.cleanupStorage(album)
}

func (a AlbumDeletedHandler) syncWithSearchEngine(album model.Album) error {
	// previously in delete album, the album was populated with songs, only if they have to be deleted too
	ids := []string{album.ToSearch().ID}
	for _, song := range album.Songs {
		ids = append(ids, song.ToSearch().ID)
	}

	err := a.messagePublisherService.Publish(topics.DeleteFromSearchEngineTopic, ids)
	if err != nil {
		return err
	}

	// if the album already has songs populated, there is no need to update the songs, as they will be deleted
	if len(album.Songs) > 0 {
		return nil
	}

	// get the songs based on the album and delete their album
	filter := fmt.Sprintf("type = %s AND album.id = %s", enums.Song, album.ID)
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

func (a AlbumDeletedHandler) cleanupStorage(album model.Album) error {
	var directoryPaths []string
	albumDirectoryPath := a.storageFilePathProvider.GetAlbumDirectoryPath(album)
	directoryPaths = append(directoryPaths, albumDirectoryPath)

	// previously in delete album, the album was populated with songs, only if they have to be deleted too
	for _, song := range album.Songs {
		directoryPaths = append(directoryPaths, a.storageFilePathProvider.GetSongDirectoryPath(song))
	}

	return a.messagePublisherService.Publish(topics.DeleteDirectoriesStorageTopic, directoryPaths)
}

func (a AlbumDeletedHandler) GetName() string {
	return a.name
}

func (a AlbumDeletedHandler) GetTopic() topics.Topic {
	return a.topic
}
