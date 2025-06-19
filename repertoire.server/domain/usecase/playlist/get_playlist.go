package playlist

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type GetPlaylist struct {
	repository repository.PlaylistRepository
}

func NewGetPlaylist(repository repository.PlaylistRepository) GetPlaylist {
	return GetPlaylist{
		repository: repository,
	}
}

func (g GetPlaylist) Handle(request requests.GetPlaylistRequest) (playlist model.Playlist, e *wrapper.ErrorCode) {
	err := g.repository.Get(&playlist, request.ID)
	if err != nil {
		return playlist, wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(playlist).IsZero() {
		return playlist, wrapper.NotFoundError(errors.New("playlist not found"))
	}
	return playlist, nil
}
