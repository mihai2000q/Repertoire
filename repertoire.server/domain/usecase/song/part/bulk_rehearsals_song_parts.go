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
	"slices"
	"time"

	"github.com/google/uuid"
)

type BulkRehearsalsSongParts struct {
	songRepository     repository.SongRepository
	transactionManager transaction.Manager
	progressProcessor  processor.ProgressProcessor
}

func NewBulkRehearsalsSongParts(
	songRepository repository.SongRepository,
	transactionManager transaction.Manager,
	progressProcessor processor.ProgressProcessor,
) BulkRehearsalsSongParts {
	return BulkRehearsalsSongParts{
		songRepository:     songRepository,
		transactionManager: transactionManager,
		progressProcessor:  progressProcessor,
	}
}

func (b BulkRehearsalsSongParts) Handle(request requests.BulkRehearsalsSongPartsRequest) *wrapper.ErrorCode {
	var song model.Song
	err := b.songRepository.GetWithParts(&song, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	// check whether all parts can be found
	partsFound := 0
	for _, part := range song.Parts {
		ind := slices.IndexFunc(request.Parts, func(sec requests.BulkRehearsalsSongPartRequest) bool {
			return sec.ID == part.ID
		})
		if ind == -1 {
			continue
		}
		partsFound++
	}

	if partsFound != len(request.Parts) {
		return wrapper.NotFoundError(errors.New("song parts not found"))
	}

	var errCode *wrapper.ErrorCode
	err = b.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		transactionSongPartRepository := factory.NewSongPartRepository()
		transactionSongRepository := factory.NewSongRepository()

		totalOldRehearsals := uint(0)
		totalNewRehearsals := uint(0)
		totalOldProgress := uint64(0)
		totalNewProgress := uint64(0)
		for i, part := range song.Parts {
			ind := slices.IndexFunc(request.Parts, func(sec requests.BulkRehearsalsSongPartRequest) bool {
				return sec.ID == part.ID
			})
			if ind == -1 || request.Parts[ind].Rehearsals == 0 {
				continue
			}
			oldProgress := part.Progress
			oldRehearsals := part.Rehearsals
			newRehearsals := part.Rehearsals + request.Parts[ind].Rehearsals

			// add history of the rehearsals change
			newHistory := model.SongPartHistory{
				ID:         uuid.New(),
				Property:   model.RehearsalsProperty,
				From:       oldRehearsals,
				To:         newRehearsals,
				SongPartID: part.ID,
			}
			err = transactionSongPartRepository.CreateHistory(&newHistory)
			if err != nil {
				errCode = wrapper.InternalServerError(err)
				return err
			}

			// update part's rehearsals score based on the history changes
			var history []model.SongPartHistory
			err = transactionSongPartRepository.GetHistory(&history, part.ID, model.RehearsalsProperty)
			if err != nil {
				errCode = wrapper.InternalServerError(err)
				return err
			}
			song.Parts[i].RehearsalsScore = b.progressProcessor.ComputeRehearsalsScore(history)

			// update part's progress (depends on the rehearsals score)
			newProgress := b.progressProcessor.ComputeProgress(part)
			song.Parts[i].Progress = newProgress

			song.Parts[i].Rehearsals = newRehearsals
			totalOldRehearsals += oldRehearsals
			totalNewRehearsals += newRehearsals
			totalOldProgress += oldProgress
			totalNewProgress += newProgress
		}

		// means that no part got updated (because if it did, the total would be at least 1)
		if totalNewRehearsals == 0 {
			return nil
		}

		// update song's new rehearsals and progress medians
		partsLength := len(song.Parts)
		song.Rehearsals =
			(song.Rehearsals*float64(partsLength) + float64(totalNewRehearsals) - float64(totalOldRehearsals)) /
				float64(partsLength)
		song.Progress =
			(song.Progress*float64(partsLength) + float64(totalNewProgress) - float64(totalOldProgress)) /
				float64(partsLength)
		song.LastTimePlayed = &[]time.Time{time.Now().UTC()}[0]

		err = transactionSongRepository.UpdateWithAssociations(&song)
		if err != nil {
			errCode = wrapper.InternalServerError(err)
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
