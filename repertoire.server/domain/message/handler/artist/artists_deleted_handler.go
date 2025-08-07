package artist

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

type ArtistsDeletedHandler struct {
	name                    string
	topic                   topics.Topic
	messagePublisherService service.MessagePublisherService
	searchEngineService     service.SearchEngineService
	storageFilePathProvider provider.StorageFilePathProvider
}

func NewArtistsDeletedHandler(
	messagePublisherService service.MessagePublisherService,
	searchEngineService service.SearchEngineService,
	storageFilePathProvider provider.StorageFilePathProvider,
) ArtistsDeletedHandler {
	return ArtistsDeletedHandler{
		name:                    "artist_deleted_handler",
		topic:                   topics.ArtistsDeletedTopic,
		messagePublisherService: messagePublisherService,
		searchEngineService:     searchEngineService,
		storageFilePathProvider: storageFilePathProvider,
	}
}

func (a ArtistsDeletedHandler) Handle(msg *watermillMessage.Message) error {
	var artists []model.Artist
	err := json.Unmarshal(msg.Payload, &artists)
	if err != nil {
		return err
	}

	err = a.syncWithSearchEngine(artists)
	if err != nil {
		return err
	}

	return a.cleanupStorage(artists)
}

func (a ArtistsDeletedHandler) syncWithSearchEngine(artists []model.Artist) error {
	// previously in delete artist, the artist was populated with songs and albums,
	// but only if they have to be deleted too
	var artistIDs []string
	var artistSearchIDs []string
	var albumIDs []string
	var songIDs []string
	for _, art := range artists {
		artistIDs = append(artistIDs, art.ID.String())
		artistSearchIDs = append(artistSearchIDs, art.ToSearch().ID)
		for _, song := range art.Songs {
			songIDs = append(songIDs, song.ToSearch().ID)
		}
		for _, album := range art.Albums {
			albumIDs = append(albumIDs, album.ToSearch().ID)
		}
	}

	ids := slices.Concat(artistSearchIDs, albumIDs, songIDs)
	err := a.messagePublisherService.Publish(topics.DeleteFromSearchEngineTopic, ids)
	if err != nil {
		return err
	}

	// if the artist already has songs and albums populated, there is no need to update them, as they have been deleted
	if len(songIDs) > 0 && len(albumIDs) > 0 {
		return nil
	}

	// get the songs and albums based on the artist and delete their artist
	filter := fmt.Sprintf(
		"(type = %s OR type = %s) AND artist.id IN [%s]",
		enums.Song,
		enums.Album,
		strings.Join(artistIDs, ","),
	)
	searches, err := a.searchEngineService.GetDocuments(filter)
	if err != nil {
		return err
	}
	if len(searches) == 0 {
		return nil
	}

	var documentsToUpdate []any
	for _, search := range searches {
		search["artist"] = nil
		documentsToUpdate = append(documentsToUpdate, search)
	}

	return a.messagePublisherService.Publish(topics.UpdateFromSearchEngineTopic, documentsToUpdate)
}

func (a ArtistsDeletedHandler) cleanupStorage(artists []model.Artist) error {
	var directoryPaths []string

	for _, art := range artists {
		directoryPaths = append(directoryPaths, a.storageFilePathProvider.GetArtistDirectoryPath(art))
		// previously in delete artist, the artist was populated with songs and albums,
		// only if they have to be deleted too
		for _, album := range art.Albums {
			directoryPaths = append(directoryPaths, a.storageFilePathProvider.GetAlbumDirectoryPath(album))
		}
		for _, song := range art.Songs {
			directoryPaths = append(directoryPaths, a.storageFilePathProvider.GetSongDirectoryPath(song))
		}
	}

	return a.messagePublisherService.Publish(topics.DeleteDirectoriesStorageTopic, directoryPaths)
}

func (a ArtistsDeletedHandler) GetName() string {
	return a.name
}

func (a ArtistsDeletedHandler) GetTopic() topics.Topic {
	return a.topic
}
