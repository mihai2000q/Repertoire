package playlist

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type UpdatePlaylist struct {
	repository repository.PlaylistRepository
}

func NewUpdatePlaylist(repository repository.PlaylistRepository) UpdatePlaylist {
	return UpdatePlaylist{repository: repository}
}

func (u UpdatePlaylist) Handle(request requests.UpdatePlaylistRequest) *wrapper.ErrorCode {
	var playlist model.Playlist
	err := u.repository.Get(&playlist, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(playlist).IsZero() {
		return wrapper.NotFoundError(errors.New("playlist not found"))
	}

	playlist.Title = request.Title
	playlist.Description = request.Description

	err = u.repository.Update(&playlist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
