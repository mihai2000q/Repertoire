package song

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/wrapper"
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

func (a AddPerfectSongRehearsal) Handle(request requests.AddPerfectSongRehearsalRequest) *wrapper.ErrorCode {
	var song model.Song
	err := a.repository.GetWithSections(&song, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	var errCode *wrapper.ErrorCode
	err = a.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		transactionSongRepository := factory.NewSongRepository()

		errC, isUpdated := a.songProcessor.AddPerfectRehearsal(&song, transactionSongRepository)
		if errC != nil {
			errCode = errC
			return errCode.Error
		}

		if isUpdated {
			err := transactionSongRepository.UpdateWithAssociations(&song)
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
