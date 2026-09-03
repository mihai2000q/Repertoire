package album

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type AddPerfectRehearsalsToAlbums struct {
	repository         repository.AlbumRepository
	songProcessor      processor.SongProcessor
	transactionManager transaction.Manager
}

func NewAddPerfectRehearsalsToAlbums(
	repository repository.AlbumRepository,
	songProcessor processor.SongProcessor,
	transactionManager transaction.Manager,
) AddPerfectRehearsalsToAlbums {
	return AddPerfectRehearsalsToAlbums{
		repository:         repository,
		songProcessor:      songProcessor,
		transactionManager: transactionManager,
	}
}

func (a AddPerfectRehearsalsToAlbums) Handle(request requests.AddPerfectRehearsalsToAlbumsRequest) *httperror.ErrorCode {
	var albums []model.Album
	err := a.repository.GetAllByIDsWithSongPartsAndDefaultOccurrences(&albums, request.IDs)
	if err != nil {
		return httperror.DatabaseError(err)
	}
	if len(albums) != len(request.IDs) {
		return httperror.NotFoundError(errors.New("albums not found"))
	}

	var errCode *httperror.ErrorCode
	err = a.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txSongPartRepository := factory.NewSongPartRepository()
		txSongRepository := factory.NewSongRepository()

		var newSongs []model.Song
		for _, album := range albums {
			for _, song := range album.Songs {
				errC, isUpdated := a.songProcessor.AddPerfectRehearsal(&song, txSongPartRepository)
				if errC != nil {
					errCode = errC
					return errCode.Error
				}
				if isUpdated {
					newSongs = append(newSongs, song)
				}
			}
		}

		if len(newSongs) > 0 {
			if err = txSongRepository.UpdateAllWithAssociations(&newSongs); err != nil {
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
