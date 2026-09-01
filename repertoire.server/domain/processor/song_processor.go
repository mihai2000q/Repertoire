package processor

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/internal/deduplicate"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"
	"time"

	"github.com/google/uuid"
)

type SongProcessor interface {
	AddCustomRehearsal(
		song *model.Song,
		songPartRepository repository.SongPartRepository,
		arrangementID *uuid.UUID,
	) (errCode *wrapper.ErrorCode, updatedSong bool)
	AddPerfectRehearsal(
		song *model.Song,
		songPartRepository repository.SongPartRepository,
	) (errCode *wrapper.ErrorCode, updatedSong bool)
	UpdateSongAfterPartsDeletion(
		songRepository repository.SongRepository,
		songID uuid.UUID,
		partIDs []uuid.UUID,
	) *wrapper.ErrorCode
}

type songProcessor struct {
	progressProcessor ProgressProcessor
}

func NewSongProcessor(progressProcessor ProgressProcessor) SongProcessor {
	return &songProcessor{progressProcessor: progressProcessor}
}

func (s *songProcessor) AddCustomRehearsal(
	song *model.Song,
	songPartRepository repository.SongPartRepository,
	arrangementID *uuid.UUID,
) (*wrapper.ErrorCode, bool) {
	if len(song.Parts) == 0 || (arrangementID == nil && len(song.Parts[0].ArrangementOccurrences) == 0) {
		return nil, false
	}
	if arrangementID != nil {
		index := slices.IndexFunc(song.Parts[0].ArrangementOccurrences, func(o model.SongPartOccurrences) bool {
			return o.ArrangementID == *arrangementID
		})
		if index == -1 {
			return wrapper.NotFoundError(errors.New("song arrangement not found")), false
		}
	}

	return s.addRehearsal(song, songPartRepository, arrangementID)
}

func (s *songProcessor) AddPerfectRehearsal(
	song *model.Song,
	songPartRepository repository.SongPartRepository,
) (*wrapper.ErrorCode, bool) {
	if song.DefaultArrangementID == nil {
		return nil, false
	}
	return s.addRehearsal(song, songPartRepository, nil)
}

func (s *songProcessor) addRehearsal(
	song *model.Song,
	songPartRepository repository.SongPartRepository,
	arrangementID *uuid.UUID,
) (*wrapper.ErrorCode, bool) {
	var totalRehearsals float64 = 0
	var totalProgress float64 = 0
	for i, part := range song.Parts {
		var arrangementOccurrence model.SongPartOccurrences
		if arrangementID != nil {
			index := slices.IndexFunc(part.ArrangementOccurrences, func(o model.SongPartOccurrences) bool {
				return o.ArrangementID == *arrangementID
			})
			arrangementOccurrence = part.ArrangementOccurrences[index]
		} else {
			arrangementOccurrence = part.ArrangementOccurrences[0]
		}

		if arrangementOccurrence.Occurrences == 0 {
			continue
		}

		newRehearsals := part.Rehearsals + arrangementOccurrence.Occurrences
		// add history of the rehearsals change
		newHistory := model.SongPartHistory{
			ID:        uuid.New(),
			Property:  model.RehearsalsProperty,
			From:      part.Rehearsals,
			To:        newRehearsals,
			PartID:    part.ID,
			CreatedAt: time.Now().UTC(),
		}
		err := songPartRepository.CreateHistory(&newHistory)
		if err != nil {
			return wrapper.InternalServerError(err), false
		}

		// update part's rehearsals score based on the history changes and update the rehearsals and progress too
		var history []model.SongPartHistory
		err = songPartRepository.GetHistory(&history, part.ID, model.RehearsalsProperty)
		if err != nil {
			return wrapper.InternalServerError(err), false
		}

		song.Parts[i].Rehearsals = newRehearsals
		song.Parts[i].RehearsalsScore = s.progressProcessor.ComputeRehearsalsScore(history)
		song.Parts[i].Progress = s.progressProcessor.ComputeProgress(song.Parts[i])

		// add to the total for the median
		totalProgress += float64(song.Parts[i].Progress)
		totalRehearsals += float64(song.Parts[i].Rehearsals)
	}

	// means that no part got updated (because if it did, the total would be at least 1 from an occurrence)
	if totalRehearsals == 0 {
		return nil, false
	}

	// update song media progress and rehearsals + update last time played
	partsCount := len(song.Parts)
	song.Rehearsals = totalRehearsals / float64(partsCount)
	song.Progress = totalProgress / float64(partsCount)
	song.LastTimePlayed = &[]time.Time{time.Now().UTC()}[0]

	return nil, true
}

func (s *songProcessor) UpdateSongAfterPartsDeletion(
	songRepository repository.SongRepository,
	songID uuid.UUID,
	partIDs []uuid.UUID,
) *wrapper.ErrorCode {
	// Fetch the song with its parts
	var song model.Song
	if err := songRepository.GetWithParts(&song, songID); err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	// Reorder remaining parts and accumulate deleted stats
	var totalConfidence, totalRehearsals uint
	var totalProgress uint64
	idsSet := deduplicate.Deduplicate(partIDs)
	partsToDeleteCount := 0
	remainingPartsCount := 0

	for i, part := range song.Parts {
		if idsSet[part.ID] {
			partsToDeleteCount++
			totalConfidence += part.Confidence
			totalRehearsals += part.Rehearsals
			totalProgress += part.Progress
			continue
		}
		// Shift SongOrder down by the number of deleted parts before it
		song.Parts[i].SongOrder -= uint(partsToDeleteCount)
		remainingPartsCount++
	}

	// Validate that all unique part IDs were found
	if partsToDeleteCount != len(idsSet) {
		return wrapper.NotFoundError(errors.New("song parts not found"))
	}

	// Recalculate song stats
	if remainingPartsCount == 0 {
		song.Confidence = 0
		song.Rehearsals = 0
		song.Progress = 0
	} else {
		oldPartsLen := remainingPartsCount + partsToDeleteCount
		newPartsLen := float64(remainingPartsCount)
		song.Confidence = (song.Confidence*float64(oldPartsLen) - float64(totalConfidence)) / newPartsLen
		song.Rehearsals = (song.Rehearsals*float64(oldPartsLen) - float64(totalRehearsals)) / newPartsLen
		song.Progress = (song.Progress*float64(oldPartsLen) - float64(totalProgress)) / newPartsLen
	}

	if err := songRepository.UpdateWithAssociations(&song); err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
