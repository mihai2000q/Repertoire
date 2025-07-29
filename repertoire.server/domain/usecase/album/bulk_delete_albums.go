package album

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type BulkDeleteAlbums struct {
	repository              repository.AlbumRepository
	messagePublisherService service.MessagePublisherService
}

func NewBulkDeleteAlbums(
	repository repository.AlbumRepository,
	messagePublisherService service.MessagePublisherService,
) BulkDeleteAlbums {
	return BulkDeleteAlbums{
		repository:              repository,
		messagePublisherService: messagePublisherService,
	}
}

func (b BulkDeleteAlbums) Handle(request requests.BulkDeleteAlbumsRequest) *wrapper.ErrorCode {
	var albums []model.Album
	var err error
	if request.WithSongs {
		err = b.repository.GetAllByIDsWithSongs(&albums, request.IDs)
	} else {
		err = b.repository.GetAllByIDs(&albums, request.IDs)
	}
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if len(albums) == 0 {
		return wrapper.NotFoundError(errors.New("albums not found"))
	}

	if request.WithSongs {
		err = b.repository.DeleteWithSongs(request.IDs)
	} else {
		err = b.repository.Delete(request.IDs)
	}
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = b.messagePublisherService.Publish(topics.AlbumsDeletedTopic, albums)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
