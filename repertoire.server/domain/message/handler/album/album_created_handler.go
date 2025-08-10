package album

import (
	"encoding/json"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"

	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
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

func (a AlbumCreatedHandler) Handle(msg *watermillMessage.Message) error {
	var album model.Album
	err := json.Unmarshal(msg.Payload, &album)
	if err != nil {
		return err
	}

	var searches []any
	albumSearch := album.ToSearch()

	if album.ArtistID != nil {
		var artist model.Artist
		err = a.artistRepository.Get(&artist, *album.ArtistID)
		if err != nil {
			return err
		}
		albumSearch.Artist = artist.ToAlbumSearch()
	}
	searches = append(searches, albumSearch)
	if album.Artist != nil {
		searches = append(searches, album.Artist.ToSearch())
	}

	err = a.messagePublisherService.Publish(topics.AddToSearchEngineTopic, searches)
	if err != nil {
		return err
	}
	return nil
}

func (a AlbumCreatedHandler) GetName() string {
	return a.name
}

func (a AlbumCreatedHandler) GetTopic() topics.Topic {
	return a.topic
}
