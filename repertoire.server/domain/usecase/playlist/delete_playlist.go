package playlist

import (
	"errors"
	"net/http"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type DeletePlaylist struct {
	repository              repository.PlaylistRepository
	storageService          service.StorageService
	storageFilePathProvider provider.StorageFilePathProvider
}

func NewDeletePlaylist(
	repository repository.PlaylistRepository,
	storageService service.StorageService,
	storageFilePathProvider provider.StorageFilePathProvider,
) DeletePlaylist {
	return DeletePlaylist{
		repository:              repository,
		storageService:          storageService,
		storageFilePathProvider: storageFilePathProvider,
	}
}

func (d DeletePlaylist) Handle(id uuid.UUID) *wrapper.ErrorCode {
	var playlist model.Playlist
	err := d.repository.Get(&playlist, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(playlist).IsZero() {
		return wrapper.NotFoundError(errors.New("playlist not found"))
	}

	directoryPath := d.storageFilePathProvider.GetPlaylistDirectoryPath(playlist)
	errCode := d.storageService.DeleteDirectory(directoryPath)
	if errCode != nil && errCode.Code != http.StatusNotFound {
		return errCode
	}

	err = d.repository.Delete(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
