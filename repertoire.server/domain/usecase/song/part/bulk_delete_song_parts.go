package part

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type BulkDeleteSongParts struct {
	transactionManager transaction.Manager
	songProcessor      processor.SongProcessor

	txSongRepository        repository.SongRepository
	txSongSectionRepository repository.SongSectionRepository
}

func NewBulkDeleteSongParts(
	transactionManager transaction.Manager,
	songProcessor processor.SongProcessor,
) BulkDeleteSongParts {
	return BulkDeleteSongParts{
		transactionManager: transactionManager,
		songProcessor:      songProcessor,
	}
}

func (b BulkDeleteSongParts) Handle(request requests.BulkDeleteSongPartsRequest) *wrapper.ErrorCode {
	var errCode *wrapper.ErrorCode
	err := b.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		b.txSongSectionRepository = factory.NewSongSectionRepository()
		txSongRepository := factory.NewSongRepository()
		txSongPartRepository := factory.NewSongPartRepository()

		idsSet := make(map[uuid.UUID]bool, len(request.IDs))
		for _, id := range request.IDs {
			idsSet[id] = true
		}

		errCode = b.songProcessor.UpdateSongAfterPartsDeletion(txSongRepository, request.SongID, request.IDs)
		if errCode != nil {
			return errCode.Error
		}
		if errCode = b.reorderSectionParts(idsSet, request); errCode != nil {
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

func (b BulkDeleteSongParts) reorderSectionParts(
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
