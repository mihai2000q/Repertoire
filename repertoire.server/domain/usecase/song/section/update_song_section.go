package section

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
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
	if reflect.ValueOf(section).IsZero() {
		return wrapper.NotFoundError(errors.New("song section not found"))
	}

	section.Name = request.Name
	section.Confidence = request.Confidence
	section.Rehearsals = request.Rehearsals
	section.SongSectionTypeID = request.TypeID

	err = c.songRepository.UpdateSection(&section)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
