package song

import (
	"encoding/json"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"

	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

type SongsUpdatedHandler struct {
	name                    string
	topic                   topics.Topic
	songRepository          repository.SongRepository
	messagePublisherService service.MessagePublisherService
}

func NewSongsUpdatedHandler(
	songsRepository repository.SongRepository,
	messagePublisherService service.MessagePublisherService,
) SongsUpdatedHandler {
	return SongsUpdatedHandler{
		name:                    "songs_updated_handler",
		topic:                   topics.SongsUpdatedTopic,
		songRepository:          songsRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (s SongsUpdatedHandler) Handle(msg *watermillMessage.Message) error {
	var ids []uuid.UUID
	err := json.Unmarshal(msg.Payload, &ids)
	if err != nil {
		return err
	}

	var songs []model.Song
	err = s.songRepository.GetAllByIDsWithArtistAndAlbum(&songs, ids)
	if err != nil {
		return err
	}
	if len(songs) == 0 {
		return nil
	}

	var songSearches []any
	for _, song := range songs {
		songSearches = append(songSearches, song.ToSearch())
	}

	err = s.messagePublisherService.Publish(topics.UpdateFromSearchEngineTopic, songSearches)
	if err != nil {
		return err
	}
	return nil
}

func (s SongsUpdatedHandler) GetName() string {
	return s.name
}

func (s SongsUpdatedHandler) GetTopic() topics.Topic {
	return s.topic
}
