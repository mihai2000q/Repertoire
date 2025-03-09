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

func (a AlbumDeletedHandler) Handle(msg *watermillMessage.Message) error {
	var album model.Album
	err := json.Unmarshal(msg.Payload, &album)
	if err != nil {
		return err
	}

	ids := []string{album.ToSearch().ID}
	for _, song := range album.Songs {
		ids = append(ids, song.ToSearch().ID)
	}

	err = a.messagePublisherService.Publish(topics.DeleteFromSearchEngineTopic, ids)
	if err != nil {
		return err
	}
	// TODO: UPDATE SONGS WITH ALBUM NULL ON MEILI
	return nil
}

func (a AlbumDeletedHandler) GetName() string {
	return a.name
}

func (a AlbumDeletedHandler) GetTopic() topics.Topic {
	return a.topic
}
