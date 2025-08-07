package song

import (
	"encoding/json"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"

	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
)

type SongDeletedHandler struct {
	name                    string
	topic                   topics.Topic
	storageFilePathProvider provider.StorageFilePathProvider
	messagePublisherService service.MessagePublisherService
}

func NewSongDeletedHandler(
	storageFilePathProvider provider.StorageFilePathProvider,
	messagePublisherService service.MessagePublisherService,
) SongDeletedHandler {
	return SongDeletedHandler{
		name:                    "song_deleted_handler",
		topic:                   topics.SongDeletedTopic,
		storageFilePathProvider: storageFilePathProvider,
		messagePublisherService: messagePublisherService,
	}
}

func (s SongDeletedHandler) Handle(msg *watermillMessage.Message) error {
	var song model.Song
	err := json.Unmarshal(msg.Payload, &song)
	if err != nil {
		return err
	}

	err = s.messagePublisherService.Publish(topics.DeleteFromSearchEngineTopic, []string{song.ToSearch().ID})
	if err != nil {
		return err
	}

	directoryPath := s.storageFilePathProvider.GetSongDirectoryPath(song)
	return s.messagePublisherService.Publish(topics.DeleteDirectoriesStorageTopic, []string{directoryPath})
}

func (s SongDeletedHandler) GetName() string {
	return s.name
}

func (s SongDeletedHandler) GetTopic() topics.Topic {
	return s.topic
}
