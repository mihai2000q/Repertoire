package song

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type AddPerfectPlaylistSongRehearsals struct {
	playlistRepository repository.PlaylistRepository
	songProcessor      processor.SongProcessor
	transactionManager transaction.Manager
}

func NewAddPerfectPlaylistSongRehearsals(
	playlistRepository repository.PlaylistRepository,
	songProcessor processor.SongProcessor,
	transactionManager transaction.Manager,
) AddPerfectPlaylistSongRehearsals {
	return AddPerfectPlaylistSongRehearsals{
		playlistRepository: playlistRepository,
		songProcessor:      songProcessor,
		transactionManager: transactionManager,
	}
}

func (a AddPerfectPlaylistSongRehearsals) Handle(request requests.AddPerfectPlaylistSongRehearsalsRequest) *httperror.ErrorCode {
	var playlistSongs []model.PlaylistSong
	err := a.playlistRepository.GetPlaylistSongsByIDsWithPartsAndDefaultOccurrences(&playlistSongs, request.IDs, request.PlaylistID)
	if err != nil {
		return httperror.DatabaseError(err)
	}
	if len(playlistSongs) != len(request.IDs) {
		return httperror.NotFoundError(errors.New("playlist songs not found"))
	}

	var errCode *httperror.ErrorCode
	err = a.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txSongPartRepo := factory.NewSongPartRepository()
		txSongRepo := factory.NewSongRepository()

		var newSongs []model.Song
		for _, playlistSong := range playlistSongs {
			errC, isUpdated := a.songProcessor.AddPerfectRehearsal(&playlistSong.Song, txSongPartRepo)
			if errC != nil {
				errCode = errC
				return errCode.Error
			}
			if isUpdated {
				newSongs = append(newSongs, playlistSong.Song)
			}
		}

		if len(newSongs) > 0 {
			err = txSongRepo.UpdateAllWithAssociations(&newSongs)
			if err != nil {
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
