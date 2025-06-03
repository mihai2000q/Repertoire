package playlist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"
)

type RemoveSongsFromPlaylist struct {
	repository  repository.PlaylistRepository
	transaction transaction.Manager
}

func NewRemoveSongsFromPlaylist(
	repository repository.PlaylistRepository,
	transaction transaction.Manager,
) RemoveSongsFromPlaylist {
	return RemoveSongsFromPlaylist{
		repository:  repository,
		transaction: transaction,
	}
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
		if slices.Contains(request.PlaylistSongIDs, playlistSong.ID) {
			songsToDelete = append(songsToDelete, playlistSong)
		} else {
			// reorder preserved songs
			playlistSong.SongTrackNo = songTrackNo
			songsToPreserve = append(songsToPreserve, playlistSong)
			songTrackNo++
		}
	}

	if len(songsToDelete) != len(request.PlaylistSongIDs) {
		return wrapper.NotFoundError(errors.New("could not find all songs"))
	}

	err = r.transaction.Execute(func(factory transaction.RepositoryFactory) error {
		playlistRepo := factory.NewPlaylistRepository()

		if err := playlistRepo.RemoveSongs(&songsToDelete); err != nil {
			return err
		}
		if err := playlistRepo.UpdateAllPlaylistSongs(&songsToPreserve); err != nil {
			return err
		} // preserver order
		return nil
	})
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
