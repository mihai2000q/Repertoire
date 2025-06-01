package playlist

import (
	"github.com/google/uuid"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type AddSongsToPlaylist struct {
	repository repository.PlaylistRepository
}

func NewAddSongsToPlaylist(repository repository.PlaylistRepository) AddSongsToPlaylist {
	return AddSongsToPlaylist{repository: repository}
}

func (a AddSongsToPlaylist) Handle(request requests.AddSongsToPlaylistRequest) *wrapper.ErrorCode {
	var count int64
	err := a.repository.CountSongs(&count, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	songsLength := int(count + 1)
	var playlistSongs []model.PlaylistSong
	for i, songID := range request.SongIDs {
		playlistSong := model.PlaylistSong{
			ID:          uuid.New(),
			PlaylistID:  request.ID,
			SongID:      songID,
			SongTrackNo: uint(songsLength + i),
		}
		playlistSongs = append(playlistSongs, playlistSong)
	}

	err = a.repository.AddSongs(&playlistSongs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
