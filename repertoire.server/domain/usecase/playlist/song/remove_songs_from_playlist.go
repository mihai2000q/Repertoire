package song

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
	"slices"
)

type RemoveSongsFromPlaylist struct {
	playlistRepository repository.PlaylistRepository
	transaction        transaction.Manager
}

func NewRemoveSongsFromPlaylist(
	playlistRepository repository.PlaylistRepository,
	transaction transaction.Manager,
) RemoveSongsFromPlaylist {
	return RemoveSongsFromPlaylist{
		playlistRepository: playlistRepository,
		transaction:        transaction,
	}
}

func (r RemoveSongsFromPlaylist) Handle(request requests.RemoveSongsFromPlaylistRequest) *httperror.ErrorCode {
	var playlistSongs []model.PlaylistSong
	if err := r.playlistRepository.GetPlaylistSongs(&playlistSongs, request.ID); err != nil {
		return httperror.DatabaseError(err)
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
		return httperror.NotFoundError(errors.New("could not find all songs"))
	}

	err := r.transaction.Execute(func(factory transaction.RepositoryFactory) error {
		txPlaylistRepo := factory.NewPlaylistRepository()

		if err := txPlaylistRepo.RemoveSongs(&songsToDelete); err != nil {
			return err
		}
		if err := txPlaylistRepo.UpdateAllPlaylistSongs(&songsToPreserve); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
