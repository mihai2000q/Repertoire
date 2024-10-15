package artist

import (
	"errors"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils"

	"github.com/google/uuid"
)

type UpdateArtist struct {
	repository repository.ArtistRepository
}

func NewUpdateArtist(repository repository.ArtistRepository) UpdateArtist {
	return UpdateArtist{repository: repository}
}

func (u UpdateArtist) Handle(request requests.UpdateArtistRequest) *utils.ErrorCode {
	var artist models.Artist
	err := u.repository.Get(&artist, request.ID)
	if err != nil {
		return utils.InternalServerError(err)
	}
	if artist.ID == uuid.Nil {
		return utils.NotFoundError(errors.New("artist not found"))
	}

	artist.Name = request.Name

	err = u.repository.Update(&artist)
	if err != nil {
		return utils.InternalServerError(err)
	}

	return nil
}
