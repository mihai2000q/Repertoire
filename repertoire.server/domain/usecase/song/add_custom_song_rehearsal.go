package song

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type AddCustomSongRehearsal struct {
	repository         repository.SongRepository
	songProcessor      processor.SongProcessor
	transactionManager transaction.Manager
}

func NewAddCustomSongRehearsal(
	repository repository.SongRepository,
	songProcessor processor.SongProcessor,
	transactionManager transaction.Manager,
) AddCustomSongRehearsal {
	return AddCustomSongRehearsal{
		repository:         repository,
		songProcessor:      songProcessor,
		transactionManager: transactionManager,
	}
}

func (a AddCustomSongRehearsal) Handle(request requests.AddCustomSongRehearsalRequest) *httperror.ErrorCode {
	var song model.Song
	err := a.repository.GetWithPartsAndArrangementOccurrences(&song, request.ID, request.ArrangementID)
	if err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}
	if len(song.Parts[0].ArrangementOccurrences) == 0 {
		return httperror.NotFoundError(errors.New("song arrangement not found"))
	}

	var errCode *httperror.ErrorCode
	err = a.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txSongPartRepository := factory.NewSongPartRepository()
		txSongRepository := factory.NewSongRepository()

		errC, isUpdated := a.songProcessor.AddCustomRehearsal(&song, txSongPartRepository, nil)
		if errC != nil {
			errCode = errC
			return errCode.Error
		}

		if isUpdated {
			if err := txSongRepository.UpdateWithAssociations(&song); err != nil {
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
