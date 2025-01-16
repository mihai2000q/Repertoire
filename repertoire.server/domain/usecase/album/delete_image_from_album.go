package album

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type DeleteImageFromAlbum struct {
	repository     repository.AlbumRepository
	storageService service.StorageService
}

func NewDeleteImageFromAlbum(
	repository repository.AlbumRepository,
	storageService service.StorageService,
) DeleteImageFromAlbum {
	return DeleteImageFromAlbum{
		repository:     repository,
		storageService: storageService,
	}
}

func (d DeleteImageFromAlbum) Handle(id uuid.UUID) *wrapper.ErrorCode {
	var album model.Album
	err := d.repository.Get(&album, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return wrapper.NotFoundError(errors.New("album not found"))
	}
	if album.ImageURL == nil {
		return wrapper.BadRequestError(errors.New("album does not have an image"))
	}

	err = d.storageService.DeleteFile(*album.ImageURL)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	album.ImageURL = nil
	err = d.repository.Update(&album)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
