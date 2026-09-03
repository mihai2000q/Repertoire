package section

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/deduplicate"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type UpdateSongSection struct {
	songSectionRepository repository.SongSectionRepository
	songPartRepository    repository.SongPartRepository
	progressProcessor     processor.ProgressProcessor
	transactionManager    transaction.Manager
}

func NewUpdateSongSection(
	songSectionRepository repository.SongSectionRepository,
	songPartRepository repository.SongPartRepository,
	progressProcessor processor.ProgressProcessor,
	transactionManager transaction.Manager,
) UpdateSongSection {
	return UpdateSongSection{
		songSectionRepository: songSectionRepository,
		songPartRepository:    songPartRepository,
		progressProcessor:     progressProcessor,
		transactionManager:    transactionManager,
	}
}

func (u UpdateSongSection) Handle(request requests.UpdateSongSectionRequest) *httperror.ErrorCode {
	var section model.SongSection
	if err := u.songSectionRepository.GetWithSectionParts(&section, request.ID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(section).IsZero() {
		return httperror.NotFoundError(errors.New("song section not found"))
	}

	if len(request.PartIDs) > 0 {
		errCode := u.ensurePartsBelongToSameSong(request, section.SongID)
		if errCode != nil {
			return errCode
		}
	}

	err := u.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txSongSectionRepo := factory.NewSongSectionRepository()

		section.Name = request.Name
		section.SongSectionTypeID = request.TypeID
		if err := txSongSectionRepo.Update(&section); err != nil {
			return err
		}

		if err := u.updateSectionParts(txSongSectionRepo, &section, request.PartIDs); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}

func (u UpdateSongSection) ensurePartsBelongToSameSong(
	request requests.UpdateSongSectionRequest,
	songID uuid.UUID,
) *httperror.ErrorCode {
	var parts []model.SongPart
	if err := u.songPartRepository.GetAllByIDs(&parts, request.PartIDs); err != nil {
		return httperror.DatabaseError(err)
	}
	partIDSet := deduplicate.Deduplicate(request.PartIDs)
	if len(parts) != len(partIDSet) {
		return httperror.NotFoundError(errors.New("some parts not found"))
	}
	for _, p := range parts {
		if p.SongID != songID {
			return httperror.ConflictError(errors.New("song part does not belong to the same song as the section"))
		}
	}
	return nil
}

func (u UpdateSongSection) updateSectionParts(
	txSongSectionRepo repository.SongSectionRepository,
	section *model.SongSection,
	partIDs []uuid.UUID,
) error {
	// Build maps for lookup
	oldParts := make(map[uuid.UUID]model.SongSectionPart)
	for _, sp := range section.SectionParts {
		oldParts[sp.PartID] = sp
	}
	newParts := deduplicate.Deduplicate(partIDs)

	// Identify parts to delete: those in oldMap but not in newSet
	var partsToDelete []model.SongSectionPart
	for pid, sp := range oldParts {
		if !newParts[pid] {
			partsToDelete = append(partsToDelete, sp)
		}
	}
	if len(partsToDelete) > 0 {
		if err := txSongSectionRepo.DeleteSectionParts(&partsToDelete); err != nil {
			return err
		}
	}

	// Identify parts to update (existing) and parts to create (new)
	var partsToUpdate []model.SongSectionPart
	var partsToCreate []model.SongSectionPart
	order := 0
	for pid := range newParts {
		if sp, exists := oldParts[pid]; exists {
			sp.Order = uint(order)
			partsToUpdate = append(partsToUpdate, sp)
		} else {
			partsToCreate = append(partsToCreate, model.SongSectionPart{
				PartID:    pid,
				SectionID: section.ID,
				Order:     uint(order),
			})
		}
		order++
	}

	if len(partsToUpdate) > 0 {
		if err := txSongSectionRepo.UpdateAllSectionParts(&partsToUpdate); err != nil {
			return err
		}
	}

	if len(partsToCreate) > 0 {
		if err := txSongSectionRepo.CreateAllSectionParts(&partsToCreate); err != nil {
			return err
		}
	}

	return nil
}
