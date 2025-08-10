package song

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"time"

	"github.com/google/uuid"
)

type AddPartialSongRehearsal struct {
	repository        repository.SongRepository
	progressProcessor processor.ProgressProcessor
}

func NewAddPartialSongRehearsal(
	repository repository.SongRepository,
	progressProcessor processor.ProgressProcessor,
) AddPartialSongRehearsal {
	return AddPartialSongRehearsal{
		repository:        repository,
		progressProcessor: progressProcessor,
	}
}

func (a AddPartialSongRehearsal) Handle(request requests.AddPartialSongRehearsalRequest) *wrapper.ErrorCode {
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
		if section.PartialOccurrences == 0 {
			continue
		}

		newRehearsals := section.Rehearsals + section.PartialOccurrences
		// add history of the rehearsals change
		newHistory := model.SongSectionHistory{
			ID:            uuid.New(),
			Property:      model.RehearsalsProperty,
			From:          section.Rehearsals,
			To:            newRehearsals,
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

		song.Sections[i].Rehearsals = newRehearsals
		song.Sections[i].RehearsalsScore = a.progressProcessor.ComputeRehearsalsScore(history)
		song.Sections[i].Progress = a.progressProcessor.ComputeProgress(song.Sections[i])

		// add to the total for the median
		totalProgress += float64(song.Sections[i].Progress)
		totalRehearsals += float64(song.Sections[i].Rehearsals)
	}

	// means that no section got updated (because if it did, the total would be at least 1 from an occurrence)
	if totalRehearsals == 0 {
		return nil
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
