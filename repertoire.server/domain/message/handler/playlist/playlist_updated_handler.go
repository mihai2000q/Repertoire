package playlist

import (
	"encoding/json"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"

	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
)

type PlaylistUpdatedHandler struct {
	name                    string
	topic                   topics.Topic
	messagePublisherService service.MessagePublisherService
}

func NewPlaylistUpdatedHandler(messagePublisherService service.MessagePublisherService) PlaylistUpdatedHandler {
	return PlaylistUpdatedHandler{
		name:                    "playlist_updated_handler",
		topic:                   topics.PlaylistUpdatedTopic,
		messagePublisherService: messagePublisherService,
	}
}

func (p PlaylistUpdatedHandler) Handle(msg *watermillMessage.Message) error {
	var playlist model.Playlist
	err := json.Unmarshal(msg.Payload, &playlist)
	if err != nil {
		return err
	}

	err = p.messagePublisherService.Publish(topics.UpdateFromSearchEngineTopic, []any{playlist.ToSearch()})
	if err != nil {
		return err
	}
	return nil
}

func (p PlaylistUpdatedHandler) GetName() string {
	return p.name
}

func (p PlaylistUpdatedHandler) GetTopic() topics.Topic {
	return p.topic
}
