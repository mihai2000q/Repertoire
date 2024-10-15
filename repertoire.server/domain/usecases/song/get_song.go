package song

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils"
)

type GetSong struct {
	repository repository.SongRepository
}

func NewGetSong(repository repository.SongRepository) *GetSong {
	return &GetSong{
		repository: repository,
	}
}

func (g GetSong) Handle(id uuid.UUID) (song models.Song, e *utils.ErrorCode) {
	err := g.repository.Get(&song, id)
	if err != nil {
		return song, utils.InternalServerError(err)
	}
	if song.ID == uuid.Nil {
		return song, utils.NotFoundError(errors.New("song not found"))
	}
	return song, nil
}
