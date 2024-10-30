package section

import (
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type DeleteSongSection struct {
	songRepository repository.SongRepository
}

func NewDeleteSongSection(repository repository.SongRepository) DeleteSongSection {
	return DeleteSongSection{
		songRepository: repository,
	}
}

func (d DeleteSongSection) Handle(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	var song model.Song
	err := d.songRepository.GetWithSections(&song, songID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = d.songRepository.DeleteSection(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	errCode := d.reorder(id, song.Sections)
	if errCode != nil {
		return errCode
	}

	return nil
}

func (d DeleteSongSection) reorder(id uuid.UUID, sections []model.SongSection) *wrapper.ErrorCode {
	indexOfSection := -1
	for i, section := range sections {
		if indexOfSection != -1 {
			section.Order = section.Order - 1
			err := d.songRepository.UpdateSection(&section)
			if err != nil {
				return wrapper.InternalServerError(err)
			}
		}
		if section.ID == id {
			indexOfSection = i
		}
	}

	return nil
}
