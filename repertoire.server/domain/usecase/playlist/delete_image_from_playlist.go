package playlist

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type DeleteImageFromPlaylist struct {
	repository     repository.PlaylistRepository
	storageService service.StorageService
}

func NewDeleteImageFromPlaylist(
	repository repository.PlaylistRepository,
	storageService service.StorageService,
) DeleteImageFromPlaylist {
	return DeleteImageFromPlaylist{
		repository:     repository,
		storageService: storageService,
	}
}

func (d DeleteImageFromPlaylist) Handle(id uuid.UUID) *wrapper.ErrorCode {
	var playlist model.Playlist
	err := d.repository.Get(&playlist, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(playlist).IsZero() {
		return wrapper.NotFoundError(errors.New("playlist not found"))
	}

	err = d.storageService.Delete(string(*playlist.ImageURL))
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	playlist.ImageURL = nil
	err = d.repository.Update(&playlist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
