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

type AlbumsUpdatedHandler struct {
	name                    string
	topic                   topics.Topic
	albumRepository         repository.AlbumRepository
	messagePublisherService service.MessagePublisherService
}

func NewAlbumsUpdatedHandler(
	albumRepository repository.AlbumRepository,
	messagePublisherService service.MessagePublisherService,
) AlbumsUpdatedHandler {
	return AlbumsUpdatedHandler{
		name:                    "albums_updated_handler",
		topic:                   topics.AlbumsUpdatedTopic,
		albumRepository:         albumRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (a AlbumsUpdatedHandler) Handle(msg *watermillMessage.Message) error {
	var ids []uuid.UUID
	err := json.Unmarshal(msg.Payload, &ids)
	if err != nil {
		return err
	}

	var albums []model.Album
	err = a.albumRepository.GetAllByIDsWithSongsAndArtist(&albums, ids)
	if err != nil {
		return err
	}

	var documentsToUpdate []any
	for _, album := range albums {
		documentsToUpdate = append(documentsToUpdate, album.ToSearch())
		for _, song := range album.Songs {
			songSearch := song.ToSearch()
			songSearch.Album = album.ToSongSearch()
			documentsToUpdate = append(documentsToUpdate, songSearch)
		}
	}

	err = a.messagePublisherService.Publish(topics.UpdateFromSearchEngineTopic, documentsToUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (a AlbumsUpdatedHandler) GetName() string {
	return a.name
}

func (a AlbumsUpdatedHandler) GetTopic() topics.Topic {
	return a.topic
}
