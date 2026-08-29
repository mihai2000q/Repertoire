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
	songSectionRepository repository.SongSectionRepository
	songRepository        repository.SongRepository
}

func NewBulkDeleteSongSections(
	songSectionRepository repository.SongSectionRepository,
	songRepository repository.SongRepository,
) BulkDeleteSongSections {
	return BulkDeleteSongSections{
		songSectionRepository: songSectionRepository,
		songRepository:        songRepository,
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

	// reorder the other sections
	// TODO: REMOVE PARTS TOO?
	sectionsFound := uint(0)
	for i, section := range song.Sections {
		if slices.ContainsFunc(request.IDs, func(id uuid.UUID) bool {
			return id == section.ID
		}) {
			sectionsFound++
			continue
		}
		song.Sections[i].Order = song.Sections[i].Order - sectionsFound
	}

	if int(sectionsFound) != len(request.IDs) {
		return wrapper.NotFoundError(errors.New("song sections not found"))
	}

	err = b.songRepository.UpdateWithAssociations(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	err = b.songSectionRepository.Delete(request.IDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
