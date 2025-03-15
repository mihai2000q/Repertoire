package artist

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type DeleteArtist struct {
	repository              repository.ArtistRepository
	messagePublisherService service.MessagePublisherService
}

func NewDeleteArtist(
	repository repository.ArtistRepository,
	messagePublisherService service.MessagePublisherService,
) DeleteArtist {
	return DeleteArtist{
		repository:              repository,
		messagePublisherService: messagePublisherService,
	}
}

func (d DeleteArtist) Handle(request requests.DeleteArtistRequest) *wrapper.ErrorCode {
	var artist model.Artist
	var err error
	if request.WithAlbums && request.WithSongs {
		err = d.repository.GetWithAlbumsAndSongs(&artist, request.ID)
	} else if request.WithSongs {
		err = d.repository.GetWithSongs(&artist, request.ID)
	} else if request.WithAlbums {
		err = d.repository.GetWithAlbums(&artist, request.ID)
	} else {
		err = d.repository.Get(&artist, request.ID)
	}
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return wrapper.NotFoundError(errors.New("artist not found"))
	}
	
	if request.WithAlbums {
		err = d.repository.DeleteAlbums(request.ID)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
	}
	if request.WithSongs {
		err = d.repository.DeleteSongs(request.ID)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
	}

	err = d.repository.Delete(request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = d.messagePublisherService.Publish(topics.ArtistDeletedTopic, artist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
