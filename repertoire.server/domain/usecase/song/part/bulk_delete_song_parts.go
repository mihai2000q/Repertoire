package part

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type BulkDeleteSongParts struct {
	transactionManager transaction.Manager

	txSongRepository        repository.SongRepository
	txSongSectionRepository repository.SongSectionRepository
}

func NewBulkDeleteSongParts(
	transactionManager transaction.Manager,
) BulkDeleteSongParts {
	return BulkDeleteSongParts{
		transactionManager: transactionManager,
	}
}

func (b BulkDeleteSongParts) Handle(request requests.BulkDeleteSongPartsRequest) *wrapper.ErrorCode {
	var errCode *wrapper.ErrorCode
	err := b.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		b.txSongRepository = factory.NewSongRepository()
		b.txSongSectionRepository = factory.NewSongSectionRepository()
		txSongPartRepository := factory.NewSongPartRepository()

		idsSet := make(map[uuid.UUID]bool, len(request.IDs))
		for _, id := range request.IDs {
			idsSet[id] = true
		}

		if errCode = b.updateSong(idsSet, request); errCode != nil {
			return errCode.Error
		}
		if errCode = b.reorderSections(idsSet, request); errCode != nil {
			return errCode.Error
		}
		if err := txSongPartRepository.Delete(request.IDs); err != nil {
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

func (b BulkDeleteSongParts) updateSong(
	idsSet map[uuid.UUID]bool,
	request requests.BulkDeleteSongPartsRequest,
) *wrapper.ErrorCode {
	var song model.Song
	if err := b.txSongRepository.GetWithParts(&song, request.SongID); err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	var remainingParts []model.SongPart
	var totalConfidence, totalRehearsals uint
	var totalProgress uint64
	var partsFound uint

	for _, part := range song.Parts {
		if idsSet[part.ID] {
			partsFound++
			totalConfidence += part.Confidence
			totalRehearsals += part.Rehearsals
			totalProgress += part.Progress
			continue
		}
		part.SongOrder -= partsFound
		remainingParts = append(remainingParts, part)
	}

	if int(partsFound) != len(request.IDs) {
		return wrapper.NotFoundError(errors.New("song parts not found"))
	}

	song.Parts = remainingParts

	partsLength := len(song.Parts) + int(partsFound)
	deletedCount := int(partsFound)
	if partsLength == deletedCount {
		song.Confidence = 0
		song.Rehearsals = 0
		song.Progress = 0
	} else {
		newPartsLength := float64(partsLength - deletedCount)
		song.Confidence = (song.Confidence*float64(partsLength) - float64(totalConfidence)) / newPartsLength
		song.Rehearsals = (song.Rehearsals*float64(partsLength) - float64(totalRehearsals)) / newPartsLength
		song.Progress = (song.Progress*float64(partsLength) - float64(totalProgress)) / newPartsLength
	}

	if err := b.txSongRepository.UpdateWithAssociations(&song); err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}

func (b BulkDeleteSongParts) reorderSections(
	idsSet map[uuid.UUID]bool,
	request requests.BulkDeleteSongPartsRequest,
) *wrapper.ErrorCode {
	var sections []model.SongSection
	err := b.txSongSectionRepository.GetAllByPartIDsWithSectionParts(&sections, request.IDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if len(sections) == 0 {
		return nil
	}

	var sectionPartsToUpdate []model.SongSectionPart
	for _, section := range sections {
		var shift uint
		for _, sp := range section.SectionParts {
			if idsSet[sp.PartID] {
				shift++ // this entry will be deleted; shift subsequent ones down
				continue
			}
			sp.Order -= shift
			sectionPartsToUpdate = append(sectionPartsToUpdate, sp)
		}
	}

	if len(sectionPartsToUpdate) > 0 {
		if err = b.txSongSectionRepository.UpdateAllSectionParts(&sectionPartsToUpdate); err != nil {
			return wrapper.InternalServerError(err)
		}
	}

	return nil
}
