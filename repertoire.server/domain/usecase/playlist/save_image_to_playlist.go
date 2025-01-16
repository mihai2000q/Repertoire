package playlist

import (
	"errors"
	"mime/multipart"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type SaveImageToPlaylist struct {
	repository              repository.PlaylistRepository
	storageFilePathProvider provider.StorageFilePathProvider
	storageService          service.StorageService
}

func NewSaveImageToPlaylist(
	repository repository.PlaylistRepository,
	storageFilePathProvider provider.StorageFilePathProvider,
	storageService service.StorageService,
) SaveImageToPlaylist {
	return SaveImageToPlaylist{
		repository:              repository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
	}
}

func (s SaveImageToPlaylist) Handle(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode {
	var playlist model.Playlist
	err := s.repository.Get(&playlist, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(playlist).IsZero() {
		return wrapper.NotFoundError(errors.New("playlist not found"))
	}

	imagePath := s.storageFilePathProvider.GetPlaylistImagePath(file, playlist)

	err = s.storageService.Upload(file, imagePath)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	playlist.ImageURL = (*internal.FilePath)(&imagePath)
	err = s.repository.Update(&playlist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
