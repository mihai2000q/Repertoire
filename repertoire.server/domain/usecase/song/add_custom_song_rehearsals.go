package song

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type AddCustomSongRehearsals struct {
	songRepository     repository.SongRepository
	songProcessor      processor.SongProcessor
	transactionManager transaction.Manager
}

func NewAddCustomSongRehearsals(
	songRepository repository.SongRepository,
	songProcessor processor.SongProcessor,
	transactionManager transaction.Manager,
) AddCustomSongRehearsals {
	return AddCustomSongRehearsals{
		songRepository:     songRepository,
		songProcessor:      songProcessor,
		transactionManager: transactionManager,
	}
}

func (a AddCustomSongRehearsals) Handle(request requests.AddCustomSongRehearsalsRequest) *httperror.ErrorCode {
	var ids []uuid.UUID
	idsCount := 0
	for _, r := range request.Requests {
		ids = append(ids, r.ID)
		idsCount++
	}

	var songs []model.Song
	err := a.songRepository.GetAllByIDsWithPartsAndArrangementOccurrences(&songs, ids)
	if err != nil {
		return httperror.DatabaseError(err)
	}
	if len(songs) != idsCount {
		return httperror.NotFoundError(errors.New("songs not found"))
	}

	// use a map in order to not go through the songs, but through the ids that were sent
	// so the order is preserved
	// bonus, it handles duplicates this way too
	songsMap := make(map[uuid.UUID]model.Song)
	for _, song := range songs {
		songsMap[song.ID] = song
	}

	var errCode *httperror.ErrorCode
	err = a.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txSongPartRepo := factory.NewSongPartRepository()
		txSongRepo := factory.NewSongRepository()

		var newSongs []model.Song
		for _, r := range request.Requests {
			song, ok := songsMap[r.ID]
			if !ok {
				errCode = httperror.NotFoundError(errors.New("songs not found"))
				return errCode.Error
			}

			errC, isUpdated := a.songProcessor.AddCustomRehearsal(&song, txSongPartRepo, &r.ArrangementID)
			if errC != nil {
				errCode = errC
				return errCode.Error
			}
			if isUpdated {
				newSongs = append(newSongs, song)
			}
		}

		if len(newSongs) > 0 {
			if err := txSongRepo.UpdateAllWithAssociations(&newSongs); err != nil {
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
