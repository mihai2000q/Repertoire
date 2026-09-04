package artist

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type DeleteArtist struct {
	artistRepository        repository.ArtistRepository
	messagePublisherService service.MessagePublisherService
	transaction             transaction.Manager
}

func NewDeleteArtist(
	artistRepository repository.ArtistRepository,
	messagePublisherService service.MessagePublisherService,
	transaction transaction.Manager,
) DeleteArtist {
	return DeleteArtist{
		artistRepository:        artistRepository,
		messagePublisherService: messagePublisherService,
		transaction:             transaction,
	}
}

func (d DeleteArtist) Handle(request requests.DeleteArtistRequest) *httperror.ErrorCode {
	var artist model.Artist
	if err := d.artistRepository.GetWithSongsOrAlbums(&artist, request.ID, request.WithSongs, request.WithAlbums); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return httperror.NotFoundError(errors.New("artist not found"))
	}

	err := d.transaction.Execute(func(factory transaction.RepositoryFactory) error {
		txArtistRepo := factory.NewArtistRepository()

		if request.WithAlbums {
			if err := txArtistRepo.DeleteAlbums([]uuid.UUID{request.ID}); err != nil {
				return err
			}
		}
		if request.WithSongs {
			if err := txArtistRepo.DeleteSongs([]uuid.UUID{request.ID}); err != nil {
				return err
			}
		}

		if err := txArtistRepo.Delete([]uuid.UUID{request.ID}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	if err = d.messagePublisherService.Publish(topics.ArtistsDeletedTopic, artist); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
