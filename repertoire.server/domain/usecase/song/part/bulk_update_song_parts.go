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

type BulkUpdateSongParts struct {
	transactionManager transaction.Manager
	progressProcessor  processor.ProgressProcessor
	txSongPartRepo     repository.SongPartRepository
}

func NewBulkUpdateSongParts(
	transactionManager transaction.Manager,
	progressProcessor processor.ProgressProcessor,
) BulkUpdateSongParts {
	return BulkUpdateSongParts{
		transactionManager: transactionManager,
		progressProcessor:  progressProcessor,
	}
}

func (b BulkUpdateSongParts) Handle(request requests.BulkUpdateSongPartsRequest) *httperror.ErrorCode {
	var errCode *httperror.ErrorCode
	err := b.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		b.txSongPartRepo = factory.NewSongPartRepository()
		txSongRepo := factory.NewSongRepository()

		var totalOldRehearsals, totalNewRehearsals uint
		var totalOldConfidence, totalNewConfidence uint
		var totalOldProgress, totalNewProgress uint64
		var partsToUpdate []model.SongPart
		rehearsalsChanged := false
		partsFound := 0

		var song model.Song
		if err := txSongRepo.GetWithParts(&song, request.SongID); err != nil {
			return err
		}
		if reflect.ValueOf(song).IsZero() {
			errCode = httperror.NotFoundError(errors.New("song not found"))
			return errCode.Error
		}

		// map for easy lookup
		requestMap := make(map[uuid.UUID]requests.BulkUpdateSongPartRequest, len(request.Requests))
		for _, req := range request.Requests {
			requestMap[req.ID] = req
		}

		// validate non-decrease rehearsals on parts
		for _, part := range song.Parts {
			req, ok := requestMap[part.ID]
			if !ok {
				continue
			}
			partsFound++
			if req.Rehearsals < part.Rehearsals {
				errCode = httperror.ConflictError(errors.New("rehearsals can only be increased"))
				return errCode.Error
			}
		}

		// validate parts' existence
		if partsFound != len(request.Requests) {
			errCode = httperror.NotFoundError(errors.New("song parts not found"))
			return errCode.Error
		}

		for i, part := range song.Parts {
			req, ok := requestMap[part.ID]
			if !ok {
				continue // part not requested, skip
			}

			hasRehearsalsChanged := req.Rehearsals != part.Rehearsals
			hasConfidenceChanged := req.Confidence != part.Confidence

			if !hasRehearsalsChanged && !hasConfidenceChanged {
				continue
			}

			if hasRehearsalsChanged {
				rehearsalsChanged = true
			}

			oldPart := model.SongPart{
				Rehearsals: part.Rehearsals,
				Confidence: part.Confidence,
				Progress:   part.Progress,
			}

			if hasRehearsalsChanged {
				if errCode = b.updateRehearsals(&part, req.Rehearsals); errCode != nil {
					return errCode.Error
				}
			}

			if hasConfidenceChanged {
				if errCode = b.updateConfidence(&part, req.Confidence); errCode != nil {
					return errCode.Error
				}
			}

			// Compute new progress
			part.Progress = b.progressProcessor.ComputeProgress(part)

			// Accumulate totals
			totalOldRehearsals += oldPart.Rehearsals
			totalOldConfidence += oldPart.Confidence
			totalOldProgress += oldPart.Progress
			totalNewRehearsals += part.Rehearsals
			totalNewConfidence += part.Confidence
			totalNewProgress += part.Progress

			song.Parts[i] = part
			partsToUpdate = append(partsToUpdate, part)
		}

		if len(partsToUpdate) == 0 {
			return nil
		}

		// Update song medians
		partsLength := len(song.Parts)
		song.Confidence = (song.Confidence*float64(partsLength) + float64(totalNewConfidence) - float64(totalOldConfidence)) / float64(partsLength)
		song.Rehearsals = (song.Rehearsals*float64(partsLength) + float64(totalNewRehearsals) - float64(totalOldRehearsals)) / float64(partsLength)
		song.Progress = (song.Progress*float64(partsLength) + float64(totalNewProgress) - float64(totalOldProgress)) / float64(partsLength)

		if rehearsalsChanged {
			song.LastTimePlayed = &[]time.Time{time.Now().UTC()}[0]
		}

		if err := txSongRepo.UpdateWithAssociations(&song); err != nil {
			return err
		}

		if err := b.txSongPartRepo.UpdateAll(&partsToUpdate); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errCode != nil {
			return errCode
		}
		return httperror.DatabaseError(err)
	}

	return nil
}

func (b BulkUpdateSongParts) updateRehearsals(part *model.SongPart, newRehearsals uint) *httperror.ErrorCode {
	// add history of the rehearsals change
	newHistory := model.SongPartHistory{
		ID:       uuid.New(),
		Property: model.RehearsalsProperty,
		From:     part.Rehearsals,
		To:       newRehearsals,
		PartID:   part.ID,
	}
	if err := b.txSongPartRepo.CreateHistory(&newHistory); err != nil {
		return httperror.DatabaseError(err)
	}

	// update part's rehearsals score based on the history changes
	var history []model.SongPartHistory
	if err := b.txSongPartRepo.GetHistory(&history, part.ID, model.RehearsalsProperty); err != nil {
		return httperror.DatabaseError(err)
	}
	part.RehearsalsScore = b.progressProcessor.ComputeRehearsalsScore(history)
	part.Rehearsals = newRehearsals

	return nil
}

func (b BulkUpdateSongParts) updateConfidence(part *model.SongPart, newConfidence uint) *httperror.ErrorCode {
	// add history of the confidence change
	newHistory := model.SongPartHistory{
		ID:       uuid.New(),
		Property: model.ConfidenceProperty,
		From:     part.Confidence,
		To:       newConfidence,
		PartID:   part.ID,
	}
	if err := b.txSongPartRepo.CreateHistory(&newHistory); err != nil {
		return httperror.DatabaseError(err)
	}

	// update part's confidence score based on the history changes
	var history []model.SongPartHistory
	if err := b.txSongPartRepo.GetHistory(&history, part.ID, model.ConfidenceProperty); err != nil {
		return httperror.DatabaseError(err)
	}
	part.ConfidenceScore = b.progressProcessor.ComputeConfidenceScore(history)
	part.Confidence = newConfidence

	return nil
}
