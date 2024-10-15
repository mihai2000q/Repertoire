package artist

import (
	"errors"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils"

	"github.com/google/uuid"
)

type GetArtist struct {
	repository repository.ArtistRepository
}

func NewGetArtist(repository repository.ArtistRepository) GetArtist {
	return GetArtist{
		repository: repository,
	}
}

func (g GetArtist) Handle(id uuid.UUID) (artist models.Artist, e *utils.ErrorCode) {
	err := g.repository.Get(&artist, id)
	if err != nil {
		return artist, utils.InternalServerError(err)
	}
	if artist.ID == uuid.Nil {
		return artist, utils.NotFoundError(errors.New("artist not found"))
	}
	return artist, nil
}
