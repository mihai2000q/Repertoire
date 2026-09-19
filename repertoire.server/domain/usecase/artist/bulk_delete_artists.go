package artist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
)

type BulkDeleteArtists struct {
	artistRepository        repository.ArtistRepository
	messagePublisherService service.MessagePublisherService
	transaction             transaction.Manager
}

func NewBulkDeleteArtists(
	artistRepository repository.ArtistRepository,
	messagePublisherService service.MessagePublisherService,
	transaction transaction.Manager,
) BulkDeleteArtists {
	return BulkDeleteArtists{
		artistRepository:        artistRepository,
		messagePublisherService: messagePublisherService,
		transaction:             transaction,
	}
}

func (b BulkDeleteArtists) Handle(request requests.BulkDeleteArtistsRequest) *httperror.ErrorCode {
	var artists []model.Artist
	err := b.artistRepository.GetAllByIDs(&artists, request.IDs, request.WithSongs, request.WithAlbums)
	if err != nil {
		return httperror.DatabaseError(err)
	}
	if len(artists) != len(request.IDs) {
		return httperror.NotFoundError(errors.New("artists not found"))
	}

	err = b.transaction.Execute(func(factory transaction.RepositoryFactory) error {
		txArtistRepo := factory.NewArtistRepository()

		if request.WithAlbums {
			if err = txArtistRepo.DeleteAlbums(request.IDs); err != nil {
				return err
			}
		}
		if request.WithSongs {
			if err = txArtistRepo.DeleteSongs(request.IDs); err != nil {
				return err
			}
		}

		if err = txArtistRepo.Delete(request.IDs); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	err = b.messagePublisherService.Publish(topics.ArtistsDeletedTopic, artists)
	if err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
