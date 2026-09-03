package song

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type GetSong struct {
	repository repository.SongRepository
}

func NewGetSong(repository repository.SongRepository) GetSong {
	return GetSong{
		repository: repository,
	}
}

func (g GetSong) Handle(id uuid.UUID) (song model.Song, e *httperror.ErrorCode) {
	if err := g.repository.GetWithAssociations(&song, id); err != nil {
		return song, httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return song, httperror.NotFoundError(errors.New("song not found"))
	}
	return song, nil
}
