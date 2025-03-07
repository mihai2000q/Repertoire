package song

import (
	"encoding/json"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
)

type SongDeletedHandler struct {
	name                    string
	topic                   topics.Topic
	messagePublisherService service.MessagePublisherService
}

func NewSongDeletedHandler(messagePublisherService service.MessagePublisherService) SongDeletedHandler {
	return SongDeletedHandler{
		name:                    "song_deleted_handler",
		topic:                   topics.SongDeletedTopic,
		messagePublisherService: messagePublisherService,
	}
}

func (h SongDeletedHandler) Handle(msg *watermillMessage.Message) error {
	var song model.Song
	err := json.Unmarshal(msg.Payload, &song)
	if err != nil {
		return err
	}

	err = h.messagePublisherService.Publish(topics.DeleteFromSearchEngineTopic, []string{song.ToSearch().ID})
	if err != nil {
		return err
	}
	return nil
}

func (h SongDeletedHandler) GetName() string {
	return h.name
}

func (h SongDeletedHandler) GetTopic() topics.Topic {
	return h.topic
}
