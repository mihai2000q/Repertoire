package arrangement

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"

	"github.com/google/uuid"
)

type DeleteSongArrangement struct {
	songArrangementRepository repository.SongArrangementRepository
	songRepository            repository.SongRepository
}

func NewDeleteSongArrangement(
	songArrangementRepository repository.SongArrangementRepository,
	songRepository repository.SongRepository,
) DeleteSongArrangement {
	return DeleteSongArrangement{
		songArrangementRepository: songArrangementRepository,
		songRepository:            songRepository,
	}
}

func (d DeleteSongArrangement) Handle(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	var song model.Song
	err := d.songRepository.GetWithArrangements(&song, songID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	index := slices.IndexFunc(song.Arrangements, func(a model.SongArrangement) bool {
		return a.ID == id
	})
	if index == -1 {
		return wrapper.NotFoundError(errors.New("song arrangement not found"))
	}

	// reorder the other arrangements
	for i := index + 1; i < len(song.Arrangements); i++ {
		song.Arrangements[i].Order = song.Arrangements[i].Order - 1
	}

	err = d.songRepository.UpdateWithAssociations(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	err = d.songArrangementRepository.Delete(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
