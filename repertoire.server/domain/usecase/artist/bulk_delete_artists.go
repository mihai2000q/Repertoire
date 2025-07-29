package artist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type BulkDeleteArtists struct {
	repository              repository.ArtistRepository
	messagePublisherService service.MessagePublisherService
	transaction             transaction.Manager
}

func NewBulkDeleteArtists(
	repository repository.ArtistRepository,
	messagePublisherService service.MessagePublisherService,
	transaction transaction.Manager,
) BulkDeleteArtists {
	return BulkDeleteArtists{
		repository:              repository,
		messagePublisherService: messagePublisherService,
		transaction:             transaction,
	}
}

func (b BulkDeleteArtists) Handle(request requests.BulkDeleteArtistsRequest) *wrapper.ErrorCode {
	var artists []model.Artist
	err := b.repository.GetAllByIDs(&artists, request.IDs, request.WithSongs, request.WithAlbums)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if len(artists) == 0 {
		return wrapper.NotFoundError(errors.New("artists not found"))
	}

	err = b.transaction.Execute(func(factory transaction.RepositoryFactory) error {
		artistRepo := factory.NewArtistRepository()

		if request.WithAlbums {
			err = artistRepo.DeleteAlbums(request.IDs)
			if err != nil {
				return err
			}
		}
		if request.WithSongs {
			err = artistRepo.DeleteSongs(request.IDs)
			if err != nil {
				return err
			}
		}

		err = artistRepo.Delete(request.IDs)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = b.messagePublisherService.Publish(topics.ArtistsDeletedTopic, artists)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
