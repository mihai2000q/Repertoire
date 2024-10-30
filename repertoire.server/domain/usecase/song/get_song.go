package song

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type GetSong struct {
	repository repository.SongRepository
}

func NewGetSong(repository repository.SongRepository) GetSong {
	return GetSong{
		repository: repository,
	}
}

func (g GetSong) Handle(id uuid.UUID) (song model.Song, e *wrapper.ErrorCode) {
	err := g.repository.GetWithAssociations(&song, id)
	if err != nil {
		return song, wrapper.InternalServerError(err)
	}
	if song.ID == uuid.Nil {
		return song, wrapper.NotFoundError(errors.New("song not found"))
	}
	return song, nil
}
