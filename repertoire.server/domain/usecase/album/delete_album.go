package album

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type DeleteAlbum struct {
	repository              repository.AlbumRepository
	messagePublisherService service.MessagePublisherService
}

func NewDeleteAlbum(
	repository repository.AlbumRepository,
	messagePublisherService service.MessagePublisherService,
) DeleteAlbum {
	return DeleteAlbum{
		repository:              repository,
		messagePublisherService: messagePublisherService,
	}
}

func (d DeleteAlbum) Handle(request requests.DeleteAlbumRequest) *wrapper.ErrorCode {
	var album model.Album
	var err error
	if request.WithSongs {
		err = d.repository.GetWithSongs(&album, request.ID)
	} else {
		err = d.repository.Get(&album, request.ID)
	}
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return wrapper.NotFoundError(errors.New("album not found"))
	}

	if request.WithSongs {
		err = d.repository.DeleteWithSongs([]uuid.UUID{request.ID})
	} else {
		err = d.repository.Delete([]uuid.UUID{request.ID})
	}
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = d.messagePublisherService.Publish(topics.AlbumsDeletedTopic, []model.Album{album})
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
