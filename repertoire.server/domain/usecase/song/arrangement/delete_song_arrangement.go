package arrangement

import (
	"errors"
	"reflect"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
	"slices"

	"github.com/google/uuid"
)

type DeleteSongArrangement struct {
	songRepository     repository.SongRepository
	transactionManager transaction.Manager
}

func NewDeleteSongArrangement(
	songRepository repository.SongRepository,
	transactionManager transaction.Manager,
) DeleteSongArrangement {
	return DeleteSongArrangement{
		songRepository:     songRepository,
		transactionManager: transactionManager,
	}
}

func (d DeleteSongArrangement) Handle(id uuid.UUID, songID uuid.UUID) *httperror.ErrorCode {
	var song model.Song
	if err := d.songRepository.GetWithArrangements(&song, songID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}

	index := slices.IndexFunc(song.Arrangements, func(a model.SongArrangement) bool {
		return a.ID == id
	})
	if index == -1 {
		return httperror.NotFoundError(errors.New("song arrangement not found"))
	}

	// reorder the other arrangements
	for i := index + 1; i < len(song.Arrangements); i++ {
		song.Arrangements[i].Order = song.Arrangements[i].Order - 1
	}

	err := d.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txSongRepo := factory.NewSongRepository()
		txSongArrangementRepo := factory.NewSongArrangementRepository()
		if err := txSongRepo.UpdateWithAssociations(&song); err != nil {
			return err
		}
		if err := txSongArrangementRepo.Delete(id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
