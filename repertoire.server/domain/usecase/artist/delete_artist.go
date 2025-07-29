package artist

import (
	"errors"
	"github.com/google/uuid"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type DeleteArtist struct {
	repository              repository.ArtistRepository
	messagePublisherService service.MessagePublisherService
	transaction             transaction.Manager
}

func NewDeleteArtist(
	repository repository.ArtistRepository,
	messagePublisherService service.MessagePublisherService,
	transaction transaction.Manager,
) DeleteArtist {
	return DeleteArtist{
		repository:              repository,
		messagePublisherService: messagePublisherService,
		transaction:             transaction,
	}
}

func (d DeleteArtist) Handle(request requests.DeleteArtistRequest) *wrapper.ErrorCode {
	var artist model.Artist
	err := d.repository.GetWithSongsOrAlbums(&artist, request.ID, request.WithSongs, request.WithAlbums)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return wrapper.NotFoundError(errors.New("artist not found"))
	}

	err = d.transaction.Execute(func(factory transaction.RepositoryFactory) error {
		artistRepo := factory.NewArtistRepository()

		if request.WithAlbums {
			err = artistRepo.DeleteAlbums([]uuid.UUID{request.ID})
			if err != nil {
				return err
			}
		}
		if request.WithSongs {
			err = artistRepo.DeleteSongs([]uuid.UUID{request.ID})
			if err != nil {
				return err
			}
		}

		err = artistRepo.Delete([]uuid.UUID{request.ID})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = d.messagePublisherService.Publish(topics.ArtistsDeletedTopic, artist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
