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

type AddPerfectSongRehearsal struct {
	repository         repository.SongRepository
	songProcessor      processor.SongProcessor
	transactionManager transaction.Manager
}

func NewAddPerfectSongRehearsal(
	repository repository.SongRepository,
	songProcessor processor.SongProcessor,
	transactionManager transaction.Manager,
) AddPerfectSongRehearsal {
	return AddPerfectSongRehearsal{
		repository:         repository,
		songProcessor:      songProcessor,
		transactionManager: transactionManager,
	}
}

func (a AddPerfectSongRehearsal) Handle(request requests.AddPerfectSongRehearsalRequest) *httperror.ErrorCode {
	var song model.Song
	err := a.repository.GetWithPartsAndDefaultOccurrences(&song, request.ID)
	if err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}
	if song.DefaultArrangementID == nil {
		return httperror.BadRequestError(errors.New("song has no default arrangement set"))
	}

	var errCode *httperror.ErrorCode
	err = a.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txSongPartRepository := factory.NewSongPartRepository()
		txSongRepository := factory.NewSongRepository()

		errC, isUpdated := a.songProcessor.AddPerfectRehearsal(&song, txSongPartRepository)
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
