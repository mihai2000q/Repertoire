package playlist

import (
	"errors"
	"repertoire/server/data/repository"
	"repertoire/server/model"
	"repertoire/server/utils/wrapper"

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

func (g GetPlaylist) Handle(id uuid.UUID) (playlist model.Playlist, e *wrapper.ErrorCode) {
	err := g.repository.GetWithAssociations(&playlist, id)
	if err != nil {
		return playlist, wrapper.InternalServerError(err)
	}
	if playlist.ID == uuid.Nil {
		return playlist, wrapper.NotFoundError(errors.New("playlist not found"))
	}
	return playlist, nil
}
