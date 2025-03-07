package artist

import (
	"encoding/json"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
)

type ArtistCreatedHandler struct {
	name                    string
	topic                   topics.Topic
	messagePublisherService service.MessagePublisherService
}

func NewArtistCreatedHandler(messagePublisherService service.MessagePublisherService) ArtistCreatedHandler {
	return ArtistCreatedHandler{
		name:                    "artist_created_handler",
		topic:                   topics.ArtistCreatedTopic,
		messagePublisherService: messagePublisherService,
	}
}

func (h ArtistCreatedHandler) Handle(msg *watermillMessage.Message) error {
	var artist model.Artist
	err := json.Unmarshal(msg.Payload, &artist)
	if err != nil {
		return err
	}

	err = h.messagePublisherService.Publish(topics.AddToSearchEngineTopic, []any{artist.ToSearch()})
	if err != nil {
		return err
	}
	return nil
}

func (h ArtistCreatedHandler) GetName() string {
	return h.name
}

func (h ArtistCreatedHandler) GetTopic() topics.Topic {
	return h.topic
}
