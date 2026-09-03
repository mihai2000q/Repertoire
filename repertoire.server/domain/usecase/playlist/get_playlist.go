package playlist

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
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

func (g GetPlaylist) Handle(request requests.GetPlaylistRequest) (playlist model.Playlist, e *httperror.ErrorCode) {
	if err := g.repository.Get(&playlist, request.ID); err != nil {
		return playlist, httperror.DatabaseError(err)
	}
	if reflect.ValueOf(playlist).IsZero() {
		return playlist, httperror.NotFoundError(errors.New("playlist not found"))
	}
	return playlist, nil
}
