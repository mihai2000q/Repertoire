package album

import (
	"encoding/json"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
)

type AlbumCreatedHandler struct {
	name                    string
	topic                   topics.Topic
	artistRepository        repository.ArtistRepository
	messagePublisherService service.MessagePublisherService
}

func NewAlbumCreatedHandler(
	artistRepository repository.ArtistRepository,
	messagePublisherService service.MessagePublisherService,
) AlbumCreatedHandler {
	return AlbumCreatedHandler{
		name:                    "album_created_handler",
		topic:                   topics.AlbumCreatedTopic,
		artistRepository:        artistRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (h AlbumCreatedHandler) Handle(msg *watermillMessage.Message) error {
	var album model.Album
	err := json.Unmarshal(msg.Payload, &album)
	if err != nil {
		return err
	}

	var searches []any
	albumSearch := album.ToSearch()

	if album.ArtistID != nil {
		var artist model.Artist
		err = h.artistRepository.Get(&artist, *album.ArtistID)
		if err != nil {
			return err
		}
		albumSearch.Artist = artist.ToAlbumSearch()
	}
	searches = append(searches, albumSearch)
	if album.Artist != nil {
		searches = append(searches, album.Artist.ToSearch())
	}

	err = h.messagePublisherService.Publish(topics.AddToSearchEngineTopic, searches)
	if err != nil {
		return err
	}
	return nil
}

func (h AlbumCreatedHandler) GetName() string {
	return h.name
}

func (h AlbumCreatedHandler) GetTopic() topics.Topic {
	return h.topic
}
