package processor

import (
	"errors"
	"repertoire/server/data/repository"
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
			ID:         uuid.New(),
			Property:   model.RehearsalsProperty,
			From:       part.Rehearsals,
			To:         newRehearsals,
			SongPartID: part.ID,
			CreatedAt:  time.Now().UTC(),
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
