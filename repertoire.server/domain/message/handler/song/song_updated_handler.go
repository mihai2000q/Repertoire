package song

import (
	"encoding/json"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
)

type SongUpdatedHandler struct {
	name                    string
	topic                   topics.Topic
	songRepository          repository.SongRepository
	messagePublisherService service.MessagePublisherService
}

func NewSongUpdatedHandler(
	songRepository repository.SongRepository,
	messagePublisherService service.MessagePublisherService,
) SongUpdatedHandler {
	return SongUpdatedHandler{
		name:                    "song_updated_handler",
		topic:                   topics.SongUpdatedTopic,
		songRepository:          songRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (s SongUpdatedHandler) Handle(msg *watermillMessage.Message) error {
	var song model.Song
	err := json.Unmarshal(msg.Payload, &song)
	if err != nil {
		return err
	}

	// making sure that the song artist and album are present for mapping
	err = s.songRepository.GetWithArtistAndAlbum(&song, song.ID)
	if err != nil {
		return err
	}

	err = s.messagePublisherService.Publish(topics.UpdateFromSearchEngineTopic, []any{song.ToSearch()})
	if err != nil {
		return err
	}
	return nil
}

func (s SongUpdatedHandler) GetName() string {
	return s.name
}

func (s SongUpdatedHandler) GetTopic() topics.Topic {
	return s.topic
}
