package artist

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type DeleteImageFromArtist struct {
	repository     repository.ArtistRepository
	storageService service.StorageService
}

func NewDeleteImageFromArtist(
	repository repository.ArtistRepository,
	storageService service.StorageService,
) DeleteImageFromArtist {
	return DeleteImageFromArtist{
		repository:     repository,
		storageService: storageService,
	}
}

func (d DeleteImageFromArtist) Handle(id uuid.UUID) *wrapper.ErrorCode {
	var artist model.Artist
	err := d.repository.Get(&artist, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return wrapper.NotFoundError(errors.New("artist not found"))
	}

	err = d.storageService.Delete(string(*artist.ImageURL))
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	artist.ImageURL = nil
	err = d.repository.Update(&artist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}