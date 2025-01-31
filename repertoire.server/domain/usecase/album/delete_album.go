package album

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
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

func (d DeleteAlbum) Handle(request requests.DeleteAlbumRequest) *wrapper.ErrorCode {
	var album model.Album
	err := d.repository.Get(&album, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return wrapper.NotFoundError(errors.New("album not found"))
	}

	if d.storageFilePathProvider.HasAlbumFiles(album) {
		directoryPath := d.storageFilePathProvider.GetAlbumDirectoryPath(album)
		errCode := d.storageService.DeleteDirectory(directoryPath)
		if errCode != nil {
			return errCode
		}
	}

	if request.WithSongs {
		err = d.repository.DeleteWithSongs(request.ID)
	} else {
		err = d.repository.Delete(request.ID)
	}
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
