package playlist

import (
	"encoding/json"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
)

type PlaylistDeletedHandler struct {
	name                    string
	topic                   topics.Topic
	storageFilePathProvider provider.StorageFilePathProvider
	messagePublisherService service.MessagePublisherService
}

func NewPlaylistDeletedHandler(
	storageFilePathProvider provider.StorageFilePathProvider,
	messagePublisherService service.MessagePublisherService,
) PlaylistDeletedHandler {
	return PlaylistDeletedHandler{
		name:                    "playlist_deleted_handler",
		topic:                   topics.PlaylistDeletedTopic,
		storageFilePathProvider: storageFilePathProvider,
		messagePublisherService: messagePublisherService,
	}
}

func (p PlaylistDeletedHandler) Handle(msg *watermillMessage.Message) error {
	var playlist model.Playlist
	err := json.Unmarshal(msg.Payload, &playlist)
	if err != nil {
		return err
	}

	err = p.messagePublisherService.Publish(topics.DeleteFromSearchEngineTopic, []string{playlist.ToSearch().ID})
	if err != nil {
		return err
	}

	directoryPath := p.storageFilePathProvider.GetPlaylistDirectoryPath(playlist)
	return p.messagePublisherService.Publish(topics.DeleteDirectoriesStorageTopic, []string{directoryPath})
}

func (p PlaylistDeletedHandler) GetName() string {
	return p.name
}

func (p PlaylistDeletedHandler) GetTopic() topics.Topic {
	return p.topic
}
