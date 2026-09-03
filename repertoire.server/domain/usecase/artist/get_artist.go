package artist

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

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

func (g GetArtist) Handle(id uuid.UUID) (artist model.Artist, e *httperror.ErrorCode) {
	if err := g.repository.GetWithAssociations(&artist, id); err != nil {
		return artist, httperror.DatabaseError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return artist, httperror.NotFoundError(errors.New("artist not found"))
	}
	return artist, nil
}
