package playlist

import (
	"encoding/json"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
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

func (a PlaylistCreatedHandler) Handle(msg *watermillMessage.Message) error {
	var playlist model.Playlist
	err := json.Unmarshal(msg.Payload, &playlist)
	if err != nil {
		return err
	}

	err = a.messagePublisherService.Publish(topics.AddToSearchEngineTopic, []any{playlist.ToSearch()})
	if err != nil {
		return err
	}
	return nil
}

func (a PlaylistCreatedHandler) GetName() string {
	return a.name
}

func (a PlaylistCreatedHandler) GetTopic() topics.Topic {
	return a.topic
}
