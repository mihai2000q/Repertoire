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

func (a AddCustomSongRehearsal) Handle(request requests.AddCustomSongRehearsalRequest) *wrapper.ErrorCode {
	var song model.Song
	err := a.repository.GetWithSectionsAndArrangementOccurrences(&song, request.ID, request.ArrangementID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}
	if len(song.Sections[0].ArrangementOccurrences) == 0 {
		return wrapper.NotFoundError(errors.New("song arrangement not found"))
	}

	var errCode *wrapper.ErrorCode
	err = a.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		transactionSongSectionRepository := factory.NewSongSectionRepository()
		transactionSongRepository := factory.NewSongRepository()

		errC, isUpdated := a.songProcessor.AddCustomRehearsal(&song, transactionSongSectionRepository)
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
