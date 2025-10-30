package section

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"

	"github.com/google/uuid"
)

type BulkDeleteSongSections struct {
	songRepository repository.SongRepository
}

func NewBulkDeleteSongSections(repository repository.SongRepository) BulkDeleteSongSections {
	return BulkDeleteSongSections{
		songRepository: repository,
	}
}

func (b BulkDeleteSongSections) Handle(request requests.BulkDeleteSongSectionsRequest) *wrapper.ErrorCode {
	var song model.Song
	err := b.songRepository.GetWithSections(&song, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	// reorder the other sections and gather total values from deleted sections
	sectionsFound := uint(0)
	totalConfidence := uint(0)
	totalRehearsals := uint(0)
	totalProgress := uint64(0)
	for i, section := range song.Sections {
		if slices.ContainsFunc(request.IDs, func(id uuid.UUID) bool {
			return id == section.ID
		}) {
			sectionsFound++
			totalConfidence += section.Confidence
			totalRehearsals += section.Rehearsals
			totalProgress += section.Progress
			continue
		}
		song.Sections[i].Order = song.Sections[i].Order - sectionsFound
	}

	if sectionsFound == 0 {
		return wrapper.NotFoundError(errors.New("song sections not found"))
	}

	// update song's new confidence, rehearsals and progress medians
	sectionsLength := len(song.Sections)
	sectionsDeletedLength := len(request.IDs)
	if sectionsLength == sectionsDeletedLength {
		song.Confidence = 0
		song.Rehearsals = 0
		song.Progress = 0
	} else {
		song.Confidence = (song.Confidence*float64(sectionsLength) - float64(totalConfidence)) / float64(sectionsLength-sectionsDeletedLength)
		song.Rehearsals = (song.Rehearsals*float64(sectionsLength) - float64(totalRehearsals)) / float64(sectionsLength-sectionsDeletedLength)
		song.Progress = (song.Progress*float64(sectionsLength) - float64(totalProgress)) / float64(sectionsLength-sectionsDeletedLength)
	}

	err = b.songRepository.UpdateWithAssociations(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	err = b.songRepository.DeleteSections(request.IDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
