package playlist

import (
	"errors"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils"

	"github.com/google/uuid"
)

type UpdatePlaylist struct {
	repository repository.PlaylistRepository
}

func NewUpdatePlaylist(repository repository.PlaylistRepository) UpdatePlaylist {
	return UpdatePlaylist{repository: repository}
}

func (u UpdatePlaylist) Handle(request requests.UpdatePlaylistRequest) *utils.ErrorCode {
	var playlist models.Playlist
	err := u.repository.Get(&playlist, request.ID)
	if err != nil {
		return utils.InternalServerError(err)
	}
	if playlist.ID == uuid.Nil {
		return utils.NotFoundError(errors.New("playlist not found"))
	}

	playlist.Title = request.Title
	playlist.Description = request.Description

	err = u.repository.Update(&playlist)
	if err != nil {
		return utils.InternalServerError(err)
	}

	return nil
}
