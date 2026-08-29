package section

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type UpdateSongSection struct {
	songSectionRepository repository.SongSectionRepository
	songRepository        repository.SongRepository
	progressProcessor     processor.ProgressProcessor
}

func NewUpdateSongSection(
	songSectionRepository repository.SongSectionRepository,
	songRepository repository.SongRepository,
	progressProcessor processor.ProgressProcessor,
) UpdateSongSection {
	return UpdateSongSection{
		songSectionRepository: songSectionRepository,
		songRepository:        songRepository,
		progressProcessor:     progressProcessor,
	}
}

func (u UpdateSongSection) Handle(request requests.UpdateSongSectionRequest) *wrapper.ErrorCode {
	var section model.SongSection
	err := u.songSectionRepository.GetWithParts(&section, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(section).IsZero() {
		return wrapper.NotFoundError(errors.New("song section not found"))
	}

	section.Name = request.Name
	section.SongSectionTypeID = request.TypeID

	err = u.songSectionRepository.Update(&section)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
