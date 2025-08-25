package playlist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/wrapper"
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

func (a AddPerfectRehearsalsToPlaylists) Handle(request requests.AddPerfectRehearsalsToPlaylistsRequest) *wrapper.ErrorCode {
	var playlists []model.Playlist
	err := a.repository.GetAllByIDsWithSongSections(&playlists, request.IDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if len(playlists) == 0 {
		return wrapper.NotFoundError(errors.New("playlists not found"))
	}

	var errCode *wrapper.ErrorCode
	err = a.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		transactionSongRepository := factory.NewSongRepository()

		var newSongs []model.Song
		for _, playlist := range playlists {
			for _, playlistSong := range playlist.PlaylistSongs {
				errC, isUpdated := a.songProcessor.AddPerfectRehearsal(&playlistSong.Song, transactionSongRepository)
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
			err = transactionSongRepository.UpdateAllWithAssociations(&newSongs)
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
