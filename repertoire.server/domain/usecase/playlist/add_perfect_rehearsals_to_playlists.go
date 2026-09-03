package playlist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type AddPerfectRehearsalsToPlaylists struct {
	playlistRepository repository.PlaylistRepository
	songProcessor      processor.SongProcessor
	transactionManager transaction.Manager
}

func NewAddPerfectRehearsalsToPlaylists(
	playlistRepository repository.PlaylistRepository,
	songProcessor processor.SongProcessor,
	transactionManager transaction.Manager,
) AddPerfectRehearsalsToPlaylists {
	return AddPerfectRehearsalsToPlaylists{
		playlistRepository: playlistRepository,
		songProcessor:      songProcessor,
		transactionManager: transactionManager,
	}
}

func (a AddPerfectRehearsalsToPlaylists) Handle(request requests.AddPerfectRehearsalsToPlaylistsRequest) *httperror.ErrorCode {
	var playlists []model.Playlist
	err := a.playlistRepository.GetAllByIDsWithSongPartsAndDefaultOccurrences(&playlists, request.IDs)
	if err != nil {
		return httperror.DatabaseError(err)
	}
	if len(playlists) != len(request.IDs) {
		return httperror.NotFoundError(errors.New("playlists not found"))
	}

	var errCode *httperror.ErrorCode
	err = a.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txSongPartRepo := factory.NewSongPartRepository()
		txSongRepo := factory.NewSongRepository()

		var newSongs []model.Song
		for _, playlist := range playlists {
			for _, playlistSong := range playlist.PlaylistSongs {
				errC, isUpdated := a.songProcessor.AddPerfectRehearsal(&playlistSong.Song, txSongPartRepo)
				if errC != nil {
					errCode = errC
					return errCode.Error
				}
				if isUpdated {
					newSongs = append(newSongs, playlistSong.Song)
				}
			}
		}

		if len(newSongs) > 0 {
			if err = txSongRepo.UpdateAllWithAssociations(&newSongs); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		if errCode != nil {
			return errCode
		}
		return httperror.DatabaseError(err)
	}

	return nil
}
