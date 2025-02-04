package song

import (
	"errors"
	"github.com/google/uuid"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"time"
)

type AddPerfectSongRehearsal struct {
	repository        repository.SongRepository
	progressProcessor processor.ProgressProcessor
}

func NewAddPerfectSongRehearsal(
	repository repository.SongRepository,
	progressProcessor processor.ProgressProcessor,
) AddPerfectSongRehearsal {
	return AddPerfectSongRehearsal{
		repository:        repository,
		progressProcessor: progressProcessor,
	}
}

func (a AddPerfectSongRehearsal) Handle(request requests.AddPerfectSongRehearsalRequest) *wrapper.ErrorCode {
	var song model.Song
	err := a.repository.GetWithSections(&song, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	var totalRehearsals float64 = 0
	var totalProgress float64 = 0
	for i, section := range song.Sections {
		// add history of the rehearsals change
		newHistory := model.SongSectionHistory{
			ID:            uuid.New(),
			Property:      model.RehearsalsProperty,
			From:          section.Rehearsals,
			To:            section.Occurrences,
			SongSectionID: section.ID,
		}
		err = a.repository.CreateSongSectionHistory(&newHistory)
		if err != nil {
			return wrapper.InternalServerError(err)
		}

		// update section's rehearsals score based on the history changes and update the rehearsals and progress too
		var history []model.SongSectionHistory
		err = a.repository.GetSongSectionHistory(&history, section.ID, model.RehearsalsProperty)
		if err != nil {
			return wrapper.InternalServerError(err)
		}

		song.Sections[i].Rehearsals += section.Occurrences
		song.Sections[i].RehearsalsScore = a.progressProcessor.ComputeRehearsalsScore(history)
		song.Sections[i].Progress = a.progressProcessor.ComputeProgress(song.Sections[i])

		// add to the total for the median
		totalProgress += float64(song.Sections[i].Progress)
		totalRehearsals += float64(song.Sections[i].Rehearsals)
	}
	
	// update song media progress and rehearsals + update last time played
	sectionsCount := len(song.Sections)
	song.Rehearsals = totalRehearsals / float64(sectionsCount)
	song.Progress = totalProgress / float64(sectionsCount)
	song.LastTimePlayed = &[]time.Time{time.Now().UTC()}[0]

	err = a.repository.UpdateWithAssociations(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
