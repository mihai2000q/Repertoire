package section

import (
	"errors"
	"reflect"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
	"slices"

	"github.com/google/uuid"
)

type DeleteSongSection struct {
	songRepository     repository.SongRepository
	transactionManager transaction.Manager
	songProcessor      processor.SongProcessor

	txSongSectionRepo repository.SongSectionRepository
	txSongPartRepo    repository.SongPartRepository
	txSongRepo        repository.SongRepository
}

func NewDeleteSongSection(
	songRepository repository.SongRepository,
	transactionManager transaction.Manager,
	songProcessor processor.SongProcessor,
) DeleteSongSection {
	return DeleteSongSection{
		songRepository:     songRepository,
		transactionManager: transactionManager,
		songProcessor:      songProcessor,
	}
}

func (d DeleteSongSection) Handle(id uuid.UUID, songID uuid.UUID, withParts bool) *httperror.ErrorCode {
	var song model.Song
	if err := d.songRepository.GetWithSections(&song, songID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}

	index := slices.IndexFunc(song.Sections, func(a model.SongSection) bool {
		return a.ID == id
	})
	if index == -1 {
		return httperror.NotFoundError(errors.New("song section not found"))
	}

	var errCode *httperror.ErrorCode
	err := d.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		d.txSongRepo = factory.NewSongRepository()
		d.txSongSectionRepo = factory.NewSongSectionRepository()

		// reorder the other sections
		sectionsLength := len(song.Sections)
		for i := index + 1; i < sectionsLength; i++ {
			song.Sections[i].Order = song.Sections[i].Order - 1
		}
		if err := d.txSongRepo.UpdateWithAssociations(&song); err != nil {
			return err
		}

		if withParts {
			d.txSongPartRepo = factory.NewSongPartRepository()
			if errCode = d.deleteParts(id); errCode != nil {
				return errCode.Error
			}
		}

		if err := d.txSongSectionRepo.Delete([]uuid.UUID{id}); err != nil {
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

func (d DeleteSongSection) deleteParts(id uuid.UUID) *httperror.ErrorCode {
	var section model.SongSection
	if err := d.txSongSectionRepo.GetWithSectionParts(&section, id); err != nil {
		return httperror.DatabaseError(err)
	}

	var partIDsToDelete []uuid.UUID
	partsCount := 0
	for _, sp := range section.SectionParts {
		partIDsToDelete = append(partIDsToDelete, sp.PartID)
		partsCount++
	}

	if partsCount == 0 {
		return nil
	}

	errCode := d.songProcessor.UpdateSongAfterPartsDeletion(d.txSongRepo, section.SongID, partIDsToDelete)
	if errCode != nil {
		return errCode
	}

	if err := d.txSongPartRepo.Delete(partIDsToDelete); err != nil {
		return httperror.DatabaseError(err)
	}
	return nil
}
