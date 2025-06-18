package section

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

type UpdateSongSection struct {
	songRepository    repository.SongRepository
	progressProcessor processor.ProgressProcessor
}

func NewUpdateSongSection(
	repository repository.SongRepository,
	progressProcessor processor.ProgressProcessor,
) UpdateSongSection {
	return UpdateSongSection{
		songRepository:    repository,
		progressProcessor: progressProcessor,
	}
}

func (u UpdateSongSection) Handle(request requests.UpdateSongSectionRequest) *wrapper.ErrorCode {
	var section model.SongSection
	err := u.songRepository.GetSection(&section, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(section).IsZero() {
		return wrapper.NotFoundError(errors.New("song section not found"))
	}
	if section.Rehearsals > request.Rehearsals {
		return wrapper.ConflictError(errors.New("rehearsals can only be increased"))
	}

	hasRehearsalsChanged := section.Rehearsals != request.Rehearsals
	hasConfidenceChanged := section.Confidence != request.Confidence
	hasBandMemberChanged := section.BandMemberID != nil && request.BandMemberID == nil ||
		section.BandMemberID == nil && request.BandMemberID != nil ||
		section.BandMemberID != nil && request.BandMemberID != nil && *section.BandMemberID != *request.BandMemberID

	var song model.Song
	var sectionsCount int64
	if hasRehearsalsChanged || hasConfidenceChanged {
		err = u.songRepository.Get(&song, section.SongID)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
		err = u.songRepository.CountSectionsBySong(&sectionsCount, section.SongID)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
	}

	if hasBandMemberChanged && request.BandMemberID != nil {
		res, err := u.songRepository.IsBandMemberAssociatedWithSong(section.SongID, *request.BandMemberID)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
		if !res {
			return wrapper.ConflictError(errors.New("band member is not part of the artist associated with this song"))
		}
	}

	if hasRehearsalsChanged {
		errCode := u.rehearsalsHasChanged(&section, request.Rehearsals, sectionsCount, &song)
		if errCode != nil {
			return errCode
		}
	}

	if hasConfidenceChanged {
		errCode := u.confidenceHasChanged(&section, request.Confidence, sectionsCount, &song)
		if errCode != nil {
			return errCode
		}
	}

	section.Name = request.Name
	section.Confidence = request.Confidence
	section.Rehearsals = request.Rehearsals
	section.SongSectionTypeID = request.TypeID
	section.BandMemberID = request.BandMemberID
	section.InstrumentID = request.InstrumentID

	if hasRehearsalsChanged || hasConfidenceChanged {
		err = u.songRepository.Update(&song)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
	}
	err = u.songRepository.UpdateSection(&section)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (u UpdateSongSection) rehearsalsHasChanged(
	section *model.SongSection,
	newRehearsals uint,
	sectionsCount int64,
	song *model.Song,
) *wrapper.ErrorCode {
	// add history of the rehearsals change
	newHistory := model.SongSectionHistory{
		ID:            uuid.New(),
		Property:      model.RehearsalsProperty,
		From:          section.Rehearsals,
		To:            newRehearsals,
		SongSectionID: section.ID,
	}
	err := u.songRepository.CreateSongSectionHistory(&newHistory)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	// remove section's rehearsals and progress from the song's median
	song.Rehearsals = song.Rehearsals*float64(sectionsCount) - float64(section.Rehearsals)
	song.Progress = song.Progress*float64(sectionsCount) - float64(section.Progress)

	// update section's rehearsals score based on the history changes
	var history []model.SongSectionHistory
	err = u.songRepository.GetSongSectionHistory(&history, section.ID, model.RehearsalsProperty)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	section.RehearsalsScore = u.progressProcessor.ComputeRehearsalsScore(history)

	// update section's progress (dependent on the rehearsals score)
	section.Progress = u.progressProcessor.ComputeProgress(*section)

	// update the song's rehearsals and progress median with new section values
	song.Rehearsals = (song.Rehearsals + float64(newRehearsals)) / float64(sectionsCount)
	song.Progress = (song.Progress + float64(section.Progress)) / float64(sectionsCount)

	// update song's last time played
	song.LastTimePlayed = &[]time.Time{time.Now().UTC()}[0]

	return nil
}

func (u UpdateSongSection) confidenceHasChanged(
	section *model.SongSection,
	newConfidence uint,
	sectionsCount int64,
	song *model.Song,
) *wrapper.ErrorCode {
	// add history of the confidence change
	newHistory := model.SongSectionHistory{
		ID:            uuid.New(),
		Property:      model.ConfidenceProperty,
		From:          section.Confidence,
		To:            newConfidence,
		SongSectionID: section.ID,
	}
	err := u.songRepository.CreateSongSectionHistory(&newHistory)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	// remove section's confidence and progress from the song's median
	song.Confidence = song.Confidence*float64(sectionsCount) - float64(section.Confidence)
	song.Progress = song.Progress*float64(sectionsCount) - float64(section.Progress)

	// update section's confidence score based on the history changes
	var history []model.SongSectionHistory
	err = u.songRepository.GetSongSectionHistory(&history, section.ID, model.ConfidenceProperty)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	section.ConfidenceScore = u.progressProcessor.ComputeConfidenceScore(history)

	// update section's progress (dependent on the confidence score)
	section.Progress = u.progressProcessor.ComputeProgress(*section)

	// update the song's confidence and progress median with new section values
	song.Confidence = (song.Confidence + float64(newConfidence)) / float64(sectionsCount)
	song.Progress = (song.Progress + float64(section.Progress)) / float64(sectionsCount)

	return nil
}
