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
	repository         repository.PlaylistRepository
	songProcessor      processor.SongProcessor
	transactionManager transaction.Manager
}

func NewAddPerfectRehearsalsToPlaylists(
	repository repository.PlaylistRepository,
	songProcessor processor.SongProcessor,
	transactionManager transaction.Manager,
) AddPerfectRehearsalsToPlaylists {
	return AddPerfectRehearsalsToPlaylists{
		repository:         repository,
		songProcessor:      songProcessor,
		transactionManager: transactionManager,
	}
}

func (a AddPerfectRehearsalsToPlaylists) Handle(request requests.AddPerfectRehearsalsToPlaylistsRequest) *httperror.ErrorCode {
	var playlists []model.Playlist
	err := a.repository.GetAllByIDsWithSongPartsAndDefaultOccurrences(&playlists, request.IDs)
	if err != nil {
		return httperror.DatabaseError(err)
	}
	if len(playlists) != len(request.IDs) {
		return httperror.NotFoundError(errors.New("playlists not found"))
	}

	var errCode *httperror.ErrorCode
	err = a.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txSongPartRepository := factory.NewSongPartRepository()
		txSongRepository := factory.NewSongRepository()

		var newSongs []model.Song
		for _, playlist := range playlists {
			for _, playlistSong := range playlist.PlaylistSongs {
				errC, isUpdated := a.songProcessor.AddPerfectRehearsal(&playlistSong.Song, txSongPartRepository)
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
			err = txSongRepository.UpdateAllWithAssociations(&newSongs)
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
