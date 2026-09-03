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
	artistRepository repository.ArtistRepository
}

func NewGetArtist(artistRepository repository.ArtistRepository) GetArtist {
	return GetArtist{
		artistRepository: artistRepository,
	}
}

func (g GetArtist) Handle(id uuid.UUID) (artist model.Artist, e *httperror.ErrorCode) {
	if err := g.artistRepository.GetWithAssociations(&artist, id); err != nil {
		return artist, httperror.DatabaseError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return artist, httperror.NotFoundError(errors.New("artist not found"))
	}
	return artist, nil
}
