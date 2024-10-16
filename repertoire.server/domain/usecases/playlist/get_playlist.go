package playlist

import (
	"errors"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils"

	"github.com/google/uuid"
)

type GetPlaylist struct {
	repository repository.PlaylistRepository
}

func NewGetPlaylist(repository repository.PlaylistRepository) GetPlaylist {
	return GetPlaylist{
		repository: repository,
	}
}

func (g GetPlaylist) Handle(id uuid.UUID) (playlist models.Playlist, e *utils.ErrorCode) {
	err := g.repository.Get(&playlist, id)
	if err != nil {
		return playlist, utils.InternalServerError(err)
	}
	if playlist.ID == uuid.Nil {
		return playlist, utils.NotFoundError(errors.New("playlist not found"))
	}
	return playlist, nil
}
