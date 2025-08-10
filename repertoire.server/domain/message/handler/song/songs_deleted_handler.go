package song

import (
	"encoding/json"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"

	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
)

type SongsDeletedHandler struct {
	name                    string
	topic                   topics.Topic
	storageFilePathProvider provider.StorageFilePathProvider
	messagePublisherService service.MessagePublisherService
}

func NewSongsDeletedHandler(
	storageFilePathProvider provider.StorageFilePathProvider,
	messagePublisherService service.MessagePublisherService,
) SongsDeletedHandler {
	return SongsDeletedHandler{
		name:                    "songs_deleted_handler",
		topic:                   topics.SongsDeletedTopic,
		storageFilePathProvider: storageFilePathProvider,
		messagePublisherService: messagePublisherService,
	}
}

func (s SongsDeletedHandler) Handle(msg *watermillMessage.Message) error {
	var songs []model.Song
	err := json.Unmarshal(msg.Payload, &songs)
	if err != nil {
		return err
	}

	var ids []string
	for _, song := range songs {
		ids = append(ids, song.ToSearch().ID)
	}
	err = s.messagePublisherService.Publish(topics.DeleteFromSearchEngineTopic, ids)
	if err != nil {
		return err
	}

	var directoryPaths []string
	for _, song := range songs {
		directoryPath := s.storageFilePathProvider.GetSongDirectoryPath(song)
		directoryPaths = append(directoryPaths, directoryPath)
	}
	return s.messagePublisherService.Publish(topics.DeleteDirectoriesStorageTopic, directoryPaths)
}

func (s SongsDeletedHandler) GetName() string {
	return s.name
}

func (s SongsDeletedHandler) GetTopic() topics.Topic {
	return s.topic
}
