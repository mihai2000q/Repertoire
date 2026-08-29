package part

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"

	"github.com/google/uuid"
)

type DeleteSongPart struct {
	songPartRepository    repository.SongPartRepository
	songRepository        repository.SongRepository
	songSectionRepository repository.SongSectionRepository
}

func NewDeleteSongPart(
	songPartRepository repository.SongPartRepository,
	songRepository repository.SongRepository,
	songSectionRepository repository.SongSectionRepository,
) DeleteSongPart {
	return DeleteSongPart{
		songPartRepository:    songPartRepository,
		songRepository:        songRepository,
		songSectionRepository: songSectionRepository,
	}
}

func (d DeleteSongPart) Handle(id uuid.UUID, songID uuid.UUID, sectionID uuid.UUID) *wrapper.ErrorCode {
	var song model.Song
	errCode := d.updateSong(&song, id, songID)
	if errCode != nil {
		return errCode
	}
	var section model.SongSection
	if sectionID != uuid.Nil {
		errCode = d.updateSection(&section, id, sectionID)
		if errCode != nil {
			return errCode
		}
	}

	err := d.songRepository.UpdateWithAssociations(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if sectionID != uuid.Nil {
		err = d.songSectionRepository.UpdateWithAssociations(&section)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
	}

	err = d.songPartRepository.Delete([]uuid.UUID{id})
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (d DeleteSongPart) updateSong(song *model.Song, id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	err := d.songRepository.GetWithParts(song, songID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	index := slices.IndexFunc(song.Parts, func(a model.SongPart) bool {
		return a.ID == id
	})
	if index == -1 {
		return wrapper.NotFoundError(errors.New("song part not found"))
	}

	// reorder the other parts in song
	partsLength := len(song.Parts)
	for i := index + 1; i < partsLength; i++ {
		song.Parts[i].SongOrder = song.Parts[i].SongOrder - 1
	}

	// update song's new confidence, rehearsals and progress medians
	if partsLength == 1 {
		song.Confidence = 0
		song.Rehearsals = 0
		song.Progress = 0
	} else {
		song.Confidence = (song.Confidence*float64(partsLength) - float64(song.Parts[index].Confidence)) / float64(partsLength-1)
		song.Rehearsals = (song.Rehearsals*float64(partsLength) - float64(song.Parts[index].Rehearsals)) / float64(partsLength-1)
		song.Progress = (song.Progress*float64(partsLength) - float64(song.Parts[index].Progress)) / float64(partsLength-1)
	}

	return nil
}

func (d DeleteSongPart) updateSection(section *model.SongSection, id uuid.UUID, sectionID uuid.UUID) *wrapper.ErrorCode {
	err := d.songSectionRepository.GetWithParts(section, sectionID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(section).IsZero() {
		return wrapper.NotFoundError(errors.New("song section not found"))
	}

	index := slices.IndexFunc(section.Parts, func(a model.SongPart) bool {
		return a.ID == id
	})
	if index == -1 {
		return wrapper.NotFoundError(errors.New("song part not found"))
	}

	// reorder the other parts in section
	partsLength := len(section.Parts)
	for i := index + 1; i < partsLength; i++ {
		section.Parts[i].SongOrder = section.Parts[i].SongOrder - 1
	}

	// update section's new confidence, rehearsals and progress medians
	if partsLength == 1 {
		section.Confidence = 0
		section.Rehearsals = 0
		section.Progress = 0
	} else {
		section.Confidence = (section.Confidence*float64(partsLength) - float64(section.Parts[index].Confidence)) / float64(partsLength-1)
		section.Rehearsals = (section.Rehearsals*float64(partsLength) - float64(section.Parts[index].Rehearsals)) / float64(partsLength-1)
		section.Progress = (section.Progress*float64(partsLength) - float64(section.Parts[index].Progress)) / float64(partsLength-1)
	}

	return nil
}
