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

type AddPerfectSongRehearsals struct {
	repository         repository.SongRepository
	songProcessor      processor.SongProcessor
	transactionManager transaction.Manager
}

func NewAddPerfectSongRehearsals(
	repository repository.SongRepository,
	songProcessor processor.SongProcessor,
	transactionManager transaction.Manager,
) AddPerfectSongRehearsals {
	return AddPerfectSongRehearsals{
		repository:         repository,
		songProcessor:      songProcessor,
		transactionManager: transactionManager,
	}
}

func (a AddPerfectSongRehearsals) Handle(request requests.AddPerfectSongRehearsalsRequest) *wrapper.ErrorCode {
	var songs []model.Song
	err := a.repository.GetAllByIDsWithSections(&songs, request.IDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if len(songs) == 0 {
		return wrapper.NotFoundError(errors.New("songs not found"))
	}

	var errCode *wrapper.ErrorCode
	err = a.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		transactionSongRepository := factory.NewSongRepository()

		var newSongs []model.Song
		for _, song := range songs {
			errC, isUpdated := a.songProcessor.AddPerfectRehearsal(&song, transactionSongRepository)
			if errC != nil {
				errCode = errC
				return errCode.Error
			}
			if isUpdated {
				newSongs = append(newSongs, song)
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
