package song

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/reorder"
	"repertoire/server/model"
)

type MoveSongFromPlaylist struct {
	repository repository.PlaylistRepository
}

func NewMoveSongFromPlaylist(repository repository.PlaylistRepository) MoveSongFromPlaylist {
	return MoveSongFromPlaylist{repository: repository}
}

func (m MoveSongFromPlaylist) Handle(request requests.MoveSongFromPlaylistRequest) *httperror.ErrorCode {
	var playlistSongs []model.PlaylistSong
	if err := m.repository.GetPlaylistSongs(&playlistSongs, request.ID); err != nil {
		return httperror.DatabaseError(err)
	}
	if len(playlistSongs) == 0 {
		return httperror.NotFoundError(errors.New("playlist not found"))
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

	if err := m.repository.UpdateAllPlaylistSongs(&playlistSongs); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
