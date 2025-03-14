package playlist

import (
	"encoding/json"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
)

type PlaylistDeletedHandler struct {
	name                    string
	topic                   topics.Topic
	messagePublisherService service.MessagePublisherService
}

func NewPlaylistDeletedHandler(messagePublisherService service.MessagePublisherService) PlaylistDeletedHandler {
	return PlaylistDeletedHandler{
		name:                    "playlist_deleted_handler",
		topic:                   topics.PlaylistDeletedTopic,
		messagePublisherService: messagePublisherService,
	}
}

func (s PlaylistDeletedHandler) Handle(msg *watermillMessage.Message) error {
	var playlist model.Playlist
	err := json.Unmarshal(msg.Payload, &playlist)
	if err != nil {
		return err
	}

	err = s.messagePublisherService.Publish(topics.DeleteFromSearchEngineTopic, []string{playlist.ToSearch().ID})
	if err != nil {
		return err
	}
	return nil
}

func (s PlaylistDeletedHandler) GetName() string {
	return s.name
}

func (s PlaylistDeletedHandler) GetTopic() topics.Topic {
	return s.topic
}
