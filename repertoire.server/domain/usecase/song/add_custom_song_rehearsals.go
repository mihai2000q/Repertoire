package song

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type AddCustomSongRehearsals struct {
	repository         repository.SongRepository
	songProcessor      processor.SongProcessor
	transactionManager transaction.Manager
}

func NewAddCustomSongRehearsals(
	repository repository.SongRepository,
	songProcessor processor.SongProcessor,
	transactionManager transaction.Manager,
) AddCustomSongRehearsals {
	return AddCustomSongRehearsals{
		repository:         repository,
		songProcessor:      songProcessor,
		transactionManager: transactionManager,
	}
}

func (a AddCustomSongRehearsals) Handle(request requests.AddCustomSongRehearsalsRequest) *wrapper.ErrorCode {
	var ids []uuid.UUID
	for _, r := range request.Requests {
		ids = append(ids, r.ID)
	}

	var songs []model.Song
	err := a.repository.GetAllByIDsWithPartsAndArrangementOccurrences(&songs, ids)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if len(songs) == 0 {
		return wrapper.NotFoundError(errors.New("songs not found"))
	}

	// use a map in order to not go through the songs, but through the ids that were sent
	// so the order is preserved
	// bonus, it handles duplicates this way too
	songsMap := make(map[uuid.UUID]model.Song)
	for _, song := range songs {
		songsMap[song.ID] = song
	}

	var errCode *wrapper.ErrorCode
	err = a.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txSongPartRepository := factory.NewSongPartRepository()
		transactionSongRepository := factory.NewSongRepository()

		var newSongs []model.Song
		for _, r := range request.Requests {
			song, ok := songsMap[r.ID]
			if !ok {
				errCode = wrapper.NotFoundError(errors.New("songs not found"))
				return errCode.Error
			}

			errC, isUpdated := a.songProcessor.AddCustomRehearsal(&song, txSongPartRepository, &r.ArrangementID)
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
