package section

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/deduplicate"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type BulkDeleteSongSections struct {
	transactionManager transaction.Manager
	songProcessor      processor.SongProcessor

	txSongRepository        repository.SongRepository
	txSongSectionRepository repository.SongSectionRepository
	txSongPartRepository    repository.SongPartRepository
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

func (b BulkDeleteSongSections) Handle(request requests.BulkDeleteSongSectionsRequest) *wrapper.ErrorCode {
	var errCode *wrapper.ErrorCode
	err := b.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		b.txSongRepository = factory.NewSongRepository()
		b.txSongSectionRepository = factory.NewSongSectionRepository()

		var song model.Song
		if err := b.txSongRepository.GetWithSections(&song, request.SongID); err != nil {
			return err
		}
		if reflect.ValueOf(song).IsZero() {
			errCode = wrapper.NotFoundError(errors.New("song not found"))
			return errCode.Error
		}

		idsSet := deduplicate.Deduplicate(request.IDs)

		// Reorder the remaining sections
		sectionsFound := uint(0)
		for i := range song.Sections {
			if idsSet[song.Sections[i].ID] {
				sectionsFound++
				continue
			}
			song.Sections[i].Order = song.Sections[i].Order - sectionsFound
		}

		// Validate all requested section IDs were found
		if int(sectionsFound) != len(idsSet) {
			errCode = wrapper.NotFoundError(errors.New("song sections not found"))
			return errCode.Error
		}

		if err := b.txSongRepository.UpdateWithAssociations(&song); err != nil {
			return err
		}

		if request.WithParts {
			b.txSongPartRepository = factory.NewSongPartRepository()
			if errCode = b.deleteParts(request.IDs); errCode != nil {
				return errCode.Error
			}
		}

		if err := b.txSongSectionRepository.Delete(request.IDs); err != nil {
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

func (b BulkDeleteSongSections) deleteParts(ids []uuid.UUID) *wrapper.ErrorCode {
	var sections []model.SongSection
	if err := b.txSongSectionRepository.GetAllByIDsWithSectionParts(&sections, ids); err != nil {
		return wrapper.InternalServerError(err)
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

	errCode := b.songProcessor.UpdateSongAfterPartsDeletion(b.txSongRepository, sections[0].SongID, partIDsToDelete)
	if errCode != nil {
		return errCode
	}

	if err := b.txSongPartRepository.Delete(partIDsToDelete); err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
