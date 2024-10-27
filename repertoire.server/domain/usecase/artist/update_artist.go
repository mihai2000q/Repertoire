package artist

import (
	"errors"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/model"
	"repertoire/utils/wrapper"

	"github.com/google/uuid"
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
	if artist.ID == uuid.Nil {
		return wrapper.NotFoundError(errors.New("artist not found"))
	}

	artist.Name = request.Name

	err = u.repository.Update(&artist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
