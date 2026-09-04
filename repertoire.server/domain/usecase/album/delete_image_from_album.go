package album

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type DeleteImageFromAlbum struct {
	albumRepository         repository.AlbumRepository
	storageService          service.StorageService
	messagePublisherService service.MessagePublisherService
}

func NewDeleteImageFromAlbum(
	albumRepository repository.AlbumRepository,
	storageService service.StorageService,
	messagePublisherService service.MessagePublisherService,
) DeleteImageFromAlbum {
	return DeleteImageFromAlbum{
		albumRepository:         albumRepository,
		storageService:          storageService,
		messagePublisherService: messagePublisherService,
	}
}

func (d DeleteImageFromAlbum) Handle(id uuid.UUID) *httperror.ErrorCode {
	var album model.Album
	if err := d.albumRepository.Get(&album, id); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return httperror.NotFoundError(errors.New("album not found"))
	}
	if album.ImageURL == nil {
		return nil
	}

	errCode := d.storageService.DeleteFile(*album.ImageURL)
	if errCode != nil {
		return errCode
	}

	album.ImageURL = nil
	if err := d.albumRepository.Update(&album); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := d.messagePublisherService.Publish(topics.AlbumsUpdatedTopic, []uuid.UUID{album.ID}); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
