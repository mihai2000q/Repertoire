package section

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

type BulkRehearsalsSongSections struct {
	songRepository     repository.SongRepository
	transactionManager transaction.Manager
	progressProcessor  processor.ProgressProcessor
}

func NewBulkRehearsalsSongSections(
	songRepository repository.SongRepository,
	transactionManager transaction.Manager,
	progressProcessor processor.ProgressProcessor,
) BulkRehearsalsSongSections {
	return BulkRehearsalsSongSections{
		songRepository:     songRepository,
		transactionManager: transactionManager,
		progressProcessor:  progressProcessor,
	}
}

func (b BulkRehearsalsSongSections) Handle(request requests.BulkRehearsalsSongSectionsRequest) *wrapper.ErrorCode {
	var song model.Song
	err := b.songRepository.GetWithSections(&song, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	// check whether all sections can be found
	sectionsFound := 0
	for _, section := range song.Sections {
		ind := slices.IndexFunc(request.Sections, func(sec requests.BulkRehearsalsSongSectionRequest) bool {
			return sec.ID == section.ID
		})
		if ind == -1 {
			continue
		}
		sectionsFound++
	}

	if sectionsFound != len(request.Sections) {
		return wrapper.NotFoundError(errors.New("song sections not found"))
	}

	var errCode *wrapper.ErrorCode
	err = b.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		transactionSongSectionRepository := factory.NewSongSectionRepository()
		transactionSongRepository := factory.NewSongRepository()

		totalOldRehearsals := uint(0)
		totalNewRehearsals := uint(0)
		totalOldProgress := uint64(0)
		totalNewProgress := uint64(0)
		for i, section := range song.Sections {
			ind := slices.IndexFunc(request.Sections, func(sec requests.BulkRehearsalsSongSectionRequest) bool {
				return sec.ID == section.ID
			})
			if ind == -1 || request.Sections[ind].Rehearsals == 0 {
				continue
			}
			oldProgress := section.Progress
			oldRehearsals := section.Rehearsals
			newRehearsals := section.Rehearsals + request.Sections[ind].Rehearsals

			// add history of the rehearsals change
			newHistory := model.SongSectionHistory{
				ID:            uuid.New(),
				Property:      model.RehearsalsProperty,
				From:          oldRehearsals,
				To:            newRehearsals,
				SongSectionID: section.ID,
			}
			err = transactionSongSectionRepository.CreateHistory(&newHistory)
			if err != nil {
				errCode = wrapper.InternalServerError(err)
				return err
			}

			// update section's rehearsals score based on the history changes
			var history []model.SongSectionHistory
			err = transactionSongSectionRepository.GetHistory(&history, section.ID, model.RehearsalsProperty)
			if err != nil {
				errCode = wrapper.InternalServerError(err)
				return err
			}
			song.Sections[i].RehearsalsScore = b.progressProcessor.ComputeRehearsalsScore(history)

			// update section's progress (depends on the rehearsals score)
			newProgress := b.progressProcessor.ComputeProgress(section)
			song.Sections[i].Progress = newProgress

			song.Sections[i].Rehearsals = newRehearsals
			totalOldRehearsals += oldRehearsals
			totalNewRehearsals += newRehearsals
			totalOldProgress += oldProgress
			totalNewProgress += newProgress
		}

		// means that no section got updated (because if it did, the total would be at least 1)
		if totalNewRehearsals == 0 {
			return nil
		}

		// update song's new rehearsals and progress medians
		sectionsLength := len(song.Sections)
		song.Rehearsals =
			(song.Rehearsals*float64(sectionsLength) + float64(totalNewRehearsals) - float64(totalOldRehearsals)) /
				float64(sectionsLength)
		song.Progress =
			(song.Progress*float64(sectionsLength) + float64(totalNewProgress) - float64(totalOldProgress)) /
				float64(sectionsLength)
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
