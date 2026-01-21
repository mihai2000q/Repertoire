package section

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreateSongSection struct {
	songSectionRepository repository.SongSectionRepository
	songRepository        repository.SongRepository
}

func NewCreateSongSection(
	songSectionRepository repository.SongSectionRepository,
	songRepository repository.SongRepository,
) CreateSongSection {
	return CreateSongSection{
		songSectionRepository: songSectionRepository,
		songRepository:        songRepository,
	}
}

func (c CreateSongSection) Handle(request requests.CreateSongSectionRequest) *wrapper.ErrorCode {
	var sectionsCount int64
	err := c.songSectionRepository.CountAllBySong(&sectionsCount, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	var song model.Song
	err = c.songRepository.Get(&song, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}
	if request.BandMemberID != nil {
		res, err := c.songRepository.IsBandMemberAssociatedWithSong(request.SongID, *request.BandMemberID)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
		if !res {
			return wrapper.ConflictError(errors.New("band member is not part of the artist associated with this song"))
		}
	}

	section := model.SongSection{
		ID:                uuid.New(),
		Name:              request.Name,
		Confidence:        model.DefaultSongSectionConfidence,
		SongSectionTypeID: request.TypeID,
		Order:             uint(sectionsCount),
		SongID:            request.SongID,
		BandMemberID:      request.BandMemberID,
		InstrumentID:      request.InstrumentID,
	}
	err = c.songSectionRepository.Create(&section)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	// update song's new confidence, rehearsals and progress medians
	song.Confidence = (song.Confidence*float64(sectionsCount) + float64(section.Confidence)) / float64(sectionsCount+1)
	song.Rehearsals = (song.Rehearsals*float64(sectionsCount) + float64(section.Rehearsals)) / float64(sectionsCount+1)
	song.Progress = (song.Progress*float64(sectionsCount) + float64(section.Progress)) / float64(sectionsCount+1)

	err = c.songRepository.Update(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
