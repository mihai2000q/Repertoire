package album

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type DeleteAlbum struct {
	repository              repository.AlbumRepository
	storageService          service.StorageService
	storageFilePathProvider provider.StorageFilePathProvider
}

func NewDeleteAlbum(
	repository repository.AlbumRepository,
	storageService service.StorageService,
	storageFilePathProvider provider.StorageFilePathProvider,
) DeleteAlbum {
	return DeleteAlbum{
		repository:              repository,
		storageService:          storageService,
		storageFilePathProvider: storageFilePathProvider,
	}
}

func (d DeleteAlbum) Handle(id uuid.UUID) *wrapper.ErrorCode {
	var album model.Album
	err := d.repository.Get(&album, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return wrapper.NotFoundError(errors.New("album not found"))
	}

	if d.storageFilePathProvider.HasAlbumFiles(album) {
		directoryPath := d.storageFilePathProvider.GetAlbumDirectoryPath(album)
		err = d.storageService.DeleteDirectory(directoryPath)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
	}

	err = d.repository.Delete(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
