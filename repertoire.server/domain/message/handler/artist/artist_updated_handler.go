package artist

import (
	"encoding/json"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
)

type ArtistUpdatedHandler struct {
	name                    string
	topic                   topics.Topic
	artistRepository        repository.ArtistRepository
	messagePublisherService service.MessagePublisherService
}

func NewArtistUpdatedHandler(
	artistRepository repository.ArtistRepository,
	messagePublisherService service.MessagePublisherService,
) ArtistUpdatedHandler {
	return ArtistUpdatedHandler{
		name:                    "artist_updated_handler",
		topic:                   topics.AlbumUpdatedTopic,
		artistRepository:        artistRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (a ArtistUpdatedHandler) Handle(msg *watermillMessage.Message) error {
	var artistID uuid.UUID
	err := json.Unmarshal(msg.Payload, &artistID)
	if err != nil {
		return err
	}

	var artist model.Artist
	err = a.artistRepository.GetWithAlbumsAndSongs(&artist, artistID)
	if err != nil {
		return err
	}

	documentsToUpdate := []any{artist.ToSearch()}
	for _, song := range artist.Songs {
		songSearch := song.ToSearch()
		songSearch.Artist = artist.ToSongSearch()
		documentsToUpdate = append(documentsToUpdate, songSearch)
	}
	for _, album := range artist.Albums {
		albumSearch := album.ToSearch()
		albumSearch.Artist = artist.ToAlbumSearch()
		documentsToUpdate = append(documentsToUpdate, albumSearch)
	}

	err = a.messagePublisherService.Publish(topics.UpdateFromSearchEngineTopic, documentsToUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (a ArtistUpdatedHandler) GetName() string {
	return a.name
}

func (a ArtistUpdatedHandler) GetTopic() topics.Topic {
	return a.topic
}
