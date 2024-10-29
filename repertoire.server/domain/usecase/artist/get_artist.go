package artist

import (
	"errors"
	"repertoire/data/repository"
	"repertoire/model"
	"repertoire/utils/wrapper"

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

func (g GetArtist) Handle(id uuid.UUID) (artist model.Artist, e *wrapper.ErrorCode) {
	err := g.repository.GetWithAssociations(&artist, id)
	if err != nil {
		return artist, wrapper.InternalServerError(err)
	}
	if artist.ID == uuid.Nil {
		return artist, wrapper.NotFoundError(errors.New("artist not found"))
	}
	return artist, nil
}
