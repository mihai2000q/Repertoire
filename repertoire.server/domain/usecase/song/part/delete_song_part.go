package part

import (
	"errors"
	"reflect"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
	"slices"

	"github.com/google/uuid"
)

type DeleteSongPart struct {
	transactionManager      transaction.Manager
	txSongRepository        repository.SongRepository
	txSongSectionRepository repository.SongSectionRepository
	txSongPartRepository    repository.SongPartRepository
}

func NewDeleteSongPart(
	transactionManager transaction.Manager,
) DeleteSongPart {
	return DeleteSongPart{
		transactionManager: transactionManager,
	}
}

func (d DeleteSongPart) Handle(id uuid.UUID, songID uuid.UUID) *httperror.ErrorCode {
	var errCode *httperror.ErrorCode
	err := d.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		d.txSongRepository = factory.NewSongRepository()
		d.txSongSectionRepository = factory.NewSongSectionRepository()
		d.txSongPartRepository = factory.NewSongPartRepository()

		if errCode = d.updateSong(id, songID); errCode != nil {
			return errCode.Error
		}
		if errCode = d.updateSections(id); errCode != nil {
			return errCode.Error
		}
		if err := d.txSongPartRepository.Delete([]uuid.UUID{id}); err != nil {
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

func (d DeleteSongPart) updateSong(id uuid.UUID, songID uuid.UUID) *httperror.ErrorCode {
	var song model.Song
	if err := d.txSongRepository.GetWithParts(&song, songID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}

	index := slices.IndexFunc(song.Parts, func(a model.SongPart) bool {
		return a.ID == id
	})
	if index == -1 {
		return httperror.NotFoundError(errors.New("song part not found"))
	}

	// reorder the other parts in song
	partsLength := len(song.Parts)
	for i := index + 1; i < partsLength; i++ {
		song.Parts[i].SongOrder = song.Parts[i].SongOrder - 1
	}

	// update song's new confidence, rehearsals and progress medians
	if partsLength == 1 {
		song.Confidence = 0
		song.Rehearsals = 0
		song.Progress = 0
	} else {
		song.Confidence = (song.Confidence*float64(partsLength) - float64(song.Parts[index].Confidence)) / float64(partsLength-1)
		song.Rehearsals = (song.Rehearsals*float64(partsLength) - float64(song.Parts[index].Rehearsals)) / float64(partsLength-1)
		song.Progress = (song.Progress*float64(partsLength) - float64(song.Parts[index].Progress)) / float64(partsLength-1)
	}

	if err := d.txSongRepository.UpdateWithAssociations(&song); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}

func (d DeleteSongPart) updateSections(id uuid.UUID) *httperror.ErrorCode {
	var sections []model.SongSection
	if err := d.txSongSectionRepository.GetAllByPartWithSectionParts(&sections, id); err != nil {
		return httperror.DatabaseError(err)
	}

	var sectionPartsToUpdate []model.SongSectionPart
	for _, section := range sections {
		index := slices.IndexFunc(section.SectionParts, func(sp model.SongSectionPart) bool {
			return sp.PartID == id
		})

		// reorder the other parts in section
		for i := index + 1; i < len(section.SectionParts); i++ {
			section.SectionParts[i].Order = section.SectionParts[i].Order - 1
			sectionPartsToUpdate = append(sectionPartsToUpdate, section.SectionParts[i])
		}
	}

	if len(sectionPartsToUpdate) > 0 {
		if err := d.txSongSectionRepository.UpdateAllSectionParts(&sectionPartsToUpdate); err != nil {
			return httperror.DatabaseError(err)
		}
	}

	return nil
}
