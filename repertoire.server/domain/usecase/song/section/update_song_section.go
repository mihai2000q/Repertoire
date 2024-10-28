package section

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/model"
	"repertoire/utils/wrapper"
)

type UpdateSongSection struct {
	songRepository repository.SongRepository
}

func NewUpdateSongSection(repository repository.SongRepository) UpdateSongSection {
	return UpdateSongSection{
		songRepository: repository,
	}
}

func (c UpdateSongSection) Handle(request requests.UpdateSongSectionRequest) *wrapper.ErrorCode {
	var section model.SongSection
	err := c.songRepository.GetSection(&section, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if section.ID == uuid.Nil {
		return wrapper.NotFoundError(errors.New("song section not found"))
	}

	section.Name = request.Name
	section.Rehearsals = request.Rehearsals
	section.SongSectionTypeID = request.TypeID

	err = c.songRepository.UpdateSection(&section)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
