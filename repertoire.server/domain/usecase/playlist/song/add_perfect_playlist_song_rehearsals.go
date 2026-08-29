package song

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type AddPerfectPlaylistSongRehearsals struct {
	repository         repository.PlaylistRepository
	songProcessor      processor.SongProcessor
	transactionManager transaction.Manager
}

func NewAddPerfectPlaylistSongRehearsals(
	repository repository.PlaylistRepository,
	songProcessor processor.SongProcessor,
	transactionManager transaction.Manager,
) AddPerfectPlaylistSongRehearsals {
	return AddPerfectPlaylistSongRehearsals{
		repository:         repository,
		songProcessor:      songProcessor,
		transactionManager: transactionManager,
	}
}

func (a AddPerfectPlaylistSongRehearsals) Handle(request requests.AddPerfectPlaylistSongRehearsalsRequest) *wrapper.ErrorCode {
	var playlistSongs []model.PlaylistSong
	err := a.repository.GetPlaylistSongsByIDsWithPartsAndDefaultOccurrences(&playlistSongs, request.IDs, request.PlaylistID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if len(playlistSongs) != len(request.IDs) {
		return wrapper.NotFoundError(errors.New("playlist songs not found"))
	}

	var errCode *wrapper.ErrorCode
	err = a.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txSongPartRepository := factory.NewSongPartRepository()
		txSongRepository := factory.NewSongRepository()

		var newSongs []model.Song
		for _, playlistSong := range playlistSongs {
			errC, isUpdated := a.songProcessor.AddPerfectRehearsal(&playlistSong.Song, txSongPartRepository)
			if errC != nil {
				errCode = errC
				return errCode.Error
			}
			if isUpdated {
				newSongs = append(newSongs, playlistSong.Song)
			}
		}

		if len(newSongs) > 0 {
			err = txSongRepository.UpdateAllWithAssociations(&newSongs)
			if err != nil {
				errCode = wrapper.InternalServerError(err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		if errCode != nil {
			return errCode
		}
		return wrapper.InternalServerError(err)
	}

	return nil
}
