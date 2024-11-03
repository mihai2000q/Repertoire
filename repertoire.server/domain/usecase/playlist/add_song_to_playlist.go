package playlist

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type AddSongToPlaylist struct {
	songRepository repository.SongRepository
	repository     repository.PlaylistRepository
}

func NewAddSongToPlaylist(
	songRepository repository.SongRepository,
	repository repository.PlaylistRepository,
) AddSongToPlaylist {
	return AddSongToPlaylist{
		songRepository: songRepository,
		repository:     repository,
	}
}

func (a AddSongToPlaylist) Handle(request requests.AddSongToPlaylistRequest) *wrapper.ErrorCode {
	var playlist model.Playlist
	err := a.repository.Get(&playlist, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	var song model.Song
	err = a.songRepository.Get(&song, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	var count int64
	err = a.repository.CountSongs(&count, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = a.repository.AddSong(&playlist, &song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
