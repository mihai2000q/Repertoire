package part

import (
	"errors"
	"reflect"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
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

func (d DeleteSongPart) Handle(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	var errCode *wrapper.ErrorCode
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
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (d DeleteSongPart) updateSong(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	var song model.Song
	err := d.txSongRepository.GetWithParts(&song, songID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	index := slices.IndexFunc(song.Parts, func(a model.SongPart) bool {
		return a.ID == id
	})
	if index == -1 {
		return wrapper.NotFoundError(errors.New("song part not found"))
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

	err = d.txSongRepository.UpdateWithAssociations(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (d DeleteSongPart) updateSections(id uuid.UUID) *wrapper.ErrorCode {
	var sections []model.SongSection
	err := d.txSongSectionRepository.GetAllByPartWithSectionParts(&sections, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	if len(sections) == 0 {
		return nil
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
		err = d.txSongSectionRepository.UpdateAllSectionParts(&sectionPartsToUpdate)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
	}

	return nil
}
