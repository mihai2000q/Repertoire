package part

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

type BulkDeleteSongParts struct {
	songPartRepository repository.SongPartRepository
	songRepository     repository.SongRepository
}

func NewBulkDeleteSongParts(
	songPartRepository repository.SongPartRepository,
	songRepository repository.SongRepository,
) BulkDeleteSongParts {
	return BulkDeleteSongParts{
		songPartRepository: songPartRepository,
		songRepository:     songRepository,
	}
}

func (b BulkDeleteSongParts) Handle(request requests.BulkDeleteSongPartsRequest) *wrapper.ErrorCode {
	var song model.Song
	err := b.songRepository.GetWithParts(&song, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	// reorder the other parts and gather total values from deleted parts
	partsFound := uint(0)
	totalConfidence := uint(0)
	totalRehearsals := uint(0)
	totalProgress := uint64(0)
	for i, part := range song.Parts {
		if slices.ContainsFunc(request.IDs, func(id uuid.UUID) bool {
			return id == part.ID
		}) {
			partsFound++
			totalConfidence += part.Confidence
			totalRehearsals += part.Rehearsals
			totalProgress += part.Progress
			continue
		}
		song.Parts[i].SongOrder = song.Parts[i].SongOrder - partsFound
	}

	if int(partsFound) != len(request.IDs) {
		return wrapper.NotFoundError(errors.New("song parts not found"))
	}

	// update song's new confidence, rehearsals and progress medians
	partsLength := len(song.Parts)
	partsDeletedLength := len(request.IDs)
	if partsLength == partsDeletedLength {
		song.Confidence = 0
		song.Rehearsals = 0
		song.Progress = 0
	} else {
		song.Confidence = (song.Confidence*float64(partsLength) - float64(totalConfidence)) / float64(partsLength-partsDeletedLength)
		song.Rehearsals = (song.Rehearsals*float64(partsLength) - float64(totalRehearsals)) / float64(partsLength-partsDeletedLength)
		song.Progress = (song.Progress*float64(partsLength) - float64(totalProgress)) / float64(partsLength-partsDeletedLength)
	}

	err = b.songRepository.UpdateWithAssociations(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	err = b.songPartRepository.Delete(request.IDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
