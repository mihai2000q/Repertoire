package artist

import (
	"encoding/json"
	"fmt"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"repertoire/server/data/service"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
)

type ArtistDeletedHandler struct {
	name                    string
	topic                   topics.Topic
	messagePublisherService service.MessagePublisherService
	searchEngineService     service.SearchEngineService
}

func NewArtistDeletedHandler(
	messagePublisherService service.MessagePublisherService,
	searchEngineService service.SearchEngineService,
) ArtistDeletedHandler {
	return ArtistDeletedHandler{
		name:                    "artist_deleted_handler",
		topic:                   topics.ArtistDeletedTopic,
		messagePublisherService: messagePublisherService,
		searchEngineService:     searchEngineService,
	}
}

func (a ArtistDeletedHandler) Handle(msg *watermillMessage.Message) error {
	var artist model.Artist
	err := json.Unmarshal(msg.Payload, &artist)
	if err != nil {
		return err
	}

	// previously in delete artist, the artist was populated with songs and albums,
	// but only if they have to be deleted too
	ids := []string{artist.ToSearch().ID}
	for _, song := range artist.Songs {
		ids = append(ids, song.ToSearch().ID)
	}
	for _, album := range artist.Albums {
		ids = append(ids, album.ToSearch().ID)
	}

	err = a.messagePublisherService.Publish(topics.DeleteFromSearchEngineTopic, ids)
	if err != nil {
		return err
	}

	// if the artist already has songs and albums populated, there is no need to update them, as they have been deleted
	if len(artist.Songs) > 0 && len(artist.Albums) > 0 {
		return nil
	}

	// get the songs and albums based on the artist and delete their artist
	filter := fmt.Sprintf("(type = %s OR type = %s) AND artist.id = %s", enums.Song, enums.Album, artist.ID)
	searches, err := a.searchEngineService.GetDocuments(filter)
	if err != nil {
		return err
	}
	if len(searches) == 0 {
		return nil
	}

	var documentsToUpdate []any
	for _, search := range searches {
		search["artist"] = nil
		documentsToUpdate = append(documentsToUpdate, search)
	}

	err = a.messagePublisherService.Publish(topics.UpdateFromSearchEngineTopic, documentsToUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (a ArtistDeletedHandler) GetName() string {
	return a.name
}

func (a ArtistDeletedHandler) GetTopic() topics.Topic {
	return a.topic
}
