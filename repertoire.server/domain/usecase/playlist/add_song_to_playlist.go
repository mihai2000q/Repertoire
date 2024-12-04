package playlist

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type AddSongToPlaylist struct {
	repository repository.PlaylistRepository
}

func NewAddSongToPlaylist(repository repository.PlaylistRepository) AddSongToPlaylist {
	return AddSongToPlaylist{repository: repository}
}

func (a AddSongToPlaylist) Handle(request requests.AddSongToPlaylistRequest) *wrapper.ErrorCode {
	var count int64
	err := a.repository.CountSongs(&count, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	playlistSong := model.PlaylistSong{
		PlaylistID:  request.ID,
		SongID:      request.SongID,
		SongTrackNo: uint(count),
	}
	err = a.repository.AddSong(&playlistSong)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
