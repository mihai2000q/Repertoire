package playlist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"
)

type RemoveSongsFromPlaylist struct {
	repository repository.PlaylistRepository
}

func NewRemoveSongsFromPlaylist(repository repository.PlaylistRepository) RemoveSongsFromPlaylist {
	return RemoveSongsFromPlaylist{repository: repository}
}

func (r RemoveSongsFromPlaylist) Handle(request requests.RemoveSongsFromPlaylistRequest) *wrapper.ErrorCode {
	var playlistSongs []model.PlaylistSong
	err := r.repository.GetPlaylistSongs(&playlistSongs, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	var songsToDelete []model.PlaylistSong
	var songsToPreserve []model.PlaylistSong

	songTrackNo := uint(1)
	for _, playlistSong := range playlistSongs {
		if slices.Contains(request.SongIDs, playlistSong.SongID) {
			songsToDelete = append(songsToDelete, playlistSong)
		} else {
			// reorder preserved songs
			playlistSong.SongTrackNo = songTrackNo
			songsToPreserve = append(songsToPreserve, playlistSong)
			songTrackNo++
		}
	}

	if len(songsToDelete) != len(request.SongIDs) {
		return wrapper.NotFoundError(errors.New("could not find all songs"))
	}
	// TODO: TRANSACTION!!!
	err = r.repository.RemoveSongs(&songsToDelete)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = r.repository.UpdateAllPlaylistSongs(&songsToPreserve)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
