package processor

import (
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"time"

	"github.com/google/uuid"
)

type SongProcessor interface {
	AddPerfectRehearsal(
		song *model.Song,
		repository repository.SongRepository,
	) (errCode *wrapper.ErrorCode, updatedSong bool)
}

type songProcessor struct {
	progressProcessor ProgressProcessor
}

func NewSongProcessor(progressProcessor ProgressProcessor) SongProcessor {
	return &songProcessor{progressProcessor: progressProcessor}
}

func (s *songProcessor) AddPerfectRehearsal(
	song *model.Song,
	repository repository.SongRepository,
) (*wrapper.ErrorCode, bool) {
	var totalRehearsals float64 = 0
	var totalProgress float64 = 0
	for i, section := range song.Sections {
		if section.Occurrences == 0 {
			continue
		}

		newRehearsals := section.Rehearsals + section.Occurrences
		// add history of the rehearsals change
		newHistory := model.SongSectionHistory{
			ID:            uuid.New(),
			Property:      model.RehearsalsProperty,
			From:          section.Rehearsals,
			To:            newRehearsals,
			SongSectionID: section.ID,
		}
		err := repository.CreateSongSectionHistory(&newHistory)
		if err != nil {
			return wrapper.InternalServerError(err), false
		}

		// update section's rehearsals score based on the history changes and update the rehearsals and progress too
		var history []model.SongSectionHistory
		err = repository.GetSongSectionHistory(&history, section.ID, model.RehearsalsProperty)
		if err != nil {
			return wrapper.InternalServerError(err), false
		}

		song.Sections[i].Rehearsals = newRehearsals
		song.Sections[i].RehearsalsScore = s.progressProcessor.ComputeRehearsalsScore(history)
		song.Sections[i].Progress = s.progressProcessor.ComputeProgress(song.Sections[i])

		// add to the total for the median
		totalProgress += float64(song.Sections[i].Progress)
		totalRehearsals += float64(song.Sections[i].Rehearsals)
	}

	// means that no section got updated (because if it did, the total would be at least 1 from an occurrence)
	if totalRehearsals == 0 {
		return nil, false
	}

	// update song media progress and rehearsals + update last time played
	sectionsCount := len(song.Sections)
	song.Rehearsals = totalRehearsals / float64(sectionsCount)
	song.Progress = totalProgress / float64(sectionsCount)
	song.LastTimePlayed = &[]time.Time{time.Now().UTC()}[0]

	return nil, true
}
