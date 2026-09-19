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

		// map for easy lookup
		idsMap := make(map[uuid.UUID]bool)
		for _, iD := range request.IDs {
			idsMap[iD] = true
		}

		// Reorder the remaining sections
		sectionsFound := 0
		for i := range song.Sections {
			if idsMap[song.Sections[i].ID] {
				sectionsFound++
				continue
			}
			song.Sections[i].Order = song.Sections[i].Order - uint(sectionsFound)
		}

		// Validate all requested section IDs were found
		if sectionsFound != len(request.IDs) {
			errCode = httperror.NotFoundError(errors.New("song sections not found"))
			return errCode.Error
		}

		if err := b.txSongRepo.UpdateWithAssociations(&song); err != nil {
			return err
		}

		if request.WithParts {
			b.txSongPartRepo = factory.NewSongPartRepository()
			if errCode = b.deleteParts(request.IDs); errCode != nil {
				return errCode.Error
			}
		}

		if err := b.txSongSectionRepo.Delete(request.IDs); err != nil {
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

func (b BulkDeleteSongSections) deleteParts(ids []uuid.UUID) *httperror.ErrorCode {
	var sections []model.SongSection
	if err := b.txSongSectionRepo.GetAllByIDsWithSectionParts(&sections, ids); err != nil {
		return httperror.DatabaseError(err)
	}

	// Collect all part IDs from all sections (deduplicate)
	var partIDsToDelete []uuid.UUID
	partSet := make(map[uuid.UUID]bool)
	partsCount := 0
	for _, sec := range sections {
		for _, sp := range sec.SectionParts {
			if !partSet[sp.PartID] {
				partSet[sp.PartID] = true
				partIDsToDelete = append(partIDsToDelete, sp.PartID)
				partsCount++
			}
		}
	}

	if partsCount == 0 {
		return nil
	}

	errCode := b.songProcessor.UpdateSongAfterPartsDeletion(b.txSongRepo, sections[0].SongID, partIDsToDelete)
	if errCode != nil {
		return errCode
	}

	if err := b.txSongPartRepo.Delete(partIDsToDelete); err != nil {
		return httperror.DatabaseError(err)
	}
	return nil
}
