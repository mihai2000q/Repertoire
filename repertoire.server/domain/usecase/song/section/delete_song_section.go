package section

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"

	"github.com/google/uuid"
)

type DeleteSongSection struct {
	songRepository repository.SongRepository
}

func NewDeleteSongSection(repository repository.SongRepository) DeleteSongSection {
	return DeleteSongSection{
		songRepository: repository,
	}
}

func (d DeleteSongSection) Handle(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	var song model.Song
	err := d.songRepository.GetWithSections(&song, songID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	index := slices.IndexFunc(song.Sections, func(a model.SongSection) bool {
		return a.ID == id
	})

	for i := index + 1; i < len(song.Sections); i++ {
		song.Sections[i].Order = song.Sections[i].Order - 1
	}

	err = d.songRepository.UpdateWithAssociations(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	err = d.songRepository.DeleteSection(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
