package song

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/deduplicate"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type AddPerfectSongRehearsals struct {
	songRepository     repository.SongRepository
	songProcessor      processor.SongProcessor
	transactionManager transaction.Manager
}

func NewAddPerfectSongRehearsals(
	songRepository repository.SongRepository,
	songProcessor processor.SongProcessor,
	transactionManager transaction.Manager,
) AddPerfectSongRehearsals {
	return AddPerfectSongRehearsals{
		songRepository:     songRepository,
		songProcessor:      songProcessor,
		transactionManager: transactionManager,
	}
}

func (a AddPerfectSongRehearsals) Handle(request requests.AddPerfectSongRehearsalsRequest) *httperror.ErrorCode {
	var songs []model.Song
	err := a.songRepository.GetAllByIDsWithPartsAndDefaultOccurrences(&songs, request.IDs)
	if err != nil {
		return httperror.DatabaseError(err)
	}
	ids := deduplicate.Deduplicate(request.IDs)
	if len(songs) != len(ids) {
		return httperror.NotFoundError(errors.New("songs not found"))
	}

	var errCode *httperror.ErrorCode
	err = a.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txSongPartRepo := factory.NewSongPartRepository()
		txSongRepo := factory.NewSongRepository()

		var newSongs []model.Song
		for _, song := range songs {
			errC, isUpdated := a.songProcessor.AddPerfectRehearsal(&song, txSongPartRepo)
			if errC != nil {
				errCode = errC
				return errCode.Error
			}
			if isUpdated {
				newSongs = append(newSongs, song)
			}
		}

		if len(newSongs) > 0 {
			err = txSongRepo.UpdateAllWithAssociations(&newSongs)
			if err != nil {
				errCode = httperror.DatabaseError(err)
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
