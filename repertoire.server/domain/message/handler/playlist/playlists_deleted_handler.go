package playlist

import (
	"encoding/json"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"

	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
)

type PlaylistsDeletedHandler struct {
	name                    string
	topic                   topics.Topic
	storageFilePathProvider provider.StorageFilePathProvider
	messagePublisherService service.MessagePublisherService
}

func NewPlaylistsDeletedHandler(
	storageFilePathProvider provider.StorageFilePathProvider,
	messagePublisherService service.MessagePublisherService,
) PlaylistsDeletedHandler {
	return PlaylistsDeletedHandler{
		name:                    "playlists_deleted_handler",
		topic:                   topics.PlaylistsDeletedTopic,
		storageFilePathProvider: storageFilePathProvider,
		messagePublisherService: messagePublisherService,
	}
}

func (p PlaylistsDeletedHandler) Handle(msg *watermillMessage.Message) error {
	var playlists []model.Playlist
	err := json.Unmarshal(msg.Payload, &playlists)
	if err != nil {
		return err
	}

	var playlistIDs []string
	for _, playlist := range playlists {
		playlistIDs = append(playlistIDs, playlist.ToSearch().ID)
	}
	err = p.messagePublisherService.Publish(topics.DeleteFromSearchEngineTopic, playlistIDs)
	if err != nil {
		return err
	}

	var directoryPaths []string
	for _, playlist := range playlists {
		directoryPath := p.storageFilePathProvider.GetPlaylistDirectoryPath(playlist)
		directoryPaths = append(directoryPaths, directoryPath)
	}
	return p.messagePublisherService.Publish(topics.DeleteDirectoriesStorageTopic, directoryPaths)
}

func (p PlaylistsDeletedHandler) GetName() string {
	return p.name
}

func (p PlaylistsDeletedHandler) GetTopic() topics.Topic {
	return p.topic
}
