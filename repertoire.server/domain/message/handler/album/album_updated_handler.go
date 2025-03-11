package album

import (
	"encoding/json"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
)

type AlbumUpdatedHandler struct {
	name                    string
	topic                   topics.Topic
	albumRepository         repository.AlbumRepository
	messagePublisherService service.MessagePublisherService
}

func NewAlbumUpdatedHandler(
	albumRepository repository.AlbumRepository,
	messagePublisherService service.MessagePublisherService,
) AlbumUpdatedHandler {
	return AlbumUpdatedHandler{
		name:                    "album_updated_handler",
		topic:                   topics.AlbumUpdatedTopic,
		albumRepository:         albumRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (a AlbumUpdatedHandler) Handle(msg *watermillMessage.Message) error {
	var albumID uuid.UUID
	err := json.Unmarshal(msg.Payload, &albumID)
	if err != nil {
		return err
	}

	var album model.Album
	err = a.albumRepository.GetWithSongsAndArtist(&album, albumID)
	if err != nil {
		return err
	}

	documentsToUpdate := []any{album.ToSearch()}
	for _, song := range album.Songs {
		songSearch := song.ToSearch()
		songSearch.Album = album.ToSongSearch()
		documentsToUpdate = append(documentsToUpdate, songSearch)
	}

	err = a.messagePublisherService.Publish(topics.UpdateFromSearchEngineTopic, documentsToUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (a AlbumUpdatedHandler) GetName() string {
	return a.name
}

func (a AlbumUpdatedHandler) GetTopic() topics.Topic {
	return a.topic
}
