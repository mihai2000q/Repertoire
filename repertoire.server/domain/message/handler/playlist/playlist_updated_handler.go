package playlist

import (
	"encoding/json"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
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

func (s PlaylistUpdatedHandler) Handle(msg *watermillMessage.Message) error {
	var playlist model.Playlist
	err := json.Unmarshal(msg.Payload, &playlist)
	if err != nil {
		return err
	}

	err = s.messagePublisherService.Publish(topics.UpdateFromSearchEngineTopic, []any{playlist.ToSearch()})
	if err != nil {
		return err
	}
	return nil
}

func (s PlaylistUpdatedHandler) GetName() string {
	return s.name
}

func (s PlaylistUpdatedHandler) GetTopic() topics.Topic {
	return s.topic
}
