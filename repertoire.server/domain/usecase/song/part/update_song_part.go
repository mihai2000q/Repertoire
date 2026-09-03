package part

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
	"time"

	"github.com/google/uuid"
)

type UpdateSongPart struct {
	songPartRepository repository.SongPartRepository
	songRepository     repository.SongRepository
	progressProcessor  processor.ProgressProcessor
	transactionManager transaction.Manager

	txSongPartRepository repository.SongPartRepository
	txSongRepository     repository.SongRepository
}

func NewUpdateSongPart(
	songPartRepository repository.SongPartRepository,
	songRepository repository.SongRepository,
	progressProcessor processor.ProgressProcessor,
	transactionManager transaction.Manager,
) UpdateSongPart {
	return UpdateSongPart{
		songPartRepository: songPartRepository,
		songRepository:     songRepository,
		progressProcessor:  progressProcessor,
		transactionManager: transactionManager,
	}
}

func (u UpdateSongPart) Handle(request requests.UpdateSongPartRequest) *httperror.ErrorCode {
	var part model.SongPart
	if err := u.songPartRepository.Get(&part, request.ID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(part).IsZero() {
		return httperror.NotFoundError(errors.New("song part not found"))
	}
	if part.Rehearsals > request.Rehearsals {
		return httperror.ConflictError(errors.New("rehearsals can only be increased"))
	}

	hasRehearsalsChanged := part.Rehearsals != request.Rehearsals
	hasConfidenceChanged := part.Confidence != request.Confidence
	hasBandMemberChanged := part.BandMemberID != nil && request.BandMemberID == nil ||
		part.BandMemberID == nil && request.BandMemberID != nil ||
		part.BandMemberID != nil && request.BandMemberID != nil && *part.BandMemberID != *request.BandMemberID

	if hasBandMemberChanged && request.BandMemberID != nil {
		res, err := u.songRepository.IsBandMemberAssociatedWithSong(part.SongID, *request.BandMemberID)
		if err != nil {
			return httperror.DatabaseError(err)
		}
		if !res {
			return httperror.ConflictError(errors.New("band member is not part of the artist associated with this song"))
		}
	}

	err := u.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		u.txSongRepository = factory.NewSongRepository()
		u.txSongPartRepository = factory.NewSongPartRepository()

		// Store old stats before modifications
		oldPart := model.SongPart{
			Rehearsals: part.Rehearsals,
			Confidence: part.Confidence,
			Progress:   part.Progress,
		}

		// update parts' fields
		part.Name = request.Name
		part.BandMemberID = request.BandMemberID
		part.InstrumentID = request.InstrumentID

		// complex update of rehearsals and/or confidence
		if hasRehearsalsChanged {
			if err := u.updateRehearsals(&part, request.Rehearsals); err != nil {
				return err
			}
		}

		if hasConfidenceChanged {
			if err := u.updateConfidence(&part, request.Confidence); err != nil {
				return err
			}
		}

		// compute new progress and update song's stats
		if hasRehearsalsChanged || hasConfidenceChanged {
			part.Progress = u.progressProcessor.ComputeProgress(part)

			if err := u.updateSongStats(oldPart, part); err != nil {
				return err
			}
		}

		// finally update part
		if err := u.txSongPartRepository.Update(&part); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}

func (u UpdateSongPart) updateRehearsals(
	part *model.SongPart,
	newRehearsals uint,
) error {
	// add history of the rehearsals change
	newHistory := model.SongPartHistory{
		ID:       uuid.New(),
		Property: model.RehearsalsProperty,
		From:     part.Rehearsals,
		To:       newRehearsals,
		PartID:   part.ID,
	}
	if err := u.txSongPartRepository.CreateHistory(&newHistory); err != nil {
		return err
	}

	// update part's rehearsals score based on the history changes
	var history []model.SongPartHistory
	if err := u.txSongPartRepository.GetHistory(&history, part.ID, model.RehearsalsProperty); err != nil {
		return err
	}
	part.RehearsalsScore = u.progressProcessor.ComputeRehearsalsScore(history)
	part.Rehearsals = newRehearsals

	return nil
}

func (u UpdateSongPart) updateConfidence(
	part *model.SongPart,
	newConfidence uint,
) error {
	// add history of the confidence change
	newHistory := model.SongPartHistory{
		ID:       uuid.New(),
		Property: model.ConfidenceProperty,
		From:     part.Confidence,
		To:       newConfidence,
		PartID:   part.ID,
	}
	if err := u.txSongPartRepository.CreateHistory(&newHistory); err != nil {
		return err
	}

	// update part's confidence score based on the history changes
	var history []model.SongPartHistory
	if err := u.txSongPartRepository.GetHistory(&history, part.ID, model.ConfidenceProperty); err != nil {
		return err
	}
	part.ConfidenceScore = u.progressProcessor.ComputeConfidenceScore(history)
	part.Confidence = newConfidence

	return nil
}

func (u UpdateSongPart) updateSongStats(oldPart model.SongPart, newPart model.SongPart) error {
	// fetch song and the number of parts in it
	var song model.Song
	var songPartsCount int64
	if err := u.txSongRepository.Get(&song, newPart.SongID); err != nil {
		return err
	}
	if err := u.txSongPartRepository.CountAllBySong(&songPartsCount, newPart.SongID); err != nil {
		return err
	}

	// recalculate song medians using old and new values
	totalConfidence := song.Confidence*float64(songPartsCount) - float64(oldPart.Confidence) + float64(newPart.Confidence)
	totalRehearsals := song.Rehearsals*float64(songPartsCount) - float64(oldPart.Rehearsals) + float64(newPart.Rehearsals)
	totalProgress := song.Progress*float64(songPartsCount) - float64(oldPart.Progress) + float64(newPart.Progress)

	song.Confidence = totalConfidence / float64(songPartsCount)
	song.Rehearsals = totalRehearsals / float64(songPartsCount)
	song.Progress = totalProgress / float64(songPartsCount)

	// update song's last time played
	if oldPart.Rehearsals != newPart.Rehearsals {
		song.LastTimePlayed = &[]time.Time{time.Now().UTC()}[0]
	}

	// update song
	err := u.txSongRepository.Update(&song)
	if err != nil {
		return err
	}

	return nil
}
