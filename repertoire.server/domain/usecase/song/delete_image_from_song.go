package song

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type DeleteImageFromSong struct {
	repository     repository.SongRepository
	storageService service.StorageService
}

func NewDeleteImageFromSong(
	repository repository.SongRepository,
	storageService service.StorageService,
) DeleteImageFromSong {
	return DeleteImageFromSong{
		repository:     repository,
		storageService: storageService,
	}
}

func (d DeleteImageFromSong) Handle(id uuid.UUID) *wrapper.ErrorCode {
	var song model.Song
	err := d.repository.Get(&song, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	err = d.storageService.Delete(string(*song.ImageURL))
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	song.ImageURL = nil
	err = d.repository.Update(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
