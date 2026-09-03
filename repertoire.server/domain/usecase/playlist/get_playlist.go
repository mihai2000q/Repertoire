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
	playlistRepository repository.PlaylistRepository
}

func NewGetPlaylist(playlistRepository repository.PlaylistRepository) GetPlaylist {
	return GetPlaylist{
		playlistRepository: playlistRepository,
	}
}

func (g GetPlaylist) Handle(request requests.GetPlaylistRequest) (playlist model.Playlist, e *httperror.ErrorCode) {
	if err := g.playlistRepository.Get(&playlist, request.ID); err != nil {
		return playlist, httperror.DatabaseError(err)
	}
	if reflect.ValueOf(playlist).IsZero() {
		return playlist, httperror.NotFoundError(errors.New("playlist not found"))
	}
	return playlist, nil
}
