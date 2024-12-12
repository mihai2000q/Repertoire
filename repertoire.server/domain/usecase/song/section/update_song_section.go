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
		return wrapper.BadRequestError(errors.New("rehearsals can only be increased"))
	}

	hasRehearsalsChanged := section.Rehearsals != request.Rehearsals
	hasConfidenceChanged := section.Confidence != request.Confidence
	if hasRehearsalsChanged {
		errCode := u.addRehearsalsHistory(section, request.Rehearsals)
		if errCode != nil {
			return errCode
		}
		errCode = u.updateRehearsalsScore(&section)
		if errCode != nil {
			return errCode
		}
		section.Progress = uint64(section.ConfidenceScore) * section.RehearsalsScore
	}

	if hasConfidenceChanged {
		errCode := u.addConfidenceHistory(section, request.Confidence)
		if errCode != nil {
			return errCode
		}
		errCode = u.updateConfidenceScore(&section)
		if errCode != nil {
			return errCode
		}
		section.Progress = uint64(section.ConfidenceScore) * section.RehearsalsScore
	}

	section.Name = request.Name
	section.Confidence = request.Confidence
	section.Rehearsals = request.Rehearsals
	section.SongSectionTypeID = request.TypeID

	err = u.songRepository.UpdateSection(&section)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (u UpdateSongSection) addRehearsalsHistory(section model.SongSection, newRehearsals uint) *wrapper.ErrorCode {
	history := model.SongSectionHistory{
		ID:            uuid.New(),
		Property:      model.RehearsalsProperty,
		From:          section.Rehearsals,
		To:            newRehearsals,
		SongSectionID: section.ID,
	}
	err := u.songRepository.CreateSongSectionHistory(&history)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (u UpdateSongSection) addConfidenceHistory(section model.SongSection, newConfidence uint) *wrapper.ErrorCode {
	history := model.SongSectionHistory{
		ID:            uuid.New(),
		Property:      model.ConfidenceProperty,
		From:          section.Confidence,
		To:            newConfidence,
		SongSectionID: section.ID,
	}
	err := u.songRepository.CreateSongSectionHistory(&history)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (u UpdateSongSection) updateRehearsalsScore(section *model.SongSection) *wrapper.ErrorCode {
	var history []model.SongSectionHistory
	err := u.songRepository.GetSongSectionHistory(&history, section.ID, model.RehearsalsProperty)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	section.RehearsalsScore = u.progressProcessor.ComputeRehearsalsScore(history)

	return nil
}

func (u UpdateSongSection) updateConfidenceScore(section *model.SongSection) *wrapper.ErrorCode {
	var history []model.SongSectionHistory
	err := u.songRepository.GetSongSectionHistory(&history, section.ID, model.ConfidenceProperty)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	section.ConfidenceScore = u.progressProcessor.ComputeConfidenceScore(history)

	return nil
}
