package song

import (
	"encoding/json"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"

	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
)

type SongCreatedHandler struct {
	name                    string
	topic                   topics.Topic
	artistRepository        repository.ArtistRepository
	albumRepository         repository.AlbumRepository
	messagePublisherService service.MessagePublisherService
}

func NewSongCreatedHandler(
	artistRepository repository.ArtistRepository,
	albumRepository repository.AlbumRepository,
	messagePublisherService service.MessagePublisherService,
) SongCreatedHandler {
	return SongCreatedHandler{
		name:                    "song_created_handler",
		topic:                   topics.SongCreatedTopic,
		artistRepository:        artistRepository,
		albumRepository:         albumRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (s SongCreatedHandler) Handle(msg *watermillMessage.Message) error {
	var song model.Song
	err := json.Unmarshal(msg.Payload, &song)
	if err != nil {
		return err
	}

	var searches []any
	songSearch := song.ToSearch()

	if song.ArtistID != nil {
		var artist model.Artist
		err = s.artistRepository.Get(&artist, *song.ArtistID)
		if err != nil {
			return err
		}
		songSearch.Artist = artist.ToSongSearch()
	}
	if song.AlbumID != nil {
		var album model.Album
		err = s.albumRepository.Get(&album, *song.AlbumID)
		if err != nil {
			return err
		}
		songSearch.Album = album.ToSongSearch()
	}

	searches = append(searches, songSearch)

	if song.Artist != nil {
		searches = append(searches, song.Artist.ToSearch())
	}
	if song.Album != nil {
		searches = append(searches, song.Album.ToSearch())
	}

	err = s.messagePublisherService.Publish(topics.AddToSearchEngineTopic, searches)
	if err != nil {
		return err
	}
	return nil
}

func (s SongCreatedHandler) GetName() string {
	return s.name
}

func (s SongCreatedHandler) GetTopic() topics.Topic {
	return s.topic
}
