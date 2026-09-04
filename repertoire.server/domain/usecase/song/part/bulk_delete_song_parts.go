package part

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type BulkDeleteSongParts struct {
	transactionManager transaction.Manager
	songProcessor      processor.SongProcessor

	txSongRepo        repository.SongRepository
	txSongSectionRepo repository.SongSectionRepository
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

func (b BulkDeleteSongParts) Handle(request requests.BulkDeleteSongPartsRequest) *httperror.ErrorCode {
	var errCode *httperror.ErrorCode
	err := b.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		b.txSongSectionRepo = factory.NewSongSectionRepository()
		txSongRepo := factory.NewSongRepository()
		txSongPartRepo := factory.NewSongPartRepository()

		errCode = b.songProcessor.UpdateSongAfterPartsDeletion(txSongRepo, request.SongID, request.IDs)
		if errCode != nil {
			return errCode.Error
		}
		if errCode = b.reorderSectionParts(request); errCode != nil {
			return errCode.Error
		}
		if err := txSongPartRepo.Delete(request.IDs); err != nil {
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

func (b BulkDeleteSongParts) reorderSectionParts(
	request requests.BulkDeleteSongPartsRequest,
) *httperror.ErrorCode {
	var sections []model.SongSection
	err := b.txSongSectionRepo.GetAllByPartIDsWithSectionParts(&sections, request.IDs)
	if err != nil {
		return httperror.DatabaseError(err)
	}
	if len(sections) == 0 {
		return nil
	}

	// map for easy lookup
	idsMap := make(map[uuid.UUID]bool, len(request.IDs))
	for _, id := range request.IDs {
		idsMap[id] = true
	}

	var sectionPartsToUpdate []model.SongSectionPart
	for _, section := range sections {
		var shift uint
		for _, sp := range section.SectionParts {
			if idsMap[sp.PartID] {
				shift++ // this entry will be deleted; shift subsequent ones down
				continue
			}
			sp.Order -= shift
			sectionPartsToUpdate = append(sectionPartsToUpdate, sp)
		}
	}

	if len(sectionPartsToUpdate) > 0 {
		if err = b.txSongSectionRepo.UpdateAllSectionParts(&sectionPartsToUpdate); err != nil {
			return httperror.DatabaseError(err)
		}
	}

	return nil
}
