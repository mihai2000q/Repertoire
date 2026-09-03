package album

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
)

type BulkDeleteAlbums struct {
	albumRepository         repository.AlbumRepository
	messagePublisherService service.MessagePublisherService
}

func NewBulkDeleteAlbums(
	albumRepository repository.AlbumRepository,
	messagePublisherService service.MessagePublisherService,
) BulkDeleteAlbums {
	return BulkDeleteAlbums{
		albumRepository:         albumRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (b BulkDeleteAlbums) Handle(request requests.BulkDeleteAlbumsRequest) *httperror.ErrorCode {
	var albums []model.Album
	var err error
	if request.WithSongs {
		err = b.albumRepository.GetAllByIDsWithSongs(&albums, request.IDs)
	} else {
		err = b.albumRepository.GetAllByIDs(&albums, request.IDs)
	}
	if err != nil {
		return httperror.DatabaseError(err)
	}
	if len(albums) == 0 {
		return httperror.NotFoundError(errors.New("albums not found"))
	}

	if request.WithSongs {
		err = b.albumRepository.DeleteWithSongs(request.IDs)
	} else {
		err = b.albumRepository.Delete(request.IDs)
	}
	if err != nil {
		return httperror.DatabaseError(err)
	}

	if err = b.messagePublisherService.Publish(topics.AlbumsDeletedTopic, albums); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
