package playlist

import (
	"encoding/json"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"

	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
)

type PlaylistCreatedHandler struct {
	name                    string
	topic                   topics.Topic
	messagePublisherService service.MessagePublisherService
}

func NewPlaylistCreatedHandler(messagePublisherService service.MessagePublisherService) PlaylistCreatedHandler {
	return PlaylistCreatedHandler{
		name:                    "playlist_created_handler",
		topic:                   topics.PlaylistCreatedTopic,
		messagePublisherService: messagePublisherService,
	}
}

func (p PlaylistCreatedHandler) Handle(msg *watermillMessage.Message) error {
	var playlist model.Playlist
	err := json.Unmarshal(msg.Payload, &playlist)
	if err != nil {
		return err
	}

	err = p.messagePublisherService.Publish(topics.AddToSearchEngineTopic, []any{playlist.ToSearch()})
	if err != nil {
		return err
	}
	return nil
}

func (p PlaylistCreatedHandler) GetName() string {
	return p.name
}

func (p PlaylistCreatedHandler) GetTopic() topics.Topic {
	return p.topic
}
