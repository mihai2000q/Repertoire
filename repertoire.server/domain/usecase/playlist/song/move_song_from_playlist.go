package song

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/reorder"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type MoveSongFromPlaylist struct {
	repository repository.PlaylistRepository
}

func NewMoveSongFromPlaylist(repository repository.PlaylistRepository) MoveSongFromPlaylist {
	return MoveSongFromPlaylist{repository: repository}
}

func (m MoveSongFromPlaylist) Handle(request requests.MoveSongFromPlaylistRequest) *wrapper.ErrorCode {
	var playlistSongs []model.PlaylistSong
	err := m.repository.GetPlaylistSongs(&playlistSongs, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	errCode := reorder.MoveEntity(
		playlistSongs,
		request.PlaylistSongID,
		request.OverPlaylistSongID,
		&reorder.Config{
			OrderField:            "SongTrackNo",
			StartOffset:           1,
			EntityNotFoundMsg:     "song not found",
			OverEntityNotFoundMsg: "over song not found",
		},
	)
	if errCode != nil {
		return errCode
	}

	err = m.repository.UpdateAllPlaylistSongs(&playlistSongs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
