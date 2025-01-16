package artist

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type UpdateArtist struct {
	repository repository.ArtistRepository
}

func NewUpdateArtist(repository repository.ArtistRepository) UpdateArtist {
	return UpdateArtist{repository: repository}
}

func (u UpdateArtist) Handle(request requests.UpdateArtistRequest) *wrapper.ErrorCode {
	var artist model.Artist
	err := u.repository.Get(&artist, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return wrapper.NotFoundError(errors.New("artist not found"))
	}

	artist.Name = request.Name

	err = u.repository.Update(&artist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
