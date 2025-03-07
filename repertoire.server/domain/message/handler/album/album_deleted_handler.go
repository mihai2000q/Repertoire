package album

import (
	"encoding/json"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
)

type AlbumDeletedHandler struct {
	name                    string
	topic                   topics.Topic
	messagePublisherService service.MessagePublisherService
}

func NewAlbumDeletedHandler(messagePublisherService service.MessagePublisherService) AlbumDeletedHandler {
	return AlbumDeletedHandler{
		name:                    "album_deleted_handler",
		topic:                   topics.AlbumDeletedTopic,
		messagePublisherService: messagePublisherService,
	}
}

func (h AlbumDeletedHandler) Handle(msg *watermillMessage.Message) error {
	var album model.Album
	err := json.Unmarshal(msg.Payload, &album)
	if err != nil {
		return err
	}

	ids := []string{album.ToSearch().ID}
	for _, song := range album.Songs {
		ids = append(ids, song.ToSearch().ID)
	}

	err = h.messagePublisherService.Publish(topics.DeleteFromSearchEngineTopic, ids)
	if err != nil {
		return err
	}
	// TODO: UPDATE SONGS WITH ALBUM NULL ON MEILI
	return nil
}

func (h AlbumDeletedHandler) GetName() string {
	return h.name
}

func (h AlbumDeletedHandler) GetTopic() topics.Topic {
	return h.topic
}
