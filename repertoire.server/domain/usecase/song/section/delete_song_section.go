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
	songSectionRepository repository.SongSectionRepository
	songRepository        repository.SongRepository
}

func NewDeleteSongSection(
	songSectionRepository repository.SongSectionRepository,
	songRepository repository.SongRepository,
) DeleteSongSection {
	return DeleteSongSection{
		songSectionRepository: songSectionRepository,
		songRepository:        songRepository,
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
	if index == -1 {
		return wrapper.NotFoundError(errors.New("song section not found"))
	}

	// reorder the other sections
	sectionsLength := len(song.Sections)
	for i := index + 1; i < sectionsLength; i++ {
		song.Sections[i].Order = song.Sections[i].Order - 1
	}

	// update song's new confidence, rehearsals and progress medians
	if sectionsLength == 1 {
		song.Confidence = 0
		song.Rehearsals = 0
		song.Progress = 0
	} else {
		song.Confidence = (song.Confidence*float64(sectionsLength) - float64(song.Sections[index].Confidence)) / float64(sectionsLength-1)
		song.Rehearsals = (song.Rehearsals*float64(sectionsLength) - float64(song.Sections[index].Rehearsals)) / float64(sectionsLength-1)
		song.Progress = (song.Progress*float64(sectionsLength) - float64(song.Sections[index].Progress)) / float64(sectionsLength-1)
	}

	err = d.songRepository.UpdateWithAssociations(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	err = d.songSectionRepository.Delete([]uuid.UUID{id})
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
