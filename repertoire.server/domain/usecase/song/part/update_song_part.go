package part

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"time"

	"github.com/google/uuid"
)

type UpdateSongPart struct {
	songPartRepository    repository.SongPartRepository
	songRepository        repository.SongRepository
	songSectionRepository repository.SongSectionRepository
	progressProcessor     processor.ProgressProcessor
	transactionManager    transaction.Manager

	txSongPartRepository repository.SongPartRepository
	txSongRepository     repository.SongRepository
}

func NewUpdateSongPart(
	songPartRepository repository.SongPartRepository,
	songRepository repository.SongRepository,
	songSectionRepository repository.SongSectionRepository,
	progressProcessor processor.ProgressProcessor,
	transactionManager transaction.Manager,
) UpdateSongPart {
	return UpdateSongPart{
		songPartRepository:    songPartRepository,
		songRepository:        songRepository,
		songSectionRepository: songSectionRepository,
		progressProcessor:     progressProcessor,
		transactionManager:    transactionManager,
	}
}

func (u UpdateSongPart) Handle(request requests.UpdateSongPartRequest) *wrapper.ErrorCode {
	var part model.SongPart
	err := u.songPartRepository.Get(&part, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(part).IsZero() {
		return wrapper.NotFoundError(errors.New("song part not found"))
	}
	if part.Rehearsals > request.Rehearsals {
		return wrapper.ConflictError(errors.New("rehearsals can only be increased"))
	}

	hasRehearsalsChanged := part.Rehearsals != request.Rehearsals
	hasConfidenceChanged := part.Confidence != request.Confidence
	hasBandMemberChanged := part.BandMemberID != nil && request.BandMemberID == nil ||
		part.BandMemberID == nil && request.BandMemberID != nil ||
		part.BandMemberID != nil && request.BandMemberID != nil && *part.BandMemberID != *request.BandMemberID

	if hasBandMemberChanged && request.BandMemberID != nil {
		res, err := u.songRepository.IsBandMemberAssociatedWithSong(part.SongID, *request.BandMemberID)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
		if !res {
			return wrapper.ConflictError(errors.New("band member is not part of the artist associated with this song"))
		}
	}

	var song model.Song
	var songPartsCount int64
	var section *model.SongSection
	var sectionPartsCount int64
	if hasRehearsalsChanged || hasConfidenceChanged {
		err = u.songRepository.Get(&song, part.SongID)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
		err = u.songPartRepository.CountAllBySong(&songPartsCount, part.SongID)
		if err != nil {
			return wrapper.InternalServerError(err)
		}

		if part.SectionID != nil {
			err = u.songSectionRepository.Get(section, *part.SectionID)
			if err != nil {
				return wrapper.InternalServerError(err)
			}
			err = u.songPartRepository.CountAllBySection(&sectionPartsCount, *part.SectionID)
			if err != nil {
				return wrapper.InternalServerError(err)
			}
		}
	}

	var errCode *wrapper.ErrorCode
	err = u.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		u.txSongRepository = factory.NewSongRepository()
		u.txSongPartRepository = factory.NewSongPartRepository()
		txSongSectionRepository := factory.NewSongSectionRepository()

		if hasRehearsalsChanged {
			errCode = u.rehearsalsHasChanged(
				&part,
				request.Rehearsals,
				&song,
				songPartsCount,
				section,
				sectionPartsCount,
			)
			if errCode != nil {
				return errCode.Error
			}
		}

		if hasConfidenceChanged {
			errCode = u.confidenceHasChanged(
				&part,
				request.Confidence,
				&song,
				songPartsCount,
				section,
				sectionPartsCount,
			)
			if errCode != nil {
				return errCode.Error
			}
		}

		part.Name = request.Name
		part.Confidence = request.Confidence
		part.Rehearsals = request.Rehearsals
		part.BandMemberID = request.BandMemberID
		part.InstrumentID = request.InstrumentID

		if hasRehearsalsChanged || hasConfidenceChanged {
			err = u.txSongRepository.Update(&song)
			if err != nil {
				return err
			}

			if section != nil {
				err = txSongSectionRepository.Update(section)
				if err != nil {
					return err
				}
			}
		}
		err = u.txSongPartRepository.Update(&part)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if errCode != nil {
			return errCode
		}
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (u UpdateSongPart) rehearsalsHasChanged(
	part *model.SongPart,
	newRehearsals uint,
	song *model.Song,
	songPartsCount int64,
	section *model.SongSection,
	sectionPartsCount int64,
) *wrapper.ErrorCode {
	// add history of the rehearsals change
	newHistory := model.SongPartHistory{
		ID:         uuid.New(),
		Property:   model.RehearsalsProperty,
		From:       part.Rehearsals,
		To:         newRehearsals,
		SongPartID: part.ID,
	}
	err := u.txSongPartRepository.CreateHistory(&newHistory)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	// remove part's rehearsals and progress from the song's median
	song.Rehearsals = song.Rehearsals*float64(songPartsCount) - float64(part.Rehearsals)
	song.Progress = song.Progress*float64(songPartsCount) - float64(part.Progress)

	// update part's rehearsals score based on the history changes
	var history []model.SongPartHistory
	err = u.txSongPartRepository.GetHistory(&history, part.ID, model.RehearsalsProperty)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	part.RehearsalsScore = u.progressProcessor.ComputeRehearsalsScore(history)

	// update part's progress (depends on the rehearsals score)
	part.Progress = u.progressProcessor.ComputeProgress(*part)

	if section != nil {
		// update the section's rehearsals and progress median with new part values
		section.Rehearsals = (section.Rehearsals + float64(newRehearsals)) / float64(sectionPartsCount)
		section.Progress = (section.Progress + float64(part.Progress)) / float64(sectionPartsCount)
	}

	// update the song's rehearsals and progress median with new part values
	song.Rehearsals = (song.Rehearsals + float64(newRehearsals)) / float64(songPartsCount)
	song.Progress = (song.Progress + float64(part.Progress)) / float64(songPartsCount)

	// update song's last time played
	song.LastTimePlayed = &[]time.Time{time.Now().UTC()}[0]

	return nil
}

func (u UpdateSongPart) confidenceHasChanged(
	part *model.SongPart,
	newConfidence uint,
	song *model.Song,
	songPartsCount int64,
	section *model.SongSection,
	sectionPartsCount int64,
) *wrapper.ErrorCode {
	// add history of the confidence change
	newHistory := model.SongPartHistory{
		ID:         uuid.New(),
		Property:   model.ConfidenceProperty,
		From:       part.Confidence,
		To:         newConfidence,
		SongPartID: part.ID,
	}
	err := u.txSongPartRepository.CreateHistory(&newHistory)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	// remove part's confidence and progress from the song's median
	song.Confidence = song.Confidence*float64(songPartsCount) - float64(part.Confidence)
	song.Progress = song.Progress*float64(songPartsCount) - float64(part.Progress)

	// update part's confidence score based on the history changes
	var history []model.SongPartHistory
	err = u.txSongPartRepository.GetHistory(&history, part.ID, model.ConfidenceProperty)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	part.ConfidenceScore = u.progressProcessor.ComputeConfidenceScore(history)

	// update part's progress (depends on the confidence score)
	part.Progress = u.progressProcessor.ComputeProgress(*part)

	if section != nil {
		// update the section's rehearsals and progress median with new part values
		section.Confidence = (section.Confidence + float64(newConfidence)) / float64(sectionPartsCount)
		section.Progress = (section.Progress + float64(part.Progress)) / float64(sectionPartsCount)
	}

	// update the song's confidence and progress median with new part values
	song.Confidence = (song.Confidence + float64(newConfidence)) / float64(songPartsCount)
	song.Progress = (song.Progress + float64(part.Progress)) / float64(songPartsCount)

	return nil
}
