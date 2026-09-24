package section

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type BulkDeleteSongSections struct {
	transactionManager transaction.Manager
	songProcessor      processor.SongProcessor

	txSongRepo        repository.SongRepository
	txSongSectionRepo repository.SongSectionRepository
	txSongPartRepo    repository.SongPartRepository
}

func NewBulkDeleteSongSections(
	transactionManager transaction.Manager,
	songProcessor processor.SongProcessor,
) BulkDeleteSongSections {
	return BulkDeleteSongSections{
		transactionManager: transactionManager,
		songProcessor:      songProcessor,
	}
}

func (b BulkDeleteSongSections) Handle(request requests.BulkDeleteSongSectionsRequest) *httperror.ErrorCode {
	var errCode *httperror.ErrorCode
	err := b.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		b.txSongRepo = factory.NewSongRepository()
		b.txSongSectionRepo = factory.NewSongSectionRepository()

		var song model.Song
		if err := b.txSongRepo.GetWithSections(&song, request.SongID); err != nil {
			return err
		}
		if reflect.ValueOf(song).IsZero() {
			errCode = httperror.NotFoundError(errors.New("song not found"))
			return errCode.Error
		}

		idsMap := make(map[uuid.UUID]bool, len(request.IDs))
		for _, id := range request.IDs {
			idsMap[id] = true
		}

		// Reorder the remaining sections
		sectionsFound := 0
		for i := range song.Sections {
			if idsMap[song.Sections[i].ID] {
				sectionsFound++
				continue
			}
			song.Sections[i].Order -= uint(sectionsFound)
		}

		// Validate all requested section IDs were found
		if sectionsFound != len(request.IDs) {
			errCode = httperror.NotFoundError(errors.New("song sections not found"))
			return errCode.Error
		}
		// Validate parts
		if len(request.PartIDs) > 0 {
			if errCode = b.validatePartsBelongToSections(request.IDs, request.PartIDs); errCode != nil {
				return errCode.Error
			}
		}

		// Update sections
		if err := b.txSongRepo.UpdateWithAssociations(&song); err != nil {
			return err
		}

		if len(request.PartIDs) > 0 {
			b.txSongPartRepo = factory.NewSongPartRepository()
			// Update song
			errCode = b.songProcessor.UpdateSongAfterPartsDeletion(b.txSongRepo, request.SongID, request.PartIDs)
			if errCode != nil {
				return errCode.Error
			}
			// Delete parts
			if err := b.txSongPartRepo.Delete(request.PartIDs); err != nil {
				return err
			}
		}

		// Delete sections
		return b.txSongSectionRepo.Delete(request.IDs)
	})
	if err != nil {
		if errCode != nil {
			return errCode
		}
		return httperror.DatabaseError(err)
	}
	return nil
}

func (b BulkDeleteSongSections) validatePartsBelongToSections(
	sectionIDs []uuid.UUID,
	partIDs []uuid.UUID,
) *httperror.ErrorCode {
	var sections []model.SongSection
	if err := b.txSongSectionRepo.GetAllByIDsWithSectionParts(&sections, sectionIDs); err != nil {
		return httperror.DatabaseError(err)
	}

	partsInSections := make(map[uuid.UUID]bool)
	for _, sec := range sections {
		for _, sp := range sec.SectionParts {
			partsInSections[sp.PartID] = true
		}
	}

	seen := make(map[uuid.UUID]bool, len(partIDs))
	for _, partID := range partIDs {
		if !partsInSections[partID] {
			return httperror.ConflictError(errors.New("song parts don't belong to section"))
		}
		seen[partID] = true
	}

	return nil
}
